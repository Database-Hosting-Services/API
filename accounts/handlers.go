package accounts

import (
	"DBHS/config"
	"net/http"
)

func signUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User signed up"))
	}
}
