package courses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"skillswap/backend/internal/utils"
)

// GetAllCourses returns all published courses
func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	query := courseWithRatingsQuery + `
		WHERE c.status = 'Published'
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	courses, err := fetchCourses(query)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}

// GetCourseByID returns a single course with its modules and reviews
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	courseIDStr := r.URL.Query().Get("id")
	courseID, err := validateAndParseInt64(courseIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid course ID"})
		return
	}

	// Fetch course details
	course, err := fetchCourseByID(courseID)
	if err == sql.ErrNoRows {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
		return
	}
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch course"})
		return
	}

	// Fetch course modules
	modules, err := fetchCourseModules(courseID)
	if err != nil {
		utils.HandleError(err)
		modules = []CourseModule{}
	}

	// Fetch course reviews
	reviews, err := fetchCourseReviews(courseID)
	if err != nil {
		utils.HandleError(err)
		reviews = []CourseReview{}
	}

	courseDetail := CourseDetail{
		Course:  *course,
		Modules: modules,
		Reviews: reviews,
	}

	utils.SendJSONResponse(w, http.StatusOK, courseDetail)
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

	searchPattern := "%" + requestBody.Query + "%"
	query := courseWithRatingsQuery + `
		WHERE c.status = 'Published'
		AND (c.title LIKE ? OR c.description LIKE ? OR s.name LIKE ? OR u.username LIKE ?)
		GROUP BY c.id
		ORDER BY c.created_at DESC
		LIMIT 20
	`

	courses, err := fetchCourses(query, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to search courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}

// GetCoursesByInstructor returns all courses by a specific instructor
func GetCoursesByInstructor(w http.ResponseWriter, r *http.Request) {
	instructorIDStr := r.URL.Query().Get("instructor_id")
	instructorID, err := validateAndParseInt64(instructorIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid instructor ID"})
		return
	}

	query := courseWithRatingsQuery + `
		WHERE c.instructor_id = ? AND c.status = 'Published'
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	courses, err := fetchCourses(query, instructorID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, courses)
}

// AddCourse handles the creation of a new course with modules
func AddCourse(w http.ResponseWriter, r *http.Request) {
	// Limit upload size (500 MB for videos)
	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)
	// utils.DebugPrint("Course upload started")

	if err := r.ParseMultipartForm(500 << 20); err != nil {
		// utils.DebugPrint("File too large or invalid form data")
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "File too large or invalid form data"})
		return
	}

	// Validate and extract course form data
	formData, err := validateCourseFormData(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Parse module data
	modules, err := parseModulesFromForm(r)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid modules data format"})
		return
	}

	// Create course with modules
	courseID, thumbnailURL, err := createCourseWithModules(r, formData, modules)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]any{
		"message":    "Course uploaded successfully",
		"course_id":  courseID,
		"thumbnail":  thumbnailURL,
		"skill_name": formData.SkillName,
	})
}

// ServeModuleVideo serves course module videos
func ServeModuleVideo(w http.ResponseWriter, r *http.Request) {
	videoPath := r.URL.Query().Get("path")
	if videoPath == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Video path required"})
		return
	}

	// Security check: only allow videos from the courses directory
	if !strings.HasPrefix(videoPath, "/uploads/courses/videos/") {
		utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": "Invalid path"})
		return
	}

	filePath := "." + videoPath
	http.ServeFile(w, r, filePath)
}
