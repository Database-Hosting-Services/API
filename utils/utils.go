package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func GenerateOID() string {
	return uuid.New().String() // 36 character
}

func GenerateVerficationCode() string {
	return rand.Text()[:6]
}

func HashedPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

// returns the authToken in the Authorization header
func ExtractToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return ""
	}

	// Split into exactly 2 parts: [scheme] [token]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	// Validate Bearer scheme
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func AddToContext(ctx context.Context, data map[string]interface{}) context.Context {
	for k, v := range data {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// Querier defines an interface for executing a single-row query.
// Both *pgxpool.Pool and pgx.Tx implement this interface through the QueryRow method.
type Querier interface {
	// Querier should implement both Query and Query row for pgxscan package
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// if you have request struct like UpdateUser or UpdatePassword
// after decoding the request into the struct
// you can fetch the non zero fields from the struct via this function
// you can see the UpdateUser handler in accounts.handlers.go file
// if err != nil, it means that there is no non zero fields
// IMPORTANT NOTE : YOU SHOULD PASS THE DATA STRCUT AS A POINTER TO THIS FUNCTION
// FOR EX : GetNonZeroFieldsFromStruct(&data)
func GetNonZeroFieldsFromStruct(data interface{}) ([]string, []interface{}, error) {
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	updatedFields := []string{}
	newValues := []interface{}{}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// Check if the field has a non-zero value (works with all premitive datatypes)
		if !field.IsZero() {
			updatedFields = append(updatedFields, strings.ToLower(fieldName))
			newValues = append(newValues, field.Interface())
		}
	}

	if len(updatedFields) == 0 {
		return nil, nil, errors.New("no non-zero fields found")
	}

	return updatedFields, newValues, nil
}

// Fucntion to replce the whitespaces in the string with the underscore
func ReplaceWhiteSpacesWithUnderscore(str string) string {
	replaced := strings.ReplaceAll(str, " ", "_")
	return replaced
}
