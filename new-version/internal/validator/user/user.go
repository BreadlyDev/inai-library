package user

import (
	"fmt"
	"net/mail"
	"new-version/internal/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Messages
func WrongEmailFormat(email string) string {
	return fmt.Sprintf("wrong email format: %s", email)
}

func TooShort(pass string) string {
	return fmt.Sprintf("password is too short: %s", pass)
}

func HasNoNumber(pass string) string {
	return fmt.Sprintf("password has no number: %s", pass)
}

func HasNoCapitalizedLetter(pass string) string {
	return fmt.Sprintf("password has no capitalized letter: %s", pass)
}

func HasNoSpecialSymbol(pass string) string {
	return fmt.Sprintf("password has no special symbol: %s", pass)
}

// Validators
func RightEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsShort(pass string, minLen int) bool {
	return len(pass) < minLen
}

func HasNumber(pass string) bool {
	return strings.ContainsAny(pass, "1234567890")
}

func HasCapitalizedLetter(pass string) bool {
	return strings.ContainsAny(pass, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func HasSpecialSymbol(pass string) bool {
	return strings.ContainsAny(pass, "@#$%&/?.,-_+=)(^;:!")
}

func ValidatePassword(pass string, passMinLen int) string {
	if IsShort(pass, passMinLen) {
		return fmt.Sprintf("%s", TooShort(pass))
	}

	if !HasNumber(pass) {
		return fmt.Sprintf("%s", HasNoNumber(pass))
	}

	if !HasCapitalizedLetter(pass) {
		return fmt.Sprintf("%s", HasNoCapitalizedLetter(pass))
	}

	if !HasSpecialSymbol(pass) {
		return fmt.Sprintf("%s", HasNoSpecialSymbol(pass))
	}

	return ""
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
