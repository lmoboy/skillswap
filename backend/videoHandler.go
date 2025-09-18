package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	videoUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true 
		},
	}

	
	peers     = make(map[string]*webrtc.PeerConnection)
	peersLock sync.RWMutex
)


type SignalMessage struct {
	Type string          `json:"type"` 
	Data json.RawMessage `json:"data"`
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := videoUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		return
	}
	defer peerConnection.Close()

	
	peerID := fmt.Sprintf("peer-%d", time.Now().UnixNano())
	peersLock.Lock()
	peers[peerID] = peerConnection
	peersLock.Unlock()

	
	defer func() {
		peersLock.Lock()
		delete(peers, peerID)
		peersLock.Unlock()
	}()

	
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidateJSON, err := json.Marshal(c.ToJSON())
		if err != nil {
			log.Printf("Failed to marshal ICE candidate: %v", err)
			return
		}

		msg := SignalMessage{
			Type: "ice-candidate",
			Data: candidateJSON,
		}

		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Failed to send ICE candidate: %v", err)
		}
	})

	
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer %s connection state: %s", peerID, s.String())

		switch s {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Printf("Failed to close failed peer connection: %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			if err := conn.Close(); err != nil {
				log.Printf("Failed to close WebSocket: %v", err)
			}
		}
	})

	
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track received: %s, codec: %s", track.Kind().String(), track.Codec().MimeType)

		
		localTrack, err := webrtc.NewTrackLocalStaticRTP(
			track.Codec().RTPCodecCapability,
			track.ID(),
			track.StreamID(),
		)
		if err != nil {
			log.Printf("Failed to create local track: %v", err)
			return
		}

		
		buf := make([]byte, 1500) 
		for {
			n, _, err := track.Read(buf)
			if err != nil {
				log.Printf("Failed to read from track: %v", err)
				return
			}

			if _, err = localTrack.Write(buf[:n]); err != nil {
				log.Printf("Failed to write to local track: %v", err)
				return
			}
		}
	})

	
	for {
		var msg SignalMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		switch msg.Type {
		case "offer":
			
			var offer webrtc.SessionDescription
			if err := json.Unmarshal(msg.Data, &offer); err != nil {
				log.Printf("Failed to unmarshal offer: %v", err)
				continue
			}

			
			if err := peerConnection.SetRemoteDescription(offer); err != nil {
				log.Printf("Failed to set remote description: %v", err)
				continue
			}

			
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				log.Printf("Failed to create answer: %v", err)
				continue
			}

			
			if err := peerConnection.SetLocalDescription(answer); err != nil {
				log.Printf("Failed to set local description: %v", err)
				continue
			}

			
			answerJSON, err := json.Marshal(answer)
			if err != nil {
				log.Printf("Failed to marshal answer: %v", err)
				continue
			}

			if err := conn.WriteJSON(SignalMessage{
				Type: "answer",
				Data: answerJSON,
			}); err != nil {
				log.Printf("Failed to send answer: %v", err)
			}

		case "ice-candidate":
			
			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(msg.Data, &candidate); err != nil {
				log.Printf("Failed to unmarshal ICE candidate: %v", err)
				continue
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Printf("Failed to add ICE candidate: %v", err)
			}
		}
	}
}
