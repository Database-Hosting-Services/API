package accounts

import (
	"DBHS/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
sign up logic steps :

	1- Check if the email already exists in the database.
	2- Hash the password using bcrypt before storing.
	3- Generate a verification code (e.g., a 6-digit number).
	4- Store user data in the database with is_verified = false.
	5- Send the verification code via email.
	6- Respond with success message.
*/
func signUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := checkPasswordStrength(user.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		conflicField, err := checkUserExists(r.Context(), config.DB, user.Username, user.Email)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
		}

		if conflicField != "" {
			http.Error(w, fmt.Sprintf("Invalid input %s must be unique", conflicField), http.StatusBadRequest)
			return
		}

		err = SignupUser(context.Background(), config.DB, &user)
		if err != nil {
			http.Error(w, "error creating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
