package authentication

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
	"time"

	"github.com/gorilla/sessions"
)

// Use a simple key for session store
var store = sessions.NewCookieStore([]byte("simple-session-key-12345"))

func ApplySession(w http.ResponseWriter, req *http.Request, userInfo *structs.UserInfo) error {
	if userInfo == nil || userInfo.Email == "" {
		return fmt.Errorf("user info and email are required")
	}

	utils.DebugPrint("Applying session for user", userInfo)
	session, err := store.New(req, "authentication")
	if err != nil {
		utils.DebugPrint("Create session failed", err)
		return err
	}
	session.ID = fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Email+time.Now().String())))
	session.Values["authenticated"] = true
	session.Values["email"] = userInfo.Email

	if err := session.Save(req, w); err != nil {
		utils.DebugPrint(err)
		return err
	}

	return nil
}

func CheckSession(w http.ResponseWriter, req *http.Request) {
	values, err := store.Get(req, "authentication")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}

	if values.Values["authenticated"] != true {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	row, err := database.Query("SELECT username, email, id, COALESCE(profile_picture, '') FROM users WHERE email = ?", values.Values["email"])
	var username string = ""
	var email = ""
	var id = 0
	var profilePicture = ""
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check session"})
		return
	}
	row.Next()
	err = row.Scan(&username, &email, &id, &profilePicture)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check session"})
		return
	}
	defer row.Close()
	// utils.DebugPrint(values.Values)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"user": username, "email": email, "id": fmt.Sprintf("%d", id), "profile_picture": profilePicture})
}

// RemoveSession invalidates and removes the current authentication session.
//
// If the session cannot be retrieved, it responds with 400 Bad Request and returns an error.
// If saving the invalidated session fails, it responds with 500 Internal Server Error and returns an error.
// On success, the session is invalidated and the function returns nil.
func RemoveSession(w http.ResponseWriter, req *http.Request) error {

	session, err := store.Get(req, "authentication")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid session"})
		return fmt.Errorf("invalid session")
	}
	session.Options.MaxAge = -1
	if err := session.Save(req, w); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
		return fmt.Errorf("failed to save session")
	}
	// utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "BAIIII"})

	return nil
}
