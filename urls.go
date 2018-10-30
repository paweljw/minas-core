package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

// URL stores data about a single URL.
type URL struct {
	ID       int64  `json:"id,omitempty"`
	Location string `json:"location,omitempty"`
	Active   bool   `json:"active,omitempty"`
}

// Validate validates a single URL.
func (u URL) Validate() (string, bool) {
	var count int
	if err := db.Get(&count, "SELECT count(*) FROM urls WHERE location = $1", u.Location); err != nil {
		log.Print(err)
		return "Database error", false
	}

	if count > 0 {
		return "URL is not unique", false
	}

	re := regexp.MustCompile(`https?:\/\/.*`)
	if !re.MatchString(u.Location) {
		return "URL is malformed", false
	}

	return "", true
}

// CreateURLHandler allows creation of URLs.
var CreateURLHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var url URL

	if err := json.NewDecoder(req.Body).Decode(&url); err != nil {
		RespondBadRequest(w, Response{Message: "Error", Status: "Could not parse body"})
		return
	}

	message, valid := url.Validate()

	if valid {
		tx := db.MustBegin()

		if _, err := tx.NamedExec("INSERT INTO urls (location, active) VALUES (:location, :active)", &url); err != nil {
			log.Print(err)
			RespondError(w, Response{Status: "Error", Message: "Database error"})
		}

		if err := tx.Commit(); err != nil {
			log.Print(err)
			RespondError(w, Response{Status: "Error", Message: "Database error"})
		}

		RespondOk(w, Response{Status: "OK", Message: "OK"})
	} else {
		RespondBadRequest(w, Response{Status: "Error", Message: message})
	}
})

// DeactivateURLHandler allows deactivation of a single URL.
var DeactivateURLHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	tx := db.MustBegin()

	id, _ := strconv.Atoi(params["id"])

	if _, err := tx.NamedExec("UPDATE urls SET active = false WHERE id = :id", &URL{ID: int64(id)}); err != nil {
		log.Print(err)
		RespondError(w, Response{Status: "Error", Message: "Database error"})
	}

	if err := tx.Commit(); err != nil {
		log.Print(err)
		RespondError(w, Response{Status: "Error", Message: "Database error"})
	}

	RespondOk(w, Response{Status: "OK", Message: "OK"})
})

// ActivateURLHandler allows deactivation of a single URL.
var ActivateURLHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	tx := db.MustBegin()

	id, _ := strconv.Atoi(params["id"])

	if _, err := tx.NamedExec("UPDATE urls SET active = true WHERE id = :id", &URL{ID: int64(id)}); err != nil {
		log.Print(err)
		RespondError(w, Response{Status: "Error", Message: "Database error"})
	}

	if err := tx.Commit(); err != nil {
		log.Print(err)
		RespondError(w, Response{Status: "Error", Message: "Database error"})
	}

	RespondOk(w, Response{Status: "OK", Message: "OK"})
})
