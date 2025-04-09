package main

import (
	"DBHS/accounts"
	"DBHS/indexes"
	"DBHS/projects"
)

func defineURLs() {
	accounts.DefineURLs()
	indexes.DefineURLs()
	projects.DefineURLs()
}
