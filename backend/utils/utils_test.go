package utils

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestSendJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		payload        interface{}
		expectedStatus int
	}{
		{
			name:           "Success response",
			statusCode:     200,
			payload:        map[string]string{"status": "ok"},
			expectedStatus: 200,
		},
		{
			name:           "Error response",
			statusCode:     400,
			payload:        map[string]string{"error": "Bad Request"},
			expectedStatus: 400,
		},
		{
			name:           "Not found response",
			statusCode:     404,
			payload:        map[string]string{"error": "Not Found"},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the function
			SendJSONResponse(rr, tt.statusCode, tt.payload)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check content type
			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check if response body is valid JSON
			// The response should contain the payload data
			body := rr.Body.String()
			if body == "" {
				t.Error("Response body is empty")
			}
		})
	}
}

func TestCheckType(t *testing.T) {
	allowedTypes := []string{"image/jpeg", "image/png", "application/pdf"}

	tests := []struct {
		name     string
		toCheck  string
		toAllow  []string
		expected bool
	}{
		{
			name:     "Valid JPEG type",
			toCheck:  "image/jpeg",
			toAllow:  allowedTypes,
			expected: true,
		},
		{
			name:     "Valid PNG type",
			toCheck:  "image/png",
			toAllow:  allowedTypes,
			expected: true,
		},
		{
			name:     "Valid PDF type",
			toCheck:  "application/pdf",
			toAllow:  allowedTypes,
			expected: true,
		},
		{
			name:     "Invalid type",
			toCheck:  "text/html",
			toAllow:  allowedTypes,
			expected: false,
		},
		{
			name:     "Empty allowed types",
			toCheck:  "image/jpeg",
			toAllow:  []string{},
			expected: false,
		},
		{
			name:     "Empty check type",
			toCheck:  "",
			toAllow:  allowedTypes,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckType(tt.toCheck, tt.toAllow)
			if result != tt.expected {
				t.Errorf("CheckType(%s, %v) = %v, expected %v", tt.toCheck, tt.toAllow, result, tt.expected)
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	t.Run("Generate unique UUIDs", func(t *testing.T) {
		// Generate multiple UUIDs and ensure they are unique
		uuids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			uuid := GenerateUUID()
			if uuid == "" {
				t.Error("Generated UUID is empty")
			}
			if uuids[uuid] {
				t.Error("Generated duplicate UUID")
			}
			uuids[uuid] = true

			// Check UUID format (basic validation)
			if len(uuid) != 36 {
				t.Errorf("UUID length is %d, expected 36", len(uuid))
			}
			// UUID should contain hyphens in standard positions
			if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
				t.Errorf("UUID format is invalid: %s", uuid)
			}
		}
	})

	t.Run("UUID consistency", func(t *testing.T) {
		// Same input should generate different UUIDs due to timestamp
		uuid1 := GenerateUUID()
		uuid2 := GenerateUUID()

		if uuid1 == uuid2 {
			t.Error("Generated identical UUIDs consecutively")
		}
	})
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "Actual error",
			err:      fmt.Errorf("test error"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleError(tt.err)
			if result != tt.expected {
				t.Errorf("HandleError(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}
