package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, req *http.Request) {
	var userInfo models.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if userInfo.Email == "" || userInfo.Password == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
		return
	}

	row := database.QueryRow("SELECT username, email, id, password_hash FROM users WHERE email = ?", userInfo.Email)
	var storedID int
	var storedUsername, storedEmail, storedHash string
	if err := row.Scan(&storedUsername, &storedEmail, &storedID, &storedHash); err != nil {
		if err == sql.ErrNoRows {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
			return
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "An error occurred. Please try again later"})
		return
	}

	// Verify password with bcrypt (fallback to MD5 for legacy hashes)
	if err := verifyPassword(storedHash, userInfo.Password); err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		return
	}

	if err := ApplySession(w, req, &models.UserInfo{Username: storedUsername, Email: storedEmail, ID: storedID}); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Login successful"})
}

// verifyPassword checks a password against a bcrypt hash.
func verifyPassword(storedHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}
