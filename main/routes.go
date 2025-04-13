package main

import (
	"DBHS/accounts"
	"DBHS/projects"
	"DBHS/tables"
)

func defineURLs() {
	accounts.DefineURLs()
	projects.DefineURLs()
	tables.DefineURLs()
}
