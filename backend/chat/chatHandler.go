package chat

import (
	"net/http"

	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/websocket"
)

func SimpleWebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.DebugPrint(err)
		return
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			utils.DebugPrint(err)
			break
		}
		ws.WriteMessage(websocket.TextMessage, message)
	}
}

func CreateChat(w http.ResponseWriter, req *http.Request) {
	var user1_id, user2_id string
	user1_id = req.URL.Query().Get("u1")
	user2_id = req.URL.Query().Get("u2")
	utils.DebugPrint(user1_id, user2_id)
	res := database.QueryRow("SELECT * FROM chats WHERE user1_id = ? AND user2_id = ?", user1_id, user2_id)
	utils.DebugPrint(res)
	if res != nil {
		utils.SendJSONResponse(w, http.StatusOK, res)
		return
	}

	utils.DebugPrint(req.URL.RawQuery)

	database.Execute("INSERT INTO chats (user1_id, user2_id) VALUES (?, ?)", user1_id, user2_id)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "Created a new chat with users " + user1_id + " and " + user2_id})
}
