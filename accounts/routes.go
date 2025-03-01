package accounts

import "DBHS/config"

func DefineURLs() {
	config.Mux.HandleFunc("POST /api/user/sign-up", signUp(config.App))
	config.Mux.HandleFunc("POST /api/user/sign-in", signIn(config.App))
}
