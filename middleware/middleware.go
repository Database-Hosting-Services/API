package middleware

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
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

func CheckOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectId := urlVariables["project_id"]
		userId := r.Context().Value("user-id").(int)
		ok, err := utils.CheckOwnershipQuery(r.Context(), projectId, userId, config.DB)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			return
		}

		if !ok {
			response.UnAuthorized(w, "UnAuthorized", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
