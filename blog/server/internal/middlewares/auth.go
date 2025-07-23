package middlewares

import (
	"main/internal/contexter"
	"main/internal/env"
	"main/pkg/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	verifyTokenHandler := func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")
		if authorization == "" {
			response.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authorization, " ")[1]

		type ParsedClaims struct {
			contexter.TokenUser
			jwt.RegisteredClaims
		}
		var parsedClaims ParsedClaims
		token, err := jwt.ParseWithClaims(bearerToken, &parsedClaims, func(t *jwt.Token) (any, error) {
			return env.JwtSecret, nil
		})
		if err != nil {
			response.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			response.Error(w, "Invalid Token Provided", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = contexter.ContextWithUser(ctx, parsedClaims.TokenUser)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(verifyTokenHandler)
}
