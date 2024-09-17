package middleware

import (
	"log"
	"net/http"
	"time"
)

// UseLogger is a middleware that logs HTTP requests with response time
func UseLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqName := r.Method + " " + r.RequestURI
		log.Printf(reqName)
		next.ServeHTTP(w, r)
		log.Printf("Took %v\n", time.Since(start))
	}
	return http.HandlerFunc(fn)
}
