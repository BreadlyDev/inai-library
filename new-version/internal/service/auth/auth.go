package auth

import (
	"fmt"
	"log/slog"
	"new-version/internal/config"
	"new-version/internal/contract/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	HashPassword(pass string) (string, error)
	ComparePassword(hashPass string, pass string) (bool, error)
	GenerateJwtToken(userInfo user.Model) (string, error)
}

type JwtService struct {
	log *slog.Logger
	cfg *config.Security
}

func New(log *slog.Logger, cfg *config.Security) *JwtService {
	return &JwtService{log: log, cfg: cfg}
}

func (j *JwtService) HashPassword(pass string) (string, error) {
	const op = "modules.user.service.HashPass"

	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(h), nil
}

func (j *JwtService) ComparePassword(hashPass string, pass string) (bool, error) {
	const op = "service.auth.ComparePassword"

	if err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("%s: wrong password", op)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (j *JwtService) GenerateJwtToken(userInfo user.Model) (string, error) {
	const op = "service.auth.GenJwtToken"

	claims := jwt.MapClaims{
		"sub":          userInfo.Email,
		"access_level": userInfo.AccessLevel,
		"exp":          time.Now().Add(j.cfg.AccessTokenExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.cfg.JwtSecret))

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedToken, nil
}
