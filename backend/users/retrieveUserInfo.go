package users

import (
	"encoding/json"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
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
  COALESCE(u.profile_picture, "noPicture") as profile_picture,
  COALESCE(u.aboutme, "") as aboutme,
  COALESCE(u.location, "") as location,
  COALESCE(u.profession, "") as profession,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", s.name, "verified", us.verified)
    ), JSON_ARRAY()
  ) AS skills,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", up.name, "description", COALESCE(up.description, ""), "link", COALESCE(up.link, ""))
    ), JSON_ARRAY()
  ) AS projects,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", uc.name, "link", COALESCE(uc.link, ""), "icon", uc.icon)
    ), JSON_ARRAY()
  ) AS contacts

FROM users AS u

LEFT JOIN user_skills AS us ON u.id = us.user_id
LEFT JOIN skills AS s ON us.skill_id = s.id

LEFT JOIN user_projects AS up ON u.id = up.user_id
LEFT JOIN user_contacts AS uc ON u.id = uc.user_id

WHERE u.id = ?
GROUP BY u.id, u.username, u.email, u.profile_picture, u.aboutme, u.location, u.profession;
`, req.URL.Query().Get("q"))
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user info"})
		return
	}

	var user structs.UserInfo
	var skillsJSON []byte
	var projectsJSON []byte
	var contactsJSON []byte
	if !rows.Next() {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilePicture, &user.AboutMe, &user.Location, &user.Professions, &skillsJSON, &projectsJSON, &contactsJSON)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to parse user data"})
		return
	}

	// Unmarshal JSON arrays into structs
	if err := json.Unmarshal(skillsJSON, &user.Skills); err != nil {
		utils.HandleError(err)
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

	utils.SendJSONResponse(w, http.StatusOK, user)
}
