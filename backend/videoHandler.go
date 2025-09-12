package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected")
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received message type: %d, length: %d bytes", messageType, len(p))

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
