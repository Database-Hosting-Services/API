package main

import (
	"DBHS/accounts"
	"DBHS/indexes"
	"DBHS/projects"
	"DBHS/schemas"
	"DBHS/tables"
)

func defineURLs() {
	accounts.DefineURLs()
	indexes.DefineURLs()
	projects.DefineURLs()
	schemas.DefineURLs()
	tables.DefineURLs()
}
