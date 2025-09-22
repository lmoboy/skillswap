package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func sendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
func handleError(err error) (b bool) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d (%s)\n", filename, line, err)
		b = true
	}
	return
}

// GetUserIDFromEmail retrieves a user's ID from their email
func GetUserIDFromEmail(email string) (int64, error) {
	db, err := getDatabase()
	if err != nil {
		return 0, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found with email: %s", email)
		}
		return 0, fmt.Errorf("error querying user ID: %v", err)
	}

	return userID, nil
}

// GetEmailFromUserID retrieves a user's email from their ID
func GetEmailFromUserID(userID int64) (string, error) {
	db, err := getDatabase()
	if err != nil {
		return "", fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	var email string
	err = db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found with ID: %d", userID)
		}
		return "", fmt.Errorf("error querying user email: %v", err)
	}

	return email, nil
}

// GetChatsByUserEmail retrieves all chats for a user given their email
func GetChatsByUserEmail(email string) ([]map[string]interface{}, error) {
	db, err := getDatabase()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	// First, get the user ID from email
	userID, err := GetUserIDFromEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error getting user ID: %v", err)
	}

	// Query to get all chats where the user is either user1 or user2
	query := `
		SELECT 
			c.id as chat_id,
			c.user1_id,
			c.user2_id,
			c.initiated_by,
			c.created_at,
			u1.username as user1_username,
			u1.email as user1_email,
			u2.username as user2_username,
			u2.email as user2_email
		FROM chats c
		JOIN users u1 ON c.user1_id = u1.id
		JOIN users u2 ON c.user2_id = u2.id
		WHERE c.user1_id = ? OR c.user2_id = ?
		ORDER BY c.created_at DESC
	`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chats: %v", err)
	}
	defer rows.Close()

	var chats []map[string]interface{}
	for rows.Next() {
		var (
			chatID, user1ID, user2ID, initiatedBy int64
			createdAt                              string
			user1Username, user1Email              string
			user2Username, user2Email              string
		)

		err := rows.Scan(
			&chatID, &user1ID, &user2ID, &initiatedBy,
			&createdAt, &user1Username, &user1Email,
			&user2Username, &user2Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat row: %v", err)
		}

		// Determine the other user's info
		var otherUserID int64
		var otherUsername, otherEmail string
		if user1ID == userID {
			otherUserID = user2ID
			otherUsername = user2Username
			otherEmail = user2Email
		} else {
			otherUserID = user1ID
			otherUsername = user1Username
			otherEmail = user1Email
		}

		chat := map[string]interface{}{
			"chat_id":        chatID,
			"other_user_id":  otherUserID,
			"username":       otherUsername,
			"email":          otherEmail,
			"initiated_by":   initiatedBy == userID,
			"created_at":     createdAt,
		}

		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat rows: %v", err)
	}

	return chats, nil
}
