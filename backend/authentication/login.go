package authentication

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, req *http.Request) {
	var userInfo Info
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-227"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		handleError(err)
		SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-233"})
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT username, email, id FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))
	var storedID int
	var storedUsername, storedEmail string
	if err := row.Scan(&storedUsername, &storedEmail, &storedID); err != nil {
		if err == sql.ErrNoRows {
			SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}
		SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-246"})
		return
	}

	applySession(w, req, &Info{Username: storedUsername, Email: storedEmail, id: storedID})

	SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}
