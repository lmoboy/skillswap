package websocket

import (
	"fmt"
	"skillswap/backend/structs"
	"skillswap/backend/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender    structs.UserInfo `json:"sender"`
	TimeStamp string           `json:"timestamp"`
	Content   string           `json:"content"`
}

type MessageHub struct {
	Clients    map[*websocket.Conn]bool
	Broadcasts chan Message
	Mutex      sync.Mutex
	Messages   []Message
}

// Run starts the hub loop to broadcast messages to all connected clients.
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
