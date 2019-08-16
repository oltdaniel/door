# door :door:

through door to numbers

# about

Count your views on your website without JavaScript and privacy focused. No
tracking just counting.

```bash
# Clone repository.
$ git clone https://github.com/oltdaniel/door.git
# Enter repository.
$ cd door
# Install dependencies.
$ ./deps.sh
# Start server.
$ go run main.go
# Run to store stats as txt
$ ./store.sh
```

To call this server to count the views, you can simply add a new css file to
your existing website. The average size this request will consume is under 200
bytes.

```html
<link rel="styelsheet" href="http://localhost:8080/style.css">
```

To view your stats call [localhost:8080/stats](http://localhost:8080/stats).

# how

We remove the data that allows us to track down the user from the IP in order
to make it more private. In order to gain some unqiness we add the user agent
of the browser and the path that has been called. All values will be then hashed
with the <code>blake2b</code> hash function.

**The result**: We can distinguish wether this call already has been made, or
if the same 'user' called a path twice, but not if a certain user called another
path.

# license

_just do what you want_

MIT License
