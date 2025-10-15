package middleware

import (
	"net/http"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
)

// AuthMiddleware checks if the user is authenticated before allowing access to protected routes
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session from store
		session, err := auth.Store.Get(r, "authentication")
		if err != nil {
			// utils.DebugPrint("AuthMiddleware: Failed to get session:", err)
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Authentication required",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Check if authenticated flag is set
		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			// utils.DebugPrint("AuthMiddleware: User not authenticated")
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Authentication required",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Verify email exists in session
		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			// utils.DebugPrint("AuthMiddleware: Invalid session - no email")
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "Invalid session",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// Verify user exists in database
		row, err := database.Query("SELECT id FROM users WHERE email = ?", email)
		if err != nil {
			utils.HandleError(err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to verify authentication",
			})
			return
		}
		defer row.Close()

		if !row.Next() {
			// utils.DebugPrint("AuthMiddleware: User not found in database")
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"error":        "User not found",
				"redirect":     "/auth/login",
				"previousPath": r.URL.Path,
			})
			return
		}

		// User is authenticated, proceed to next handler
		// utils.DebugPrint(fmt.Sprintf("AuthMiddleware: User %s authenticated successfully", email))
		next(w, r)
	}
}

// OptionalAuthMiddleware adds user information to request if authenticated, but doesn't block unauthenticated users
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get session, but don't fail if it doesn't exist
		session, err := auth.Store.Get(r, "authentication")
		if err == nil {
			if authenticated, ok := session.Values["authenticated"].(bool); ok && authenticated {
				// utils.DebugPrint("OptionalAuthMiddleware: User authenticated")
			}
		}
		// Always proceed to next handler
		next(w, r)
	}
}

