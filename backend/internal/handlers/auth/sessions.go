package auth

import (
	"fmt"
	"net/http"
	"os"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"
	"time"

	"github.com/gorilla/sessions"
)

// Use a simple key for session Store
var Authenticated = false

var Store = sessions.NewCookieStore([]byte(getSessionKey()))

func getSessionKey() string {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		// Generate a random key at startup if not configured
		return "skillswap-dev-session-key-change-in-production"
	}
	return key
}

func init() {
	isDev := os.Getenv("ENVIRONMENT") != "production"
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   !isDev,
		SameSite: http.SameSiteLaxMode,
	}
}

func ApplySession(w http.ResponseWriter, req *http.Request, userInfo *models.UserInfo) error {
	if userInfo == nil || userInfo.Email == "" {
		return fmt.Errorf("user info and email are required")
	}

	// utils.DebugPrint("Applying session for user", userInfo)
	session, err := Store.New(req, "authentication")
	if err != nil {
		// utils.DebugPrint("Create session failed", err)
		Authenticated = false
		return err
	}
	session.ID = fmt.Sprintf("%x", time.Now().UnixNano())
	session.Values["authenticated"] = true
	session.Values["email"] = userInfo.Email

	if err := session.Save(req, w); err != nil {
		// utils.DebugPrint(err)
		Authenticated = false
		return err
	}

	return nil
}

func CheckSession(w http.ResponseWriter, req *http.Request) {
	values, err := Store.Get(req, "authentication")
	if err != nil {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}

	if values.Values["authenticated"] != true {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	row, err := database.Query("SELECT username, email, id, COALESCE(profile_picture, ''), is_admin, swaps FROM users WHERE email = ?", values.Values["email"])
	var username string = ""
	var email = ""
	var id = 0
	var profilePicture = ""
	var isAdmin = 0
	var swaps = 0
	if err != nil {
		Authenticated = false
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check session"})
		return
	}
	defer row.Close()
	if !row.Next() {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "User not found"})
		return
	}
	err = row.Scan(&username, &email, &id, &profilePicture, &isAdmin, &swaps)
	if err != nil {
		Authenticated = false
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check session"})
		return
	}
	// utils.DebugPrint(values.Values)
	Authenticated = true
	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user":            username,
		"email":           email,
		"id":              fmt.Sprintf("%d", id),
		"profile_picture": profilePicture,
		"is_admin":        isAdmin == 1,
		"swaps":           swaps,
	})
}

// RemoveSession invalidates and removes the current authentication session.
//
// If the session cannot be retrieved, it responds with 400 Bad Request and returns an error.
// If saving the invalidated session fails, it responds with 500 Internal Server Error and returns an error.
// On success, the session is invalidated and the function returns nil.
func RemoveSession(w http.ResponseWriter, req *http.Request) error {

	session, err := Store.Get(req, "authentication")
	if err != nil {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid session"})
		return fmt.Errorf("invalid session")
	}
	session.Options.MaxAge = -1
	if err := session.Save(req, w); err != nil {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
		return fmt.Errorf("failed to save session")
	}
	// utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "BAIIII"})
	Authenticated = false
	return nil
}
