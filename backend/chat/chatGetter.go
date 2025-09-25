package chat

import (
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/utils"
)

func GetMessagesFromUID(w http.ResponseWriter, req *http.Request) {
	chatId := req.URL.Query().Get("cid")
	res, err := database.Query(`
	SELECT u.id, u.username, u.email, u.profile_picture, u.aboutme, u.profession, u.location, m.content, m.created_at
	FROM messages AS m
	JOIN chats AS c ON m.chat_id = c.id
	JOIN users AS u ON m.sender_id = u.id
	WHERE c.id = ? ORDER BY m.id DESC`, chatId)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get chat messages"})
		return
	}

	var contents []Message

	for res.Next() {
		var content Message
		err := res.Scan(&content.Sender.ID, &content.Sender.Username, &content.Sender.Email, &content.Sender.ProfilePicture, &content.Sender.AboutMe, &content.Sender.Professions, &content.Sender.Location, &content.Content, &content.TimeStamp)
		if err != nil {
			return
		}
		contents = append(contents, content)
	}
	utils.DebugPrint("gotta send the messages from: ", chatId)

	utils.SendJSONResponse(w, http.StatusOK, map[string]interface{}{"messages": contents})
}

func GetChatsFromUserID(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get("uid")
	res, err := database.Query(`
	SELECT DISTINCT c.* 
	FROm chats AS c, users AS u 
	WHERE c.user1_id = ? 
	AND c.user1_id = u.id 
	OR c.user2_id = u.id`, userId)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get chat messages"})
		return
	}
	var contents []Chat
	for res.Next() {
		var content Chat
		err := res.Scan(
			&content.Id, &content.Initiator, &content.Responder, &content.Created_at,
		)
		if err != nil {
			utils.HandleError(err)
			return
		}
		contents = append(contents, content)
	}
	// utils.DebugPrint(contents)
	utils.DebugPrint("so we got the messages here: ", userId)
	utils.SendJSONResponse(w, http.StatusOK, contents)
}
