package chat

import (
	"encoding/json"
	"net/http"

	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/websocket"
)

// WebSocketMessage represents the incoming WebSocket message structure
type WebSocketMessage struct {
	Type    string `json:"type"`
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

func SimpleWebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.DebugPrint(err)
		return
	}
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			utils.DebugPrint(err)
			break
		}

		// Parse the incoming message
		var wsMessage WebSocketMessage
		messageStr := string(message)

		// Handle both array format [{"type":"post",...}] and single object format {"type":"post",...}
		if len(messageStr) > 0 && messageStr[0] == '[' && messageStr[len(messageStr)-1] == ']' {
			// Array format - take the first element
			var messages []WebSocketMessage
			if err := json.Unmarshal(message, &messages); err != nil {
				utils.DebugPrint("Error parsing array message:", err)
				continue
			}
			if len(messages) > 0 {
				wsMessage = messages[0]
			}
		} else {
			// Single object format
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				utils.DebugPrint("Error parsing message:", err)
				continue
			}
		}

		utils.DebugPrint("Parsed message:", wsMessage)

		// Handle different message types
		switch wsMessage.Type {
		case "post":
			utils.DebugPrint("Handling POST message with ID:", wsMessage.ID)
			response := map[string]interface{}{
				"type":    wsMessage.Type,
				"id":      wsMessage.ID,
				"user_id": wsMessage.UserID,
				"content": wsMessage.Content,
				"status":  "processed",
			}
			_, err = database.Execute("INSERT INTO messages (chat_id, sender_id, content) VALUES (?, ?, ?)", wsMessage.ID, wsMessage.UserID, wsMessage.Content)
			utils.DebugPrint("INSERT INTO messages (chat_id, sender_id, content) VALUES (%v, %v, %v)", wsMessage.ID, wsMessage.UserID, wsMessage.Content)

			if err != nil {
				utils.DebugPrint(err)
				response["status"] = "error"
				response["error"] = err.Error()
			}
			responseBytes, _ := json.Marshal(response)
			ws.WriteMessage(websocket.TextMessage, responseBytes)

		case "update":
			utils.DebugPrint("Handling UPDATE message with ID:", wsMessage.ID)
			response := map[string]interface{}{
				"type":    wsMessage.Type,
				"id":      wsMessage.ID,
				"user_id": wsMessage.UserID,
				"content": wsMessage.Content,
				"status":  "processed",
			}
			responseBytes, _ := json.Marshal(response)
			ws.WriteMessage(websocket.TextMessage, responseBytes)
		default:
			utils.DebugPrint("Unknown message type:", wsMessage.Type)
		}

		// Echo back the processed message
	}
	ws.WriteMessage(websocket.TextMessage, []byte("Connection closed"))
}

func CreateChat(w http.ResponseWriter, req *http.Request) {
	var user1_id, user2_id string
	user1_id = req.URL.Query().Get("u1")
	user2_id = req.URL.Query().Get("u2")
	utils.DebugPrint(user1_id, user2_id)
	res := database.QueryRow("SELECT * FROM chats WHERE user1_id = ? AND user2_id = ?", user1_id, user2_id)
	utils.DebugPrint(res.Err())
	if res.Err() != nil {
		utils.DebugPrint(res.Err())
		utils.SendJSONResponse(w, http.StatusNotFound, "somethng happned")
		return
	}

	utils.DebugPrint(req.URL.RawQuery)

	database.Execute("INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)", user1_id, user2_id)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "Created a new chat with users " + user1_id + " and " + user2_id})
}
