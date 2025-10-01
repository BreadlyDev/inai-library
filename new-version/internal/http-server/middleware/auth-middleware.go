package middleware

import (
	"net/http"
	"strings"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			authParts := strings.SplitN(auth, " ", 2)
			if len(authParts) != 2 || authParts[0] != "Bearer" {
				http.Error(w, "invalid header authorization format", http.StatusUnauthorized)
				return
			}

			if authParts[0] == "" {
				http.Error(w, "empty token", http.StatusUnauthorized)
				return
			}

			// TODO: finish the function

			next.ServeHTTP(w, r)
		})
	}
}
