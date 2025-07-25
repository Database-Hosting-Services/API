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
			response.MethodNotAllowed(w, r, strings.Join(methods, ","), "", nil)
		})
	}
}

func Route(hundlers map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, ok := hundlers[r.Method]
		if !ok {
			response.MethodNotAllowed(w, r, strings.Join(slices.Collect(maps.Keys(hundlers)), ","), "", nil)
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
		// check if the project exists
		exists, err := utils.CheckProjectExist(r.Context(), projectOID, config.DB)
		if err != nil {
			response.InternalServerError(w, r, "Failed to check project existence", err)
			return
		}

		if !exists {
			config.App.ErrorLog.Printf("Project %s does not exist", projectOID)
			response.NotFound(w, r, "Project not found", nil)
			return
		}
		// check if the user is the owner of the project
		ok, err := utils.CheckOwnershipQuery(r.Context(), projectOID, userId, config.DB)
		if err != nil {
			response.InternalServerError(w, r, err.Error(), err)
			return
		}
		config.App.InfoLog.Printf("Ownership check for project %s by user %d: %t", projectOID, userId, ok)
		if !ok {
			response.UnAuthorized(w, r, "UnAuthorized", nil)
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
				response.InternalServerError(w, r, "Failed to get project ID", err)
				return
			}
			// check if the table exists
			exists, err := utils.CheckTableExist(r.Context(), tableOID, config.DB)
			if err != nil {
				response.InternalServerError(w, r, "Failed to check table existence", err)
				return
			}

			if !exists {
				config.App.ErrorLog.Printf("Table %s does not exist in project %s", tableOID, projectOID)
				response.NotFound(w, r, "Table not found", nil)
				return
			}
			//check if the table belongs to the project
			ok, err := utils.CheckOwnershipQueryTable(r.Context(), tableOID, projectId.(int64), config.DB)
			if err != nil {
				response.InternalServerError(w, r, err.Error(), err)
				return
			}

			if !ok {
				response.NotFound(w, r, "Table not found", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow all common HTTP methods used by the API
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		// Allow common headers used by the API
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")

		// Note: When using "*" for Allow-Origin, we cannot use Allow-Credentials: true
		// If you need credentials, you'll need to specify specific origins instead of "*"
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
