package auth

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"
	"time"

	"github.com/gorilla/sessions"
)

// Use a simple key for session Store
var Authenticated = false
var Store = sessions.NewCookieStore([]byte("simple-session-key-12345"))

func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func ApplySession(w http.ResponseWriter, req *http.Request, userInfo *models.UserInfo) error {
	if userInfo == nil || userInfo.Email == "" {
		return fmt.Errorf("user info and email are required")
	}

	// utils.DebugPrint("Applying session for user", userInfo)
	session, err := Store.Get(req, "authentication")
	if err != nil {
		// utils.DebugPrint("Get session failed", err)
		Authenticated = false
		return err
	}

	// Only create new session ID if one doesn't exist or if it's expired
	if session.IsNew {
		session.ID = fmt.Sprintf("%x", md5.Sum([]byte(userInfo.Email+time.Now().String())))
	}

	session.Values["authenticated"] = true
	session.Values["email"] = userInfo.Email
	session.Values["is_admin"] = userInfo.IsAdmin
	session.Values["lastAccess"] = time.Now().Unix()

	if err := session.Save(req, w); err != nil {
		// utils.DebugPrint(err)
		Authenticated = false
		return err
	}

	Authenticated = true
	return nil
}

func CheckSession(w http.ResponseWriter, req *http.Request) {
	session, err := Store.Get(req, "authentication")
	if err != nil {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}

	// Check if session is authenticated
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}

	// Check if session has timed out (optional additional security)
	if lastAccess, ok := session.Values["lastAccess"].(int64); ok {
		// Session expires after 7 days of inactivity
		if time.Now().Unix()-lastAccess > 86400*7 {
			Authenticated = false
			session.Options.MaxAge = -1 // Expire the session
			session.Save(req, w)
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Session expired"})
			return
		}
		// Update last access time
		session.Values["lastAccess"] = time.Now().Unix()
	}

	// Get user email from session
	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session data"})
		return
	}

	// Verify user exists in database
	row, err := database.Query("SELECT username, email, id, COALESCE(profile_picture, ''), is_admin FROM users WHERE email = ?", email)
	var username string = ""
	var userEmail = ""
	var id = 0
	var profilePicture = ""
	var isAdmin bool = false
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
	err = row.Scan(&username, &userEmail, &id, &profilePicture, &isAdmin)
	if err != nil {
		Authenticated = false
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to check session"})
		return
	}

	// utils.DebugPrint(session.Values)
	Authenticated = true

	// Update session data and save once
	session.Values["is_admin"] = isAdmin
	session.Save(req, w)

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user":            username,
		"email":           userEmail,
		"id":              fmt.Sprintf("%d", id),
		"profile_picture": profilePicture,
		"is_admin":        isAdmin,
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
		// Session doesn't exist, that's fine for logout
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Already logged out"})
		return nil
	}

	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	session.Options.MaxAge = -1
	if err := session.Save(req, w); err != nil {
		// Even if we can't save, we still consider the user logged out
		Authenticated = false
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful (session cleanup warning)"})
		return nil
	}

	Authenticated = false
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})
	return nil
}
