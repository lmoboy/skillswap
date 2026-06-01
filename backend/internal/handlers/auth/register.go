package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, req *http.Request) {
	var userInfo models.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	// Validate required fields
	if userInfo.Username == "" || userInfo.Email == "" || userInfo.Password == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username, email, and password are required"})
		return
	}

	// Validate username length and format
	if len(userInfo.Username) < 3 {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username must be at least 3 characters long"})
		return
	}

	if len(userInfo.Username) > 50 {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username must not exceed 50 characters"})
		return
	}

	// Validate password length
	if len(userInfo.Password) < 8 {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters long"})
		return
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
		return
	}
	passwordHash := string(passwordHashBytes)

	_, err = database.Execute("INSERT INTO users ( username, email, password_hash) VALUES (?, ?, ?)", userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		utils.HandleError(err)
		errStr := err.Error()
		// Check for MySQL duplicate entry (Error 1062)
		if strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "1062") ||
			strings.Contains(errStr, "UNIQUE constraint failed") {
			if strings.Contains(errStr, "email") {
				utils.SendJSONResponse(w, http.StatusConflict, map[string]string{"error": "Email already registered"})
				return
			}
			if strings.Contains(errStr, "username") {
				utils.SendJSONResponse(w, http.StatusConflict, map[string]string{"error": "Username already taken"})
				return
			}
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create account. Please try again"})
		return
	}

	// Get the newly created user ID
	row := database.QueryRow("SELECT id FROM users WHERE email = ?", userInfo.Email)
	var userID int
	if err := row.Scan(&userID); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user information"})
		return
	}

	userInfo.ID = userID

	// Apply session automatically after registration
	if err := ApplySession(w, req, &userInfo); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Account created but failed to log in. Please log in manually"})
		return
	}
	
	// Return success response with user data
	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "Registration successful",
		"user":    userInfo.Username,
		"email":   userInfo.Email,
		"id":      userID,
	})
}
