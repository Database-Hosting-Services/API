package utils

import (
	"context"
	"crypto/rand"
	"net/http"
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
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
