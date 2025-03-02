package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func GenerateOID() string {
	return uuid.New().String() // 36 character
}

// returns the authToken in the Authorization header
func ExtractToken(r *http.Request) string {
	authheader := r.Header.Get("Authorization")
	header := strings.Split(authheader, " ")
	if len(header) != 2 {
		return ""
	}

	return header[1]
}

func AddToContext(ctx context.Context, data map[string]interface{}) context.Context {
	for k, v := range data {
		ctx = context.WithValue(ctx, k , v)
	}
	return ctx
}