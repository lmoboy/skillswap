package video

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// VideoUpgrader configures WebSocket upgrade settings for video signaling.
var VideoUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
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
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan Message
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
	// We don't defer conn.Close() here because room management handles it

	// Parse room ID from query parameters
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		log.Println("Room ID is required")
		conn.Close()
		return
	}

	// Get or create room
	RoomsMutex.Lock()
	room, exists := Rooms[roomID]
	if !exists {
		room = &Room{
			ID:         roomID,
			Clients:    make(map[*websocket.Conn]bool),
			Register:   make(chan *websocket.Conn),
			Unregister: make(chan *websocket.Conn),
			Broadcast:  make(chan Message),
		}
		Rooms[roomID] = room
		go room.Run()
	}
	RoomsMutex.Unlock()

	// Register client in room
	room.Register <- conn

	// Handle incoming messages
	go func() {
		defer func() {
			room.Unregister <- conn
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
			msg.From = conn.RemoteAddr().String()

			// Broadcast message to other clients in the room
			room.Broadcast <- msg
		}
	}()
}

// Run handles the room's message broadcasting and client management.
func (room *Room) Run() {
	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
			log.Printf("Client connected to room %s. Total clients: %d", room.ID, len(room.Clients))

		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				client.Close()
				log.Printf("Client disconnected from room %s. Total clients: %d", room.ID, len(room.Clients))
			}

			// Clean up empty rooms
			if len(room.Clients) == 0 {
				RoomsMutex.Lock()
				delete(Rooms, room.ID)
				RoomsMutex.Unlock()
				return
			}

		case message := <-room.Broadcast:
			// Broadcast message to all clients except sender
			senderAddr := message.From
			for client := range room.Clients {
				if client.RemoteAddr().String() != senderAddr {
					if err := client.WriteJSON(message); err != nil {
						log.Println("Write error:", err)
						client.Close()
						delete(room.Clients, client)
					}
				}
			}
		}
	}
}
