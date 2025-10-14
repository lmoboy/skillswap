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

func Login(w http.ResponseWriter, req *http.Request) {
	var userInfo structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if userInfo.Email == "" || userInfo.Password == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
		return
	}

	row := database.QueryRow("SELECT username, email, id FROM users WHERE email = ? AND password_hash = ?", userInfo.Email, fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password))))
	var storedID int
	var storedUsername, storedEmail string
	if err := row.Scan(&storedUsername, &storedEmail, &storedID); err != nil {
		if err == sql.ErrNoRows {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
			return
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "An error occurred. Please try again later"})
		return
	}

	if err := ApplySession(w, req, &structs.UserInfo{Username: storedUsername, Email: storedEmail, ID: storedID}); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}
