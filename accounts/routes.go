package accounts

import "DBHS/config"

func DefineURLs() {
	config.Mux.HandleFunc("POST /api/sign-up", signUp(config.App))
}
