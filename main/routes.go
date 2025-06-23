package main

import (
	"DBHS/accounts"
	"DBHS/AI"
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
	ai.DefineURLs()
}
