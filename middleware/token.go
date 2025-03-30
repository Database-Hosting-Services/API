package middleware

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"DBHS/utils/token"
	"context"
	"errors"
	"fmt"
	"net/http"
)

func GetUserByOid(ctx context.Context, oid string) (int, error) {
	var id int
	err := config.DB.QueryRow(
		ctx,
		`SELECT id FROM "users" WHERE oid= $1`,
		oid,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := utils.ExtractToken(r)
		if authToken == "" {
			response.UnAuthorized(w, "Authorization failed", fmt.Errorf("JWT token is empty"))
			return
		}

		err := token.IsAuthorized(authToken)
		if err != nil {
			response.UnAuthorized(w, "Authorization failed", err)
			return
		}

		fields, err := token.GetData(authToken, "oid", "username")
		if err != nil {
			response.UnAuthorized(w, "Authorization failed", err)
			return
		}

		if len(fields) >= 2 {
			userID, err := GetUserByOid(r.Context(), fields[0].(string))
			if err != nil {
				response.UnAuthorized(w, "Authorization failed", errors.New("No user found for this token"))
				return
			}
			ctx := utils.AddToContext(r.Context(), map[string]interface{}{
				"user-id":   userID,
				"user-oid":  fields[0],
				"user-name": fields[1],
			})
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
