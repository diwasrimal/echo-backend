package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/jwt"
	"github.com/diwasrimal/echo-backend/types"

	"golang.org/x/crypto/bcrypt"
)

func LoginPost(w http.ResponseWriter, r *http.Request) api.Response {
	body := r.Context().Value("body").(types.Json)
	log.Printf("Hit LoginPost() with body: %v\n", body)

	// Ensure both username and password are given
	username, usernameOk := body["username"].(string)
	password, passwordOk := body["password"].(string)
	if !usernameOk || !passwordOk {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Missing data"},
		}
	}
	username = strings.Trim(username, " \t\n\r")
	if len(username) == 0 {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "Username is empty"},
		}
	}

	// Retreive user details from username and check password
	user, err := db.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user details from db: %v\n", err)
		return api.Response{
			Status:  http.StatusInternalServerError,
			Payload: types.Json{"message": "Error logging in"},
		}
	}
	if user == nil {
		return api.Response{
			Status:  http.StatusBadRequest,
			Payload: types.Json{"message": "No such username exists"},
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return api.Response{
			Status:  http.StatusUnauthorized,
			Payload: types.Json{"message": "Incorrect password"},
		}
	}

	token := jwt.Create(types.Json{"userId": user.Id})
	log.Println("Logged in!")
	return api.Response{
		Status:  http.StatusOK,
		Payload: types.Json{"jwt": token, "userId": user.Id},
	}
}
