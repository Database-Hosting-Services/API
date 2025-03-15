package middleware

import (
	"DBHS/response"
	"net/http"
)

func MethodsAllowed(methods ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { 
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, method := range methods {
				if r.Method == method {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.MethodNotAllowed(w,"",nil)
		})
	}
}