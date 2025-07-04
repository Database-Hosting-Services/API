package sqleditor

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/sqlEditor").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)

	router.Handle("/run-query", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: RunSqlQuery(config.App),
	}))
}
