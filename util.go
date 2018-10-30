package main

import (
	"encoding/json"
	"net/http"
)

// Response is a generic HTTP response.
type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// Respond is a generic helper to set up content type.
func Respond(w http.ResponseWriter, r Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&r)
}

// RespondBadRequest responds with a HTTP status of Bad Request and a given Response.
func RespondBadRequest(w http.ResponseWriter, r Response) {
	Respond(w, r, http.StatusBadRequest)
}

// RespondOk responds with a HTTP status of OK and a given Response.
func RespondOk(w http.ResponseWriter, r Response) {
	Respond(w, r, http.StatusOK)
}

// RespondError responds with a HTTP status of Error and a given Response.
func RespondError(w http.ResponseWriter, r Response) {
	Respond(w, r, http.StatusInternalServerError)
}
