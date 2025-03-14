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

    userDynamic.HandleFunc("/sign-up", signUp(config.App)).Methods("POST")
    userDynamic.HandleFunc("/sign-in", SignIn(config.App)).Methods("POST")
    userDynamic.HandleFunc("/verify", Verify(config.App)).Methods("POST")
    userDynamic.HandleFunc("/resend-code", resendCode(config.App)).Methods("POST")
    userDynamic.HandleFunc("/forget-password", ForgetPassword(config.App)).Methods("POST")
    userDynamic.HandleFunc("/forget-password/verify", ForgetPasswordVerify(config.App)).Methods("POST")
}

func protectedRoutes() {
    userProtected := config.Router.PathPrefix("/api/users").Subrouter()

    userProtected.Use(middleware.JwtAuthMiddleware)
    userProtected.HandleFunc("/update-password", UpdatePassword(config.App)).Methods("POST")
    userProtected.HandleFunc("/{id}", UpdateUser(config.App)).Methods("PATCH")
}
