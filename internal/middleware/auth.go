package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/theashgen/url-short/internal/auth"
)

const UsernameKey = "username"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access_token")
		// fmt.Println(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// fmt.Print(cookie.Value)	
		if cookie.Value == "" {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), UsernameKey, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
