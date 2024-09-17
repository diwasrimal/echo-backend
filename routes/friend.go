package routes

import (
	"log"
	"net/http"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/types"
)

// Records mutual friendship among requesting user and given
// user. Accepts json payload with field "targetId",
// which is the user that will be befriended.
// Should be used with auth middleware.
func FriendPost(w http.ResponseWriter, r *http.Request) api.Response {
	body := r.Context().Value("body").(types.Json)
	log.Printf("Hit FriendPost() with body: %v\n", body)

	userId := r.Context().Value("userId").(uint64)
	tid, ok := body["targetId"].(float64)
	if !ok {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Missing/Invalid targetId in body"},
		}
	}
	targetId := uint64(tid)

	// Can only be friends if a request was received before from the other user
	status, err := db.GetFriendshipStatus(userId, targetId)
	if err != nil {
		log.Printf("Error checking friendship status while creating new friend: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	if status != "req-received" {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "No request was received from other user"},
		}
	}

	// Record friend and delete friend request that other user sent before
	err = db.RecordFriendship(userId, targetId)
	if err != nil {
		log.Printf("Error recording friendship among (%v, %v) in db: %v\n", userId, targetId, err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	err = db.DeleteFriendRequest(targetId, userId)
	if err != nil {
		log.Printf("Error deleting prev friend request after being friends in db: %v\n", err)
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

func FriendDelete(w http.ResponseWriter, r *http.Request) api.Response {
	body := r.Context().Value("body").(types.Json)
	log.Printf("Hit FriendDelete() with body: %v\n", body)

	userId := r.Context().Value("userId").(uint64)
	tid, ok := body["targetId"].(float64)
	if !ok {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Missing/Invalid targetId in body"},
		}
	}
	targetId := uint64(tid)

	err := db.DeleteFriendship(userId, targetId)
	if err != nil {
		log.Printf("Error deleting friendship among (%v, %v) in db: %v\n", userId, targetId, err)
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

// Returns list of users that are friends of requesting user
func FriendsGet(w http.ResponseWriter, r *http.Request) api.Response {
	userId := r.Context().Value("userId").(uint64)
	log.Printf("Hit FriendsGet() with userId: %v\n", userId)
	friends, err := db.GetFriends(userId)
	if err != nil {
		log.Printf("Error getting friends from db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{"friends": friends},
	}
}
