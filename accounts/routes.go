package accounts

import (
	"DBHS/config"
	"DBHS/middleware"
)

func DefineURLs() {
	dynamicRoutes()
	protectedRoutes()
}

func dynamicRoutes() {
	userDynamic := config.Router.PathPrefix("/api/user").Subrouter()

	userDynamic.HandleFunc("POST /sign-up", signUp(config.App))
	userDynamic.HandleFunc("POST /sign-in", SignIn(config.App))
	userDynamic.HandleFunc("POST /verify", Verify(config.App))
	userDynamic.HandleFunc("POST /resend-code", resendCode(config.App))
	userDynamic.HandleFunc("POST /forget-password", ForgetPassword(config.App))
	userDynamic.HandleFunc("POST /forget-password/verify", ForgetPasswordVerify(config.App))
}

func protectedRoutes() {
	userProtected := config.Router.PathPrefix("/api/users").Subrouter()

	userProtected.Use(middleware.JwtAuthMiddleware)
	userProtected.HandleFunc("POST /update-password", UpdatePassword(config.App))
	userProtected.HandleFunc("PATCH /{id}", UpdateUser(config.App))
}
