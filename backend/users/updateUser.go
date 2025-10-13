package users

import (
	"encoding/json"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
)

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	var payload structs.UserInfo
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		utils.HandleError(err)

		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-146"})
		return
	}

	for _, project := range payload.Projects {
		if project.Description == "" || project.Name == "" || project.Link == "" {
			continue
		}
		addProject(payload, project)
	}
	for _, skill := range payload.Skills {
		if skill.Name == "" {
			continue
		}
		addSkill(payload, skill)
	}
	for _, contact := range payload.Contacts {
		if contact.Name == "" || contact.Link == "" {
			continue
		}
		addContacts(payload, contact)
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "received"})
}

func addProject(user structs.UserInfo, project structs.UserProject) {
	_, err := database.Execute("INSERT INTO user_projects (user_id, name,description,link) VALUES (?, ?, ?, ?)", user.ID, project.Name, project.Description, project.Link)
	if err != nil {
		utils.HandleError(err)
	}
}

func addSkill(user structs.UserInfo, skill structs.UserSkill) {
	skillId, err := database.GetSkillIDFromName(skill.Name)
	if err != nil {
		utils.HandleError(err)
		return
	}
	database.Debug("INSERT INTO user_skills (user_id, skill_id, verified) VALUES (%v, %v, %v)", user.ID, skillId, skill.Verified)

	_, err = database.Query("INSERT INTO user_skills (user_id, skill_id, verified) VALUES (?,?,?)", user.ID, skillId, bool(skill.Verified))
	if err != nil {
		utils.HandleError(err)
		return
	}
}

func addContacts(user structs.UserInfo, contact structs.UserContact) {
	_, err := database.Execute("INSERT INTO user_contacts (user_id, name, link, icon) VALUES (?, ?, ?, ?)", user.ID, contact.Name, contact.Link, contact.Icon)
	if err != nil {
		utils.HandleError(err)
		return
	}
}
