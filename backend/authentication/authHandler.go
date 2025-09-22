package authentication

import (
	"encoding/json"
	"net/http"

	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func IsEmailUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-186"})
		return
	}

	db, err := database.GetDatabase()
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-192"})
		return
	}

	var exists int
	qerr := db.QueryRow("SELECT 1 FROM users WHERE email = ? LIMIT 1", userInfo.Email).Scan(&exists)
	if qerr == nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email already in use"})
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Email available"})
	return
}

func IsUsernameUsed(w http.ResponseWriter, req *http.Request) {
	var userInfo structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-206"})
		return
	}
	db, err := database.GetDatabase()
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-211"})
		return
	}
	var exists int
	qerr := db.QueryRow("SELECT 1 FROM users WHERE username = ? LIMIT 1", userInfo.Username).Scan(&exists)
	if qerr == nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username already in use"})
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Username available"})
	return
}
