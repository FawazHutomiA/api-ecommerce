package middleware

import (
	"context"
	"encoding/json"
	"example/pkg/exception"
	jwtValidate "example/pkg/jwt"
	"example/pkg/response"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")

		token := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))

		tokenValidate, err := jwtValidate.ValidateToken(token)
		if err != nil {
			resp := response.Error(response.StatusForbiddend, "Unauthorized", exception.ErrUnauthorized)
			resp.JSON(w)
			return
		}

		// Extract claims
		var claims Claims
		claimsBytes, err := json.Marshal(tokenValidate.Claims)
		if err != nil {
			resp := response.Error(response.StatusForbiddend, "Unauthorized", exception.ErrUnauthorized)
			resp.JSON(w)
			return
		}
		json.Unmarshal(claimsBytes, &claims)

		userID := claims.Data.UserID
		role := claims.Data.Role

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
