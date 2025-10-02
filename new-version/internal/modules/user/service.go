package user

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"new-version/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	HashPass(pass string) string
	ComparePassword(hashPass string, pass string) (bool, error)
	GenJwtToken(userIn UserLogin) (string, error)
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

func (j *JwtAuthService) GenJwtToken(userIn UserLogin) (string, error) {
	const op = "modules.user.service.GenJwtToken"

	claims := jwt.MapClaims{
		"sub":          userIn.Email,
		"access_level": userIn.AccessLevel,
		"exp":          time.Now().Add(j.cfg.AccessTokenExpire).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.cfg.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return signedToken, nil
}

func ValidateJwt(cfg *config.Security, signedToken string) (jwt.Claims, error) {
	const op = "modules.user.service.ValidateJwt"

	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: unexpected signing method", op)
		}

		return cfg.JwtSecret, nil
	})

	if err != nil {
		log.Println("error checking: ", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("%s: %w", op, ok)
	}

	return claims, nil
}

type UserService interface {
	Login(ctx context.Context) (string, error)
	Logout(ctx context.Context) error
	Register(ctx context.Context, userIn UserCreate) (uuid.UUID, error)
}

type SqliteUserService struct {
	log  *slog.Logger
	repo UserRepo
	auth AuthService
	cfg  *config.Security
}

func NewUserService(log *slog.Logger, repo UserRepo, auth AuthService, cfg *config.Security) *SqliteUserService {
	return &SqliteUserService{
		log:  log,
		repo: repo,
		auth: auth,
	}
}

func (u *SqliteUserService) Register(ctx context.Context, userIn UserCreate) (uuid.UUID, error) {
	const op = "modules.user.service.Register"

	if !ValidateEmailFormat(userIn.Email) {
		return uuid.UUID{}, fmt.Errorf("%s", WrongEmailFormat(userIn.Email))
	}

	if res := ValidatePassword(userIn.Pass, u.cfg.PasswordMinLen); res != "" {
		return uuid.UUID{}, fmt.Errorf("%s", res)
	}

	userIn.Pass = u.auth.HashPass(userIn.Pass)

	id, err := u.repo.Create(ctx, userIn)
	if err != nil {
		u.log.Error("%s: %w", op, err)
		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (u *SqliteUserService) Login(ctx context.Context, userIn UserLogin) (string, error) {
	const op = "modules.user.service.Login"

	user, err := u.repo.GetByEmail(ctx, userIn.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	valid, err := u.auth.ComparePassword(user.Pass, userIn.Pass)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !valid {
		return "", fmt.Errorf("wrong password")
	}

	token, err := u.auth.GenJwtToken(userIn)

	return token, nil
}

func (u *SqliteUserService) Logout(ctx context.Context) error {
	const op = "modules.user.service.Logout"

	// Add logout logic

	return nil
}
