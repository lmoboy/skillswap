package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/websocket"
)

// NewHub creates a new MessageHub with an initialized client map (capacity 2), a broadcasts channel, and an empty messages slice.
func NewHub() *MessageHub {
	return &MessageHub{
		Clients:    make(map[*websocket.Conn]bool, 2),
		Broadcasts: make(chan Message),
		Messages:   make([]Message, 0),
	}
}

func (h *MessageHub) Run() {
	for {
		msg := <-h.Broadcasts
		h.Mutex.Lock()
		h.Messages = append(h.Messages, msg)
		for client := range h.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				utils.HandleError(err)
				fmt.Println("Error writing message:", err)
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mutex.Unlock()
	}
}

// WSEndpoint is the websocket endpoint handler
func WSEndpoint(w http.ResponseWriter, req *http.Request, hub *MessageHub) {
	MessageUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	webSocket, err := MessageUpgrader.Upgrade(w, req, nil)
	if err != nil {
		utils.HandleError(err)
		fmt.Println("Websocket connection error :", err)
		return
	}

	hub.Mutex.Lock()
	hub.Clients[webSocket] = true
	hub.Mutex.Unlock()

	// send all previous Messages to the new client
	for _, msg := range hub.Messages {
		err := webSocket.WriteJSON(msg)
		if err != nil {
			utils.HandleError(err)
			fmt.Println("Error writing message:", err)
			hub.Mutex.Lock()
			delete(hub.Clients, webSocket)
			hub.Mutex.Unlock()
			webSocket.Close()
			return
		}
	}

	for {
		var p Message
		err := webSocket.ReadJSON(&p)
		if err != nil {
			utils.HandleError(err)
			fmt.Println("Websocket reading error:", err)
			hub.Mutex.Lock()
			delete(hub.Clients, webSocket)
			hub.Mutex.Unlock()
			webSocket.Close()
			return
		}
		hub.Broadcasts <- p
	}
}

// JoinWebSocket upgrades the HTTP request to a WebSocket connection and attaches it to the provided MessageHub for receiving and broadcasting messages.
func JoinWebSocket(w http.ResponseWriter, req *http.Request, hub *MessageHub) {
	WSEndpoint(w, req, hub)
}

// SaveToDBLink decodes a JSON Message from the request body and stores it in the database.
// On JSON decode failure it responds with HTTP 400 Bad Request.
// On database save failure it responds with HTTP 500 Internal Server Error.
// Successful requests produce no response body.
func SaveToDBLink(w http.ResponseWriter, req *http.Request) {
	var message Message

	err := json.NewDecoder(req.Body).Decode(&message)

	utils.DebugPrint(message)
	if err != nil {
		utils.HandleError(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = SaveMessageToDatabase(message)
	if err != nil {
		utils.HandleError(err)
		http.Error(w, "Failed to save message to database", http.StatusInternalServerError)
		return
	}
}

// SaveMessageToDatabase inserts msg into the `messages` table, persisting its chat_id, sender_id, and content.
// It returns an error if the database insert fails.
func SaveMessageToDatabase(msg Message) error {
	_, err := database.Execute(`
		INSERT INTO messages (chat_id, sender_id, content)
		VALUES (?, ?, ?)
	`, msg.ChatID, msg.Sender.ID, msg.Content)
	if err != nil {
		utils.HandleError(err)
		return err
	}
	return nil
}

// LoadMessagesFromDatabase loads Messages from the database

// RunWebsocket upgrades an HTTP request to a websocket, initializes and starts a MessageHub,
// loads persisted messages into the hub, and attaches the requesting client to that hub.
// It sends HTTP 400 if the request body is missing or contains invalid JSON; if loading
// messages from the database fails the error is logged and the handler returns without
// attaching the client.
func RunWebsocket(w http.ResponseWriter, req *http.Request) {

	if req.Body == nil {
		http.Error(w, "No request body", http.StatusBadRequest)
		return
	}
	var p struct {
		UsrToken  string `json:"usrtokn"`
		JoinsRoom string `json:"joinroom"`
	}
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		utils.HandleError(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hub := NewHub()
	go hub.Run()

	// load Messages from database to the hub
	Messages, err := LoadMessagesFromDatabase()
	if err != nil {
		utils.HandleError(err)
		fmt.Println("Error loading Messages from database:", err)
		return
	}
	hub.Messages = Messages

	JoinWebSocket(w, req, hub)
}
