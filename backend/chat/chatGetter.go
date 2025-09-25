package chat

import (
	"net/http"
	"skillswap/backend/database"
	"skillswap/backend/utils"
)

// GetMessagesFromUID retrieves all messages for the chat specified by the "cid" query parameter and writes them as JSON.
// Messages include sender details (id, username, email, profile picture, about me, professions, location), message content, and timestamp.
// Messages are returned ordered by message ID in descending order.
// On database query error the handler sends HTTP 500 with `{"error": "Failed to get chat messages"}`.
// If scanning a result row fails the handler returns early and does not write a response.
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

// GetChatsFromUserID reads the "uid" query parameter, retrieves all chats involving that user from the database, and writes the resulting chats as JSON to the response.
// GetChatsFromUserID retrieves all chats involving the user specified by the "uid" query parameter and writes them as JSON to the response.
// It queries the database for distinct chats for that user and returns HTTP 200 with the slice of Chat objects on success.
// If the initial database query fails it logs the error and sends HTTP 500 with `{"error": "Failed to get chat messages"}`.
// If scanning a result row fails it logs the error and returns immediately without writing a response.
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



// LoadMessagesFromDatabase loads the latest 100 messages along with each sender's details from the database.
// It returns a slice of Message populated with sender fields, content, and timestamp, or an error if the query or row scanning fails.
func LoadMessagesFromDatabase() ([]Message, error) {
	res, err := database.Query(`
		SELECT u.id, u.username, u.email, u.profile_picture, u.aboutme, u.profession, u.location, m.content, m.created_at
		FROM messages AS m
		JOIN users AS u ON m.sender_id = u.id
		ORDER BY m.id DESC
		LIMIT 100`)
	if err != nil {
		utils.HandleError(err)
		return nil, err
	}
	defer res.Close()

	var messages []Message
	for res.Next() {
		var msg Message
		err := res.Scan(&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.Email, &msg.Sender.ProfilePicture, &msg.Sender.AboutMe, &msg.Sender.Professions, &msg.Sender.Location, &msg.Content, &msg.TimeStamp)
		if err != nil {
			utils.HandleError(err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}