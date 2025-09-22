package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// +++++++++++++++++structs+++++++++++++++++
type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	id       int    `json:"id"`
}

// User represents the data returned to the client
type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Id       int    `json:"id"`
}

// Global session store
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

// +++++++++++++++++ Helper function for JSON responses +++++++++++++++++

// +++++++++++++++++ Helper function for applying sessions +++++++++++++++++
func applySession(w http.ResponseWriter, req *http.Request, user *Info) {
	// First check if there's an existing session
	session, err := store.Get(req, "authentication")
	if err != nil {
		handleError(err)
		// If we can't get the session, create a new one
		session, err = store.New(req, "authentication")
		if err != nil {
			handleError(err)
		}
	}

	// Check if there's an existing authentication cookie
	cookie, err := req.Cookie("authentication")
	if err == nil && cookie != nil && cookie.Raw != "" {
		// If cookie exists, check if the session is valid
		row, err := findValues("sessions", []string{"session_token"}, map[string]string{"session_token": cookie.Raw})
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-66"})
			return
		}
		for _, row := range row {
			if row["session_token"] == cookie.Raw {
				sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-71"})
				return
			}
		}
	}

	if err != nil && !session.IsNew {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-78"})
		return
	}

	session.Options = &sessions.Options{HttpOnly: false, SameSite: http.SameSiteDefaultMode, Secure: false, MaxAge: 2628000, Path: "/"}

	// Use a less expensive unique ID for the session
	session.ID = fmt.Sprintf("%x", md5.Sum([]byte(user.Email+time.Now().String()+string(securecookie.GenerateRandomKey(16)))))
	session.Values["Authenticated"] = true
	session.Values["email"] = user.Email

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-91"})
		return
	}
	defer db.Close()
	// fmt.Printf("session: %v\n", session)
	if err := session.Save(req, w); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-97"})
	}

	// Get the session ID that was just saved
	sessionToken := session.ID

	// First, try to update existing session
	updateStmt, err := db.Prepare(
		`UPDATE sessions 
		SET session_token = ?, expires_at = DATE_ADD(NOW(), INTERVAL 1 MONTH) 
		WHERE user_id = (SELECT id FROM users WHERE email = ?)`)
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-109"})
		return
	}
	defer updateStmt.Close()

	result, err := updateStmt.Exec(sessionToken, user.Email)
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-116"})
		return
	}

	// If no rows were updated, insert a new session
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		insertStmt, err := db.Prepare(
			`INSERT INTO sessions (user_id, session_token, expires_at) 
			VALUES ((SELECT id FROM users WHERE email = ?), ?, DATE_ADD(NOW(), INTERVAL 1 MONTH))`)
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-127"})
			return
		}
		defer insertStmt.Close()

		_, err = insertStmt.Exec(user.Email, sessionToken)
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-134"})
			fmt.Println("here?")
			return
		}
	}

}

// ++++++++++++++ Register Handler ++++++++++++++
func register(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-146"})
		return
	}

	if userInfo.Username == "" || userInfo.Email == "" || userInfo.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-151"})
		return
	}

	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password)))

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-159"})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users ( username, email, password_hash) VALUES (?, ?, ?)")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-166"})
		return
	}
	defer stmt.Close()

	// fmt.Println(userInfo)
	_, err = stmt.Exec(userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusConflict, map[string]string{"error": "AH-174", "message": err.Error()})
		return
	}

	applySession(w, req, &userInfo)

	sendJSONResponse(w, http.StatusCreated, map[string]string{"status": "ok", "message": "Registration successful"})
}

func isEmailUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-186"})
		return
	}
	rows, err := findValues("users", []string{"email"}, map[string]string{"email": userInfo.Email})
	// fmt.Println(rows)
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-192"})
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
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-206"})
		return
	}
	rows, err := findValues("users", []string{"name"}, map[string]string{"name": userInfo.Username})
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-211"})
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
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-227"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-233"})
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT username, email, id FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))
	var storedID int
	var storedUsername, storedEmail string
	if err := row.Scan(&storedUsername, &storedEmail, &storedID); err != nil {
		if err == sql.ErrNoRows {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-246"})
		return
	}

	applySession(w, req, &Info{Username: storedUsername, Email: storedEmail, id: storedID})

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}

func getCookieUser(w http.ResponseWriter, req *http.Request) {
	// First try to get the session
	session, err := store.Get(req, "authentication")
	if err != nil {
		handleError(err)
		// If we can't get the session, try to get the cookie
		cookie, cookieErr := req.Cookie("authentication")
		if cookieErr != nil || cookie == nil {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"authenticated": "false", "error": "No valid session or cookie found"})
			return
		}

		// If we have a cookie but no session, check if it exists in the database
		db, dbErr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
		if dbErr != nil {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database connection failed"})
			return
		}
		defer db.Close()

		// Get user data based on the session token from cookie
		var userID int
		var email string
		err = db.QueryRow(
			`SELECT u.user_id, u.email FROM users u 
			JOIN sessions s ON u.user_id = s.user_id 
			WHERE s.session_token = ? AND s.expires_at > NOW()`,
			cookie.Value,
		).Scan(&userID, &email)

		if err != nil {
			handleError(err)
			if err == sql.ErrNoRows {
				sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"authenticated": "false", "error": "Session expired or invalid"})
				return
			}
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database query failed"})
			return
		}

		// Create a new session since we have a valid cookie
		session, err = store.New(req, "authentication")
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
			return
		}

		session.Options = &sessions.Options{
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,             // Set to true in production with HTTPS
			MaxAge:   30 * 24 * 60 * 60, // 30 days
			Path:     "/",
		}

		session.Values["Authenticated"] = true
		session.Values["email"] = email

		if err = session.Save(req, w); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
			return
		}

		// Update the session token in the database with the new session ID
		_, err = db.Exec(
			`UPDATE sessions SET session_token = ?, expires_at = DATE_ADD(NOW(), INTERVAL 1 MONTH) 
			WHERE user_id = ?`,
			session.ID,
			userID,
		)
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update session"})
			return
		}

		// Get the username for the response
		var username string
		err = db.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username)
		if err != nil {
			handleError(err)
			sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get user data"})
			return
		}

		sendJSONResponse(w, http.StatusOK, UserData{
			Username: username,
			Email:    email,
			Id:       userID,
		})
		return
	}

	// If we have a valid session, return the user data
	auth, ok := session.Values["Authenticated"].(bool)
	if !ok || !auth {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"authenticated": "false", "error": "Not authenticated"})
		return
	}

	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"authenticated": "false", "error": "Invalid session data"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database connection failed"})
		return
	}
	defer db.Close()
	var userId int
	var username string
	err = db.QueryRow("SELECT username, id FROM users WHERE email = ?", email).Scan(&username, &userId)
	if err != nil {
		handleError(err)
		if err == sql.ErrNoRows {
			sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"authenticated": "false", "error": "User not found"})
			return
		}
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Database query failed"})
		return
	}
	fmt.Println(username, email, userId)
	sendJSONResponse(w, http.StatusOK, UserData{
		Username: username,
		Email:    email,
		Id:       userId,
	})
}

func logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "authentication")
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request: No session found"})
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session save failed"})
		return
	}


	err = removeValues("sessions", map[string]string{"session_token": session.ID})
	if err != nil {
		handleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session deletion failed"})
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})

}
