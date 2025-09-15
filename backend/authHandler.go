package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// +++++++++++++++++structs+++++++++++++++++
type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// User represents the data returned to the client
type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Global session store
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// +++++++++++++++++ Helper function for JSON responses +++++++++++++++++
func sendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// +++++++++++++++++ Helper function for applying sessions +++++++++++++++++
func applySession(w http.ResponseWriter, req *http.Request, user *Info) {
	session, err := store.Get(req, "authentication")

	if err != nil && !session.IsNew {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
		return
	}

	session.Options = &sessions.Options{HttpOnly: false, SameSite: http.SameSiteDefaultMode, Secure: false, MaxAge: 2628000, Path: "/"}

	// Use a less expensive unique ID for the session
	session.ID = fmt.Sprintf("%x", md5.Sum([]byte(user.Email+time.Now().String()+string(securecookie.GenerateRandomKey(16)))))
	session.Values["Authenticated"] = true
	session.Values["email"] = user.Email // Store email in the session

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer db.Close()

	if err := session.Save(req, w); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
	}

	req.Cookie("authentication")
	stmt, err := db.Prepare("INSERT INTO sessions (user_id, session_token, expires_at) VALUES ((SELECT id FROM users WHERE email = ?), ?, DATE_ADD(NOW(), INTERVAL 1 MONTH)) ON DUPLICATE KEY UPDATE expires_at = DATE_ADD(NOW(), INTERVAL 1 MONTH)")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	_, err = stmt.Exec(user.Email, session.ID)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	fmt.Println("Session:", session)
}

// ++++++++++++++ Register Handler ++++++++++++++
func register(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if userInfo.Username == "" || userInfo.Email == "" || userInfo.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Inputs can't be empty"})
		return
	}

	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password)))

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (user_id, username, email, password_hash) VALUES (?, ?, ?, ?)")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer stmt.Close()

	id := fmt.Sprintf("%s%s%s", strings.ToUpper(strings.Replace(uuid.New().String(), "-", "", -1)), strings.ToLower(strings.Replace(uuid.New().String(), "-", "", -1)), strings.Replace(uuid.New().String(), "-", "", -1))
	fmt.Println(userInfo)
	_, err = stmt.Exec(id, userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		sendJSONResponse(w, http.StatusConflict, map[string]string{"error": "Username or email already exists", "message": err.Error()})
		return
	}

	applySession(w, req, &userInfo)

	sendJSONResponse(w, http.StatusCreated, map[string]string{"status": "ok", "message": "Registration successful"})
}

func isEmailUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	rows, err := findValues("users", []string{"email"}, map[string]string{"email": userInfo.Email})
	fmt.Println(rows)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Internal server error: ERR-0010"})
		return
	}
	if len(rows) > 0 {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email already in use"})
		return
	} else {
		sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Email available"})
		return
	}
}
func isUsernameUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	rows, err := findValues("users", []string{"name"}, map[string]string{"name": userInfo.Username})
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Internal server error: ERR-0010"})
		return
	}
	if len(rows) > 0 {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username already in use"})
		return
	} else {
		sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Username available"})
		return
	}
}

// ++++++++++++++ Login Handler ++++++++++++++
func login(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT username, email FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))

	var storedUsername, storedEmail string
	if err := row.Scan(&storedUsername, &storedEmail); err != nil {
		if err == sql.ErrNoRows {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	applySession(w, req, &Info{Username: storedUsername, Email: storedEmail})

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}

func getCookieUser(w http.ResponseWriter, req *http.Request) {

	session, err := store.Get(req, "authentication")
	if err != nil {
		fmt.Println("Error getting session:", err)
		// sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
	}

	auth, ok := session.Values["Authenticated"].(bool)
	if !ok || !auth {
		fmt.Println("Unauthorized :", err)

		// sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Not authenticated"})
	}

	email, ok := session.Values["email"].(string)
	if !ok || email == "" {

		fmt.Println("Where the email:", err)
		// sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session: Missing email"})
	}

	cookie, err := req.Cookie("authentication")
	if err != nil {
		if err == http.ErrNoCookie {

			fmt.Println("Where is the cookie :", err)
			// sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: No cookie found"})
		}

		fmt.Println("Error getting session cookie:", err)
		// sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error reading cookie"})
	}
	if session != nil || auth || email != "" || cookie != nil {
		var session_value string

		// fmt.Println("Session:", session, "Auth:", auth, "Email:", email, "Cookie:", cookie)
		if session != nil {
			session_value = session.ID
		} else if auth {
			session_value = email
		} else if cookie != nil {
			session_value = cookie.String()
		} else {
			session_value = ""
		}

		// fmt.Println(cookie)
		fmt.Printf("Session: %v", session_value)
		// fmt.Println("Decided on:", session_value)
		// sendJSONResponse(w, http.StatusOK, map[string]string{"cookie": cookie.String()})
		// return
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database connection failed"})
		return
	}
	defer db.Close()

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: User not found"})
			return
		}
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database query failed"})
		return
	}

	applySession(w, req, &Info{Username: username, Email: email})
	sendJSONResponse(w, http.StatusOK, UserData{Username: username, Email: email})
}

func logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "authentication")
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request: No session found"})
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session save failed"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database connection failed"})
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE session_token = ?", session.ID)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session deletion failed"})
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})

}
