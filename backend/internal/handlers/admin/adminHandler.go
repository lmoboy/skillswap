package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
	"strconv"
	"time"
)

// AdminStats - big picture numbers for admin dashboard
type AdminStats struct {
	TotalUsers       int `json:"total_users"`
	TotalAdmins      int `json:"total_admins"`
	TotalSkills      int `json:"total_skills"`
	TotalCourses     int `json:"total_courses"`
	TotalChats       int `json:"total_chats"`
	TotalMessages    int `json:"total_messages"`
	NewUsersToday    int `json:"new_users_today"`
	NewUsersWeek     int `json:"new_users_week"`
}

// GetAllStats - return ALL the numbers (admin only)
func GetAllStats(w http.ResponseWriter, r *http.Request) {
	// Get counts from all tables
	stats := AdminStats{}

	// Total users
	row := database.QueryRow("SELECT COUNT(*) FROM users")
	row.Scan(&stats.TotalUsers)

	// Total admins
	row = database.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 1")
	row.Scan(&stats.TotalAdmins)

	// Total skills
	row = database.QueryRow("SELECT COUNT(*) FROM skills")
	row.Scan(&stats.TotalSkills)

	// Total courses
	row = database.QueryRow("SELECT COUNT(*) FROM courses")
	row.Scan(&stats.TotalCourses)

	// Total chats
	row = database.QueryRow("SELECT COUNT(*) FROM chats")
	row.Scan(&stats.TotalChats)

	// Total messages
	row = database.QueryRow("SELECT COUNT(*) FROM messages")
	row.Scan(&stats.TotalMessages)

	// New users today
	row = database.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()")
	row.Scan(&stats.NewUsersToday)

	// New users this week
	row = database.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)")
	row.Scan(&stats.NewUsersWeek)

	utils.SendJSONResponse(w, http.StatusOK, stats)
}

// GetAllUsers - return ALL users (admin only, no escape)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Optional: ?search=xxx filter
	search := r.URL.Query().Get("search")

	var rows *sql.Rows
	var err error

	if search != "" {
		searchTerm := "%" + search + "%"
		rows, err = database.Query(`
			SELECT u.id, u.username, u.email, u.profile_picture, u.aboutme,
				   u.profession, u.location, u.swaps, u.is_admin, u.created_at
			FROM users u
			WHERE u.username LIKE ? OR u.email LIKE ?
			ORDER BY u.created_at DESC
		`, searchTerm, searchTerm)
	} else {
		rows, err = database.Query(`
			SELECT u.id, u.username, u.email, u.profile_picture, u.aboutme,
				   u.profession, u.location, u.swaps, u.is_admin, u.created_at
			FROM users u
			ORDER BY u.created_at DESC
		`)
	}

	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var user struct {
			ID             int
			Username       string
			Email          string
			ProfilePicture string
			AboutMe        sql.NullString
			Profession     sql.NullString
			Location       sql.NullString
			Swaps          int
			IsAdmin        int
			CreatedAt      time.Time
		}

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.ProfilePicture,
			&user.AboutMe, &user.Profession, &user.Location, &user.Swaps,
			&user.IsAdmin, &user.CreatedAt,
		)
		if err != nil {
			utils.HandleError(err)
			continue
		}

		userMap := map[string]interface{}{
			"id":              user.ID,
			"username":        user.Username,
			"email":           user.Email,
			"profile_picture": user.ProfilePicture,
			"aboutme":         "",
			"profession":      "",
			"location":        "",
			"swaps":           user.Swaps,
			"is_admin":        user.IsAdmin == 1,
			"created_at":      user.CreatedAt.Format("2006-01-02"),
		}

		if user.AboutMe.Valid {
			userMap["aboutme"] = user.AboutMe.String
		}
		if user.Profession.Valid {
			userMap["profession"] = user.Profession.String
		}
		if user.Location.Valid {
			userMap["location"] = user.Location.String
		}

		users = append(users, userMap)
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// ToggleUserAdmin - make/unmake user admin (admin only)
func ToggleUserAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req struct {
		UserID  int  `json:"user_id"`
		SetAdmin bool `json:"set_admin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
		return
	}

	// Safety: cannot remove own admin status
	adminID, ok := r.Context().Value("userID").(int)
	if ok && adminID == req.UserID && !req.SetAdmin {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{
			"error": "Cannot remove your own admin privileges",
		})
		return
	}

	adminValue := 0
	if req.SetAdmin {
		adminValue = 1
	}

	_, err := database.Execute("UPDATE users SET is_admin = ? WHERE id = ?", adminValue, req.UserID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update user admin status",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": fmt.Sprintf("User admin status updated to %v", req.SetAdmin),
	})
}

// DeleteUser - HARD delete user (admin only, dangerous)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		utils.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}

	// Try JSON decode first, then query param
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing user_id",
			})
			return
		}
		req.UserID = userID
	}

	// Safety: cannot delete yourself
	adminID, ok := r.Context().Value("userID").(int)
	if ok && adminID == req.UserID {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{
			"error": "Cannot delete your own account",
		})
		return
	}

	// HARD delete - all cascades will fire (chats, messages, courses, etc.)
	_, err := database.Execute("DELETE FROM users WHERE id = ?", req.UserID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "User deleted successfully",
	})
}

// GetAllCourses - return ALL courses (admin only)
func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query(`
		SELECT c.id, c.title, c.description, c.instructor_id, c.skill_id,
			   c.difficulty_level, c.duration_hours, c.max_students,
			   c.current_students, c.price, c.thumbnail_url, c.status,
			   c.created_at, u.username as instructor_name
		FROM courses c
		LEFT JOIN users u ON c.instructor_id = u.id
		ORDER BY c.created_at DESC
	`)

	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve courses",
		})
		return
	}
	defer rows.Close()

	courses := []map[string]interface{}{}
	for rows.Next() {
		var course struct {
			ID              int
			Title           string
			Description     sql.NullString
			InstructorID    int
			SkillID         int
			DifficultyLevel string
			DurationHours   sql.NullInt64
			MaxStudents     int
			CurrentStudents int
			Price           float64
			ThumbnailURL    sql.NullString
			Status          string
			CreatedAt       time.Time
			InstructorName  sql.NullString
		}

		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID,
			&course.SkillID, &course.DifficultyLevel, &course.DurationHours,
			&course.MaxStudents, &course.CurrentStudents, &course.Price,
			&course.ThumbnailURL, &course.Status, &course.CreatedAt,
			&course.InstructorName,
		)
		if err != nil {
			utils.HandleError(err)
			continue
		}

		courseMap := map[string]interface{}{
			"id":               course.ID,
			"title":            course.Title,
			"description":      "",
			"instructor_id":    course.InstructorID,
			"skill_id":         course.SkillID,
			"difficulty_level": course.DifficultyLevel,
			"duration_hours":   0,
			"max_students":     course.MaxStudents,
			"current_students": course.CurrentStudents,
			"price":            course.Price,
			"thumbnail_url":    "",
			"status":           course.Status,
			"created_at":       course.CreatedAt.Format("2006-01-02"),
			"instructor_name":  "",
		}

		if course.Description.Valid {
			courseMap["description"] = course.Description.String
		}
		if course.DurationHours.Valid {
			courseMap["duration_hours"] = int(course.DurationHours.Int64)
		}
		if course.ThumbnailURL.Valid {
			courseMap["thumbnail_url"] = course.ThumbnailURL.String
		}
		if course.InstructorName.Valid {
			courseMap["instructor_name"] = course.InstructorName.String
		}

		courses = append(courses, courseMap)
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"courses": courses,
		"count":   len(courses),
	})
}

// DeleteCourse - HARD delete course (admin only)
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		utils.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req struct {
		CourseID int `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		courseIDStr := r.URL.Query().Get("course_id")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing course_id",
			})
			return
		}
		req.CourseID = courseID
	}

	_, err := database.Execute("DELETE FROM courses WHERE id = ?", req.CourseID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete course",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Course deleted successfully",
	})
}

