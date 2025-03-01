package accounts

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"regexp"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[@#$%^&*!_]`).MatchString(password)

	if !hasUpperCase {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLowerCase {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character (@#$%^&*!)")
	}

	return nil
}

func checkUserExists(ctx context.Context, db *pgxpool.Pool, username, email string) (string, error) {
	/*
		new users should sign up with unique `username` and `email`
		if one of these fields are not unique this function will return the existingField
	*/
	var existingField string

	query := `SELECT 
                CASE 
                    WHEN username = $1 THEN 'username' 
                    WHEN email = $2 THEN 'email' 
                END 
              FROM "users" 
              WHERE username = $1 OR email = $2 
              LIMIT 1`

	err := db.QueryRow(ctx, query, username, email).Scan(&existingField)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return existingField, nil
}

func CheckPasswordHash(inputPassword, storedHash string) bool {
    byteHash := []byte(storedHash)
    bytePassword := []byte(inputPassword)

    // Compare the password with the hash
    err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
    return err == nil
}