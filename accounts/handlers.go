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

func SignIn(app *config.Application) http.HandlerFunc {
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

		resp, err := SignInUser(r.Context(), config.DB, &user)
		if err != nil {
			if err.Error() == "no rows in result set" || err.Error() == "InCorrect Email or Password" {
				response.BadRequest(w, "InCorrect Email or Password", nil)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response.OK(w, "User signed in successfully", resp)
	}
}

func Verify(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserVerify
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			app.ErrorLog.Println(err.Error())
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}

		data, err := VerifyUser(r.Context(), config.DB, config.VerifyCache, &user)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			if err.Error() == "Wrong verification code" {
				response.BadRequest(w, err.Error(), err)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}
		app.InfoLog.Println("User verified successfully", user.Username)
		response.Created(w, "User verified successfully", data)
	}
}
