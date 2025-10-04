package user

import (
	"fmt"
	"log/slog"
	"new-version/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	HashPass(pass string) (string, error)
	ComparePassword(hashPass string, pass string) (bool, error)
	GenJwtToken(userIn UserLoginResp) (string, error)
}

type JwtAuthService struct {
	log *slog.Logger
	cfg *config.Security
}

func NewJwtAuthService(log *slog.Logger, cfg *config.Security) *JwtAuthService {
	return &JwtAuthService{log: log, cfg: cfg}
}

func (j *JwtAuthService) HashPass(pass string) (string, error) {
	const op = "modules.user.service.HashPass"

	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(h), nil
}

func (j *JwtAuthService) ComparePassword(hashPass string, pass string) (bool, error) {
	const op = "modules.user.service.ComparePassword"

	if err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("%s: wrong password", op)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (j *JwtAuthService) GenJwtToken(userIn UserLoginResp) (string, error) {
	const op = "modules.user.service.GenJwtToken"

	claims := jwt.MapClaims{
		"sub":          userIn.Email,
		"access_level": userIn.AccessLevel,
		"exp":          time.Now().Add(j.cfg.AccessTokenExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.cfg.JwtSecret))

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedToken, nil
}
