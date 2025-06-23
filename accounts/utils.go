package accounts

import (
	"DBHS/config"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// EmailSenderFunc defines the signature of the SendMail function for mocking
type EmailSenderFunc func(d interface{}, from, to, code, subject string) error

// Mock sender for tests
var currentEmailSender EmailSenderFunc = nil

func checkPasswordStrength(password string) error {
	/*
		The password should contains uppercase, lowercase , digits and special characters
		special characters : [@#$%^&*!_]
	*/
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

func checkUserExistsInCache(username, email string) (string, error) {
	/*
		if email or username exists in the cache
		this function will return (email, nil) or (username, null)
	*/
	for key, value := range map[string]string{"email": email, "username": username} {
		exists, err := config.VerifyCache.Exists(value)
		if err != nil {
			return "", err
		}
		if exists {
			return key, nil
		}
	}
	return "", nil
}

func CheckPasswordHash(inputPassword, storedHash string) bool {
	byteHash := []byte(storedHash)
	bytePassword := []byte(inputPassword)

	// Compare the password with the hash
	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	return err == nil
}

// SetEmailSender allows testing to set a mock email sender
func SetEmailSender(sender EmailSenderFunc) EmailSenderFunc {
	old := currentEmailSender
	currentEmailSender = sender
	return old
}

func SendMail(d *gomail.Dialer, from, to, code, Subject string) error {
	// Use mock sender if set (for testing)
	if currentEmailSender != nil {
		return currentEmailSender(d, from, to, code, Subject)
	}

	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", Subject)

	// get the directory of the current file using runtime
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	mailTemplatePath := filepath.Join(dir, "templates", "mailTemplate.html")
	data, err := os.ReadFile(mailTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read mail template: %w, path: %s", err, mailTemplatePath)
	}

	body := fmt.Sprintf(string(data), code)
	m.SetBody("text/html", body)

	// Send email using the provided dialer
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
