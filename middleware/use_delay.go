package middleware

import (
	"net/http"
	"time"
)

// Middleware function to imitate processing delay.
func UseDelay(delay time.Duration, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
