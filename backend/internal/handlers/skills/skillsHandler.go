package skills

import (
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"
)

// GetSkills retrieves all skills from the database and returns them as JSON
func GetSkills(w http.ResponseWriter, req *http.Request) {
	rows, err := database.Query(`SELECT id,name,description FROM skills`)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch skills"})
		return
	}
	defer rows.Close()
	skills := []models.Skill{}
	for rows.Next() {
		var p models.Skill
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			utils.HandleError(err)
			continue
		}
		skills = append(skills, p)
	}

	if err = rows.Err(); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process skills"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, skills)
}


