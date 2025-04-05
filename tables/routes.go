package tables

import (
	"DBHS/config"
	"DBHS/middleware"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/tables").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)


}