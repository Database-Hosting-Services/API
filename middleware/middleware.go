package middleware

import (
	"DBHS/response"
	"maps"
	"net/http"
	"slices"
	"strings"
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
			response.MethodNotAllowed(w, strings.Join(methods, ","), "", nil)
		})
	}
}

func Route(hundlers map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, ok := hundlers[r.Method]
		if !ok {
			response.MethodNotAllowed(w, strings.Join(slices.Collect(maps.Keys(hundlers)), ","), "", nil)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
