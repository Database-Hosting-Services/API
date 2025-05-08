package accounts

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// signUp godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password, sends verification code to email
// @Tags accounts
// @Accept json
// @Produce json
// @Param user body SignUpUser true "User registration information"
// @Success 201 {object} CreatedResponse "User signed up successfully, check your email for verification"
// @Failure 400 {object} ErrorResponse "Invalid input data or user already exists"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/sign-up [post]
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

// SignIn godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags accounts
// @Accept json
// @Produce json
// @Param user body UserCredentials true "User login credentials"
// @Success 200 {object} LoginResponse "User signed in successfully with JWT token and user details"
// @Success 302 {object} RedirectResponse "User redirected to verification page"
// @Failure 400 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/sign-in [post]
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
			response.Redirect(w, verification, resp)
		} else {
			response.OK(w, "User signed in successfully", resp)
		}
	}
}

// Verify godoc
// @Summary Verify user account
// @Description Verify a user account with verification code sent to email
// @Tags accounts
// @Accept json
// @Produce json
// @Param verification body VerificationRequest true "User verification information with code"
// @Success 201 {object} VerificationSuccessResponse "User verified successfully with JWT token"
// @Failure 400 {object} ErrorResponse "Invalid verification code"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/verify [post]
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
		response.Created(w, "User verified successfully", data)
	}
}

// resendCode godoc
// @Summary Resend verification code
// @Description Resend verification code to user email
// @Tags accounts
// @Accept json
// @Produce json
// @Param user body EmailRequest true "User email information"
// @Success 200 {object} SuccessMessageResponse "Verification code sent successfully"
// @Failure 400 {object} ErrorResponse "Invalid email"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/resend-code [post]
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

// UpdatePassword godoc
// @Summary Update user password
// @Description Update the password for an authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param password body PasswordUpdateRequest true "Password update information"
// @Security BearerAuth
// @Success 200 {object} SuccessMessageResponse "Password updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /users/update-password [put]
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

// UpdateUser godoc
// @Summary Update user information
// @Description Update user profile information such as username and image
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body ProfileUpdateRequest true "User information to update"
// @Security BearerAuth
// @Success 200 {object} ProfileUpdateResponse "User's data updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /users/{id} [put]
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

// ForgetPassword godoc
// @Summary Initiate password reset
// @Description Send a verification code to reset password
// @Tags accounts
// @Accept json
// @Produce json
// @Param user body EmailRequest true "User email information"
// @Success 200 {object} SuccessMessageResponse "Verification code sent"
// @Failure 400 {object} ErrorResponse "User does not exist"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/forgot-password [post]
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

		response.OK(w, "Verification Code Sent", nil)
	}
}

// ForgetPasswordVerify godoc
// @Summary Verify and reset password
// @Description Verify code and reset user password
// @Tags accounts
// @Accept json
// @Produce json
// @Param reset body PasswordResetRequest true "Password reset information with verification code"
// @Success 200 {object} SuccessMessageResponse "Password reset successfully"
// @Failure 400 {object} ErrorResponse "Invalid code or password"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /user/forget-password/verify [post]
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
