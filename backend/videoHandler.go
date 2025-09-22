package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var oogabooga = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var rtcConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{URLs: []string{
			// "stun:stun.l.google.com:19302",
			"turn:localhost:3499"},
			Username: "admin", Credential: "admin"},
	},
}

type Signal struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func handleVideo(w http.ResponseWriter, r *http.Request) {
	conn, err := oogabooga.Upgrade(w, r, nil)
	if err != nil {
		handleError(err)
		return
	}
	defer conn.Close()

	peer, err := webrtc.NewPeerConnection(rtcConfig)
	if err != nil {
		handleError(err)
		return
	}
	defer peer.Close()

	// We need to keep track of senders to send RTCP etc
	// Map remote track ID to local track sender
	senders := make(map[string]*webrtc.RTPSender)

	// Handle received tracks (echo back)
	// fmt.Println(peer)
	peer.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("OnTrack kind=%s id=%s codec=%s", remoteTrack.Kind(), remoteTrack.ID(), remoteTrack.Codec().MimeType)

		// Create local track matching remote
		localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), "echo")
		if err != nil {
			handleError(err)
			return
		}

		sender, err := peer.AddTrack(localTrack)
		if err != nil {
			handleError(err)
			return
		}
		senders[remoteTrack.ID()] = sender

		// Read from remote, write into local
		go func() {
			buf := make([]byte, 1500)
			for {
				n, _, err := remoteTrack.Read(buf)
				if err != nil {
					handleError(err)
					return
				}
				if _, err := localTrack.Write(buf[:n]); err != nil {
					handleError(err)
					return
				}
			}
		}()

		// Also handle RTCP, to respond to receiver reports or send PLI etc
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				if sender == nil {
					return
				}
				_, _, err := sender.Read(rtcpBuf)
				if err != nil {
					handleError(err)
					return
				}
				// Optionally parse rtcpBuf and take action if needed
				// e.g. on PLI request, but for echo probably not strictly needed
			}
		}()
	})

	// ICE candidates from server → client
	peer.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		jsonMsg := map[string]interface{}{
			"type": "ice-candidate",
			"data": c.ToJSON(),
		}
		msgBytes, _ := json.Marshal(jsonMsg)
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			handleError(err)
		}
	})

	// Read signaling messages from client
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			handleError(err)
			return
		}

		var sig Signal
		if err := json.Unmarshal(msgBytes, &sig); err != nil {
			handleError(err)
			continue
		}

		switch sig.Type {
		case "offer":
			var offer webrtc.SessionDescription
			if err := json.Unmarshal(sig.Data, &offer); err != nil {
				handleError(err)
				continue
			}

			if err := peer.SetRemoteDescription(offer); err != nil {
				handleError(err)
				continue
			}

			// At this point, after SetRemoteDescription, before CreateAnswer is OK
			// Create answer
			answer, err := peer.CreateAnswer(nil)
			if err != nil {
				handleError(err)
				continue
			}

			// Wait for ICE gathering to complete or send partial if using trickle ICE
			// If not using trickle, wait:
			gatherComplete := webrtc.GatheringCompletePromise(peer)
			if err := peer.SetLocalDescription(answer); err != nil {
				handleError(err)
				continue
			}

			<-gatherComplete // wait

			resp := map[string]interface{}{
				"type": "answer",
				"data": peer.LocalDescription(),
			}
			respMsg, _ := json.Marshal(resp)
			if err := conn.WriteMessage(websocket.TextMessage, respMsg); err != nil {
				handleError(err)
			}

		case "ice-candidate":
			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(sig.Data, &candidate); err != nil {
				handleError(err)
				continue
			}
			if err := peer.AddICECandidate(candidate); err != nil {
				handleError(err)
			}

		default:
			log.Println("Unknown signal type:", sig.Type)
		}
	}
}
