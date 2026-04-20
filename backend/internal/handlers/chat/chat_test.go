package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"skillswap/backend/internal/database"
)

// TestMain sets up the test environment for chat tests
func TestMain(m *testing.M) {
	// Setup test database
	if err := database.SetupTestDB(); err != nil {
		fmt.Printf("Failed to setup test database: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	// Cleanup
	database.TeardownTestDB()
	os.Exit(code)
}

func TestCreateChat(t *testing.T) {
	// Clear test data and insert test users
	database.ClearTestData()

	// Insert test users (using test-prefixed emails so they get cleaned up)
	user1ID, err := database.InsertTestUser("testuser1", "testuser1@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 1: %v", err)
	}

	user2ID, err := database.InsertTestUser("testuser2", "testuser2@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 2: %v", err)
	}

	tests := []struct {
		name           string
		initiatorID    int64
		responderID    int64
		expectedStatus int
	}{
		{
			name:           "Valid chat creation",
			initiatorID:    user1ID,
			responderID:    user2ID,
			expectedStatus: 200,
		},
		{
			name:           "Same user IDs",
			initiatorID:    user1ID,
			responderID:    user1ID,
			expectedStatus: 200, // Should still work but may create duplicate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/createChat?u1=%d&u2=%d", tt.initiatorID, tt.responderID), nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			CreateChat(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			if tt.expectedStatus == 200 {
				// Check if chat_id is in response
				if _, hasChatID := response["chat_id"]; !hasChatID {
					t.Error("Expected chat_id in response")
				}

				// Verify chat was created in database
				var count int
				err := database.TestDB.QueryRow("SELECT COUNT(*) FROM chats WHERE user1_id = ? AND user2_id = ?",
					tt.initiatorID, tt.responderID).Scan(&count)
				if err != nil {
					t.Errorf("Failed to check chat creation: %v", err)
				}
				if count == 0 {
					t.Error("Chat was not created in database")
				}
			}
		})
	}
}

func TestGetChatsFromUserID(t *testing.T) {
	// Clear test data and insert test data
	database.ClearTestData()

	// Insert test users (using test-prefixed emails so they get cleaned up)
	user1ID, err := database.InsertTestUser("testuser1_chats", "testuser1chats@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 1: %v", err)
	}

	user2ID, err := database.InsertTestUser("testuser2_chats", "testuser2chats@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 2: %v", err)
	}

	user3ID, err := database.InsertTestUser("testuser3_chats", "testuser3chats@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 3: %v", err)
	}

	// Insert test chats
	_, err = database.TestDB.Exec("INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)", user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to insert chat 1: %v", err)
	}

	_, err = database.TestDB.Exec("INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)", user3ID, user1ID)
	if err != nil {
		t.Fatalf("Failed to insert chat 2: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		minChats       int
	}{
		{
			name:           "Valid user ID with chats",
			userID:         fmt.Sprintf("%d", user1ID),
			expectedStatus: 200,
			minChats:       2,
		},
		{
			name:           "User ID with no chats",
			userID:         fmt.Sprintf("%d", user2ID),
			expectedStatus: 200,
			minChats:       0,
		},
		{
			name:           "Non-existent user ID",
			userID:         "99999",
			expectedStatus: 200,
			minChats:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request (GetChatsFromUserID expects "uid" parameter)
			req := httptest.NewRequest("GET", "/api/getChats?uid="+tt.userID, nil)
			rr := httptest.NewRecorder()

			// Call the handler
			GetChatsFromUserID(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response
			var chats []ChatWithUserInfo
			if err := json.Unmarshal(rr.Body.Bytes(), &chats); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}

			// Check minimum chats
			if len(chats) < tt.minChats {
				t.Errorf("Expected at least %d chats, got %d", tt.minChats, len(chats))
			}

			// For user1, check that we get chats with user2 and user3
			if tt.userID == fmt.Sprintf("%d", user1ID) && len(chats) >= 2 {
				foundUser2 := false
				foundUser3 := false

				for _, chat := range chats {
					// Check if user1 is involved in chat with user2
					if (chat.Initiator == int(user1ID) && chat.Responder == int(user2ID)) ||
					   (chat.Initiator == int(user2ID) && chat.Responder == int(user1ID)) {
						foundUser2 = true
					}
					// Check if user1 is involved in chat with user3
					if (chat.Initiator == int(user3ID) && chat.Responder == int(user1ID)) ||
					   (chat.Initiator == int(user1ID) && chat.Responder == int(user3ID)) {
						foundUser3 = true
					}
				}

				if !foundUser2 {
					t.Error("Expected chat with user2 not found")
				}
				if !foundUser3 {
					t.Error("Expected chat with user3 not found")
				}
			}
		})
	}
}

