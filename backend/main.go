package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	server := mux.NewRouter().StrictSlash(true)
	server.HandleFunc("/chat", RunWebsocket)
	server.HandleFunc("/api/login", login).Methods("POST")
	server.HandleFunc("/api/register", register).Methods("POST")
	server.HandleFunc("/api/cookieUser", getCookieUser).Methods("GET")
	server.HandleFunc("/api/isEmailUsed", isEmailUsed).Methods("POST")
	server.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		sendJSONResponse(w, http.StatusOK, map[string]string{"status": "pong"})
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
