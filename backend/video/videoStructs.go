// Package video provides WebRTC signaling and video handling functionality.
//
// This package contains:
// - WebSocket signaling server for peer-to-peer WebRTC connections
// - Video upload and streaming capabilities
// - Room-based message routing for multiple users
package video

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// VideoUpgrader configures WebSocket upgrade settings for video signaling.
var VideoUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Message represents a WebSocket message for WebRTC signaling between peers.
type Message struct {
	Type   string      `json:"type"`   // Message type (offer, answer, candidate, etc.)
	RoomID string      `json:"roomId"` // Room identifier for multi-user signaling
	From   string      `json:"from"`   // Sender identifier
	To     string      `json:"to"`     // Recipient identifier
	Data   interface{} `json:"data"`   // Message payload (SDP, ICE candidates, etc.)
}
