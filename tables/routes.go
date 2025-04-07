package tables

import (
	"DBHS/config"
	"DBHS/middleware"
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
	
	router.Handle("",)

}