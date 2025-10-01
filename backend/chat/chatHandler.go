package chat

import (
	"encoding/json"
	"net/http"
	"time"

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

// Client represents a single user's WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages to be broadcasted to all clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// --- HUB IMPLEMENTATION ---

// NewHub creates a Hub with initialized channels and a client registry for managing registrations, unregistrations, and broadcast messages.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the Hub's main loop. This should be run in a separate goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Add new client
			h.clients[client] = true
			utils.DebugPrint("New client connected. Total clients:", len(h.clients))

		case client := <-h.unregister:
			// Remove client
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				utils.DebugPrint("Client disconnected. Total clients:", len(h.clients))
			}

		case message := <-h.broadcast:
			// Send message to all active clients
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message successfully sent to client's channel
				default:
					// Failed to send (channel blocked/full), unregister client
					close(client.send)
					delete(h.clients, client)
					client.conn.Close() // Close the connection
					utils.DebugPrint("Client unregistering due to failed send.")
				}
			}
		}
	}
}

// --- CLIENT READ/WRITE IMPLEMENTATION ---

// readPump reads messages from the WebSocket connection and processes them.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	// Set read deadline to ensure the connection eventually times out
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.DebugPrint("error reading message:", err)
			}
			break // Exit the loop and run deferred cleanup
		}

		var wsMessage WebSocketMessage
		messageStr := string(message)

		// Original array/single message parsing logic (preserved)
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
			
			// Insert message into database
			if(wsMessage.Content == "") {
				continue
			}
			_, err = database.Execute("INSERT INTO messages (chat_id, sender_id, content) VALUES (?, ?, ?)", wsMessage.ID, wsMessage.UserID, wsMessage.Content)

			if err != nil {
				utils.HandleError(err)
				errorResponse := map[string]interface{}{
					"type":   "error",
					"status": "error",
					"error":  err.Error(),
				}
				errorBytes, _ := json.Marshal(errorResponse)
				c.hub.broadcast <- errorBytes
				continue
			}

			// Fetch the complete user information for the sender
			row := database.QueryRow(`
				SELECT u.id, u.username, u.email, u.profile_picture, u.aboutme, u.profession, u.location
				FROM users AS u
				WHERE u.id = ?`, wsMessage.UserID)
			
			var sender Message
			err = row.Scan(&sender.Sender.ID, &sender.Sender.Username, &sender.Sender.Email, 
				&sender.Sender.ProfilePicture, &sender.Sender.AboutMe, &sender.Sender.Professions, 
				&sender.Sender.Location)
			
			if err != nil {
				utils.HandleError(err)
				continue
			}

			sender.Content = wsMessage.Content
			// Get current timestamp
			sender.TimeStamp = time.Now().Format("2006-01-02 15:04:05")

			response := map[string]interface{}{
				"type":    "new_message",
				"chat_id": wsMessage.ID,
				"message": sender,
			}
			responseBytes, _ := json.Marshal(response)

			// Broadcast the complete message to all connected clients
			c.hub.broadcast <- responseBytes

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

			// For updates, we'll also broadcast the status to all users
			c.hub.broadcast <- responseBytes

		default:
			utils.DebugPrint("Unknown message type:", wsMessage.Type)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				utils.HandleError(err)
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				utils.DebugPrint("Sending message:", message)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				utils.HandleError(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				utils.HandleError(err)
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var globalHub = NewHub()

// StartHub initializes and runs the global WebSocket hub
func StartHub() {
	globalHub.Run()
}

// SimpleWebSocketEndpoint upgrades the HTTP request to a WebSocket, registers a new Client with the global hub, and runs the client's I/O pumps.
// 
// On successful upgrade it creates a Client with a buffered send channel, registers the client with globalHub, starts writePump in a new goroutine and runs readPump (blocking) until the connection closes. If the WebSocket upgrade fails, the error is logged and the handler returns.
func SimpleWebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.DebugPrint("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:  globalHub,
		conn: conn,
		send: make(chan []byte, 256), // Buffered channel for outbound messages
	}
	client.hub.register <- client
	utils.DebugPrint("New client connected:", len(globalHub.clients))

	go client.writePump() // Handles outgoing broadcast messages
	client.readPump()     // Handles incoming messages (blocking until connection closes)
}

// CreateChat handles HTTP requests to create a chat between two users.
//
// It reads the user IDs from query parameters "u1" and "u2". If a database query to check the chat fails, it responds with HTTP 404 and a JSON error message. On success it inserts a new chat record and responds with HTTP 200 and a JSON object indicating the created chat and the involved user IDs.
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