// GetAllSkills - return ALL skills (admin only)
func GetAllSkills(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT id, name, description FROM skills ORDER BY name ASC")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve skills",
		})
		return
	}
	defer rows.Close()

	skills := []map[string]interface{}{}
	for rows.Next() {
		var s struct {
			ID          int
			Name        string
			Description string
		}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description); err != nil {
			utils.HandleError(err)
			continue
		}
		skills = append(skills, map[string]interface{}{
			"id":          s.ID,
			"name":        s.Name,
			"description": s.Description,
		})
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"skills": skills,
		"count":  len(skills),
	})
}

// AddSkill - create new skill (admin only)
func AddSkill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
		return
	}

	if req.Name == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Skill name is required",
		})
		return
	}

	_, err := database.Execute("INSERT INTO skills (name, description) VALUES (?, ?)", req.Name, req.Description)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to add skill",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Skill added successfully",
	})
}

// UpdateSkill - update existing skill (admin only)
func UpdateSkill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
		return
	}

	_, err := database.Execute("UPDATE skills SET name = ?, description = ? WHERE id = ?", req.Name, req.Description, req.ID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update skill",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Skill updated successfully",
	})
}

// DeleteSkill - HARD delete skill (admin only)
func DeleteSkill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SkillID int `json:"skill_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		skillIDStr := r.URL.Query().Get("skill_id")
		skillID, err := strconv.Atoi(skillIDStr)
		if err != nil {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing skill_id",
			})
			return
		}
		req.SkillID = skillID
	}

	_, err := database.Execute("DELETE FROM skills WHERE id = ?", req.SkillID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete skill",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Skill deleted successfully",
	})
}

// UpdateUserSwaps - modify user's available swaps (admin only)
func UpdateUserSwaps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req struct {
		UserID int `json:"user_id"`
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
		return
	}

	_, err := database.Execute("UPDATE users SET swaps = swaps + ? WHERE id = ?", req.Amount, req.UserID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update user swaps",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "User swaps updated successfully",
	})
}

// GetSystemHealth - check all systems go (admin only)
func GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"database":    "ok",
		"timestamp":   time.Now().Format(time.RFC3339),
		"connections": 0,
	}

	// Check database connection
	db, err := database.GetDatabase()
	if err != nil {
		health["database"] = "error"
		health["error"] = err.Error()
		utils.SendJSONResponse(w, http.StatusServiceUnavailable, health)
		return
	}

	// Get connection stats
	stats := db.Stats()
	health["connections"] = map[string]any{
		"max_open":      stats.MaxOpenConnections,
		"open":          stats.OpenConnections,
		"in_use":        stats.InUse,
		"idle":          stats.Idle,
		"wait_count":    stats.WaitCount,
		"wait_duration": int(stats.WaitDuration.Seconds()),
	}

	// Ping database
	if err := db.Ping(); err != nil {
		health["database"] = "error"
		health["error"] = err.Error()
		utils.SendJSONResponse(w, http.StatusServiceUnavailable, health)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, health)
}
