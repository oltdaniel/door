# door :door:

through door to numbers

# about

Count your views on your website without JavaScript and privacy focused. No
tracking just counting.

<pre><code>Clone repository.
$ git clone https://github.com/oltdaniel/door.git
Enter repository.
$ cd door
Install dependencies.
$ ./install.sh
Start server.
$ go run main.go</code></pre>

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
