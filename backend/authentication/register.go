package authentication

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func Register(w http.ResponseWriter, req *http.Request) {
	var userInfo structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-146"})
		return
	}

	if userInfo.Username == "" || userInfo.Email == "" || userInfo.Password == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-151"})
		return
	}

	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Password)))

	_, err := database.Execute("INSERT INTO users ( username, email, password_hash) VALUES (?, ?, ?)", userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusConflict, map[string]string{"error": "AH-174", "message": err.Error()})
		return
	}

	if err := ApplySession(w, req, &userInfo); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
		return
	}
	CheckSession(w, req)
	// utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"status": "ok", "message": "Registration successful"})
}
