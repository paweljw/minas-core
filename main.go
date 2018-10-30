package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func main() {
	log.Print("Booting")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}

	db = sqlx.MustConnect("postgres", viper.Get("connection").(string))
	db.MustExec(schema)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ping", PingEndpoint).Methods("GET")
	router.Handle("/api/v1/urls", authMiddleware(CreateURLHandler)).Methods("POST")
	router.Handle("/api/v1/urls/{id}/deactivate", authMiddleware(DeactivateURLHandler)).Methods("PUT")
	router.Handle("/api/v1/urls/{id}/activate", authMiddleware(ActivateURLHandler)).Methods("PUT")
	router.Handle("/api/v1/urls/{id}/scans", authMiddleware(ScansHandler)).Methods("GET")

	log.Print("Finished bootstrap")

	c := make(chan bool)

	go ScanURLs(c)
	go func() {
		log.Fatal(http.ListenAndServe(":9000", handlers.LoggingHandler(os.Stdout, router)))
	}()

	for range c {
		go func() {
			time.Sleep(time.Minute)
			ScanURLs(c)
		}()
	}
}
