package chat

import (
	"encoding/json"
	"time"

	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
)

// handleMessage routes message to appropriate handler based on type
func (c *Client) handleMessage(wsMessage *WebSocketMessage) {
	switch wsMessage.Type {
	case "post":
		c.handlePostMessage(wsMessage)
	case "update":
		c.handleUpdateMessage(wsMessage)
	default:
		// utils.DebugPrint("Unknown message type:", wsMessage.Type)
	}
}

// handlePostMessage handles posting a new message to the database and broadcasts it
func (c *Client) handlePostMessage(wsMessage *WebSocketMessage) {
	// utils.DebugPrint("Handling POST message with ID:", wsMessage.ID)

	if wsMessage.Content == "" {
		return
	}

	// Save message to database
	messageID, err := saveMessageToDB(wsMessage)
	if err != nil {
		c.broadcastError(err)
		return
	}

	// Fetch complete message with sender info
	message, err := fetchMessageWithSender(messageID, wsMessage)
	if err != nil {
		utils.HandleError(err)
		return
	}

	// Broadcast to all clients
	c.broadcastNewMessage(wsMessage.ID, message)
}

// handleUpdateMessage handles update message type
func (c *Client) handleUpdateMessage(wsMessage *WebSocketMessage) {
	// utils.DebugPrint("Handling UPDATE message with ID:", wsMessage.ID)

	response := map[string]interface{}{
		"type":    wsMessage.Type,
		"id":      wsMessage.ID,
		"user_id": wsMessage.UserID,
		"content": wsMessage.Content,
		"status":  "processed",
	}
	responseBytes, _ := json.Marshal(response)
	// utils.DebugPrint(response)

	c.hub.broadcast <- responseBytes
}

// saveMessageToDB inserts a message into the database and returns the message ID
func saveMessageToDB(wsMessage *WebSocketMessage) (int64, error) {
	result, err := database.Execute(
		"INSERT INTO messages (chat_id, sender_id, content) VALUES (?, ?, ?)",
		wsMessage.ID, wsMessage.UserID, wsMessage.Content,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// fetchMessageWithSender retrieves complete user information for a message sender
func fetchMessageWithSender(messageID int64, wsMessage *WebSocketMessage) (*Message, error) {
	row := database.QueryRow(`
		SELECT u.id, u.username, u.email,
			COALESCE(u.profile_picture, ''),
			COALESCE(u.aboutme, ''),
			COALESCE(u.profession, ''),
			COALESCE(u.location, '')
		FROM users AS u
		WHERE u.id = ?`, wsMessage.UserID)

	var message Message
	err := row.Scan(
		&message.Sender.ID,
		&message.Sender.Username,
		&message.Sender.Email,
		&message.Sender.ProfilePicture,
		&message.Sender.AboutMe,
		&message.Sender.Professions,
		&message.Sender.Location,
	)

	if err != nil {
		return nil, err
	}

	message.Id = int(messageID)
	message.Content = wsMessage.Content
	message.TimeStamp = time.Now().Format("2006-01-02 15:04:05")

	return &message, nil
}

// broadcastNewMessage sends a new message notification to all connected clients
func (c *Client) broadcastNewMessage(chatID int, message *Message) {
	response := map[string]any{
		"type":    "new_message",
		"chat_id": chatID,
		"message": message,
	}
	responseBytes, _ := json.Marshal(response)
	c.hub.broadcast <- responseBytes
}

// broadcastError sends an error message to all connected clients
func (c *Client) broadcastError(err error) {
	utils.HandleError(err)
	errorResponse := map[string]interface{}{
		"type":   "error",
		"status": "error",
		"error":  err.Error(),
	}
	errorBytes, _ := json.Marshal(errorResponse)
	c.hub.broadcast <- errorBytes
}


