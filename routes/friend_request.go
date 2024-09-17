package routes

import (
	"log"
	"net/http"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/types"
)

// Records a new entry into the friend requests table.
// Accepts json payload with field "targetId", which is the user
// that will receive this friend request. Requestor is the one
// who made this request, i.e. the logged in user.
// Should be used with auth middleware
func FriendRequestPost(w http.ResponseWriter, r *http.Request) api.Response {
	body := r.Context().Value("body").(types.Json)
	log.Printf("Hit FriendRequestPost() with body: %v\n", body)

	userId := r.Context().Value("userId").(uint64)
	targetId, ok := body["targetId"].(float64)
	if !ok {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Missing/Invalid targetId in body"},
		}
	}

	// Check status of friendship before making a request, if the status
	// not "unknown", i.e is "friends", "req-sent" or "req-received", then we
	// can't send friend request
	status, err := db.GetFriendshipStatus(userId, uint64(targetId))
	if err != nil {
		log.Printf("Error getting friendship status while creating friend req: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	if status != "unknown" {
		return api.Response{
			Status: http.StatusBadRequest,
			Payload: types.Json{
				"message": "Other user is either a friend or a request is sent/received already",
			},
		}
	}

	err = db.RecordFriendRequest(userId, uint64(targetId)) // from userId -> targetId
	if err != nil {
		log.Printf("Error recording friend request in db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	return api.Response{
		Status:  http.StatusCreated,
		Payload: types.Json{},
	}
}

// Deletes a friend request send from requesting user to provided user
// if the request was sent previously.
func FriendRequestDelete(w http.ResponseWriter, r *http.Request) api.Response {
	body := r.Context().Value("body").(types.Json)
	log.Printf("Hit FriendRequestDelete() with body: %v\n", body)

	userId := r.Context().Value("userId").(uint64)
	tid, ok := body["targetId"].(float64)
	if !ok {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Missing/Invalid targetId in body"},
		}
	}
	targetId := uint64(tid)

	// Check the friendship status beforing deleting request. If request is
	// not sent, we can't delete it
	status, err := db.GetFriendshipStatus(userId, uint64(targetId))
	if err != nil {
		log.Printf("Error getting friendship status while deleting friend req: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	if status != "req-sent" {
		return api.Response{
			Status: http.StatusBadRequest,
			Payload: types.Json{
				"message": "Request wasn't sent, cannnot delete!",
			},
		}
	}

	err = db.DeleteFriendRequest(userId, targetId)
	if err != nil {
		log.Printf("Error deleting friend request in db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{},
	}
}
