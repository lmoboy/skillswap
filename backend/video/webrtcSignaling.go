package video

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Room represents a signaling room for peer-to-peer WebRTC connections
var Rooms = make(map[string]*Room)
var RoomsMutex = sync.RWMutex{}

// Room represents a signaling room for peer-to-peer WebRTC connections
type Room struct {
	ID         string
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan Message
}

// Message represents a WebSocket message for signaling


// HandleWebSocket handles WebSocket connections for WebRTC signaling
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := VideoUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Parse room ID from query parameters
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		log.Println("Room ID is required")
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
		go room.run()
	}
	RoomsMutex.Unlock()

	// Register client in room
	room.Register <- conn

	// Handle incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			room.Unregister <- conn
			break
		}

		// Add sender info to message
		msg.RoomID = roomID
		msg.From = conn.RemoteAddr().String()

		// Broadcast message to other clients in the room
		select {
		case room.Broadcast <- msg:
		default:
			log.Println("Broadcast channel is full, dropping message")
			room.Unregister <- conn
			break
		}
	}
}

// run handles the room's message broadcasting
func (room *Room) run() {
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
			for client := range room.Clients {
				if client != nil && client != getConnByAddr(message.From) {
					if err := client.WriteJSON(message); err != nil {
						log.Println("Write error:", err)
						delete(room.Clients, client)
						client.Close()
					}
				}
			}
		}
	}
}

// getConnByAddr finds a connection by its remote address string
func getConnByAddr(addr string) *websocket.Conn {
	RoomsMutex.RLock()
	defer RoomsMutex.RUnlock()

	for _, room := range Rooms {
		for conn := range room.Clients {
			if conn.RemoteAddr().String() == addr {
				return conn
			}
		}
	}
	return nil
}
