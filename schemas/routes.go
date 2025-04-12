package schemas

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project-id}/schema/tables").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)

	router.Handle("", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: GetDatabaseSchema(config.App),
	}))

	router.Handle("/{table-id}", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: GetDatabaseTableSchema(config.App),
	}))
}
