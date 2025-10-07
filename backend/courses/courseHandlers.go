package courses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/utils"
	"strconv"
)

// GetAllCourses returns all published courses
func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query(`
		SELECT 
			c.id, c.title, c.description, c.instructor_id, u.username as instructor_name,
			c.skill_id, s.name as skill_name, c.difficulty_level, c.duration_hours,
			c.max_students, c.current_students, c.price, c.thumbnail_url, c.status,
			c.created_at, COALESCE(c.updated_at, ''),
			COALESCE(AVG(cr.rating), 0) as average_rating,
			COUNT(cr.id) as review_count
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
		WHERE c.status = 'Published'
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
			&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
			&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
			&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
		)
		if err != nil {
			utils.HandleError(err)
			continue
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}

// GetCourseByID returns a single course with its modules and reviews
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	courseIDStr := r.URL.Query().Get("id")
	if courseIDStr == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Course ID is required"})
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	// Get course details
	var course Course
	err = database.QueryRow(`
		SELECT 
			c.id, c.title, c.description, c.instructor_id, u.username as instructor_name,
			c.skill_id, s.name as skill_name, c.difficulty_level, c.duration_hours,
			c.max_students, c.current_students, c.price, c.thumbnail_url, c.status,
			c.created_at, COALESCE(c.updated_at, ''),
			COALESCE(AVG(cr.rating), 0) as average_rating,
			COUNT(cr.id) as review_count
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
		WHERE c.id = ? AND c.status = 'Published'
		GROUP BY c.id
	`, courseID).Scan(
		&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
		&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
		&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
		&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
	)

	if err == sql.ErrNoRows {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
		return
	}

	// Get course modules
	moduleRows, err := database.Query(`
		SELECT id, course_id, title, description, order_index, created_at
		FROM course_modules
		WHERE course_id = ?
		ORDER BY order_index ASC
	`, courseID)
	if err != nil {
		utils.HandleError(err)
	} else {
		defer moduleRows.Close()
		var modules []CourseModule
		for moduleRows.Next() {
			var module CourseModule
			err := moduleRows.Scan(&module.ID, &module.CourseID, &module.Title, &module.Description, &module.OrderIndex, &module.CreatedAt)
			if err == nil {
				modules = append(modules, module)
			}
		}
		courseDetail := CourseDetail{
			Course:  course,
			Modules: modules,
		}

		// Get course reviews
		reviewRows, err := database.Query(`
			SELECT cr.id, cr.course_id, cr.student_id, u.username as student_name, cr.rating, cr.review_text, cr.created_at
			FROM course_reviews cr
			JOIN users u ON cr.student_id = u.id
			WHERE cr.course_id = ?
			ORDER BY cr.created_at DESC
			LIMIT 10
		`, courseID)
		if err != nil {
			utils.HandleError(err)
		} else {
			defer reviewRows.Close()
			var reviews []CourseReview
			for reviewRows.Next() {
				var review CourseReview
				err := reviewRows.Scan(&review.ID, &review.CourseID, &review.StudentID, &review.StudentName, &review.Rating, &review.ReviewText, &review.CreatedAt)
				if err == nil {
					reviews = append(reviews, review)
				}
			}
			courseDetail.Reviews = reviews
		}

		utils.SendJSONResponse(w, http.StatusOK, courseDetail)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, course)
}

// SearchCourses searches for courses by title, description, or skill
func SearchCourses(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	searchQuery := "%" + requestBody.Query + "%"

	rows, err := database.Query(`
		SELECT 
			c.id, c.title, c.description, c.instructor_id, u.username as instructor_name,
			c.skill_id, s.name as skill_name, c.difficulty_level, c.duration_hours,
			c.max_students, c.current_students, c.price, c.thumbnail_url, c.status,
			c.created_at, COALESCE(c.updated_at, ''),
			COALESCE(AVG(cr.rating), 0) as average_rating,
			COUNT(cr.id) as review_count
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
		WHERE c.status = 'Published' 
		AND (c.title LIKE ? OR c.description LIKE ? OR s.name LIKE ? OR u.username LIKE ?)
		GROUP BY c.id
		ORDER BY c.created_at DESC
		LIMIT 20
	`, searchQuery, searchQuery, searchQuery, searchQuery)

	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search courses"})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
			&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
			&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
			&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
		)
		if err != nil {
			utils.HandleError(err)
			continue
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}

// GetCoursesByInstructor returns all courses by a specific instructor
func GetCoursesByInstructor(w http.ResponseWriter, r *http.Request) {
	instructorIDStr := r.URL.Query().Get("instructor_id")
	if instructorIDStr == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Instructor ID is required"})
		return
	}

	instructorID, err := strconv.ParseInt(instructorIDStr, 10, 64)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid instructor ID"})
		return
	}

	rows, err := database.Query(`
		SELECT 
			c.id, c.title, c.description, c.instructor_id, u.username as instructor_name,
			c.skill_id, s.name as skill_name, c.difficulty_level, c.duration_hours,
			c.max_students, c.current_students, c.price, c.thumbnail_url, c.status,
			c.created_at, COALESCE(c.updated_at, ''),
			COALESCE(AVG(cr.rating), 0) as average_rating,
			COUNT(cr.id) as review_count
		FROM courses c
		JOIN users u ON c.instructor_id = u.id
		JOIN skills s ON c.skill_id = s.id
		LEFT JOIN course_reviews cr ON c.id = cr.course_id
		WHERE c.instructor_id = ? AND c.status = 'Published'
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`, instructorID)

	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.InstructorID, &course.InstructorName,
			&course.SkillID, &course.SkillName, &course.DifficultyLevel, &course.DurationHours,
			&course.MaxStudents, &course.CurrentStudents, &course.Price, &course.ThumbnailURL, &course.Status,
			&course.CreatedAt, &course.UpdatedAt, &course.AverageRating, &course.ReviewCount,
		)
		if err != nil {
			utils.HandleError(err)
			continue
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to process courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}
