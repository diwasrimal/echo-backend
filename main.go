package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/diwasrimal/echo-backend/api"
	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/jwt"
	mw "github.com/diwasrimal/echo-backend/middleware"
	"github.com/diwasrimal/echo-backend/routes"
	"github.com/diwasrimal/echo-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
)

func main() {
	var (
		serverPort     = utils.MustGetEnv("SERVER_PORT")
		jwtSecret      = utils.MustGetEnv("JWT_SECRET")
		pgUrl          = utils.MustGetEnv("POSTGRES_URL")
		allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	)

	jwt.Init(jwtSecret)
	db.MustInitPostgres(pgUrl)
	defer db.Close()

	h := api.MakeHandler
	handlers := map[string]http.Handler{
		"GET /api/logout":            mw.UseJson(h(routes.LogoutGet)),
		"GET /api/auth":              mw.UseJWTAuth(h(routes.AuthGet)),
		"GET /api/users/{id}":        mw.UseJWTAuth(mw.UseJson(h(routes.UsersGet))),
		"GET /api/chat-partners":     mw.UseJWTAuth(mw.UseJson(h(routes.ChatPartnersGet))),
		"GET /api/search":            mw.UseJWTAuth(mw.UseJson(h(routes.SearchGet))),
		"GET /api/messages/{pairId}": mw.UseJWTAuth(mw.UseJson(h(routes.MessagesGet))),
		"GET /api/friends":           mw.UseJWTAuth(mw.UseJson(h(routes.FriendsGet))),
		"GET /api/friend-requestors": mw.UseJWTAuth(mw.UseJson(h(routes.FriendRequestorsGet))),
		"GET /ws":                    mw.UseJWTAuth(http.HandlerFunc(routes.WSHandleFunc)),
		"GET /api/tmp":               h(routes.TmpGet),

		"GET /api/friendship-status/{targetId}": mw.UseJWTAuth(mw.UseJson(h(routes.FriendshipStatusGet))),

		"POST /api/login":           mw.UseJson(h(routes.LoginPost)),
		"POST /api/register":        mw.UseJson(h(routes.RegisterPost)),
		"POST /api/friend-requests": mw.UseJWTAuth(mw.UseJson(h(routes.FriendRequestPost))),
		"POST /api/friends":         mw.UseJWTAuth(mw.UseJson(h(routes.FriendPost))),

		"DELETE /api/friend-requests": mw.UseJWTAuth(mw.UseJson(h(routes.FriendRequestDelete))),
		"DELETE /api/friends":         mw.UseJWTAuth(mw.UseJson(h(routes.FriendDelete))),
	}

	mux := http.NewServeMux()
	for route, handler := range handlers {
		mux.Handle(route, mw.UseLogger(handler))
	}

	_ = allowedOrigins
	handler := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(allowedOrigins, ","),
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowCredentials: true,
	}).Handler(mux)

	addr := "0.0.0.0:" + serverPort
	log.Printf("Listening on %v...\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
