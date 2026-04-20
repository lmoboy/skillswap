package users

import (
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
)

func GetSwapperList(w http.ResponseWriter, req *http.Request) {
	requestId := req.URL.Query().Get("uid")
	res, err := database.Execute(`
		SELECT u.*, c.* FROM users AS u, chats AS c WHERE c.user2_id = ?
		`, requestId)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get swapees"})
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, res)
}
