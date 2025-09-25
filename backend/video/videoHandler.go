package video

import (
	"encoding/json"
	"log"
	"net/http"
	"skillswap/backend/utils"

	"github.com/pion/webrtc/v3"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := VideoUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println(err)
			continue
		}

		switch msg["type"] {
		case "join":
			log.Printf("%s joined the session", msg["name"])
		case "offer":
			// Handle offer message
		case "answer":
			// Handle answer message
		case "candidate":
			// Handle ICE candidate message
		}
	}
}

func createPeerConnection() (*webrtc.PeerConnection, error) {
	iceServers := []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	}

	config := webrtc.Configuration{
		ICEServers: iceServers,
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		utils.DebugPrint("ICE Connection State has changed: %s\n", state.String())
	})

	return peerConnection, nil
}
