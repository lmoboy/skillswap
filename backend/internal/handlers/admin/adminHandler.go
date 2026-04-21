package admin

import (
	"encoding/json"
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/models"
	"skillswap/backend/internal/utils"
	"strconv"
	"strings"
)

// IsAdmin checks if the current session belongs to an admin
func IsAdmin(req *http.Request) bool {
	session, err := auth.Store.Get(req, "authentication")
	if err != nil {
		return false
	}
	isAdmin, ok := session.Values["is_admin"].(bool)
	return ok && isAdmin
}

// AdminOnly is a helper to wrap admin-only handlers
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !IsAdmin(req) {
			utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": "Admin access required"})
			return
		}
		next(w, req)
	}
}

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, req *http.Request) {
	rows, err := database.Query("SELECT id, username, email, is_admin, created_at FROM users")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var username, email, createdAt string
		var isAdmin bool
		if err := rows.Scan(&id, &username, &email, &isAdmin, &createdAt); err != nil {
			utils.HandleError(err)
			continue
		}
		users = append(users, map[string]interface{}{
			"id":         id,
			"username":   username,
			"email":      email,
			"is_admin":   isAdmin,
			"created_at": createdAt,
		})
	}
	utils.SendJSONResponse(w, http.StatusOK, users)
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	idStr := strings.TrimPrefix(req.URL.Path, "/api/admin/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	_, err = database.Execute("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// AddSkill adds a new skill
func AddSkill(w http.ResponseWriter, req *http.Request) {
	var skill models.Skill
	if err := json.NewDecoder(req.Body).Decode(&skill); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		return
	}

	if skill.Name == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Skill name is required"})
		return
	}

	_, err := database.Execute("INSERT INTO skills (name, description) VALUES (?, ?)", skill.Name, skill.Description)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add skill"})
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// DeleteSkill deletes a skill
func DeleteSkill(w http.ResponseWriter, req *http.Request) {
	idStr := strings.TrimPrefix(req.URL.Path, "/api/admin/skills/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid skill ID"})
		return
	}

	_, err = database.Execute("DELETE FROM skills WHERE id = ?", id)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete skill"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DeleteCourse deletes a course
func DeleteCourse(w http.ResponseWriter, req *http.Request) {
	idStr := strings.TrimPrefix(req.URL.Path, "/api/admin/courses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	_, err = database.Execute("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete course"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GetCourses returns all courses for admin
func GetCourses(w http.ResponseWriter, req *http.Request) {
	rows, err := database.Query(`
		SELECT c.id, c.title, c.description, u.username as instructor_name, s.name as skill_name, c.status, c.created_at
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
	`)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
		return
	}
	defer rows.Close()

	courses := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var title, description, instructorName, skillName, status, createdAt string
		if err := rows.Scan(&id, &title, &description, &instructorName, &skillName, &status, &createdAt); err != nil {
			utils.HandleError(err)
			continue
		}
		courses = append(courses, map[string]interface{}{
			"id":              id,
			"title":           title,
			"description":     description,
			"instructor_name": instructorName,
			"skill_name":      skillName,
			"status":          status,
			"created_at":      createdAt,
		})
	}
	utils.SendJSONResponse(w, http.StatusOK, courses)
}
