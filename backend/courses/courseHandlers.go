package courses

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"skillswap/backend/authentication"
	"skillswap/backend/database"
	"skillswap/backend/utils"
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
		SELECT id, course_id, title, description, COALESCE(video_url, ''), video_duration, COALESCE(thumbnail_url, ''), order_index, created_at
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
			err := moduleRows.Scan(&module.ID, &module.CourseID, &module.Title, &module.Description, &module.VideoURL, &module.VideoDuration, &module.ThumbnailURL, &module.OrderIndex, &module.CreatedAt)
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
func AddCourse(w http.ResponseWriter, r *http.Request) {
	// Limit upload size (e.g., 500 MB for videos)
	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)
	utils.DebugPrint("Course upload started")

	if err := r.ParseMultipartForm(500 << 20); err != nil {
		utils.DebugPrint("File too large or invalid form data")
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "File too large or invalid form data"})
		return
	}

	// Extract course basic fields
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	skillName := strings.TrimSpace(r.FormValue("skill_name"))
	durationMinutesStr := r.FormValue("duration_minutes")

	if title == "" || description == "" || skillName == "" || durationMinutesStr == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
		return
	}

	durationMinutes, err := strconv.Atoi(durationMinutesStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid duration value"})
		return
	}

	// Get user from session (middleware ensures authentication)
	session, err := authentication.Store.Get(r, "authentication")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
		return
	}

	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid session"})
		return
	}
	
	utils.DebugPrint("User creating course:", email)
	instructorID, err := database.GetUserIDFromEmail(email)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get user ID"})
		return
	}

	// Resolve skill name → skill_id
	var skillID int64
	err = database.QueryRow("SELECT id FROM skills WHERE name = ?", skillName).Scan(&skillID)
	if err == sql.ErrNoRows {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid skill name"})
		return
	} else if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch skill"})
		return
	}

	// Handle course thumbnail upload
	var thumbnailURL string
	previewFile, previewHeader, err := r.FormFile("preview_photo")
	if err == nil {
		defer previewFile.Close()

		// Validate preview photo file type
		if !utils.CheckType(filepath.Ext(previewHeader.Filename), []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}) {
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid preview photo file type. Only images are allowed."})
			return
		}

		uploadDir := "./uploads/course_thumbnails"
		os.MkdirAll(uploadDir, os.ModePerm)

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(previewHeader.Filename))
		filePath := filepath.Join(uploadDir, filename)

		out, err := os.Create(filePath)
		if err != nil {
			utils.HandleError(err)
			utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save preview photo"})
			return
		}
		defer out.Close()
		io.Copy(out, previewFile)

		thumbnailURL = "/uploads/course_thumbnails/" + filename
	}

	// Insert course into DB
	res, err := database.Execute(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, thumbnail_url, status)
		VALUES (?, ?, ?, ?, ?, ?, 'Published')
	`, title, description, instructorID, skillID, durationMinutes/60, thumbnailURL)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create course"})
		return
	}

	courseID, _ := res.LastInsertId()
	utils.DebugPrint(fmt.Sprintf("Course created with ID: %d", courseID))

	// Parse module data from form
	modulesJSON := r.FormValue("modules")
	if modulesJSON != "" {
		var modules []struct {
			Title         string `json:"title"`
			Description   string `json:"description"`
			OrderIndex    int    `json:"order_index"`
			VideoDuration int    `json:"video_duration"`
		}

		err = json.Unmarshal([]byte(modulesJSON), &modules)
		if err != nil {
			utils.HandleError(err)
			utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid modules data format"})
			return
		}

		// Create directories for uploads
		videoUploadDir := "./uploads/courses/videos"
		moduleThumbnailDir := "./uploads/courses/module_thumbnails"
		os.MkdirAll(videoUploadDir, os.ModePerm)
		os.MkdirAll(moduleThumbnailDir, os.ModePerm)

		// Process each module
		for i, module := range modules {
			var videoURL, moduleThumbnailURL string

			// Handle video file upload for this module
			videoFieldName := fmt.Sprintf("module_%d_video", i)
			videoFile, videoHeader, err := r.FormFile(videoFieldName)
			if err == nil {
				defer videoFile.Close()

				// Validate video file type
				videoExt := filepath.Ext(videoHeader.Filename)
				if !utils.CheckType(videoExt, []string{".mp4", ".avi", ".mov", ".webm", ".mkv"}) {
					utils.DebugPrint(fmt.Sprintf("Invalid video type for module %d: %s", i, videoExt))
					continue
				}

				videoFilename := fmt.Sprintf("%d_module_%d_%s", time.Now().UnixNano(), i, filepath.Base(videoHeader.Filename))
				videoPath := filepath.Join(videoUploadDir, videoFilename)

				out, err := os.Create(videoPath)
				if err != nil {
					utils.HandleError(err)
					continue
				}
				io.Copy(out, videoFile)
				out.Close()

				videoURL = "/uploads/courses/videos/" + videoFilename
			}

			// Handle optional module thumbnail
			thumbnailFieldName := fmt.Sprintf("module_%d_thumbnail", i)
			thumbFile, thumbHeader, err := r.FormFile(thumbnailFieldName)
			if err == nil {
				defer thumbFile.Close()

				thumbExt := filepath.Ext(thumbHeader.Filename)
				if utils.CheckType(thumbExt, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}) {
					thumbFilename := fmt.Sprintf("%d_module_%d_thumb_%s", time.Now().UnixNano(), i, filepath.Base(thumbHeader.Filename))
					thumbPath := filepath.Join(moduleThumbnailDir, thumbFilename)

					out, err := os.Create(thumbPath)
					if err == nil {
						io.Copy(out, thumbFile)
						out.Close()
						moduleThumbnailURL = "/uploads/courses/module_thumbnails/" + thumbFilename
					}
				}
			}

			// Insert module into database
			_, err = database.Execute(`
				INSERT INTO course_modules (course_id, title, description, video_url, video_duration, thumbnail_url, order_index)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, courseID, module.Title, module.Description, videoURL, module.VideoDuration, moduleThumbnailURL, module.OrderIndex)

			if err != nil {
				utils.HandleError(err)
				utils.DebugPrint(fmt.Sprintf("Failed to insert module %d: %v", i, err))
			} else {
				utils.DebugPrint(fmt.Sprintf("Module %d inserted successfully", i))
			}
		}
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]any{
		"message":    "Course uploaded successfully",
		"course_id":  courseID,
		"thumbnail":  thumbnailURL,
		"skill_name": skillName,
	})
}

func ServeModuleVideo(w http.ResponseWriter, r *http.Request) {
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Video path required"})
		return
	}
	if !strings.HasPrefix(videoPath, "/uploads/courses/videos/") {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": "Invalid path"})
		return
	}
	filePath := "." + videoPath
	http.ServeFile(w, r, filePath)
}
