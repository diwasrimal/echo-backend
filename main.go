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

	log.Println("Using environment vars:")
	log.Printf("SERVER_PORT: %v\n", serverPort)
	log.Printf("JWT_SECRET: %v\n", jwtSecret)
	log.Printf("POSTGRES_URL: %v\n", pgUrl)
	log.Printf("ALLOWED_ORIGINS: %v\n", allowedOrigins)

	jwt.Init(jwtSecret)
	db.MustInitPostgres(pgUrl)
	defer db.Close()

	h := api.MakeHandler
	handlers := map[string]http.Handler{
		"GET /api/auth":              mw.UseJWTAuth(h(routes.AuthGet)),
		"GET /api/users/{id}":        mw.UseJWTAuth(h(routes.UsersGet)),
		"GET /api/chat-partners":     mw.UseJWTAuth(h(routes.ChatPartnersGet)),
		"GET /api/search":            mw.UseJWTAuth(h(routes.SearchGet)),
		"GET /api/messages/{pairId}": mw.UseJWTAuth(h(routes.MessagesGet)),
		"GET /api/friends":           mw.UseJWTAuth(h(routes.FriendsGet)),
		"GET /api/friend-requestors": mw.UseJWTAuth(h(routes.FriendRequestorsGet)),
		"GET /ws":                    mw.UseWebsocketJWTAuth(http.HandlerFunc(routes.WSHandleFunc)),
		"GET /api/tmp":               h(routes.TmpGet),

		"GET /api/friendship-status/{targetId}": mw.UseJWTAuth(h(routes.FriendshipStatusGet)),

		"POST /api/login":           mw.UseJson(h(routes.LoginPost)),
		"POST /api/register":        mw.UseJson(h(routes.RegisterPost)),
		"POST /api/friend-requests": mw.UseJWTAuth(mw.UseJson(h(routes.FriendRequestPost))),
		"POST /api/friends":         mw.UseJWTAuth(mw.UseJson(h(routes.FriendPost))),

		"DELETE /api/friend-requests/{targetUserId}": mw.UseJWTAuth(h(routes.FriendRequestDelete)),
		"DELETE /api/friends/{targetUserId}":         mw.UseJWTAuth(h(routes.FriendDelete)),
	}

	mux := http.NewServeMux()
	for route, handler := range handlers {
		// handler = mw.UseDelay(time.Second*2, handler) // to imitate production delay
		handler = mw.UseLogger(handler)
		mux.Handle(route, handler)
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
