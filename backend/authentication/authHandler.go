package authentication

// import (
// 	"encoding/json"
// 	"net/http"

// 	"skillswap/backend/structs"
// 	"skillswap/backend/utils"

// 	_ "github.com/go-sql-driver/mysql"
// )

// func isEmailUsed(w http.ResponseWriter, req *http.Request) {
// 	var userInfo structs.UserInfo
// 	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-186"})
// 		return
// 	}

// 	// fmt.Println(rows)
// 	if err != nil {
// 		utils.HandleError(err)
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-192"})
// 		return
// 	}
// 	if len(rows) > 0 {
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email already in use"})
// 		return
// 	} else {
// 		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Email available"})
// 		return
// 	}
// }
// func isUsernameUsed(w http.ResponseWriter, req *http.Request) {
// 	var userInfo structs.UserInfo
// 	if err := json.NewDecoder(req.Body).Decode(&userInfo); err != nil {
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-206"})
// 		return
// 	}

// 	if err != nil {
// 		utils.HandleError(err)
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "AH-211"})
// 		return
// 	}
// 	if len(rows) > 0 {
// 		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Username already in use"})
// 		return
// 	} else {
// 		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "message": "Username available"})
// 		return
// 	}
// }
