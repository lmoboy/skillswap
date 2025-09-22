package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func search(w http.ResponseWriter, req *http.Request) {
	db, err := getDatabase()
	if err != nil {
		handleError(err)
		return
	}
	defer db.Close()

	var requestBody struct {
		Query string `json:"query"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		handleError(err)
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
		handleError(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		handleError(err)
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
			handleError(err)
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		handleError(err)
		return
	}

	sendJSONResponse(w, http.StatusOK, results)
}
func getAllSwappers(w http.ResponseWriter, req *http.Request) {

	session, err := store.Get(req, "authentication")
	if err != nil {
		handleError(err)
		return
	}
	db, err := getDatabase()
	if err != nil {
		handleError(err)
		return
	}
	fmt.Println(session)
	stmt, err := db.Prepare(`
		SELECT c.id
FROM chats c
JOIN users u1 ON c.user1_id = u1.id
JOIN users u2 ON c.user2_id = u2.id
WHERE u1.email = ? OR u2.email = ?;
	`)
	if err != nil {
		handleError(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(session.Values["email"], session.Values["email"])
	if err != nil {
		handleError(err)
		return
	}
	defer rows.Close()

	var results []struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	for rows.Next() {
		var r struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}
		if err := rows.Scan(&r.ID, &r.Username, &r.Email); err != nil {
			handleError(err)
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		handleError(err)
		return
	}

	sendJSONResponse(w, http.StatusOK, results)
}
func startChatWithUser(w http.ResponseWriter, req *http.Request) {
	// Get user data from session
	session, err := store.Get(req, "authentication")
	if err != nil {
		handleError(err)
		return
	}

	// Get user email from session
	userEmail, ok := session.Values["email"].(string)
	if !ok {
		handleError(err)
		return
	}

	// Get user who is being requested
	requestBody := struct {
		Email string `json:"email"`
	}{}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		handleError(err)
		return
	}
	db, err := getDatabase()
	// Get user IDs
	user1ID, err := GetUserIDFromEmail(userEmail)
	if err != nil {
		handleError(fmt.Errorf("error getting user1 ID: %v", err))
		return
	}

	user2ID, err := GetUserIDFromEmail(requestBody.Email)
	if err != nil {
		handleError(fmt.Errorf("error getting user2 ID: %v", err))
		return
	}

	// Create a new chat entry
	stmt, err := db.Prepare(`INSERT INTO chats (user1_id, user2_id, initiated_by) VALUES (?, ?, ?)`)
	if err != nil {
		handleError(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user1ID, user2ID, user1ID) // Using user1ID as initiated_by
	if err != nil {
		handleError(err)
		return
	}
}

func getUser(w http.ResponseWriter, req *http.Request) {
	db, err := getDatabase()
	if err != nil {
		handleError(err)
		return
	}
	defer db.Close()

	vars := req.URL.Query().Get("q")
	id, err := strconv.ParseInt(vars, 10, 64)
	if err != nil {
		handleError(err)
		return
	}

	user, err := getUserByID(db, id)
	if err != nil {
		handleError(err)
		return
	}

	sendJSONResponse(w, http.StatusOK, user)
}
func getUserByID(db *sql.DB, id int64) (map[string]interface{}, error) {
	// First, get user data (excluding password_hash)
	users, err := findValues("users", []string{"id", "skills", "location", "projects", "aboutme", "contacts", "username", "email", "created_at", "profile_picture"}, map[string]string{"id": strconv.FormatInt(id, 10)})
	if err != nil {
		handleError(err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := make(map[string]interface{})
	for k, v := range users[0] {
		user[k] = v
	}

	// Get user skills with verification status
	rows, err := db.Query(`
        SELECT s.name, us.verified 
        FROM user_skills us
        JOIN skills s ON us.skill_id = s.id 
        WHERE us.user_id = ?
    `, id)

	if err != nil {
		handleError(err)
		return nil, err
	}
	defer rows.Close()

	var skills []map[string]interface{}
	for rows.Next() {
		var name string
		var verified bool
		if err := rows.Scan(&name, &verified); err != nil {
			handleError(err)
			return nil, err
		}
		skills = append(skills, map[string]interface{}{
			"name":     name,
			"verified": verified,
		})
	}

	if err = rows.Err(); err != nil {
		handleError(err)
		return nil, err
	}

	// Add skills to user data
	user["skills_coding"] = skills
	fmt.Println(user)
	return user, nil
}

func main() {
	// Izveido jaunu rūteri ar stingru pārbaudi slīpsvītrām, kas nozīmē, ka maršruti ar un bez beigu slīpsvītras tiek uzskatīti par atšķirīgiem.
	server := mux.NewRouter().StrictSlash(true)

	// Tiek definēti API ceļi (end-points) dažādām front-end darbībām.
	// "HandleFunc" piesaista konkrētu URL ceļu noteiktai Go funkcijai.
	server.HandleFunc("/api/chat", RunWebsocket)
	server.HandleFunc("/api/login", login).Methods("POST")
	server.HandleFunc("/api/register", register).Methods("POST")
	server.HandleFunc("/api/logout", logout).Methods("POST")
	server.HandleFunc("/api/cookieUser", getCookieUser).Methods("GET")
	server.HandleFunc("/api/isEmailUsed", isEmailUsed).Methods("POST")
	server.HandleFunc("/api/search", search).Methods("POST")
	server.HandleFunc("/api/video", handleVideo)
	server.HandleFunc("/api/user", getUser).Methods("GET")
	server.HandleFunc("/api/getAllSwappers", getAllSwappers).Methods("GET")
	// Vienkārša "dummy" funkcija aizmugursistēmas (backend) darbības pārbaudei.
	// Tā atgriež JSON atbildi ar statusu "pong", kad tiek saņemts GET pieprasījums.
	server.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		sendJSONResponse(w, http.StatusOK, map[string]string{"status": "pong"})
	}).Methods("GET")

	// Tiek konfigurēts CORS (Cross-Origin Resource Sharing), lai atļautu pieprasījumus
	// no noteiktiem "oriģiniem" (front-end servera adresēm).
	c := getCors()

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
