package auth

import (
	"context"
	"net/http"
	"new-version/internal/config"
	"new-version/internal/validator/user"
	"new-version/pkg/httphelpers"
)

func AuthMiddleware(cfg *config.Security) func(next http.Handler, level httphelpers.AccessLevel) http.Handler {
	return func(next http.Handler, level httphelpers.AccessLevel) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc, err := r.Cookie("access_token")
			if err != nil || tc.Value == "" {
				http.Error(w, "missing or empty token", http.StatusUnauthorized)
				return
			}

			claims, err := user.ValidateJwt(cfg.JwtSecret, tc.Value)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			lvlFloat, ok := claims["access_level"].(float64)
			if !ok {
				http.Error(w, "invalid access level", http.StatusUnauthorized)
				return
			}

			lvl := httphelpers.AccessLevel(int(lvlFloat))

			if lvl < level {
				http.Error(w, "no permission for action", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Auth(handlerLevel httphelpers.AccessLevel) func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
		jwtSecret, ok := ctx.Value("jwt_secret").(string)
		if !ok || jwtSecret == "" {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return false
		}

		tc, err := r.Cookie("access_token")
		if err != nil || tc.Value == "" {
			http.Error(w, "missing or empty token", http.StatusUnauthorized)
			return false
		}

		claims, err := user.ValidateJwt(jwtSecret, tc.Value)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return false
		}

		res, ok := claims["user_level"].(float64)
		if !ok {
			http.Error(w, "invalid access level", http.StatusUnauthorized)
			return false
		}
		userLevel := httphelpers.AccessLevel(res)

		if userLevel < handlerLevel {
			http.Error(w, "no permission for action", http.StatusForbidden)
			return false
		}

		return true
	}
}
