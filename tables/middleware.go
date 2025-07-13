package tables

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func SyncTables(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sync table schemas between old and new
		requestVars := mux.Vars(r)
		projectOID := requestVars["project_id"]
		if projectOID == "" {
			response.NotFound(w, r, "Project ID is required", nil)
			return
		}
		// Extract user ID from context
		userId, ok := r.Context().Value("user-id").(int64)
		if !ok || userId == 0 {
			response.UnAuthorized(w, r, "Unauthorized", nil)
			return
		}

		// Get user database connection
		projectId, userDb, err := utils.ExtractDb(r.Context(), projectOID, userId, config.DB)
		if err != nil {
			response.InternalServerError(w, r, "Failed to extract database connection", err)
			return
		}
		
		if err := SyncTableSchemas(r.Context(), projectId, config.DB, userDb); err != nil {
			response.InternalServerError(w, r, err.Error(), err)
			return
		}
		next.ServeHTTP(w, r)
	})
}