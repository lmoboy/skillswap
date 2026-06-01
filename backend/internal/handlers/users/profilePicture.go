package users
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/handlers/auth"
	"skillswap/backend/internal/utils"

	"github.com/gorilla/mux"
)

func UploadProfilePicture(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseMultipartForm(4 << 20)
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "No file provided"})
		return
	}
	defer file.Close()
	if !utils.CheckType(filepath.Ext(fileHeader.Filename), []string{".jpg", ".jpeg", ".png"}) {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		return
	}

	// Get user ID from session, not from form (prevents impersonation)
	session, err := auth.Store.Get(req, "authentication")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Not authenticated"})
		return
	}
	userID, err := database.GetUserIDFromEmail(email)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "User not found"})
		return
	}
	userIDStr := fmt.Sprintf("%d", userID)
	if err := os.MkdirAll(filepath.Join("uploads", "users"), 0o755); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		return
	}
	path := filepath.Join("uploads", "users", fmt.Sprintf("%s.jpg", userIDStr))
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

	publicPath := fmt.Sprintf("/api/profile/%s/picture", userIDStr)
	database.Execute("UPDATE users SET profile_picture = ? WHERE id = ?", publicPath, userID)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"profile_picture": publicPath})
}

func GetProfilePicture(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	http.ServeFile(w, req, filepath.Join("uploads", "users", fmt.Sprintf("%s.jpg", id)))
}
