package video

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"github.com/pion/webrtc/v3"
// )

// func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := VideoUpgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	pc, err := createPeerConnection(conn)
// 	if err != nil {
// 		log.Println("PeerConnection error:", err)
// 		return
// 	}
// 	defer pc.Close()

// 	for {
// 		_, msgBytes, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("ReadMessage error:", err)
// 			break
// 		}

// 		var msg Message
// 		if err := json.Unmarshal(msgBytes, &msg); err != nil {
// 			log.Println("Invalid message:", err)
// 			continue
// 		}

// 		switch msg.Type {
// 		case "offer":
// 			var offer webrtc.SessionDescription
// 			if err := json.Unmarshal(msg.Data, &offer); err != nil {
// 				log.Println("Invalid offer:", err)
// 				continue
// 			}

// 			if err := pc.SetRemoteDescription(offer); err != nil {
// 				log.Println("SetRemoteDescription error:", err)
// 				continue
// 			}

// 			answer, err := pc.CreateAnswer(nil)
// 			if err != nil {
// 				log.Println("CreateAnswer error:", err)
// 				continue
// 			}

// 			if err := pc.SetLocalDescription(answer); err != nil {
// 				log.Println("SetLocalDescription error:", err)
// 				continue
// 			}

// 			answerMsg, _ := json.Marshal(map[string]interface{}{
// 				"type": "answer",
// 				"data": answer,
// 			})
// 			conn.WriteMessage(websocket.TextMessage, answerMsg)

// 		case "candidate":
// 			var candidate webrtc.ICECandidateInit
// 			if err := json.Unmarshal(msg.Data, &candidate); err != nil {
// 				log.Println("Invalid candidate:", err)
// 				continue
// 			}

// 			if err := pc.AddICECandidate(candidate); err != nil {
// 				log.Println("AddICECandidate error:", err)
// 			}

// 		case "join":
// 			log.Println("User joined")

// 		default:
// 			log.Println("Unknown message type:", msg.Type)
// 		}
// 	}
// }

// func createPeerConnection(conn *websocket.Conn) (*webrtc.PeerConnection, error) {
// 	config := webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{URLs: []string{"stun:stun.l.google.com:19302"}},
// 		},
// 	}

// 	pc, err := webrtc.NewPeerConnection(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Example media/data channel setup
// 	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
// 		log.Println("New DataChannel:", dc.Label())
// 		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
// 			log.Printf("Message from DataChannel %s: %s\n", dc.Label(), string(msg.Data))
// 		})
// 	})

// 	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
// 		if c == nil {
// 			return
// 		}
// 		candidate := map[string]interface{}{
// 			"type": "candidate",
// 			"data": c.ToJSON(),
// 		}
// 		data, _ := json.Marshal(candidate)
// 		conn.WriteMessage(websocket.TextMessage, data)
// 	})

// 	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
// 		log.Println("ICE state:", state.String())
// 	})

// 	return pc, nil
// }