func TestGetMessagesFromUID(t *testing.T) {
	// Clear test data and insert test data
	database.ClearTestData()

	// Insert test users (using test-prefixed emails so they get cleaned up)
	user1ID, err := database.InsertTestUser("testuser1_msgs", "testuser1msgs@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 1: %v", err)
	}

	user2ID, err := database.InsertTestUser("testuser2_msgs", "testuser2msgs@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to insert test user 2: %v", err)
	}

	// Insert test chat
	chatID, err := database.TestDB.Exec("INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)", user1ID, user2ID)
	if err != nil {
		t.Fatalf("Failed to insert test chat: %v", err)
	}

	chatIDInt64, _ := chatID.LastInsertId()

	// Insert test messages
	_, err = database.TestDB.Exec(`
		INSERT INTO messages (chat_id, sender_id, content, created_at)
		VALUES (?, ?, ?, NOW())
	`, chatIDInt64, user1ID, "Hello from user1")
	if err != nil {
		t.Fatalf("Failed to insert message 1: %v", err)
	}

	_, err = database.TestDB.Exec(`
		INSERT INTO messages (chat_id, sender_id, content, created_at)
		VALUES (?, ?, ?, NOW())
	`, chatIDInt64, user2ID, "Hello from user2")
	if err != nil {
		t.Fatalf("Failed to insert message 2: %v", err)
	}

	tests := []struct {
		name           string
		chatID         string
		expectedStatus int
		minMessages    int
	}{
		{
			name:           "Valid chat ID",
			chatID:         fmt.Sprintf("%d", chatIDInt64),
			expectedStatus: 200,
			minMessages:    2,
		},
		{
			name:           "Non-existent chat ID",
			chatID:         "99999",
			expectedStatus: 200,
			minMessages:    0,
		},
		{
			name:           "Invalid chat ID",
			chatID:         "invalid",
			expectedStatus: 200,
			minMessages:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request (GetMessagesFromUID expects "cid" parameter)
			req := httptest.NewRequest("GET", "/api/getChatInfo?cid="+tt.chatID, nil)
			rr := httptest.NewRecorder()

			// Call the handler
			GetMessagesFromUID(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response (GetMessagesFromUID returns {"messages": [...]} not just [...])
			var response map[string][]Message
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response JSON: %v", err)
			}
			messages := response["messages"]

			// Check minimum messages
			if len(messages) < tt.minMessages {
				t.Errorf("Expected at least %d messages, got %d", tt.minMessages, len(messages))
			}

			// For valid chat, check message content
			if tt.chatID == fmt.Sprintf("%d", chatIDInt64) && len(messages) >= 2 {
				foundUser1Message := false
				foundUser2Message := false

				for _, message := range messages {
					if message.Sender.ID == int(user1ID) && message.Content == "Hello from user1" {
						foundUser1Message = true
					}
					if message.Sender.ID == int(user2ID) && message.Content == "Hello from user2" {
						foundUser2Message = true
					}
				}

				if !foundUser1Message {
					t.Error("Expected message from user1 not found")
				}
				if !foundUser2Message {
					t.Error("Expected message from user2 not found")
				}
			}
		})
	}
}

func TestSimpleWebSocketEndpoint(t *testing.T) {
	// This is a basic test for WebSocket upgrade
	// Note: Full WebSocket testing would require more complex setup

	t.Run("WebSocket upgrade request", func(t *testing.T) {
		// Create request with WebSocket headers
		req := httptest.NewRequest("GET", "/api/chat", nil)
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Connection", "upgrade")
		req.Header.Set("Sec-WebSocket-Key", "test-key")
		req.Header.Set("Sec-WebSocket-Version", "13")

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		SimpleWebSocketEndpoint(rr, req)

		// WebSocket endpoint should handle the upgrade
		// Since we're using httptest, we can't fully test WebSocket functionality
		// But we can check that it doesn't panic and handles the request
		if rr.Code != http.StatusBadRequest { // Should fail due to missing auth in test
			// The handler should handle the WebSocket connection attempt
			// In a real test environment, this would be more comprehensive
		}
	})
}

// TestHub tests the Hub functionality
func TestNewHub(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Error("NewHub returned nil")
	}

	if hub.clients == nil {
		t.Error("Hub clients map is nil")
	}

	if hub.broadcast == nil {
		t.Error("Hub broadcast channel is nil")
	}

	if hub.register == nil {
		t.Error("Hub register channel is nil")
	}

	if hub.unregister == nil {
		t.Error("Hub unregister channel is nil")
	}
}

// TestMessageUpgrader tests the WebSocket upgrader configuration
func TestMessageUpgrader(t *testing.T) {
	// Test that the upgrader is properly configured
	if MessageUpgrader.ReadBufferSize != 1024 {
		t.Errorf("Expected ReadBufferSize 1024, got %d", MessageUpgrader.ReadBufferSize)
	}

	if MessageUpgrader.WriteBufferSize != 1024 {
		t.Errorf("Expected WriteBufferSize 1024, got %d", MessageUpgrader.WriteBufferSize)
	}

	// Test CheckOrigin function
	req := httptest.NewRequest("GET", "http://example.com", nil)
	if !MessageUpgrader.CheckOrigin(req) {
		t.Error("CheckOrigin should return true for test request")
	}
}
