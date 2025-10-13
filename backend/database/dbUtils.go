package database

import (
	"encoding/json"
	"net/http"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

// GetUserIDFromEmail returns the user ID for a given email
func GetUserIDFromEmail(email string) (int64, error) {
	db, err := GetDatabase()
	if err != nil {
		return -1, err
	}

	var id int64
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetSkillIDFromName returns the skill ID for a given name
func GetSkillIDFromName(name string) (int64, error) {
	db, err := GetDatabase()
	if err != nil {
		return -1, err
	}

	var id int64
	err = db.QueryRow("SELECT id FROM skills WHERE NAME LIKE ? LIMIT 1", name).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetAllSkills returns a list of all skills in the database
func GetAllSkills() ([]structs.Skill, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, name, description FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []structs.Skill
	for rows.Next() {
		var skill structs.Skill
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Description)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

// Search decodes a JSON body containing a "query" string and returns up to five users matching the query along with their aggregated skills.
//
// The request body must be JSON with the field `query`. The handler responds with a JSON array of `structs.SearchResult` entries,
// each containing the user's ID, username, email, and a comma-separated `SkillsFound` string. If the JSON body cannot be decoded,
// the handler responds with HTTP 200 and a JSON error message. On database/query errors the error is recorded and the handler returns without writing a further response.
func Search(w http.ResponseWriter, req *http.Request) {

	var requestBody struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		utils.HandleError(err)

		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"error": "Failed to get data"})
		return
	}

	searchQuery := "%" + requestBody.Query + "%"

	// Use a corrected and more efficient SQL query with LEFT JOINs to include users without skills
	rows, err := Query(`
        SELECT
            u.id,
            u.username,
            u.email,
            COALESCE(GROUP_CONCAT(s.name SEPARATOR ', '), '') AS skills_found
        FROM users AS u
        LEFT JOIN user_skills AS us ON u.id = us.user_id
        LEFT JOIN skills AS s ON us.skill_id = s.id
        WHERE u.username LIKE ? OR u.email LIKE ? OR s.name LIKE ? OR s.description LIKE ?
        GROUP BY u.id, u.username, u.email
        ORDER BY u.id
		LIMIT 5
    `, searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		utils.HandleError(err)
		return
	}
	defer rows.Close()

	// Define a new struct to hold the combined skills
	var results []structs.SearchResult

	for rows.Next() {
		var r structs.SearchResult
		if err := rows.Scan(&r.User.ID, &r.User.Username, &r.User.Email, &r.SkillsFound); err != nil {
			utils.HandleError(err)
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		utils.HandleError(err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, results)
}

// FullSearch searches users by username, email, skill name, or skill description and writes matching users
// with aggregated skill names to the response as a JSON array of structs.SearchResult.
//
// The request must provide a JSON body with a `query` string field; the server performs a wildcard search
// around that query and returns each matched user's ID, username, email, about-me text, profession,
// comma-separated skills (may be empty), and join/created timestamp.
//
// On JSON decode failure it logs the error and responds with HTTP 200 and a JSON error object
// (`{"error":"Failed to get shit"}`). On database/query errors it logs the error and returns without
// writing a success payload.
func FullSearch(w http.ResponseWriter, req *http.Request) {

	var requestBody struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		utils.HandleError(err)

		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"error": "Failed to get data"})
		return
	}

	searchQuery := "%" + requestBody.Query + "%"

	// Use a corrected and more efficient SQL query with LEFT JOINs to include users without skills
	rows, err := Query(`
        SELECT
            u.id,
            u.username,
            u.email,
            u.aboutme,
            u.profession,
            COALESCE(GROUP_CONCAT(s.name SEPARATOR ', '), '') AS skills_found,
            u.created_at
        FROM users AS u
        LEFT JOIN user_skills AS us ON u.id = us.user_id
        LEFT JOIN skills AS s ON us.skill_id = s.id
        WHERE u.username LIKE ? OR u.email LIKE ? OR s.name LIKE ? OR s.description LIKE ?
        GROUP BY u.id, u.username, u.email
        ORDER BY u.id
    `, searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		utils.HandleError(err)
		return
	}
	defer rows.Close()

	// Define a new struct to hold the combined skills
	var results []structs.SearchResult

	for rows.Next() {
		var r structs.SearchResult
		if err := rows.Scan(&r.User.ID, &r.User.Username, &r.User.Email, &r.User.AboutMe, &r.User.Professions, &r.SkillsFound, &r.User.Joined); err != nil {
			utils.HandleError(err)
			return
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		utils.HandleError(err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, results)
}
