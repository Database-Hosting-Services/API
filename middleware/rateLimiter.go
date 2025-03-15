package middleware

import (
	"DBHS/utils/rateLimiter"
	"net/http"
	"DBHS/response"
	"errors"
)

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address from the request
		ip := r.RemoteAddr

		// Get the client for this IP
		client := rateLimiter.Newlimiter.GetClient(ip)

		// Check if this request is allowed
		if !client.Limiter.Allow() {
			response.TooManyRequests(w, "Rate limit exceeded", errors.New("rate limit exceeded"))
			return
		}

		next.ServeHTTP(w, r)
	})
}