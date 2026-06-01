package video

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"

	"github.com/gorilla/mux"
)

func UploadCourseVideo(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseMultipartForm(4 << 20)
	file, fileHeader, err := req.FormFile("file") // Get fileHeader
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "No file provided"})
		return
	}
	defer file.Close()

	courseID := req.FormValue("course_id")
	ext := filepath.Ext(fileHeader.Filename)
	if !utils.CheckType(ext, []string{".mp4", ".webm", ".ogv", ".mov", ".avi", ".wmv", ".flv", ".m4v"}) {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		return
	}
	uniqueFileName := utils.GenerateUUID() + ext

	if err := os.MkdirAll(filepath.Join("uploads", "courses"), 0o755); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		return
	}
	path := filepath.Join("uploads", "courses", uniqueFileName)

	// Actually write the file
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

	database.Execute("UPDATE courses SET video_path = ? WHERE id = ?", path, courseID)
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"path": path})
}

func GetCourseVideo(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	http.ServeFile(w, req, filepath.Join("uploads", "courses", fmt.Sprintf("%s.mp4", id)))
}
