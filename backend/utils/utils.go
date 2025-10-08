package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func CheckType(toCheck string, toAllow []string) bool {
	return !strings.Contains(strings.Join(toAllow, ","), toCheck)
}

func GenerateUUID() string {
	return uuid.New().String()
}
func DebugPrint(message ...any) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d (%s)\n", filename, line, message)
}
func HandleError(err error) (b bool) {
	if err != nil {

		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d (%s)\n", filename, line, err)
		b = true
	}
	return
}
