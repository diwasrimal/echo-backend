package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/types"
)

// Gets a user's profile, should be used with auth middleware
// route: /path/to/route/{id}
func UsersGet(w http.ResponseWriter, r *http.Request) api.Response {
	log.Printf("Hit UserGet(), id: %v\n", r.PathValue("userId"))
	uid, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Invalid userId in path value"},
		}
	}
	userId := uint64(uid)
	user, err := db.GetUserById(userId)
	if err != nil {
		log.Printf("Error getting user by id in db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}

	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{"user": user},
	}
}
