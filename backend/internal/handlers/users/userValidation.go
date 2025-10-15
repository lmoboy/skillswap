package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"skillswap/backend/internal/database"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/models"
)

// getUserFromSession retrieves the user ID from the session
func getUserFromSession(req *http.Request) (int64, error) {
	session, err := auth.Store.Get(req, "authentication")
	if err != nil {
		return 0, fmt.Errorf("invalid session: %w", err)
	}

	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		return 0, fmt.Errorf("invalid session: no email found")
	}

	userID, err := database.GetUserIDFromEmail(email)
	if err != nil {
		return 0, fmt.Errorf("failed to verify user: %w", err)
	}

	return userID, nil
}

// parseUserUpdatePayload decodes and validates the user update request
func parseUserUpdatePayload(req *http.Request) (*models.UserInfo, error) {
	var payload models.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("invalid request payload: %w", err)
	}

	return &payload, nil
}

// validateUserOwnership verifies that the user is updating their own profile
func validateUserOwnership(sessionUserID int64, payloadUserID int) error {
	if int64(payloadUserID) != sessionUserID {
		return fmt.Errorf("you can only update your own profile")
	}
	return nil
}


