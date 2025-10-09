package courses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"skillswap/backend/database"
)

// TestMain sets up the test environment for course tests
func TestMain(m *testing.M) {
	// Setup test database
	if err := database.SetupTestDB(); err != nil {
		fmt.Printf("Failed to setup test database: %v\n", err)
		os.Exit(1)
	}

	// Create uploads directories for testing
	os.MkdirAll("uploads/course_thumbnails", 0755)
	os.MkdirAll("uploads/courses", 0755)

	code := m.Run()

	// Cleanup
	database.TeardownTestDB()
	os.RemoveAll("uploads")

	os.Exit(code)
}

func TestGetAllCourses(t *testing.T) {
	// Clear test data and insert test courses
	database.ClearTestData()

	// Insert test instructor and skill
	instructorID, err := database.InsertTestUser("instructor", "instructor@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test instructor: %v", err)
	}

	skillID, err := database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	// Insert test course
	courseID, err := database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "Test Course", "A test course description", instructorID, skillID, 10)
	if err != nil {
		t.Fatalf("Failed to insert test course: %v", err)
	}

	courseIDInt64, _ := courseID.LastInsertId()

	// Create request
	req := httptest.NewRequest("GET", "/api/courses", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	GetAllCourses(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Parse response
	var courses []Course
	if err := json.Unmarshal(rr.Body.Bytes(), &courses); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Check that we got courses
	if len(courses) == 0 {
		t.Error("Expected courses in response")
	}

	// Check course data
	found := false
	for _, course := range courses {
		if course.ID == courseIDInt64 && course.Title == "Test Course" {
			found = true
			if course.InstructorName != "instructor" {
				t.Errorf("Expected instructor name 'instructor', got '%s'", course.InstructorName)
			}
			if course.SkillName != "JavaScript" {
				t.Errorf("Expected skill name 'JavaScript', got '%s'", course.SkillName)
			}
		}
	}

	if !found {
		t.Error("Test course not found in response")
	}
}

func TestGetCourseByID(t *testing.T) {
	// Clear test data and insert test course
	database.ClearTestData()

	// Insert test instructor
	instructorID, err := database.InsertTestUser("instructor", "instructor@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test instructor: %v", err)
	}

	// Insert test skill
	skillID, err := database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	// Insert test course
	courseID, err := database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "Test Course", "A test course description", instructorID, skillID, 10)
	if err != nil {
		t.Fatalf("Failed to insert test course: %v", err)
	}

	courseIDInt64, _ := courseID.LastInsertId()

	// Insert test modules
	_, err = database.TestDB.Exec(`
		INSERT INTO course_modules (course_id, title, description, order_index)
		VALUES (?, ?, ?, ?)
	`, courseIDInt64, "Module 1", "First module", 1)
	if err != nil {
		t.Fatalf("Failed to insert test module: %v", err)
	}

	// Insert test review
	_, err = database.TestDB.Exec(`
		INSERT INTO course_reviews (course_id, student_id, rating, review_text)
		VALUES (?, ?, ?, ?)
	`, courseIDInt64, instructorID, 5, "Great course!")
	if err != nil {
		t.Fatalf("Failed to insert test review: %v", err)
	}

	tests := []struct {
		name           string
		courseID       string
		expectedStatus int
	}{
		{
			name:           "Valid course ID",
			courseID:       fmt.Sprintf("%d", courseIDInt64),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing course ID",
			courseID:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid course ID",
			courseID:       "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent course ID",
			courseID:       "99999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/course?id="+tt.courseID, nil)
			rr := httptest.NewRecorder()

			// Call the handler
			GetCourseByID(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var courseDetail CourseDetail
				if err := json.Unmarshal(rr.Body.Bytes(), &courseDetail); err != nil {
					t.Fatalf("Failed to parse response JSON: %v", err)
				}

				// Check course data
				if courseDetail.ID != courseIDInt64 {
					t.Errorf("Expected course ID %d, got %d", courseIDInt64, courseDetail.ID)
				}
				if courseDetail.Title != "Test Course" {
					t.Errorf("Expected title 'Test Course', got '%s'", courseDetail.Title)
				}

				// Check modules
				if len(courseDetail.Modules) == 0 {
					t.Error("Expected modules in course detail")
				}

				// Check reviews
				if len(courseDetail.Reviews) == 0 {
					t.Error("Expected reviews in course detail")
				}
			}
		})
	}
}

