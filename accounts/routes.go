package accounts

import (
	"DBHS/config"
	"DBHS/middleware"
	"net/http"
)

func DefineURLs() {
	config.Mux.HandleFunc("POST /api/user/sign-up", signUp(config.App))
	config.Mux.HandleFunc("POST /api/user/sign-in", SignIn(config.App))
	config.Mux.HandleFunc("POST /api/user/verify", Verify(config.App))
	config.Mux.HandleFunc("POST /api/user/resend-code", resendCode(config.App))

	config.Mux.Handle(
		"PATCH /api/user/update-password",
		middleware.JwtAuthMiddleware(http.HandlerFunc(UpdatePassword(config.App))),
	)

}
