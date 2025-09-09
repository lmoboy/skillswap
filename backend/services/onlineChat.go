package onlineChat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	// "github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Name    string `json: name`
	Message string `json: message`
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mutex     sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

func (h *Hub) run() {
	for {
		// once again mutex lock because of concurrent connections and broadcast the sent message to every client
		msg := <-h.broadcast
		h.mutex.Lock()
		for client := range h.clients {
			err := client.WriteJSON(msg)
			// yeet the client if the connection errored, we do not need them GET OUT
			if err != nil {
				fmt.Println("Error writing message:", err)
				client.Close()
				delete(h.clients, client)
			}

		}
		h.mutex.Unlock()
	}
}

func wsEndpoint(w http.ResponseWriter, req *http.Request, hub *Hub) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	webSocket, Err := upgrader.Upgrade(w, req, nil)
	if Err != nil {
		fmt.Println("Websocket connection error :", Err)
	}
	// lock the mutex because we have concurrent connections that want to access the client map
	hub.mutex.Lock()
	hub.clients[webSocket] = true
	hub.mutex.Unlock()

	for {
		// for each connection read the input, lock the mutex for client removal and push it to the broadcast var
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

// func ping(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")
// 	w.Header().Set("Access-Control-Allow-Credentials", "true")

// 	if req.Method == http.MethodOptions {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// 	fmt.Println("WE GOT PINGED!!!")
// 	fmt.Fprintf(w, "pong\n")
// }

func runwebsocket() {
	// Create new hub to fill up and run it
	hub := newHub()
	go hub.run()
	http.HandleFunc("/chat", func(w http.ResponseWriter, req *http.Request) {
		wsEndpoint(w, req, hub)
	})
}
