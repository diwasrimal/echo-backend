package routes

import (
	"log"
	"net/http"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/types"
)

// Should be used with auth middleware to work as expected.
// This function assumes that authentication was handled by
// middleware and hence just returns a ok status with logged in userid
func AuthGet(w http.ResponseWriter, r *http.Request) api.Response {
	userId := r.Context().Value("userId").(uint64)
	log.Printf("Login ok for userId: %v\n", userId)
	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{"userId": userId},
	}
}
