package main

import (
	"encoding/json"
	"net/http"
	"time"

	"skillswap/backend/authentication"
	"skillswap/backend/config"
	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/mux"
)

func search(w http.ResponseWriter, req *http.Request) {
	db, err := database.GetDatabase()
	if err != nil {
		utils.HandleError(err)
		return
	}

	var requestBody struct {
		Query string `json:"query"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		utils.HandleError(err)
		return
	}

	searchQuery := "%" + requestBody.Query + "%"

	// Use a corrected and more efficient SQL query with GROUP BY
	stmt, err := db.Prepare(`
        SELECT
            u.id,
            u.username,
            u.email,
            GROUP_CONCAT(s.name SEPARATOR ', ') AS skills_found
        FROM users AS u
        JOIN user_skills AS us ON u.id = us.user_id
        JOIN skills AS s ON us.skill_id = s.id
        WHERE u.username LIKE ? OR u.email LIKE ? OR s.name LIKE ? OR s.description LIKE ?
        GROUP BY u.id, u.username, u.email
    `)
	if err != nil {
		utils.HandleError(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		utils.HandleError(err)
		return
	}
	defer rows.Close()

	// Define a new struct to hold the combined skills
	var results []struct {
		ID          int64  `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		SkillsFound string `json:"skills_found"`
	}

	for rows.Next() {
		var r struct {
			ID          int64  `json:"id"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			SkillsFound string `json:"skills_found"`
		}
		if err := rows.Scan(&r.ID, &r.Username, &r.Email, &r.SkillsFound); err != nil {
			utils.HandleError(err)
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		utils.HandleError(err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, results)
}

func main() {
	database.Init()
	// Izveido jaunu rūteri ar stingru pārbaudi slīpsvītrām, kas nozīmē, ka maršruti ar un bez beigu slīpsvītras tiek uzskatīti par atšķirīgiem.
	server := mux.NewRouter().StrictSlash(true)

	// Tiek definēti API ceļi (end-points) dažādām front-end darbībām.
	// "HandleFunc" piesaista konkrētu URL ceļu noteiktai Go funkcijai.
	// server.HandleFunc("/api/chat", RunWebsocket)
	server.HandleFunc("/api/login", authentication.Login).Methods("POST")
	server.HandleFunc("/api/register", authentication.Register).Methods("POST")
	server.HandleFunc("/api/logout", authentication.Logout).Methods("POST")
	// Vienkārša "dummy" funkcija aizmugursistēmas (backend) darbības pārbaudei.
	// Tā atgriež JSON atbildi ar statusu "pong", kad tiek saņemts GET pieprasījums.
	server.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "pong"})
	}).Methods("GET")

	// Tiek konfigurēts CORS (Cross-Origin Resource Sharing), lai atļautu pieprasījumus
	// no noteiktiem "oriģiniem" (front-end servera adresēm).
	c := config.CORS()

	// Piesaista CORS apstrādātāju rūterim, lai tas darbotos ar visiem ceļiem.
	server.Use(c.Handler)

	// Reģistrē rūteri kā galveno HTTP apstrādātāju.
	http.Handle("/", server)

	// Izveido un konfigurē HTTP serveri.
	serve := &http.Server{
		Handler: server,           // Norāda, ka šis serveris izmantos mūsu iepriekš definēto rūteri.
		Addr:    "localhost:8080", // Serveris klausīsies uz šīs adreses un porta.

		// Iestata maksimālo laiku, cik ilgi serveris gaidīs pieprasījuma ķermeņa nolasīšanu un atbildes uzrakstīšanu.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Sāk servera darbību. "ListenAndServe" bloķē izpildi, līdz serveris tiek apturēts.
	serve.ListenAndServe()
}
