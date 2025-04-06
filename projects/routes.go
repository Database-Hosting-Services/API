package projects

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)


	router.Handle("", middleware.Route(map[string]http.HandlerFunc{
		http.MethodPost: CreateProject(config.App),
		http.MethodGet: GetProjects(config.App),
	}))
	
	router.Handle("/{project_id}", middleware.Route(map[string]http.HandlerFunc{
		http.MethodGet: getSpecificProject(config.App),
		http.MethodPatch: updateProject(config.App),
		http.MethodDelete: DeleteProject(config.App),
	}))
}
