package tables

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

/*
	POST 	/api/projects/{project_id}/tables
	PUT 	/api/projects/{project_id}/tables/{table_id}
	DELETE 	/api/projects/{project_id}/tables/{table_id}
	GET 	/api/projects/{project_id}/tables/{table_id}?Page=X&Limit=Y
*/



func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/tables").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)
	
	router.Handle("", middleware.MethodsAllowed(http.MethodPost)(CreateTableHandler(config.App)))
	router.Handle("/{table_id}", middleware.Route( map[string]http.HandlerFunc{
		http.MethodPut: UpdateTableHandler(config.App),

	}))
}