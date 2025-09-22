package structs

import (
	"sync"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ID             int    `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	AboutMe        string `json:"aboutme"`
	Projects       string `json:"projects"`
	Contacts       string `json:"contacts"`
	Skills         string `json:"skills"`
	Location       string `json:"location"`
}

type Message struct {
	Sender    UserInfo `json:"sender"`
	TimeStamp string   `json:"timestamp"`
	Content   string   `json:"content"`
}

type MessageHub struct {
	Clients     map[*websocket.Conn]bool
	Broadcasts chan Message
	Mutex      sync.Mutex
	Messages   []Message
}
