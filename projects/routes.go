package projects

import (
	"DBHS/config"
	"DBHS/middleware"
)

func DefineURLs() {
	router := config.Router.PathPrefix("/api/projects").Subrouter()
	router.Use(middleware.JwtAuthMiddleware)
	
	router.HandleFunc("/", CreateProject(config.App)).Methods("POST")
	router.HandleFunc("/", GetProjects(config.App)).Methods("GET")
	router.HandleFunc("/{project_id}", getSpecificProject(config.App)).Methods("GET")
}
