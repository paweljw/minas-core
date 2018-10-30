package main

import (
	"net/http"
)

// PingEndpoint responds to pings over HTTP.
func PingEndpoint(w http.ResponseWriter, req *http.Request) {
	RespondOk(w, Response{Status: "OK", Message: "OK"})
}
