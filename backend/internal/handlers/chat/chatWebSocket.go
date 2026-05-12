package chat

import (
	"encoding/json"
	"net/http"
	"time"

	"skillswap/backend/internal/database"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/utils"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	globalHub = NewHub()
)

// Client represents a single user's WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID int // Authenticated user ID associated with this connection
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients         map[*Client]bool
	clientsByUserID map[int][]*Client // Track clients by user ID for targeted messages
	broadcast       chan []byte
	register        chan *Client
	unregister      chan *Client
}

// NewHub initializes and returns a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
		clientsByUserID: make(map[int][]*Client),
	}
}

// Run starts the Hub's main loop (should be run in a goroutine)
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a new client to the hub
func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
	// Add client to user ID tracking
	h.clientsByUserID[client.userID] = append(h.clientsByUserID[client.userID], client)
	// utils.DebugPrint("New client connected. Total clients:", len(h.clients))
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		// Remove client from user ID tracking
		clients := h.clientsByUserID[client.userID]
		for i, c := range clients {
			if c == client {
				h.clientsByUserID[client.userID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(h.clientsByUserID[client.userID]) == 0 {
			delete(h.clientsByUserID, client.userID)
		}
		// utils.DebugPrint("Client disconnected. Total clients:", len(h.clients))
	}
}

// broadcastMessage sends a message to all active clients
func (h *Hub) broadcastMessage(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
			// Message successfully sent
		default:
			// Failed to send, unregister client
			h.closeClient(client)
		}
	}
}

// SendToUser sends a message to all active connections of a specific user
func (h *Hub) SendToUser(userID int, message []byte) {
	clients, ok := h.clientsByUserID[userID]
	if !ok {
		return
	}
	for _, client := range clients {
		select {
		case client.send <- message:
			// Message sent
		default:
			// Failed, unregister
			h.closeClient(client)
		}
	}
}

// closeClient closes and unregisters a client
func (h *Hub) closeClient(client *Client) {
	close(client.send)
	delete(h.clients, client)
	client.conn.Close()
	// utils.DebugPrint("Client unregistering due to failed send.")
}

// readPump reads messages from the WebSocket connection and processes them
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.configureConnection()

	for {
		message, err := c.readMessage()
		if err != nil {
			break
		}

		wsMessage, err := parseWebSocketMessage(message)
		if err != nil {
			// utils.DebugPrint("Error parsing message:", err)
			continue
		}

		// utils.DebugPrint("Parsed message:", wsMessage)
		c.handleMessage(wsMessage)
	}
}

// configureConnection sets up read limits and ping/pong handlers
func (c *Client) configureConnection() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// readMessage reads a message from the WebSocket connection
func (c *Client) readMessage() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// utils.DebugPrint("error reading message:", err)
		}
		return nil, err
	}
	return message, nil
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !c.writeMessage(message, ok) {
				return
			}
		case <-ticker.C:
			if !c.writePing() {
				return
			}
		}
	}
}

// writeMessage writes a message to the WebSocket connection
func (c *Client) writeMessage(message []byte, ok bool) bool {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !ok {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return false
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		utils.HandleError(err)
		return false
	}

	return true
}

// writePing writes a ping message to the WebSocket connection
func (c *Client) writePing() bool {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		utils.HandleError(err)
		return false
	}
	return true
}

// parseWebSocketMessage parses incoming WebSocket message (handles array or single message)
func parseWebSocketMessage(message []byte) (*WebSocketMessage, error) {
	messageStr := string(message)

	// Handle array format [message] - take first element
	if len(messageStr) > 0 && messageStr[0] == '[' && messageStr[len(messageStr)-1] == ']' {
		var messages []WebSocketMessage
		if err := json.Unmarshal(message, &messages); err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			return &messages[0], nil
		}
		return nil, json.Unmarshal([]byte("{}"), &WebSocketMessage{})
	}

	// Handle single message
	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return nil, err
	}
	return &wsMessage, nil
}

// StartHub initializes and runs the global WebSocket hub
func StartHub() {
	globalHub.Run()
}

// SimpleWebSocketEndpoint handles the WebSocket endpoint for the chat application
func SimpleWebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// utils.DebugPrint("Upgrade error:", err)
		return
	}

	// Get authenticated user's email from session
	session, err := auth.Store.Get(r, "authentication")
	if err != nil {
		conn.Close()
		return
	}
	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		conn.Close()
		return
	}

	// Get user ID from database
	var userID int
	row := database.QueryRow("SELECT id FROM users WHERE email = ?", email)
	if err := row.Scan(&userID); err != nil {
		conn.Close()
		return
	}

	client := &Client{
		hub:    globalHub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
	client.hub.register <- client
	// utils.DebugPrint("New client connected. User ID:", userID, "Total clients:", len(globalHub.clients))

	go client.writePump()
	client.readPump()
}


