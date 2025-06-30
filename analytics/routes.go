package analytics

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/analytics").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)

	router.Handle("/storage", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: CurrentStorage(config.App),
	}))

	router.Handle("/execution-time", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: ExecutionTime(config.App),
	}))

	router.Handle("/usage", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: DatabaseUsage(config.App),
	}))
}
