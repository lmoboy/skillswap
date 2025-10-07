package courses

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gorilla/mux"
    "skillswap/backend/utils"
)

func UploadCourseAsset(w http.ResponseWriter, req *http.Request) {
    _ = req.ParseMultipartForm(32 << 20)
    file, header, _ := req.FormFile("file")
    defer file.Close()

    courseID := req.FormValue("course_id")
    dir := filepath.Join("uploads", "courses", courseID)
    _ = os.MkdirAll(dir, 0o755)
    path := filepath.Join(dir, header.Filename)
    dst, _ := os.Create(path)
    defer dst.Close()
    io.Copy(dst, file)

    utils.SendJSONResponse(w, http.StatusOK, map[string]string{"path": fmt.Sprintf("/api/course/%s/stream?file=%s", courseID, header.Filename)})
}

func StreamCourseAsset(w http.ResponseWriter, req *http.Request) {
    id := mux.Vars(req)["id"]
    name := req.URL.Query().Get("file")
    http.ServeFile(w, req, filepath.Join("uploads", "courses", id, name))
}
