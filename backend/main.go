package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func search(w http.ResponseWriter, req *http.Request) {
	db, err := getDatabase()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database connection failed"})
		return
	}
	defer db.Close()

	// Parse the JSON request body into a struct for type safety
	var requestBody struct {
		Query string `json:"query"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Add wildcards to the query string for 'LIKE' search
	searchQuery := "%" + requestBody.Query + "%"

	// Use a single, correct SQL query
	stmt, err := db.Prepare(`
		SELECT u.username, u.email, s.name AS skill_name, s.description AS skill_description
		FROM users AS u
		JOIN user_skills AS us ON u.id = us.user_id
		JOIN skills AS s ON us.skill_id = s.id
		WHERE u.username LIKE ? OR u.email LIKE ? OR s.name LIKE ? OR s.description LIKE ?
	`)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Query execution failed"})
		return
	}
	defer rows.Close()

	// Use a slice of structs to store the results
	var results []struct {
		Username         string `json:"username"`
		Email            string `json:"email"`
		SkillName        string `json:"skill_name"`
		SkillDescription string `json:"skill_description"`
	}

	for rows.Next() {
		var r struct {
			Username         string `json:"username"`
			Email            string `json:"email"`
			SkillName        string `json:"skill_name"`
			SkillDescription string `json:"skill_description"`
		}
		if err := rows.Scan(&r.Username, &r.Email, &r.SkillName, &r.SkillDescription); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Failed to scan row"})
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Error reading rows"})
		return
	}

	// Send the results as a JSON response
	sendJSONResponse(w, http.StatusOK, results)
}

func main() {
	server := mux.NewRouter().StrictSlash(true)
	server.HandleFunc("/api/chat", RunWebsocket)
	server.HandleFunc("/api/login", login).Methods("POST")
	server.HandleFunc("/api/register", register).Methods("POST")
	server.HandleFunc("/api/logout", logout).Methods("POST")
	server.HandleFunc("/api/cookieUser", getCookieUser).Methods("GET")
	server.HandleFunc("/api/isEmailUsed", isEmailUsed).Methods("POST")

	server.HandleFunc("/api/search", search).Methods("POST")

	server.HandleFunc("/api/video", handleWebsocket)

	server.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
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
