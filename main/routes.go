package main

import (
	"DBHS/accounts"
	"DBHS/projects"
	"DBHS/schemas"
	"DBHS/tables"
)

func defineURLs() {
	accounts.DefineURLs()
	projects.DefineURLs()
	schemas.DefineURLs()
	tables.DefineURLs()
}
