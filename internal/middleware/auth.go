package middleware

import (
	"context"
	"net/http"
	"packhaus/internal/utils"
	"strings"
)

type contextKey string

const ContextKeyUserID contextKey = "userID"

func AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := utils.ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
