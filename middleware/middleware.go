package middleware

import (
	"DBHS/response"
	"net/http"
)

func Method(method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { 
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				response.MethodNotAllowed(w,"",nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}