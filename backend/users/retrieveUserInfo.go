package users

import (
	"net/http"
	"reflect"
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
			GROUP_CONCAT(JSON_OBJECT(
				s.name,
				s.description
			), ' ') AS skills,
			GROUP_CONCAT(
			JSON_OBJECT(
				uc.name,
				uc.link
			), ' ') AS user_contacts,
			GROUP_CONCAT(
			JSON_OBJECT(
				up.name,
				up.name
			), ' ') AS user_projects
		FROM users AS u
		JOIN user_skills AS us ON u.id = us.user_id
		JOIN skills AS s ON us.skill_id = s.id
		LEFT JOIN user_contacts AS uc ON u.id = uc.user_id
		LEFT JOIN user_contacts AS c ON uc.id = c.id
		LEFT JOIN user_projects AS up ON u.id = up.user_id
		LEFT JOIN user_projects AS p ON up.id = p.id
		WHERE u.id = 1
		GROUP BY u.id, u.username, u.email, u.profile_picture, u.aboutme, u.location
	`, req.URL.Query().Get("q"))
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user info"})
		return
	}

	var user structs.UserInfo
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilePicture, &user.AboutMe, &user.Location, &user.Skills, &user.Contacts, &user.Projects)
	// err = rows.Scan()
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user info"})
		return
	}

	v := reflect.ValueOf(user)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	utils.DebugPrint(values)
	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"username":        user.Username,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
		"aboutme":         user.AboutMe,
		"location":        user.Location,
		"skills":          user.Skills,
		"contacts":        user.Contacts,
		"projects":        user.Projects,
	})
	// utils.SendJSONResponse(w, http.StatusOK, everything)
}
