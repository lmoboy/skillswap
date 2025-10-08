package video

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"skillswap/backend/database"
	"skillswap/backend/utils"

	"github.com/gorilla/mux"
)

func UploadCourseVideo(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseMultipartForm(4 << 20)
	file, fileHeader, _ := req.FormFile("file") // Get fileHeader
	defer file.Close()

	courseID := req.FormValue("course_id")
	ext := filepath.Ext(fileHeader.Filename)
	if utils.CheckType(ext, []string{".mp4", ".webm", ".ogv", ".mov", ".avi", ".wmv", ".flv", ".m4v"}) {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type"})
		return
	}
	uniqueFileName := utils.GenerateUUID() + ext

	_ = os.MkdirAll(filepath.Join("uploads", "courses"), 0o755)
	path := filepath.Join("uploads", "courses", uniqueFileName)
	database.Execute("UPDATE courses SET video_path = ? WHERE id = ?", path, courseID)

}

func GetCourseVideo(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	http.ServeFile(w, req, filepath.Join("uploads", "courses", fmt.Sprintf("%s.mp4", id)))
}
