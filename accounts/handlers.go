package accounts

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
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

func signIn(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserSignIn
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}

		if user.Email == "" || user.Password == "" {
			response.BadRequest(w, "Email and Password are required", nil)
			return
		}

		authenticatedUser, err := GetUser(r.Context(), config.DB, user.Email)
		if err != nil {
			if err.Error() == "no rows in result set" {
				response.BadRequest(w, "InCorrect Email or Password", nil)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		if !CheckPasswordHash(user.Password, authenticatedUser.Password) {
			response.BadRequest(w, "InCorrect Email or Password", nil)
			return
		}

		tokenString := utils.NewToken()
		tokenString.AddClaim("oid", authenticatedUser.OID)
		token, err := tokenString.String()

		if err != nil {
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response := map[string]interface{}{
			"message": "User signed in successfully",
			"status":  "success",
			"Data": map[string]interface{}{
				"oid":      authenticatedUser.OID,
				"username": authenticatedUser.Username,
				"email":    authenticatedUser.Email,
				"image":    authenticatedUser.Image,
				"token":    token,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
