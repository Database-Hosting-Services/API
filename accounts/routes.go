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

    postMiddleware := middleware.Method(http.MethodPost)
    userDynamic.Use(postMiddleware)
    userDynamic.HandleFunc("/sign-up", signUp(config.App))
    userDynamic.HandleFunc("/sign-in", SignIn(config.App))
    userDynamic.HandleFunc("/verify", Verify(config.App))
    userDynamic.HandleFunc("/resend-code", resendCode(config.App))
    userDynamic.HandleFunc("/forget-password", ForgetPassword(config.App))
    userDynamic.HandleFunc("/forget-password/verify", ForgetPasswordVerify(config.App))
}

func protectedRoutes() {
    userProtected := config.Router.PathPrefix("/api/users").Subrouter()

    userProtected.Use(middleware.JwtAuthMiddleware)
    userProtected.HandleFunc("/update-password", UpdatePassword(config.App)).Methods("POST")
    userProtected.HandleFunc("/{id}", UpdateUser(config.App)).Methods("PATCH")
}
