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
		projectOID := urlVariables["project_id"]
		userId := r.Context().Value("user-id").(int64)
		config.App.InfoLog.Printf("Checking ownership for project %s and user %d", projectOID, userId)
		ok, err := utils.CheckOwnershipQuery(r.Context(), projectOID, userId, config.DB)
		if err != nil {
			response.InternalServerError(w, err.Error(), err)
			return
		}
		config.App.InfoLog.Printf("Ownership check for project %s by user %d: %t", projectOID, userId, ok)
		if !ok {
			response.UnAuthorized(w, "UnAuthorized", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckOTableExist(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOID := urlVariables["project_id"]
		tableOID := urlVariables["table_id"]
		if tableOID != "" {
			// get the project id from the database
			_, projectId, err := utils.GetProjectNameID(r.Context(), projectOID, config.DB)
			if err != nil {
				response.InternalServerError(w, "Failed to get project ID", err)
				return
			}
			//check if the table belongs to the project
			ok, err := utils.CheckOwnershipQueryTable(r.Context(), tableOID, projectId.(int64), config.DB)
			if err != nil {
				response.InternalServerError(w, err.Error(), err)
				return
			}

			if !ok {
				response.NotFound(w, "Table not found", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if strings.HasPrefix(origin, "http://localhost") {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}