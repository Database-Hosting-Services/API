package accounts

import (
	"time"
)

// This file contains enhanced model definitions for Swagger documentation

// Request Models

// User's data response (for 'users/me' endpoint)
type UserData struct {
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	Oid       string    `json:"oid"`
	Username  string    `json:"username"`
}

type UserDataResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

// SignUpUser represents the user registration request data
type SignUpUser struct {
	Username string `json:"username" example:"ragnar" binding:"required"`
	Email    string `json:"email" example:"ragnar@email.com" binding:"required,email"`
	Password string `json:"password" example:"Password123!" binding:"required"`
}

// UserCredentials represents the login request data
type UserCredentials struct {
	Email    string `json:"email" example:"ragnar@email.com" binding:"required,email"`
	Password string `json:"password" example:"Password123!" binding:"required"`
}

// VerificationRequest represents a user verification request
type VerificationRequest struct {
	Email string `json:"email" example:"ragnar@email.com" binding:"required"`
	Code  string `json:"code" example:"123456" binding:"required"`
}

// EmailRequest represents a request with just email
type EmailRequest struct {
	Email string `json:"email" example:"ragnar@email.com" binding:"required,email"`
}

// PasswordUpdateRequest represents the password update request
type PasswordUpdateRequest struct {
	CurrentPassword string `json:"current_password" example:"OldPassword123!" binding:"required"`
	Password        string `json:"password" example:"NewPassword123!" binding:"required"`
	ConfirmPassword string `json:"confirm_password" example:"NewPassword123!" binding:"required"`
}

// PasswordResetRequest represents the password reset request with verification code
type PasswordResetRequest struct {
	Code     string `json:"code" example:"123456" binding:"required"`
	Password string `json:"password" example:"NewPassword123!" binding:"required"`
	Email    string `json:"email" example:"ragnar@email.com" binding:"required,email"`
}

// ProfileUpdateRequest represents the user profile update request
type ProfileUpdateRequest struct {
	Username string `json:"username" example:"new_ragnar"`
	Image    string `json:"image" example:"profile_image_url.jpg"`
}

// Response Models

// CreatedResponse represents a successful creation response
type CreatedResponse struct {
	Message string      `json:"message" example:"User signed up successfully, check your email for verification"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Message string `json:"message" example:"User signed in successfully"`
	Data    struct {
		Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		Username string `json:"username" example:"ragnar"`
		Email    string `json:"email" example:"ragnar@email.com"`
		OID      string `json:"oid" example:"user-id-123"`
		Image    string `json:"image" example:"profile_image.jpg"`
	} `json:"data"`
}

// VerificationSuccessResponse represents a successful user verification response
type VerificationSuccessResponse struct {
	Message string `json:"message" example:"User verified successfully"`
	Data    struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
}

// SuccessMessageResponse represents a simple success message response
type SuccessMessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ProfileUpdateResponse represents a user profile update response
type ProfileUpdateResponse struct {
	Message string `json:"message" example:"User's data updated successfully"`
	Data    struct {
		Username string `json:"username,omitempty" example:"new_ragnar"`
		Image    string `json:"image,omitempty" example:"profile_image_url.jpg"`
	} `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

// ErrorResponse represents an error response
type ErrorResponse400EmailNotFound struct {
	Error string `json:"error" example:"email not found please sign up first"`
}

// RedirectResponse represents a redirect response
type RedirectResponse struct {
	Message string `json:"message" example:"verification code sent to your email"`
}

type ErrorNotAuthorized struct {
	Status  int    `json:"status" exmaple:"401"`
	Message string `json:"message" example:"Authorization failed"`
	Error   string `json:"error" exmaple:"JWT token is empty"`
}
