package middleware

import (
	"DBHS/response"
	"DBHS/utils"
	"DBHS/utils/token"
	"fmt"
	"net/http"
)

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

		fields, err := token.GetData(authToken, "id", "oid", "username")
		if err != nil {
			response.UnAuthorized(w, "Authorization failed", err)
			return
		}

		if len(fields) >= 3 {
			idFloat, ok := fields[0].(float64)
			if !ok {
				response.UnAuthorized(w, "Invalid user ID type", fmt.Errorf("expected numeric user ID, got %T", fields[0]))
				return
			}

			userID := int(idFloat)
			ctx := utils.AddToContext(r.Context(), map[string]interface{}{
				"user-id":   userID,
				"user-oid":  fields[1],
				"user-name": fields[2],
			})
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
