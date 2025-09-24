package users

import (
	"net/http"
	"skillswap/backend/database"
)

func GetSwapeeList(w http.ResponseWriter, req *http.Request) {
	database.Execute(`
		SELECT u.*, c.* FROM users AS u, chats AS c WHERE c.initiator != 10 AND c.user1_id = u.id OR c.user2_id = u.id
		`)
}
