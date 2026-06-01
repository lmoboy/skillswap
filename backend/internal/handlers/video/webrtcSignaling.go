package video

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// VideoUpgrader configures WebSocket upgrade settings for video signaling.
var VideoUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Only allow same-origin WebSocket connections
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // Allow non-browser clients
		}
		host := r.Host
		// Check if origin matches host (same origin policy)
		return strings.HasPrefix(origin, "http://"+host) ||
			strings.HasPrefix(origin, "https://"+host) ||
			strings.HasPrefix(host, "localhost")
	},
}

// Message represents a WebSocket message for WebRTC signaling between peers.
type Message struct {
	Type   string      `json:"type"`   // Message type (offer, answer, candidate, etc.)
	RoomID string      `json:"roomId"` // Room identifier for multi-user signaling
	From   string      `json:"from"`   // Sender identifier
	To     string      `json:"to"`     // Recipient identifier
	Data   interface{} `json:"data"`   // Message payload (SDP, ICE candidates, etc.)
}

// Room represents a signaling room for peer-to-peer WebRTC connections
type Room struct {
	ID         string
	Clients    map[string]*Client // map of client ID to Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

// Client represents a connected peer in a signaling room.
type Client struct {
	ID   string
	Conn *websocket.Conn
}

var Rooms = make(map[string]*Room)
var RoomsMutex = sync.RWMutex{}

// HandleWebSocket handles WebSocket connections for WebRTC signaling.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := VideoUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Parse room ID from query parameters
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		log.Println("Room ID is required")
		conn.Close()
		return
	}

	// Create a unique client ID for this connection
	clientID := conn.RemoteAddr().String() + "_" + r.Header.Get("Sec-WebSocket-Key")

	// Get or create room
	RoomsMutex.Lock()
	room, exists := Rooms[roomID]
	if !exists {
		room = &Room{
			ID:         roomID,
			Clients:    make(map[string]*Client),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
			Broadcast:  make(chan Message),
		}
		Rooms[roomID] = room
		go room.Run()
	}
	RoomsMutex.Unlock()

	client := &Client{ID: clientID, Conn: conn}

	// Register client in room
	room.Register <- client

	// Handle incoming messages
	go func() {
		defer func() {
			room.Unregister <- client
		}()

		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			// Add sender info to message
			msg.RoomID = roomID
			msg.From = clientID

			// Broadcast message to other clients in the room
			room.Broadcast <- msg
		}
	}()
}

// Run handles the room's message broadcasting and client management.
func (room *Room) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-room.Register:
			room.Clients[client.ID] = client
			log.Printf("Client %s connected to room %s. Total clients: %d", client.ID, room.ID, len(room.Clients))

		case client := <-room.Unregister:
			if _, ok := room.Clients[client.ID]; ok {
				delete(room.Clients, client.ID)
				client.Conn.Close()
				log.Printf("Client %s disconnected from room %s. Total clients: %d", client.ID, room.ID, len(room.Clients))
			}

			// Clean up empty rooms
			if len(room.Clients) == 0 {
				RoomsMutex.Lock()
				delete(Rooms, room.ID)
				RoomsMutex.Unlock()
				return
			}

		case message := <-room.Broadcast:
			senderID := message.From
			log.Printf("Signaling: Broadcasting %s from %s to others in room %s", message.Type, senderID, room.ID)
			for clientID, client := range room.Clients {
				if clientID != senderID {
					if err := client.Conn.WriteJSON(message); err != nil {
						log.Printf("Signaling: Write error to %s: %v", clientID, err)
						client.Conn.Close()
						delete(room.Clients, clientID)
					}
				}
			}

		case <-ticker.C:
			// Heartbeat: ping all clients to detect dead connections
			for clientID, client := range room.Clients {
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Signaling: Ping failed for %s: %v", clientID, err)
					client.Conn.Close()
					delete(room.Clients, clientID)
				}
			}
			// Clean up if room went empty during ping
			if len(room.Clients) == 0 {
				RoomsMutex.Lock()
				delete(Rooms, room.ID)
				RoomsMutex.Unlock()
				return
			}
		}
	}
}
