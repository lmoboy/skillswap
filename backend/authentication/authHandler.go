package authentication

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"skillswap/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// +++++++++++++++++structs+++++++++++++++++

// Global session store

// +++++++++++++++++ Helper function for JSON responses +++++++++++++++++

// +++++++++++++++++ Helper function for applying sessions +++++++++++++++++

// ++++++++++++++ Register Handler ++++++++++++++


func isEmailUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-186"})
		return
	}
	rows, err := findValues("users", []string{"email"}, map[string]string{"email": userInfo.Email})
	// fmt.Println(rows)
	if err != nil {
		HandleError(err)
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
		HandleError(err)
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
		HandleError(err)
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


func logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "authentication")
	if err != nil {
		HandleError(err)
		sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request: No session found"})
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		HandleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session save failed"})
		return
	}

	err = removeValues("sessions", map[string]string{"session_token": session.ID})
	if err != nil {
		HandleError(err)
		sendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session deletion failed"})
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})

}
