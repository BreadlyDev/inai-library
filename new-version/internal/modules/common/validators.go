package common

import (
	"fmt"
	"new-version/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateFieldNotEmpty(field string) bool {
	if field == "" {
		return false
	}
	return true
}

func ValidateJwt(cfg *config.Security, signedToken string) (jwt.MapClaims, error) {
	const op = "modules.user.service.ValidateJwt"

	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method", op)
		}
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if parsedToken == nil || !parsedToken.Valid {
		return nil, fmt.Errorf("%s: invalid token", op)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %s", op, "invalid token")
	}

	return claims, nil
}
