package users

import (
	"net/http"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"
)

func GetSwapeeList(w http.ResponseWriter, req *http.Request) {
	// utils.DebugPrint(req.URL.RawQuery)
	requestId := req.URL.Query().Get("uid")
	// utils.DebugPrint(requestId)
	res, err := database.Query(`
		SELECT * FROm chats AS c, users AS u WHERE c.initiated_by = ? AND c.user1_id = u.id OR c.user2_id = u.id;
`, requestId)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get swapees"})
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, res)
}
