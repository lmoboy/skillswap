package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"skillswap/backend/utils"

	"github.com/gorilla/websocket"
)

// Upgrader is the websocket upgrader

// NewHub creates a new hub
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

// JoinWebSocket joins the websocket
func JoinWebSocket(w http.ResponseWriter, req *http.Request, hub *MessageHub) {
	WSEndpoint(w, req, hub)
}

// SaveMessageToDatabase saves a message to the database
func SaveMessageToDatabase(msg Message) error {
	// implement your database logic here
	return nil
}

// LoadMessagesFromDatabase loads Messages from the database
func LoadMessagesFromDatabase() ([]Message, error) {
	// implement your database logic here
	return nil, nil
}

// RunWebsocket starts the websocket
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
