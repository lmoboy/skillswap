package video

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var VideoUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type   string      `json:"type"`
	RoomID string      `json:"roomId"`
	From   string      `json:"from"`
	To     string      `json:"to"`
	Data   interface{} `json:"data"`
}
