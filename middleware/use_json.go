package middleware

import (
	"context"
	"mime"
	"net/http"

	"github.com/diwasrimal/echo-backend/types"
	"github.com/diwasrimal/echo-backend/utils"
)

func UseJson(nextHandler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Check for appropriate headers
		mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			utils.SendJsonResp(w, http.StatusBadRequest, types.Json{"message": "Malformed Content-Type header"})
			return
		}
		if mt != "application/json" {
			utils.SendJsonResp(
				w,
				http.StatusUnsupportedMediaType,
				types.Json{"message": "Content-Type must be application/json"},
			)
			return
		}

		// Parse json body and set it in request context
		body, err := utils.ParseJson(r.Body)
		if err != nil {
			utils.SendJsonResp(
				w,
				http.StatusBadRequest,
				types.Json{"message": "Error parsing json data!"},
			)
			return
		}
		ctx := context.WithValue(r.Context(), "body", body)

		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
