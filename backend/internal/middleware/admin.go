package middleware

import (
	"context"
	"net/http"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
)

// AdminMiddleware checks if user is authenticated AND has admin privileges
// Token send now, need later, fix put delay = check auth FIRST, then check admin flag
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Check if authenticated (token exist check)
		session, err := auth.Store.Get(r, "authentication")
		if err != nil {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Authentication required",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Step 2: Check authenticated flag (token valid check)
		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Authentication required",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Step 3: Get email from session (who you are check)
		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Invalid session",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Step 4: Check if user is admin (admin card check)
		var isAdmin int
		var userID int
		var username string
		row := database.QueryRow("SELECT id, username, is_admin FROM users WHERE email = ?", email)
		err = row.Scan(&userID, &username, &isAdmin)
		if err != nil {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{
				"error": "User not found",
			})
			return
		}

		// Step 5: Admin gate (if not admin, HARD reject)
		if isAdmin != 1 {
			utils.SendJSONResponse(w, http.StatusForbidden, map[string]interface{}{
				"error":        "Admin access required",
				"redirect":     "/",
				"previousPath": r.URL.Path,
			})
			return
		}

		// All checks passed - user is admin
		// Add user info to context for handlers to use
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "isAdmin", true)

		next(w, r.WithContext(ctx))
	}
}
