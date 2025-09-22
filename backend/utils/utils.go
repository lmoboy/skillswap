package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
func HandleError(err error) (b bool) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d (%s)\n", filename, line, err)
		b = true
	}
	return
}
