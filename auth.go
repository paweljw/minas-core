package main

import (
	"log"
	"net/http"

	"github.com/auth0-community/auth0"
	"github.com/spf13/viper"
	jose "gopkg.in/square/go-jose.v2"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(viper.Get("secret").(string))
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{viper.Get("audience").(string)}

		configuration := auth0.NewConfiguration(secretProvider, audience, viper.Get("domain").(string), jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			log.Print(err)
			log.Print("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
