package indexes

import (
	"DBHS/config"
	"DBHS/middleware"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/indexes").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)
}
