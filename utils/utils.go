package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/diwasrimal/echo-backend/types"
)

func ParseJson(body io.ReadCloser) (types.Json, error) {
	data := make(types.Json)
	err := json.NewDecoder(body).Decode(&data)
	return data, err
}

func SendJsonResp(w http.ResponseWriter, status int, payload types.Json) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func MustGetEnv(key string) string {
	value, set := os.LookupEnv(key)
	if !set {
		panic(fmt.Sprintf("Environment variable %q not set", key))
	}
	return value
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
