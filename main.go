package main

// load external
import "github.com/valyala/bytebufferpool"
import "github.com/willf/bloom"
import "golang.org/x/crypto/blake2b"

// load standard
import "net"
import "net/http"
import "fmt"
import "time"

// address to listen on
const PORT string = "0.0.0.0:8080"
// the maximum number of calls to store
const N    uint   = 1000000

// IP masks
var IPV4Mask []byte = []byte{255,255,255,0}
var IPV6Mask []byte = []byte{255,255,255,255,255,255,255,255,0,0,0,0,0,0,0,0}

// buffered channel of 256 items
var queue chan *Call = make(chan *Call, 256)
// stats of the calls stored
var stats map[string]uint = make(map[string]uint)
var statsTxt string = ""

// abstract the call meta data
type Call struct {
    Hash [blake2b.Size256]byte // hash of the metadata
    Path string                // full refrer
}

// entry point
func main() {
    // start consumer
    go consumer()
    // print stats in interval
    go func() {
        // endless loop
        for {
            // count total calls
            var total uint = 0
            // get a new buffer
            b := bytebufferpool.Get()
            // print banner
            b.WriteString(fmt.Sprintf("\n%s\n\n", time.Now().String()))
            // for each path
            for k,v := range stats {
                // print path
                b.WriteString(fmt.Sprintf("%s has %d call(s)\n", k, v))
                // update total
                total += v
            }
            // write total stats
            b.WriteString(fmt.Sprintf("\ntotal: %d\n", total))
            // update global
            statsTxt = b.String()
            // put back buffer
            bytebufferpool.Put(b)
            // sleep for 10s
            time.Sleep(10 * time.Second)
        }
    }()
    // assign handler
    http.HandleFunc("/style.css", handle)
    http.HandleFunc("/stats", handleStats)
    // status message
    fmt.Printf("[DEBG] server started (%s)\n", PORT)
    // start server
    err := http.ListenAndServe(PORT, nil)
    // throw errir
    if err != nil {
        panic(err)
    }
}

// anonymize remote
func anonymize(ip string) string {
    // get host only
    host, _, err := net.SplitHostPort(ip)
    // check for error (should not happen)
    if err != nil {
        return ""
    }
    // parse ip
    ipaddr := net.ParseIP(host)
    // mask addresses
    if ipaddr.To4() != nil {
        // return masked ipv4 string
        return ipaddr.Mask(IPV4Mask).String()
    } else {
        // return masked ipv6 string
        return ipaddr.Mask(IPV6Mask).String()
    }
}

// handle each request
func handle(w http.ResponseWriter, r *http.Request) {
    // concat values
    b := bytebufferpool.Get()
    b.WriteString(r.Referer())
    b.WriteString(anonymize(r.RemoteAddr))
    b.WriteString(r.UserAgent())
    // hash value
    h := blake2b.Sum256(b.Bytes())
    // put buffer back
    bytebufferpool.Put(b)
    // put to queue
    queue <- &Call{h, r.Referer()}
    // set content type
    w.Header().Set("Content-Type", "text/css")
    // response with not modified status code
    w.WriteHeader(304)
}

// return stats
func handleStats(w http.ResponseWriter, r *http.Request) {
    // set content type
    w.Header().Set("Content-Type", "text/plain")
    // response with ok status code
    w.WriteHeader(200)
    // write content
    w.Write([]byte(statsTxt))
}

// consume new calls and filter them for uniqueness
func consumer() {
    // create filter
    f := bloom.New(20*N, 5)
    // consume queue
    for h := range queue {
        // check if call exists
        if !f.Test(h.Hash[:]) {
            // filter new call
            f.Add(h.Hash[:])
            // update
            stats[h.Path]++
        }
    }
}
