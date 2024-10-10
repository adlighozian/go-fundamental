package https

import (
	"context"
	"go-axiata/model"
	"go-axiata/pkg/helper"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jwtKey = []byte("my_secret_key")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.RespondJSON(w, http.StatusUnauthorized, false, "Access denied. You do not have authorization", nil, nil)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &model.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				helper.RespondJSON(w, http.StatusUnauthorized, false, "Invalid token", nil, nil)
				return
			}
			helper.RespondJSON(w, http.StatusBadRequest, false, "Bad request", nil, nil)
			return
		}
		if !token.Valid {
			helper.RespondJSON(w, http.StatusUnauthorized, false, "Invalid token", nil, nil)
			return
		}

		if requiredRole != "" {
			if claims.Role != requiredRole {
				helper.RespondJSON(w, http.StatusForbidden, false, "Forbidden", nil, nil)
				return
			}
		}

		ctx := context.WithValue(r.Context(), 1221, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
