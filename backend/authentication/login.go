package authentication

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

// Login authenticates a user from JSON credentials in the request body, creates a session on success, and writes an HTTP JSON response.
// It expects a JSON UserInfo with Email and Password, verifies the credentials against the users table (password compared as an MD5 hex hash), applies a session if the credentials match, and returns 200 on success, 400 for malformed input, 401 for invalid credentials, or 500 for server errors.
func Login(w http.ResponseWriter, req *http.Request) {
	var userInfo structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-227"})
		return
	}

	row := database.QueryRow("SELECT username, email, id FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))

	var storedUsername, storedEmail, storedID string
	if err := row.Scan(&storedUsername, &storedEmail, &storedID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Report this to the developer\nERRC : LG-29"})
		return
	}

	if err := ApplySession(w, req, &structs.UserInfo{Username: storedUsername, Email: storedEmail, ID: storedID}); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}
