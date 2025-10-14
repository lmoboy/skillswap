package users

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/mux"
)

func UploadProfilePicture(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseMultipartForm(4 << 20)
	file, fileHeader, _ := req.FormFile("file")
	defer file.Close()
	if !utils.CheckType(filepath.Ext(fileHeader.Filename), []string{".jpg", ".jpeg", ".png"}) {
		utils.DebugPrint("type not accepted")
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		return
	}
	userID := req.FormValue("user_id")
	if err := os.MkdirAll(filepath.Join("uploads", "users"), 0o755); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		return
	}
	path := filepath.Join("uploads", "users", fmt.Sprintf("%s.jpg", userID))
	dst, err := os.Create(path)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	publicPath := fmt.Sprintf("/api/profile/%s/picture", userID)
	database.Execute("UPDATE users SET profile_picture = ? WHERE id = ?", publicPath, userID)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"profile_picture": publicPath})
}

func GetProfilePicture(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	http.ServeFile(w, req, filepath.Join("uploads", "users", fmt.Sprintf("%s.jpg", id)))
}
