package accounts

import (
	"DBHS/config"
	"DBHS/response"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func signUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.BadRequest(w, "Invalid Input", err)
			return
		}

		if err := checkPasswordStrength(user.Password); err != nil {
			response.BadRequest(w, "Invalid Password", err)
			return
		}

		// this return the field that exists in the database
		conflicField, err := checkUserExists(r.Context(), config.DB, user.Username, user.Email) // we can make it more generic
		if err != nil {
			response.BadRequest(w, "Invalid Input Data", err)
			return
		}

		if conflicField != "" {
			response.BadRequest(w, fmt.Sprintf("Invalid input Data, %s must be unique", conflicField), nil)
			return
		}

		data, err := SignupUser(context.Background(), config.DB, &user)
		if err != nil {
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response.Created(w, "User signed up successfully", data)
	}
}
