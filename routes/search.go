package routes

import (
	"log"
	"net/http"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/types"
)

// Searches for a user by their username based on the search query
// provided search type. Search type can be "by-username" or "normal".
// Searching by username is exact, while normal search is fuzzy.
// Should be used with auth
func SearchGet(w http.ResponseWriter, r *http.Request) api.Response {
	searchType := r.URL.Query().Get("type")
	searchQuery := r.URL.Query().Get("query")
	log.Printf("Hit SearchGet() with type: %q, query: %q\n", searchType, searchQuery)
	if searchType != "normal" && searchType != "by-username" {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Invalid search type"},
		}
	}
	if len(searchQuery) == 0 {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Search query cannot be empty"},
		}
	}
	results, err := db.SearchUser(searchType, searchQuery)
	if err != nil {
		log.Printf("Error getting user search results from db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{},
		}
	}
	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{"results": results},
	}
}
