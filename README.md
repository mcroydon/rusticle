# Rusticle

A collection of tools for poking at data inside [Rust](http://playrust.com), starting with...

## Cacheinspect

Take a peek at all of the horrible horrible things that humans write on signs.

Includes a simple http server in go and a simple React frontend that loads data over AJAX using JQuery.

To run:

	$ go get github.com/mcroydon/rusticle/cacheinspect
	$ $GOPATH/bin/cacheinspect
	2016/01/10 23:36:08 Server running on :8888

You can then point your web browser at `http://localhost:8888`. Data is served over json at `http://localhost:8888/data` and individual images are served at `http://localhost:8888/img?entity=<entity>&crc=<crc>`.

### To Do

* Exhaustively check possible steam directory locations rather than be fancy with runtime.GOOS.
* Allow for env var or cmdline argument for Steam location in case it's not in the standard spot.
* Don't use panic, let errors bubble up and present useful/actionable information.

# License

Rusticle is released under the MIT license.
