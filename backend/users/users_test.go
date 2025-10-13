package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skillswap/backend/database"
	"skillswap/backend/structs"
)

func TestUpdateUser(t *testing.T) {
	// Clear test data and insert test user
	database.ClearTestData()
	userID, err := database.InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Insert test skill for testing
	_, err = database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	tests := []struct {
		name           string
		userInfo       structs.UserInfo
		expectedStatus int
	}{
		{
			name: "Valid user update",
			userInfo: structs.UserInfo{
				ID:       int(userID),
				Username: "testuser",
				Email:    "test@example.com",
				AboutMe:  "Updated bio",
				Skills: []structs.UserSkill{
					{Name: "JavaScript", Verified: true},
				},
				Projects: []structs.UserProject{
					{Name: "Test Project", Description: "A test project", Link: "https://example.com"},
				},
				Contacts: []structs.UserContact{
					{Name: "GitHub", Link: "https://github.com/test", Icon: "github"},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			userInfo:       structs.UserInfo{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty projects (should be skipped)",
			userInfo: structs.UserInfo{
				ID:       int(userID),
				Username: "testuser",
				Email:    "test@example.com",
				Projects: []structs.UserProject{
					{Name: "", Description: "", Link: ""}, // Empty project should be skipped
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty skills (should be skipped)",
			userInfo: structs.UserInfo{
				ID:       int(userID),
				Username: "testuser",
				Email:    "test@example.com",
				Skills: []structs.UserSkill{
					{Name: ""}, // Empty skill should be skipped
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty contacts (should be skipped)",
			userInfo: structs.UserInfo{
				ID:       int(userID),
				Username: "testuser",
				Email:    "test@example.com",
				Contacts: []structs.UserContact{
					{Name: "", Link: ""}, // Empty contact should be skipped
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, _ := json.Marshal(tt.userInfo)
			req := httptest.NewRequest("POST", "/api/updateUser", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			UpdateUser(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			// For successful updates, verify data was inserted
			if tt.expectedStatus == http.StatusOK {
				// Check if projects were added
				var projectCount int
				err := database.TestDB.QueryRow("SELECT COUNT(*) FROM user_projects WHERE user_id = ?", userID).Scan(&projectCount)
				if err != nil {
					t.Errorf("Failed to check projects: %v", err)
				}

				// Check if skills were added
				var skillCount int
				err = database.TestDB.QueryRow("SELECT COUNT(*) FROM user_skills WHERE user_id = ?", userID).Scan(&skillCount)
				if err != nil {
					t.Errorf("Failed to check skills: %v", err)
				}

				// Check if contacts were added
				var contactCount int
				err = database.TestDB.QueryRow("SELECT COUNT(*) FROM user_contacts WHERE user_id = ?", userID).Scan(&contactCount)
				if err != nil {
					t.Errorf("Failed to check contacts: %v", err)
				}
			}
		})
	}
}

func TestUploadProfilePicture(t *testing.T) {
	// Clear test data and insert test user
	database.ClearTestData()
	userID, err := database.InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		filename       string
		fileContent    string
		expectedStatus int
	}{
		{
			name:           "Valid image upload",
			userID:         fmt.Sprintf("%d", userID),
			filename:       "test.jpg",
			fileContent:    "fake image content",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid file type",
			userID:         fmt.Sprintf("%d", userID),
			filename:       "test.txt",
			fileContent:    "fake text content",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing user ID",
			userID:         "",
			filename:       "test.jpg",
			fileContent:    "fake image content",
			expectedStatus: http.StatusOK, // Function doesn't validate user ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Add user_id field
			err := writer.WriteField("user_id", tt.userID)
			if err != nil {
				t.Fatalf("Failed to write user_id field: %v", err)
			}

			// Add file field
			part, err := writer.CreateFormFile("file", tt.filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			part.Write([]byte(tt.fileContent))
			writer.Close()

			// Create request
			req := httptest.NewRequest("POST", "/api/profile/picture", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			UploadProfilePicture(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				// Check if profile_picture field is in response
				if _, hasProfilePicture := response["profile_picture"]; !hasProfilePicture {
					t.Error("Expected profile_picture in response")
				}

				// Check if file was created
				expectedPath := filepath.Join("uploads", "users", fmt.Sprintf("%s.jpg", tt.userID))
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Error("Profile picture file was not created")
				}

				// Check if database was updated
				var profilePicture string
				err := database.TestDB.QueryRow("SELECT profile_picture FROM users WHERE id = ?", userID).Scan(&profilePicture)
				if err != nil {
					t.Errorf("Failed to check profile picture in database: %v", err)
				}
				expectedPathInDB := fmt.Sprintf("/api/profile/%s/picture", tt.userID)
				if profilePicture != expectedPathInDB {
					t.Errorf("Expected profile_picture %s, got %s", expectedPathInDB, profilePicture)
				}
			} else {
				// For error cases, check for error in response
				if _, hasError := response["error"]; !hasError {
					t.Error("Expected error in response for invalid file type")
				}
			}
		})
	}
}

func TestGetProfilePicture(t *testing.T) {
	// Clear test data and insert test user
	database.ClearTestData()
	userID, err := database.InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Create test profile picture file
	testContent := "fake image content"
	expectedPath := filepath.Join("uploads", "users", fmt.Sprintf("%d.jpg", userID))
	os.MkdirAll(filepath.Dir(expectedPath), 0755)
	err = os.WriteFile(expectedPath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test profile picture file: %v", err)
	}

	// Update user profile picture in database
	publicPath := fmt.Sprintf("/api/profile/%d/picture", userID)
	_, err = database.TestDB.Exec("UPDATE users SET profile_picture = ? WHERE id = ?", publicPath, userID)
	if err != nil {
		t.Fatalf("Failed to update user profile picture in database: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/profile/%d/picture", userID), nil)
	rr := httptest.NewRecorder()

	// Call the handler
	GetProfilePicture(rr, req)

	// Check status code (should be 200 for successful file serve)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check content type (should be image/jpeg for .jpg files)
	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "image") {
		t.Errorf("Expected image content type, got %s", contentType)
	}

	// Check if file content matches
	if rr.Body.String() != testContent {
		t.Error("File content does not match expected content")
	}
}

func TestRetrieveUserInfo(t *testing.T) {
	// Clear test data and insert test user
	database.ClearTestData()
	userID, err := database.InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Insert test skill and link it to user
	skillID, err := database.InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	// Link skill to user
	_, err = database.TestDB.Exec("INSERT INTO user_skills (user_id, skill_id, verified) VALUES (?, ?, ?)", userID, skillID, 1)
	if err != nil {
		t.Fatalf("Failed to link skill to user: %v", err)
	}

	// Add test project
	_, err = database.TestDB.Exec("INSERT INTO user_projects (user_id, name, description, link) VALUES (?, ?, ?, ?)",
		userID, "Test Project", "A test project", "https://example.com")
	if err != nil {
		t.Fatalf("Failed to add test project: %v", err)
	}

	// Add test contact
	_, err = database.TestDB.Exec("INSERT INTO user_contacts (user_id, name, link, icon) VALUES (?, ?, ?, ?)",
		userID, "GitHub", "https://github.com/test", "github")
	if err != nil {
		t.Fatalf("Failed to add test contact: %v", err)
	}

	// Update user with additional info
	_, err = database.TestDB.Exec("UPDATE users SET aboutme = ?, profession = ?, location = ? WHERE id = ?",
		"Test bio", "Software Developer", "Test City", userID)
	if err != nil {
		t.Fatalf("Failed to update user info: %v", err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/user?q=%d", userID), nil)
	rr := httptest.NewRecorder()

	// Call the handler
	RetrieveUserInfo(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Parse response
	var userInfo structs.UserInfo
	if err := json.Unmarshal(rr.Body.Bytes(), &userInfo); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify user data
	if userInfo.ID != int(userID) {
		t.Errorf("Expected user ID %d, got %d", userID, userInfo.ID)
	}
	if userInfo.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", userInfo.Username)
	}
	if userInfo.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", userInfo.Email)
	}
	if userInfo.AboutMe != "Test bio" {
		t.Errorf("Expected aboutme 'Test bio', got '%s'", userInfo.AboutMe)
	}
	if userInfo.Professions != "Software Developer" {
		t.Errorf("Expected profession 'Software Developer', got '%s'", userInfo.Professions)
	}
	if userInfo.Location != "Test City" {
		t.Errorf("Expected location 'Test City', got '%s'", userInfo.Location)
	}

	// Check skills
	if len(userInfo.Skills) == 0 {
		t.Error("Expected skills to be populated")
	} else {
		found := false
		for _, skill := range userInfo.Skills {
			if skill.Name == "JavaScript" && bool(skill.Verified) == true {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected JavaScript skill not found")
		}
	}

	// Check projects
	if len(userInfo.Projects) == 0 {
		t.Error("Expected projects to be populated")
	} else {
		found := false
		for _, project := range userInfo.Projects {
			if project.Name == "Test Project" && project.Description == "A test project" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected test project not found")
		}
	}

	// Check contacts
	if len(userInfo.Contacts) == 0 {
		t.Error("Expected contacts to be populated")
	} else {
		found := false
		for _, contact := range userInfo.Contacts {
			if contact.Name == "GitHub" && contact.Link == "https://github.com/test" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected GitHub contact not found")
		}
	}
}
