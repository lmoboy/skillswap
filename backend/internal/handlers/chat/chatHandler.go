package chat

import (
	"net/http"

	"skillswap/backend/internal/utils"
)

// WebSocketMessage represents the incoming WebSocket message structure
type WebSocketMessage struct {
	Type    string `json:"type"`
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// CreateChat creates or finds a chat between two users
func CreateChat(w http.ResponseWriter, req *http.Request) {
	user1ID := req.URL.Query().Get("u1")
	user2ID := req.URL.Query().Get("u2")

	result, err := findOrCreateChat(user1ID, user2ID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if result.IsNew {
		utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "Created a new chat",
			"chat_id": result.ChatID,
			"users":   []string{result.User1ID, result.User2ID},
		})
	} else {
		utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "Chat already exists",
			"chat_id": result.ChatID,
		})
	}
}
