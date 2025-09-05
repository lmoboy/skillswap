package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	server := mux.NewRouter().StrictSlash(true)

	// server.HandleFunc("/api/login", services.login)
	// server.HandleFunc("/api/register", services.register)
	server.HandleFunc("/api/checks", func(w http.ResponseWriter, req *http.Request) {

	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	})

	server.Use(c.Handler)

	http.Handle("/", server)
	serve := &http.Server{
		Handler: server,
		Addr:    "localhost:8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	serve.ListenAndServe()
}
