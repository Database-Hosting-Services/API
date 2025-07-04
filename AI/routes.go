package ai

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	AIProtected := config.Router.PathPrefix("/api/projects/{project_id}/ai").Subrouter()
	AIProtected.Use(middleware.JwtAuthMiddleware, middleware.CheckOwnership)

	AIProtected.Handle("/report", middleware.MethodsAllowed(http.MethodGet)(Report(config.App)))
	AIProtected.Handle(("/chatbot/ask"), middleware.MethodsAllowed(http.MethodPost)(ChatBotAsk(config.App)))
	
	AIProtected.Handle("/agent", middleware.MethodsAllowed(http.MethodPost)(Agent(config.App)))
	AIProtected.Handle("/agent/accept", middleware.MethodsAllowed(http.MethodPost)(Agent(config.App)))
}