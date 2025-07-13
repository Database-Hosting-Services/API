package main

import (
	"DBHS/AI"
	sqleditor "DBHS/SqlEditor"
	"DBHS/accounts"
	"DBHS/analytics"
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
	analytics.DefineURLs()
	sqleditor.DefineURLs()
}
