package accounts

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	dynamicRoutes()
	protectedRoutes()
}

func dynamicRoutes() {
	userDynamic := config.Router.PathPrefix("/api/user").Subrouter()
	// global middleware
	userDynamic.Use(middleware.MethodsAllowed(http.MethodPost))

	userDynamic.Handle("/sign-up", signUp(config.App))
	userDynamic.Handle("/sign-in", SignIn(config.App))
	userDynamic.Handle("/verify", Verify(config.App))
	userDynamic.Handle("/resend-code", resendCode(config.App))
	userDynamic.Handle("/forget-password", ForgetPassword(config.App))
	userDynamic.Handle("/forget-password/verify", ForgetPasswordVerify(config.App))
}

func protectedRoutes() {
	userProtected := config.Router.PathPrefix("/api/users").Subrouter()
	// global middleware
	userProtected.Use(middleware.JwtAuthMiddleware)

	userProtected.Handle("/update-password", middleware.MethodsAllowed(http.MethodPost)(UpdatePassword(config.App)))
	userProtected.Handle("/{id}", middleware.MethodsAllowed(http.MethodPatch)(UpdateUser(config.App)))

}
