package chat

import (
	"database/sql"
	"fmt"

	"skillswap/backend/internal/database"

)

// ChatCreationResult represents the result of creating a chat
type ChatCreationResult struct {
	ChatID   int64
	IsNew    bool
	User1ID  string
	User2ID  string
}

// findExistingChat checks if a chat already exists between two users
func findExistingChat(user1ID, user2ID string) (int64, error) {
	var chatID int64
	err := database.QueryRow(
		"SELECT id FROM chats WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).Scan(&chatID)

	if err == sql.ErrNoRows {
		return 0, nil // No existing chat found
	}
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// createNewChat creates a new chat between two users
func createNewChat(user1ID, user2ID string) (int64, error) {
	result, err := database.Execute(
		"INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)",
		user1ID, user2ID,
	)
	if err != nil {
		return 0, err
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// findOrCreateChat finds an existing chat or creates a new one between two users
func findOrCreateChat(user1ID, user2ID string) (*ChatCreationResult, error) {
	if user1ID == "" || user2ID == "" {
		return nil, fmt.Errorf("both user IDs are required")
	}

	// utils.DebugPrint("Finding or creating chat between:", user1ID, user2ID)

	// Check for existing chat
	existingChatID, err := findExistingChat(user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing chat: %w", err)
	}

	if existingChatID != 0 {
		// Chat already exists
		return &ChatCreationResult{
			ChatID:  existingChatID,
			IsNew:   false,
			User1ID: user1ID,
			User2ID: user2ID,
		}, nil
	}

	// Create new chat
	newChatID, err := createNewChat(user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &ChatCreationResult{
		ChatID:  newChatID,
		IsNew:   true,
		User1ID: user1ID,
		User2ID: user2ID,
	}, nil
}


