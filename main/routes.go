package main

import (
	"DBHS/accounts"
	"DBHS/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

func defineURLs() {
	config.Mux.HandleFunc("GET /swagger/docs/", httpSwagger.WrapHandler)
	accounts.DefineURLs()
}
