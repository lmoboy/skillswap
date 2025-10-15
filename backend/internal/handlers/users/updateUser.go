package users

import (
	"net/http"

	"skillswap/backend/internal/utils"
)

// UpdateUser handles updating user profile information
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	// Get and validate user session
	sessionUserID, err := getUserFromSession(req)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	// Parse and validate request payload
	payload, err := parseUserUpdatePayload(req)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Verify user is updating their own profile
	if err := validateUserOwnership(sessionUserID, payload.ID); err != nil {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}

	// Perform the update
	if err := performUserUpdate(payload); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "Profile updated successfully"})
}
