package indexes

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

// prefix routes with: `/api/projects/{project_id}/indexes""`
var empty = map[string]http.HandlerFunc{
	http.MethodPost: CreateIndex(config.App),
	http.MethodGet:  ProjectIndexes(config.App),
}

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/indexes").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)

	router.Handle("", middleware.Route(empty))
}
