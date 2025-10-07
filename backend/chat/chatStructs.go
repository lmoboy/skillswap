package chat

import (
	"net/http"
	"skillswap/backend/structs"
	"sync"

	"github.com/gorilla/websocket"
)

type Chat struct {
	Id         int    `json:"id"`
	Initiator  int    `json:"user1_id"`
	Responder  int    `json:"user2_id"`
	Created_at string `json:"created_at"`
}

type ChatWithUserInfo struct {
	Id                      int    `json:"id"`
	Initiator               int    `json:"user1_id"`
	Responder               int    `json:"user2_id"`
	Created_at              string `json:"created_at"`
	InitiatorUsername       string `json:"user1_username"`
	InitiatorProfilePicture string `json:"user1_profile_picture"`
	ResponderUsername       string `json:"user2_username"`
	ResponderProfilePicture string `json:"user2_profile_picture"`
}

type Message struct {
	Sender    structs.UserInfo `json:"sender"`
	Content   string           `json:"content"`
	TimeStamp string           `json:"timestamp"`
}

type MessageHub struct {
	Clients    map[*websocket.Conn]bool
	Broadcasts chan Message
	Mutex      sync.Mutex
	Messages   []Message
}

var MessageUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}