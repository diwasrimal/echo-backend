package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/jwt"
	"github.com/diwasrimal/echo-backend/types"
	"github.com/diwasrimal/echo-backend/utils"
)

// Authorizes request by validating session id contained in the cookie
// and adds user id assosiated with the session to the request context.
func UseCookieAuth(nextHandler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionId")
		if err != nil {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Missing session cookie"})
			return
		}
		sessionId := cookie.Value
		if len(sessionId) == 0 {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Invalid session credentials"})
			return
		}
		session, err := db.GetSession(sessionId)
		if err != nil {
			log.Printf("Error getting session from db: %v\n", err)
			utils.SendJsonResp(w, http.StatusInternalServerError, types.Json{"message": "Error validating session credentials"})
			return
		}
		if session == nil {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Invalid session credentials"})
			return
		}

		ctx := context.WithValue(r.Context(), "userId", session.UserId)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Performs a JWT authentication and saves the requesting
// user id in request's context
func UseJWTAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) != 2 {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{
				"message": "Missing or invalid 'Authorization' header",
			})
			return
		}
		token := parts[1]

		// Ensure that token has valid signature and is not expired
		validSignature, jwtPayload := jwt.VerifyAndDecode(token)
		if !validSignature {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Invalid token"})
			return
		}
		expTime := int64(jwtPayload["exp"].(float64))
		if time.Now().Unix() > expTime {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Token has expired"})
			return
		}

		userId := uint64(jwtPayload["userId"].(float64))
		ctx := context.WithValue(r.Context(), "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// The browser's WebSocket API doesnot allow adding Authorization headers,
// Thus we can't get token from http headers, so get it as a query param
// and authorize the websocket connection request
// More: https://stackoverflow.com/questions/4361173/http-headers-in-websockets-client-api
func UseWebsocketJWTAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("jwt")
		if len(token) == 0 {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Include jwt token as query param jwt=token"})
			return
		}
		validSignature, jwtPayload := jwt.VerifyAndDecode(token)
		if !validSignature {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Invalid token"})
			return
		}
		expTime := int64(jwtPayload["exp"].(float64))
		if time.Now().Unix() > expTime {
			utils.SendJsonResp(w, http.StatusUnauthorized, types.Json{"message": "Token has expired"})
			return
		}
		userId := uint64(jwtPayload["userId"].(float64))
		ctx := context.WithValue(r.Context(), "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
