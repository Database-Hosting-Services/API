package accounts

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func signUp(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserUnVerified
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.BadRequest(w, "Invalid Input", err)
			return
		}

		if err := checkPasswordStrength(user.Password); err != nil {
			response.BadRequest(w, "Invalid Password", err)
			return
		}

		field, err := checkUserExistsInCache(user.Username, user.Email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Server Error", err)
			return
		}

		if field != "" {
			response.BadRequest(w,
				"Invalid User",
				errors.New(fmt.Sprintf("User with this %s already exists", field)),
			)
			return
		}

		// this return the field that exists in the database
		field, err = checkUserExists(r.Context(), config.DB, user.Username, user.Email) // we can make it more generic
		if err != nil {
			response.BadRequest(w, "Invalid Input Data", err)
			return
		}

		if field != "" {
			response.BadRequest(w, fmt.Sprintf("Invalid input Data, this %s is already exists", field), nil)
			return
		}

		err = SignupUser(context.Background(), config.DB, &user)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response.Created(w, "User signed up successfully, check your email for verification", nil)
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

		resp, err := SignInUser(r.Context(), config.DB, config.VerifyCache, &user)
		if err != nil {
			if err.Error() == "no rows in result set" || err.Error() == "InCorrect Email or Password" {
				response.BadRequest(w, "InCorrect Email or Password", nil)
				return
			}
			app.InfoLog.Println(err.Error())
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		verification, ok := resp["Verification"].(string)
		if ok {
			response.OK(w, verification, nil)
		} else {
			response.OK(w, "User signed in successfully", resp)
		}
	}
}

func Verify(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserUnVerified
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

func resendCode(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserSignIn
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}

		err := UpdateVerificationCode(config.VerifyCache, user)
		if err != nil {
			if err.Error() == "invalid email" {
				response.BadRequest(w, "Invalid Email", err)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}
		app.InfoLog.Println("Verification code sent successfully", user.Email)
		response.OK(w, "Verification code sent successfully", nil)
	}
}

func UpdatePassword(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var UserPassword UpdatePasswordModel
		if err := json.NewDecoder(r.Body).Decode(&UserPassword); err != nil {
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}
		err := UpdateUserPassword(r.Context(), config.DB, &UserPassword)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}
		response.OK(w, "Password updated successfully", nil)
	}
}

func UpdateUser(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		userOid := urlVariables["id"]

		if userOid == "" {
			response.BadRequest(w, "User Id is required", nil)
			return
		}

		var requestData UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			response.BadRequest(w, "Invalid Input Data", err)
			return
		}
		defer r.Body.Close()

		transaction, err := config.DB.Begin(r.Context())
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		fieldsToUpdate, newValues, err := utils.GetNonZeroFieldsFromStruct(&requestData)
		if err != nil {
			response.BadRequest(w, "Invalid Input Data", err)
			return
		}

		query, err := BuildUserUpdateQuery(userOid, fieldsToUpdate)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Internal Server Error", err)
			return
		}

		err = UpdateUserData(r.Context(), transaction, query, newValues)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Internal Server Error", err)
			return
		}

		Data := make(map[string]interface{})

		for idx := 0; idx < len(fieldsToUpdate); idx++ {
			Data[fieldsToUpdate[idx]] = newValues[idx]
		}

		if err := transaction.Commit(r.Context()); err != nil {
			app.ErrorLog.Println(err.Error())
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response.OK(w, "User's data updated successfully", Data)
	}
}

func ForgetPassword(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}
		err := ForgetPasswordService(r.Context(), config.DB, config.VerifyCache, user.Email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			if err.Error() == "User does not exist" {
				response.BadRequest(w, err.Error(), err)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}

		response.OK(w, "Verifacation Code Sent", nil)
	}
}

func ForgetPasswordVerify(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body ResetPasswordForm
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.BadRequest(w, "Invalid JSON body", err)
			return
		}

		err := ForgetPasswordVerifyService(r.Context(), config.DB, config.VerifyCache, &body)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			if err.Error() == "Wrong verification code" {
				response.BadRequest(w, err.Error(), err)
				return
			}
			response.InternalServerError(w, "Server Error, please try again later.", err)
			return
		}
		response.OK(w, "Password updated successfully", nil)
	}
}
