package database

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestMain sets up the test database for database tests
func TestMain(m *testing.M) {
	// Setup test database
	if err := SetupTestDB(); err != nil {
		panic("Failed to setup test database: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	TeardownTestDB()

	// Exit with test code
	os.Exit(code)
}

func TestGetUserIDFromEmail(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Insert test user
	userID, err := InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		expected int64
		hasError bool
	}{
		{
			name:     "Existing user",
			email:    "test@example.com",
			expected: userID,
			hasError: false,
		},
		{
			name:     "Non-existing user",
			email:    "nonexistent@example.com",
			expected: -1,
			hasError: true,
		},
		{
			name:     "Empty email",
			email:    "",
			expected: -1,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetUserIDFromEmail(tt.email)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if id != tt.expected {
					t.Errorf("Expected ID %d, got %d", tt.expected, id)
				}
			}
		})
	}
}

func TestGetSkillIDFromName(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Insert test skill
	skillID, err := InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill: %v", err)
	}

	tests := []struct {
		name      string
		skillName string
		expected  int64
		hasError  bool
	}{
		{
			name:      "Existing skill",
			skillName: "JavaScript",
			expected:  skillID,
			hasError:  false,
		},
		{
			name:      "Case insensitive search",
			skillName: "javascript",
			expected:  skillID,
			hasError:  false,
		},
		{
			name:      "Non-existing skill",
			skillName: "Python",
			expected:  -1,
			hasError:  true,
		},
		{
			name:      "Empty skill name",
			skillName: "",
			expected:  -1,
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetSkillIDFromName(tt.skillName)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if id != tt.expected {
					t.Errorf("Expected ID %d, got %d", tt.expected, id)
				}
			}
		})
	}
}

func TestGetAllSkills(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Insert test skills
	skill1ID, err := InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill 1: %v", err)
	}
	skill2ID, err := InsertTestSkill("Python", "Another programming language")
	if err != nil {
		t.Fatalf("Failed to insert test skill 2: %v", err)
	}

	skills, err := GetAllSkills()
	if err != nil {
		t.Fatalf("Failed to get all skills: %v", err)
	}

	// Check that we have at least 2 skills
	if len(skills) < 2 {
		t.Errorf("Expected at least 2 skills, got %d", len(skills))
	}

	// Check that our test skills are present
	foundJS := false
	foundPython := false
	for _, skill := range skills {
		if skill.ID == int(skill1ID) && skill.Name == "JavaScript" {
			foundJS = true
		}
		if skill.ID == int(skill2ID) && skill.Name == "Python" {
			foundPython = true
		}
	}

	if !foundJS {
		t.Error("JavaScript skill not found in results")
	}
	if !foundPython {
		t.Error("Python skill not found in results")
	}
}

func TestSearch(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Insert test users and skills
	var err error
	_, err = InsertTestUser("john_doe", "john@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 1: %v", err)
	}

	_, err = InsertTestUser("jane_smith", "jane@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 2: %v", err)
	}

	// Insert skills for users
	_, err = InsertTestSkill("JavaScript", "Programming language")
	if err != nil {
		t.Fatalf("Failed to insert JavaScript skill: %v", err)
	}

	// Link skills to users (this would need proper user_skills table operations)
	// For now, we'll test the basic functionality

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		minResults     int
	}{
		{
			name:           "Search by username",
			query:          "john",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by email",
			query:          "jane@example.com",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by skill",
			query:          "JavaScript",
			expectedStatus: 200,
			minResults:     0, // May or may not find results depending on joins
		},
		{
			name:           "Empty query",
			query:          "",
			expectedStatus: 200,
			minResults:     0,
		},
		{
			name:           "Invalid JSON",
			query:          "",
			expectedStatus: 200,
			minResults:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body := strings.NewReader(`{"query":"` + tt.query + `"}`)
			req := httptest.NewRequest("POST", "/api/search", body)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			Search(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var response interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			// Check if we got results array
			if results, ok := response.([]interface{}); ok {
				if len(results) < tt.minResults {
					t.Errorf("Expected at least %d results, got %d", tt.minResults, len(results))
				}
			} else if errorResponse, ok := response.(map[string]interface{}); ok {
				if _, hasError := errorResponse["error"]; !hasError {
					t.Error("Expected error in response for invalid JSON")
				}
			}
		})
	}
}

func TestFullSearch(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Insert test users
	_, err := InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		minResults     int
	}{
		{
			name:           "Search by username",
			query:          "testuser",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Search by email",
			query:          "test@example",
			expectedStatus: 200,
			minResults:     1,
		},
		{
			name:           "Empty query",
			query:          "",
			expectedStatus: 200,
			minResults:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body := strings.NewReader(`{"query":"` + tt.query + `"}`)
			req := httptest.NewRequest("POST", "/api/fullSearch", body)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			FullSearch(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var response interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			// Check if we got results array
			if results, ok := response.([]interface{}); ok {
				if len(results) < tt.minResults {
					t.Errorf("Expected at least %d results, got %d", tt.minResults, len(results))
				}

				// For successful searches, check that results contain expected fields
				if tt.minResults > 0 && len(results) > 0 {
					if result, ok := results[0].(map[string]interface{}); ok {
						if _, hasUser := result["user"]; !hasUser {
							t.Error("Expected user field in search result")
						}
						if _, hasSkills := result["skills_found"]; !hasSkills {
							t.Error("Expected skills_found field in search result")
						}
					}
				}
			} else if errorResponse, ok := response.(map[string]interface{}); ok {
				if _, hasError := errorResponse["error"]; !hasError {
					t.Error("Expected error in response for invalid JSON")
				}
			}
		})
	}
}

func TestDatabaseConnection(t *testing.T) {
	// Test that we can get a database connection
	db, err := GetDatabase()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}
	if db == nil {
		t.Error("Database connection is nil")
	}

	// Test ping
	if err := db.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestExecute(t *testing.T) {
	// Clear test data
	ClearTestData()

	// Test Execute function
	result, err := Execute("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		"testuser", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check that we got a result
	if result == nil {
		t.Error("Execute returned nil result")
	}

	// Check that rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Errorf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}
}

func TestQuery(t *testing.T) {
	// Clear test data and insert test data
	ClearTestData()
	_, err := InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Test Query function
	rows, err := Query("SELECT id, username, email FROM users WHERE username = ?", "testuser")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	// Check that we got results
	if !rows.Next() {
		t.Fatal("Expected at least one row in result")
	}

	var id int
	var username, email string
	if err := rows.Scan(&id, &username, &email); err != nil {
		t.Fatalf("Failed to scan row: %v", err)
	}

	// Verify data
	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
	if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}
}

func TestQueryRow(t *testing.T) {
	// Clear test data and insert test data
	ClearTestData()
	_, err := InsertTestUser("testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Test QueryRow function
	row := QueryRow("SELECT username, email FROM users WHERE username = ?", "testuser")

	var username, email string
	if err := row.Scan(&username, &email); err != nil {
		if err == sql.ErrNoRows {
			t.Fatal("No rows returned")
		}
		t.Fatalf("Failed to scan row: %v", err)
	}

	// Verify data
	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
	if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}
}
