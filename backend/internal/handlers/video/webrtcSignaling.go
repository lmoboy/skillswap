// Package video provides WebRTC signaling and video handling functionality.
//
// This package contains:
// - WebSocket signaling server for peer-to-peer WebRTC connections
// - Video upload and streaming capabilities
// - Room-based message routing for multiple users
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

// HandleWebSocket handles WebSocket connections for WebRTC signaling.
// It upgrades HTTP connections to WebSocket, manages room-based peer connections,
// HandleWebSocket upgrades an HTTP request to a WebSocket, places the connection into a room (created if needed), and relays signaling messages between clients in the same room.
// The room is determined by the "room" query parameter; each incoming JSON message is annotated with the RoomID and the sender's remote address before being broadcast to other room members. If the room parameter is missing the request is rejected; if the room's broadcast channel is full, the message is dropped and the sender is unregistered.
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

// run handles the room's message broadcasting and client management.
// It listens for client registration, unregistration, and message broadcasting,
// managing the lifecycle of WebSocket connections within a room.
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

// getConnByAddr finds a connection by its remote address string.
// getConnByAddr searches all rooms and returns the first websocket connection whose RemoteAddr string equals the provided addr.
// If no matching connection is found, it returns nil.
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