func TestSearchCourses(t *testing.T) {
	// Clear test data and insert test course
	database.ClearTestData()

	// Insert test instructor
	instructorID, err := database.InsertTestUser("instructor", "instructor@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test instructor: %v", err)
	}

	// Insert test skill
	skillID, err := database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	// Insert test course
	_, err = database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "JavaScript Basics", "Learn JavaScript fundamentals", instructorID, skillID, 10)
	if err != nil {
		t.Fatalf("Failed to insert test course: %v", err)
	}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		minResults     int
	}{
		{
			name:           "Search by title",
			query:          "JavaScript",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by description",
			query:          "fundamentals",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by skill",
			query:          "JavaScript",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by instructor",
			query:          "instructor",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Empty query",
			query:          "",
			expectedStatus: 200,
			minResults:     0,
		},
		{
			name:           "No results query",
			query:          "nonexistent",
			expectedStatus: 200,
			minResults:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body := strings.NewReader(`{"query":"` + tt.query + `"}`)
			req := httptest.NewRequest("POST", "/api/searchCourses", body)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			SearchCourses(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var courses []Course
			if err := json.Unmarshal(rr.Body.Bytes(), &courses); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			// Check minimum results
			if len(courses) < tt.minResults {
				t.Errorf("Expected at least %d results, got %d", tt.minResults, len(courses))
			}

			// For successful searches, verify course data
			if tt.minResults > 0 && len(courses) > 0 {
				found := false
				for _, course := range courses {
					if strings.Contains(course.Title, "JavaScript") {
						found = true
						if course.InstructorName != "instructor" {
							t.Errorf("Expected instructor name 'instructor', got '%s'", course.InstructorName)
						}
					}
				}
				if !found {
					t.Error("Expected JavaScript course not found in search results")
				}
			}
		})
	}
}

func TestGetCoursesByInstructor(t *testing.T) {
	// Clear test data and insert test courses
	database.ClearTestData()

	// Insert test instructor
	instructorID, err := database.InsertTestUser("instructor", "instructor@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test instructor: %v", err)
	}

	// Insert another instructor for comparison
	otherInstructorID, err := database.InsertTestUser("other_instructor", "other@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert other instructor: %v", err)
	}

	// Insert test skill
	skillID, err := database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	// Insert test courses
	_, err = database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "Course 1", "First course", instructorID, skillID, 10)
	if err != nil {
		t.Fatalf("Failed to insert course 1: %v", err)
	}

	_, err = database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "Course 2", "Second course", instructorID, skillID, 15)
	if err != nil {
		t.Fatalf("Failed to insert course 2: %v", err)
	}

	// Insert course by other instructor
	_, err = database.TestDB.Exec(`
		INSERT INTO courses (title, description, instructor_id, skill_id, duration_hours, status)
		VALUES (?, ?, ?, ?, ?, 'Published')
	`, "Other Course", "Course by other instructor", otherInstructorID, skillID, 8)
	if err != nil {
		t.Fatalf("Failed to insert other course: %v", err)
	}

	tests := []struct {
		name           string
		instructorID   string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Valid instructor ID",
			instructorID:   fmt.Sprintf("%d", instructorID),
			expectedStatus: 200,
			expectedCount:  2,
		},
		{
			name:           "Other instructor ID",
			instructorID:   fmt.Sprintf("%d", otherInstructorID),
			expectedStatus: 200,
			expectedCount:  1,
		},
		{
			name:           "Missing instructor ID",
			instructorID:   "",
			expectedStatus: 400,
		},
		{
			name:           "Invalid instructor ID",
			instructorID:   "invalid",
			expectedStatus: 400,
		},
		{
			name:           "Non-existent instructor ID",
			instructorID:   "99999",
			expectedStatus: 200,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/api/coursesByInstructor?instructor_id="+tt.instructorID, nil)
			rr := httptest.NewRecorder()

			// Call the handler
			GetCoursesByInstructor(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == 200 && tt.expectedCount >= 0 {
				// Parse response
				var courses []Course
				if err := json.Unmarshal(rr.Body.Bytes(), &courses); err != nil {
					t.Fatalf("Failed to parse response JSON: %v", err)
				}

				// Check course count
				if len(courses) != tt.expectedCount {
					t.Errorf("Expected %d courses, got %d", tt.expectedCount, len(courses))
				}

				// Verify instructor name for all courses
				for _, course := range courses {
					if course.InstructorID != instructorID && course.InstructorID != otherInstructorID {
						t.Errorf("Course has unexpected instructor ID: %d", course.InstructorID)
					}
				}
			}
		})
	}
}

