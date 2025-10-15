package auth

import (
	"net/http"
	"skillswap/backend/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

func Logout(w http.ResponseWriter, req *http.Request) {
	if err := RemoveSession(w, req); err != nil {
		utils.HandleError(err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})
}
