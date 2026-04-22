package video

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/utils"

	"github.com/gorilla/mux"
)

func UploadCourseVideo(w http.ResponseWriter, req *http.Request) {
	// Increase maximum upload size to 2GB
	req.ParseMultipartForm(2000 << 20)

	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Check file size
	if fileHeader.Size > 2000<<20 {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "File too large. Maximum size is 2GB."})
		return
	}

	courseID := req.FormValue("course_id")
	ext := filepath.Ext(fileHeader.Filename)

	// Fix the file type validation logic (was inverted)
	if !utils.CheckType(ext, []string{".mp4", ".webm", ".ogv", ".mov", ".avi", ".wmv", ".flv", ".m4v"}) {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type. Supported formats: mp4, webm, ogv, mov, avi, wmv, flv, m4v"})
		return
	}

	uniqueFileName := utils.GenerateUUID() + ext
	uploadPath := filepath.Join("uploads", "courses")

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
		return
	}

	// Save file to disk
	fullPath := filepath.Join(uploadPath, uniqueFileName)
	out, err := os.Create(fullPath)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
		return
	}
	defer out.Close()

	// Copy file data
	_, err = io.Copy(out, file)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	// Update database with file path
	_, err = database.Execute("UPDATE courses SET video_path = ? WHERE id = ?", fullPath, courseID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update database"})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "File uploaded successfully",
		"path":    fullPath,
	})
}

func GetCourseVideo(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	// Sanitize the file path to prevent directory traversal
	filename := filepath.Base(id)
	if filename == "" || filename == "." || filename == ".." {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", "courses", filename+".mp4")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, req, filePath)
}
