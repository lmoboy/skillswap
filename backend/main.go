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
	server.HandleFunc("/api/logout", logout).Methods("POST")
	server.HandleFunc("/api/cookieUser", getCookieUser).Methods("GET")
	server.HandleFunc("/api/isEmailUsed", isEmailUsed).Methods("POST")


	server.HandleFunc("/ws/videoCall", handleWebsocket)

	server.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		sendJSONResponse(w, http.StatusOK, map[string]string{"status": "pong"})
	}).Methods("GET")
	
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Set-Cookie",
		},
		Debug: false,
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
