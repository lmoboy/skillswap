package users

import (
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func RetrieveUserInfo(w http.ResponseWriter, req *http.Request) {

	rows, err := database.Query(`
SELECT 
  u.id, 
  u.username, 
  u.email, 
  u.profile_picture, 
  u.aboutme, 
  u.location,
  u.profession,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", s.name, "description", s.description)
    ), JSON_ARRAY()
  ) AS skills,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", up.name, "link", up.link, "description", up.description)
    ), JSON_ARRAY()
  ) AS projects,

  COALESCE(
    JSON_ARRAYAGG(
      JSON_OBJECT("name", uc.name, "link", uc.link, "icon", uc.icon)
    ), JSON_ARRAY()
  ) AS contacts

FROM users AS u

LEFT JOIN (
  SELECT DISTINCT user_id, skill_id
  FROM user_skills
) AS us ON u.id = us.user_id
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
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilePicture, &user.AboutMe, &user.Location, &user.Professions, &user.Skills, &user.Projects, &user.Contacts)
	// err = rows.Scan()
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user info"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"username":        user.Username,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
		"aboutme":         user.AboutMe,
		"location":        user.Location,
		"profession":      user.Professions,
		"skills":          user.Skills,
		"contacts":        user.Contacts,
		"projects":        user.Projects,
	})
}
