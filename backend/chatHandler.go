package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader is the websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Message is a websocket message
type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Hub is the central message broker
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mutex     sync.Mutex
	messages  []Message
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool, 2),
		broadcast: make(chan Message),
		messages:  make([]Message, 0),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		msg := <-h.broadcast
		h.mutex.Lock()
		h.messages = append(h.messages, msg)
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("Error writing message:", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}

// WSEndpoint is the websocket endpoint handler
func WSEndpoint(w http.ResponseWriter, req *http.Request, hub *Hub) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	webSocket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Println("Websocket connection error :", err)
		return
	}

	hub.mutex.Lock()
	hub.clients[webSocket] = true
	hub.mutex.Unlock()

	// send all previous messages to the new client
	for _, msg := range hub.messages {
		err := webSocket.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			hub.mutex.Lock()
			delete(hub.clients, webSocket)
			hub.mutex.Unlock()
			webSocket.Close()
			return
		}
	}

	for {
		var p Message
		err := webSocket.ReadJSON(&p)
		if err != nil {
			fmt.Println("Websocket reading error:", err)
			hub.mutex.Lock()
			delete(hub.clients, webSocket)
			hub.mutex.Unlock()
			webSocket.Close()
			return
		}
		hub.broadcast <- p
	}
}

// JoinWebSocket joins the websocket
func JoinWebSocket(w http.ResponseWriter, req *http.Request, hub *Hub) {
	WSEndpoint(w, req, hub)
}

// SaveMessageToDatabase saves a message to the database
func SaveMessageToDatabase(msg Message) error {
	// implement your database logic here
	return nil
}

// LoadMessagesFromDatabase loads messages from the database
func LoadMessagesFromDatabase() ([]Message, error) {
	// implement your database logic here
	return nil, nil
}

// RunWebsocket starts the websocket
func RunWebsocket(w http.ResponseWriter, req *http.Request) {

	/* 
	so a hypothetical is
	we send over a request to create a new hub and connect to the websocket related to that hub
	then we save the hub in a global array in the app, unless there are no connections we remove it
	we identify the hub using a map for easier lookup on request,

	if a user want to join our hub they send a body containing their token and the hub they want to join,
	if the tokens do not match up with the database 
	(aka the chat token doesn't ontain both users that are trying to join, make a new hub or lookup for old one)
	----------------------------------scenarios------------------------------------------------
	User Request -> Send CREATE HUB request with a unique ID

	Server -> Check global hub map for ID existence

		IF ID exists -> ❌ Response: Hub already exists

		IF ID does not exist -> ✅ Server Actions:

			Create new Hub struct

			Create new WebSocket (ws) connection

			Add the Hub to the global map

			Response -> Send SUCCESS message with the hub ID

	----------------------------------scenarios------------------------------------------------
	Scenario 2: User Joining an Existing Hub

	This flow is for a user connecting to a hub that has already been created.

	User Request -> Send JOIN HUB request with token and hub ID

	Server -> Validate user's token against the database

		IF token is invalid -> ❌ Response: Authentication failed

		IF token is valid -> ✅ Server Actions:

			Look up hub ID in the global map

			IF hub does not exist -> Server Actions:

				Check for an existing private chat hub between the two users

				IF private chat exists -> Redirect user to join that hub

				IF private chat does not exist -> Create a new hub -> Connect user to it -> Save it to the global map -> Response: SUCCESS with new hub ID

			IF hub exists -> Server Actions:

				Establish a new WebSocket connection for the user

				Associate the user with the existing hub

				Add user details to the hub's list of connected users

				Response -> Send SUCCESS message confirming the join
	----------------------------------scenarios------------------------------------------------

	Scenario 3: Hub Removal

	This flow describes how an empty hub is removed from the system.

	Event -> User disconnects from WebSocket

	Server -> Check the number of active connections in the hub

		IF connections > 0 -> Hub remains active

		IF connections = 0 -> Server Actions:

			Remove the hub from the global map

			Close the WebSocket connection

			Server -> Log event for monitoring

	*/
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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hub := NewHub()
	go hub.Run()

	// load messages from database to the hub
	messages, err := LoadMessagesFromDatabase()
	if err != nil {
		fmt.Println("Error loading messages from database:", err)
		return
	}
	hub.messages = messages

	JoinWebSocket(w, req, hub)
}
