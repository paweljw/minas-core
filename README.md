# minas-core

A webservice which periodically checks whether an HTTP(S) resource is up or down, and saves the result in a database. It exposes a very simple API
for interrogating the checked resources.

Note that this is a learning project and not intended for any serious use.

## Configuration

Copy `config.example.yml` to `config.yml` and fill in the blanks. You will need an Auth0 API set up for authentication.

## Running in "production"

First let's preface this section with: **DON'T**.

Hell-bent on your own destruction? All right:

* Create `config.production.yml` with appropriate values (remember to point database host to "postgresql" if using the attached `docker-compose.production.yml`)
* Run `docker-compose -f docker-compose.production.yml start`
* Point your favorite reverse proxy at port 9000

## TODO

* [ ] Detect whether the service is actually down, or just has a bad SSL cert
* [ ] Timeout for URL checks, ideally under a minute
* [ ] Configurable check interval
* [ ] Figure out a good way to test an API webapp in Go

## License

This work is shared under [DBAD License](http://dbad-license.org/). See [`LICENSE`](LICENSE).