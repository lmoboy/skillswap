package users

import (
	"encoding/json"
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"
)

// RetrieveUserInfo handles an HTTP request to fetch a user's complete profile by id (query parameter "q")
// and writes the user object as JSON to the response.
// It queries the database for user fields and JSON-encoded arrays for skills, projects, and contacts,
// unmarshals those arrays into the corresponding struct fields, and returns the assembled user.
// Responds with HTTP 200 and the user on success, HTTP 404 if the user is not found, and HTTP 500 on database
// or JSON unmarshalling errors.
func RetrieveUserInfo(w http.ResponseWriter, req *http.Request) {

	rows, err := database.Query(`
		SELECT
  u.id,
  u.username,
  u.email,
  u.is_admin,
  COALESCE(u.profile_picture, 'noPicture') as profile_picture,
  COALESCE(u.aboutme, '') as aboutme,
  COALESCE(u.location, '') as location,
  COALESCE(u.profession, '') as profession,

  -- Skills Subquery
  (SELECT JSON_GROUP_ARRAY(
      JSON_OBJECT('name', s.name, 'verified', us.verified)
    )
   FROM user_skills us
   JOIN skills s ON us.skill_id = s.id
   WHERE us.user_id = u.id) AS skills,

  -- Projects Subquery
  (SELECT JSON_GROUP_ARRAY(
      JSON_OBJECT('name', up.name, 'description', COALESCE(up.description, ''), 'link', COALESCE(up.link, ''))
    )
   FROM user_projects up
   WHERE up.user_id = u.id) AS projects,

  -- Contacts Subquery
  (SELECT JSON_GROUP_ARRAY(
      JSON_OBJECT('name', uc.name, 'link', COALESCE(uc.link, ''), 'icon', uc.icon)
    )
   FROM user_contacts uc
   WHERE uc.user_id = u.id) AS contacts

FROM users AS u
WHERE u.id = ?;`, req.URL.Query().Get("q"))
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user info"})
		return
	}

	var user models.UserInfo
	var skillsJSON []byte
	var projectsJSON []byte
	var contactsJSON []byte
	if !rows.Next() {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.ProfilePicture, &user.AboutMe, &user.Location, &user.Professions, &skillsJSON, &projectsJSON, &contactsJSON)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse user data"})
		return
	}

	// Unmarshal JSON arrays into structs
	if err := json.Unmarshal(skillsJSON, &user.Skills); err != nil {
		utils.HandleError(err)
		// Log the actual JSON data for debugging
		// utils.DebugPrint("Failed to unmarshal skills JSON: " + string(skillsJSON))
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse skills data"})
		return
	}

	if err := json.Unmarshal(projectsJSON, &user.Projects); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse projects data"})
		return
	}

	if err := json.Unmarshal(contactsJSON, &user.Contacts); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse contacts data"})
		return
	}
	if(len(user.Contacts)>0){
		for i, contact := range user.Contacts {
			if contact.Name == "" && contact.Link == "" && contact.Icon == "" {
				user.Contacts = append(user.Contacts[:i], user.Contacts[i+1:]...)

			}
		}
	}
	if(len(user.Projects)>0){
		for i, project := range user.Projects {
			if project.Name == "" && project.Link == "" && project.Description == "" {
				user.Projects = append(user.Projects[:i], user.Projects[i+1:]...)

			}
		}
	}
	if(len(user.Skills)>0){
		for i, skill := range user.Skills {
			if skill.Name == "" {
				user.Skills = append(user.Skills[:i], user.Skills[i+1:]...)

			}
		}
	}
	utils.SendJSONResponse(w, http.StatusOK, user)
}
