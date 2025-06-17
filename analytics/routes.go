package analytics

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/analytics").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)

	router.Handle("/current_storage", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: CurrentStorage(config.App),
	}))
}
