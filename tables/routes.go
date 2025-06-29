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
	GET 	/api/projects/{project_id}/tables/{table_id}?
			page=  		example: 1
			limit=		example: 10
			order_by=	example: name
			order=		example: asc
			filter=		example: name=x
*/

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects/{project_id}/tables").Subrouter()
	router.Use(middleware.JwtAuthMiddleware, middleware.CheckOwnership, SyncTables)

	router.Handle("", middleware.Route(map[string]http.HandlerFunc{
		http.MethodPost: CreateTableHandler(config.App),
		http.MethodGet:  GetAllTablesHanlder(config.App),
	}))
	router.Handle("/{table_id}", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet:    ReadTableHandler(config.App),
		http.MethodPut:    UpdateTableHandler(config.App),
		http.MethodDelete: DeleteTableHandler(config.App),
	}))
}

/*
	start = (page - 1) * limit
	end = start + limit
*/
