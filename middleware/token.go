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

		fields, err := token.GetData(authToken, "id", "username")
		if err != nil {
			response.UnAuthorized(w, "Authorization failed", err)
			return
		}

		if len(fields) >= 2 {
			ctx := utils.AddToContext(r.Context(), map[string]interface{}{
				"user-oid":  fields[0],
				"username": fields[1],
			})
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
