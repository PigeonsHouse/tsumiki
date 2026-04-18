package middleware

import (
	"context"
	"fmt"
	"net/http"
	"tsumiki/auth"
	"tsumiki/env"
	"tsumiki/helper"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value("user_id").(int)
	return id, ok
}

var RequireAuth = authMiddleware(true)
var OptionalAuth = authMiddleware(false)

func authMiddleware(optional bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				if optional {
					next.ServeHTTP(w, r)
				} else {
					helper.ResponseUnauthorized(w, "アクセストークンが見つかりません")
				}
				return
			}
			token, err := jwt.ParseWithClaims(cookie.Value, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return env.JwtSecret, nil
			})
			if err != nil {
				if optional {
					next.ServeHTTP(w, r)
				} else {
					helper.ResponseUnauthorized(w, "アクセストークンが無効です")
				}
				return
			}
			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok || !token.Valid {
				if optional {
					next.ServeHTTP(w, r)
				} else {
					helper.ResponseUnauthorized(w, "アクセストークンが無効です")
				}
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