func TestAddCourse(t *testing.T) {
	// Clear test data and insert test instructor and skill
	database.ClearTestData()

	// Insert test instructor and skill
	_, err := database.InsertTestUser("instructor", "instructor@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test instructor: %v", err)
	}

	_, err = database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	tests := []struct {
		name           string
		title          string
		description    string
		skillName      string
		durationMinutes string
		hasFile        bool
		expectedStatus int
	}{
		{
			name:           "Valid course creation",
			title:          "Test Course",
			description:    "A test course description",
			skillName:      "JavaScript",
			durationMinutes: "120",
			hasFile:        false,
			expectedStatus: 200,
		},
		{
			name:           "Missing title",
			title:          "",
			description:    "A test course description",
			skillName:      "JavaScript",
			durationMinutes: "120",
			hasFile:        false,
			expectedStatus: 400,
		},
		{
			name:           "Missing description",
			title:          "Test Course",
			description:    "",
			skillName:      "JavaScript",
			durationMinutes: "120",
			hasFile:        false,
			expectedStatus: 400,
		},
		{
			name:           "Missing skill name",
			title:          "Test Course",
			description:    "A test course description",
			skillName:      "",
			durationMinutes: "120",
			hasFile:        false,
			expectedStatus: 400,
		},
		{
			name:           "Invalid skill name",
			title:          "Test Course",
			description:    "A test course description",
			skillName:      "NonExistentSkill",
			durationMinutes: "120",
			hasFile:        false,
			expectedStatus: 400,
		},
		{
			name:           "Invalid duration",
			title:          "Test Course",
			description:    "A test course description",
			skillName:      "JavaScript",
			durationMinutes: "invalid",
			hasFile:        false,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Add form fields
			writer.WriteField("title", tt.title)
			writer.WriteField("description", tt.description)
			writer.WriteField("skill_name", tt.skillName)
			writer.WriteField("duration_minutes", tt.durationMinutes)

			writer.Close()

			// Create request
			req := httptest.NewRequest("POST", "/api/course/add", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			AddCourse(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == 200 {
				// Parse response
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response JSON: %v", err)
				}

				// Check response fields
				if _, hasMessage := response["message"]; !hasMessage {
					t.Error("Expected message in response")
				}
				if _, hasCourseID := response["course_id"]; !hasCourseID {
					t.Error("Expected course_id in response")
				}
				if _, hasSkillName := response["skill_name"]; !hasSkillName {
					t.Error("Expected skill_name in response")
				}

				// Verify course was created in database
				var count int
				err := database.TestDB.QueryRow("SELECT COUNT(*) FROM courses WHERE title = ? AND description = ?",
					tt.title, tt.description).Scan(&count)
				if err != nil {
					t.Errorf("Failed to check course creation: %v", err)
				}
				if count == 0 {
					t.Error("Course was not created in database")
				}
			}
		})
	}
}
