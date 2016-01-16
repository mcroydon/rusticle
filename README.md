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

# To Do

* Update cache file path to work with `GOOS` other than `darwin`.
* Move to a separate repo.
* Find a better (silly) name.

# License

Rusticle is released under t MIT license.