package main

import (
	"DBHS/accounts"
	"DBHS/projects"
	"DBHS/tables"
)

func defineURLs() {
	config.Mux.HandleFunc("GET /swagger/docs/", httpSwagger.WrapHandler)
	accounts.DefineURLs()
	projects.DefineURLs()
	tables.DefineURLs()
}
