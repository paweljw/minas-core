package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ScanURLs pulls out data from the database and scans active URLs.
func ScanURLs(channel chan bool) {
	var urls []URL

	if err := db.Select(&urls, "SELECT * FROM urls WHERE active = true"); err != nil {
		log.Print(err)
	}

	c := make(chan int)

	for _, url := range urls {
		go ScanURL(c, url)
	}

	for i := 0; i < len(urls); i++ {
		<-c
	}

	channel <- true
}

// Scan is a single data point about a URL.
type Scan struct {
	ID       int64     `json:"id,omitempty"`
	URLID    int64     `db:"url_id" json:"url_id,omitempty"`
	Up       bool      `json:"up,omitempty"`
	Duration int64     `json:"duration,omitempty"`
	Time     time.Time `json:"time,omitempty"`
}

// ScanURL scans a particular URL and saves response to DB.
func ScanURL(c chan int, u URL) {
	start := time.Now()
	link := u.Location

	_, err := http.Get(link)
	t := time.Now()
	elapsed := t.Sub(start)

	scan := Scan{URLID: u.ID, Up: err == nil, Duration: elapsed.Nanoseconds() / 1000000}

	tx := db.MustBegin()
	if _, err := tx.NamedExec("INSERT INTO scans (url_id, up, duration) VALUES (:url_id, :up, :duration)", &scan); err != nil {
		log.Print(err)
	}
	if err := tx.Commit(); err != nil {
		log.Print(err)
	}

	c <- 1
}

// ScansHandler serves data about scans for one URL.
var ScansHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])
	last, _ := strconv.Atoi(req.URL.Query().Get("last"))
	since := req.URL.Query().Get("since")

	if last == 0 {
		last = 60
	}

	scans := []Scan{}

	if since != "" {
		if err := db.Select(&scans, "SELECT * FROM scans WHERE url_id = $1 AND time >= $2 ORDER BY time DESC", id, since); err != nil {
			log.Print(err)
			RespondError(w, Response{Status: "Error", Message: "Database error"})
		}
	} else {
		if err := db.Select(&scans, "SELECT * FROM scans WHERE url_id = $1 ORDER BY time DESC LIMIT $2", id, last); err != nil {
			log.Print(err)
			RespondError(w, Response{Status: "Error", Message: "Database error"})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scans)
})
