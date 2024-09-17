package routes

import (
	"net/http"

	"github.com/diwasrimal/echo-backend/api"
)

func TmpGet(w http.ResponseWriter, r *http.Request) api.Response {
	return api.Response{
		Status:  200,
		Payload: map[string]any{"message": "Hello from tmp route"},
	}
}
