package main

import (
	"DBHS/accounts"
	"DBHS/projects"
	"DBHS/schemas"
)

func defineURLs() {
	accounts.DefineURLs()
	projects.DefineURLs()
	schemas.DefineURLs()
}
