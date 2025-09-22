package authentication

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func register(w http.ResponseWriter, req *http.Request) {
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

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-159"})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users ( username, email, password_hash) VALUES (?, ?, ?)")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "AH-166"})
		return
	}
	defer stmt.Close()

	// fmt.Println(userInfo)
	_, err = stmt.Exec(userInfo.Username, userInfo.Email, passwordHash)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusConflict, map[string]string{"error": "AH-174", "message": err.Error()})
		return
	}

	ApplySession(w, req, &userInfo)

	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"status": "ok", "message": "Registration successful"})
}
