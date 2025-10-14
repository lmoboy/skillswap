package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skillswap/backend/authentication"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	// Get session to verify user owns this profile
	session, err := authentication.Store.Get(req, "authentication")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}

	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}

	// Get user ID from session
	sessionUserID, err := database.GetUserIDFromEmail(email)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to verify user"})
		return
	}

	var payload structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Verify user is updating their own profile
	if int64(payload.ID) != sessionUserID {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": "You can only update your own profile"})
		return
	}

	// Delete existing data before adding new ones to avoid duplicates
	database.Execute("DELETE FROM user_projects WHERE user_id = ?", payload.ID)
	database.Execute("DELETE FROM user_skills WHERE user_id = ?", payload.ID)
	database.Execute("DELETE FROM user_contacts WHERE user_id = ?", payload.ID)

	// Update basic user info if provided
	if payload.AboutMe != "" || payload.Professions != "" || payload.Location != "" {
		_, err = database.Execute(`
			UPDATE users 
			SET aboutme = COALESCE(NULLIF(?, ''), aboutme),
			    profession = COALESCE(NULLIF(?, ''), profession),
			    location = COALESCE(NULLIF(?, ''), location)
			WHERE id = ?
		`, payload.AboutMe, payload.Professions, payload.Location, payload.ID)
		if err != nil {
			utils.HandleError(err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user info"})
			return
		}
	}

	// Add projects
	for _, project := range payload.Projects {
		if project.Description == "" || project.Name == "" || project.Link == "" {
			continue
		}
		if err := addProject(payload, project); err != nil {
			utils.DebugPrint(fmt.Sprintf("Failed to add project: %v", err))
		}
	}

	// Add skills
	for _, skill := range payload.Skills {
		if skill.Name == "" {
			continue
		}
		if err := addSkill(payload, skill); err != nil {
			utils.DebugPrint(fmt.Sprintf("Failed to add skill: %v", err))
		}
	}

	// Add contacts
	for _, contact := range payload.Contacts {
		if contact.Name == "" || contact.Link == "" {
			continue
		}
		if err := addContacts(payload, contact); err != nil {
			utils.DebugPrint(fmt.Sprintf("Failed to add contact: %v", err))
		}
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "Profile updated successfully"})
}

func addProject(user structs.UserInfo, project structs.UserProject) error {
	_, err := database.Execute("INSERT INTO user_projects (user_id, name, description, link) VALUES (?, ?, ?, ?)", user.ID, project.Name, project.Description, project.Link)
	if err != nil {
		utils.HandleError(err)
		return fmt.Errorf("failed to add project: %w", err)
	}
	return nil
}

func addSkill(user structs.UserInfo, skill structs.UserSkill) error {
	skillId, err := database.GetSkillIDFromName(skill.Name)
	if err != nil {
		utils.HandleError(err)
		return fmt.Errorf("failed to get skill ID: %w", err)
	}
	
	database.Debug("INSERT INTO user_skills (user_id, skill_id, verified) VALUES (%v, %v, %v)", user.ID, skillId, skill.Verified)

	_, err = database.Execute("INSERT INTO user_skills (user_id, skill_id, verified) VALUES (?, ?, ?)", user.ID, skillId, bool(skill.Verified))
	if err != nil {
		utils.HandleError(err)
		return fmt.Errorf("failed to add skill: %w", err)
	}
	return nil
}

func addContacts(user structs.UserInfo, contact structs.UserContact) error {
	_, err := database.Execute("INSERT INTO user_contacts (user_id, name, link, icon) VALUES (?, ?, ?, ?)", user.ID, contact.Name, contact.Link, contact.Icon)
	if err != nil {
		utils.HandleError(err)
		return fmt.Errorf("failed to add contact: %w", err)
	}
	return nil
}
