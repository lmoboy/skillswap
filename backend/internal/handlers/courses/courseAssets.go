package courses

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"skillswap/backend/internal/utils"
)

func UploadCourseAsset(w http.ResponseWriter, req *http.Request) {
    _ = req.ParseMultipartForm(32 << 20)
    file, header, err := req.FormFile("file")
    if err != nil {
        utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "No file provided"})
        return
    }
    defer file.Close()

    // Validate file type - allow common course content types
    fileExt := filepath.Ext(header.Filename)
    allowedTypes := []string{".mp4", ".avi", ".mov", ".pdf", ".doc", ".docx", ".ppt", ".pptx", ".txt", ".zip", ".jpg", ".jpeg", ".png", ".gif"}
    if !utils.CheckType(fileExt, allowedTypes) {
        utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type for course asset"})
        return
    }

    courseID := req.FormValue("course_id")
    dir := filepath.Join("uploads", "courses", courseID)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
        return
    }
    path := filepath.Join(dir, header.Filename)
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

    utils.SendJSONResponse(w, http.StatusOK, map[string]string{"path": fmt.Sprintf("/api/course/%s/stream?file=%s", courseID, header.Filename)})
}

func StreamCourseAsset(w http.ResponseWriter, req *http.Request) {
    id := mux.Vars(req)["id"]
    name := req.URL.Query().Get("file")

    // Prevent path traversal: sanitize both id and name
    cleanID := filepath.Clean(id)
    cleanName := filepath.Base(filepath.Clean(name))
    if cleanID == ".." || cleanID == "." || cleanName == ".." || cleanName == "." {
        utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid path"})
        return
    }

    targetPath := filepath.Join("uploads", "courses", cleanID, cleanName)

    // Verify the resolved path stays within uploads/courses
    absTarget, err := filepath.Abs(targetPath)
    if err != nil {
        utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Invalid path"})
        return
    }
    absBase, _ := filepath.Abs(filepath.Join("uploads", "courses"))
    if !strings.HasPrefix(absTarget, absBase) {
        utils.SendJSONResponse(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
        return
    }

    http.ServeFile(w, req, targetPath)
}
