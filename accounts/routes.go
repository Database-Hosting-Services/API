package accounts

import (
	"DBHS/config"
	"DBHS/middleware"
)

func DefineURLs() {
	config.Mux.HandleFunc("POST /api/user/sign-up", signUp(config.App))
	config.Mux.HandleFunc("POST /api/user/sign-in", SignIn(config.App))
	config.Mux.HandleFunc("POST /api/user/verify", Verify(config.App))
	config.Mux.HandleFunc("POST /api/user/resend-code", resendCode(config.App))

	userProtected := config.Router.PathPrefix("/api/users").Subrouter()
	userProtected.Use(middleware.JwtAuthMiddleware)
	userProtected.HandleFunc("POST /update-password", UpdateUser(config.App))
	userProtected.HandleFunc("PATCH /{id}", UpdateUser(config.App))

}
