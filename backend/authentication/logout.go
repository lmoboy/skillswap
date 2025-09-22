package authentication

import (
	"database/sql"
	"net/http"
	"skillswap/backend/utils"

	_ "github.com/go-sql-driver/mysql"
)

func Logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "authentication")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request: No session found"})
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Session save failed"})
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Failed to connect to database"})
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE session_token = ?", session.ID)
	if err != nil {
		utils.HandleError(err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error: Failed to remove session from database"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Logout successful"})
}
