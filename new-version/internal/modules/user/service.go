package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	HashPass(pass string) string
}

type JwtAuthService struct {
	log *slog.Logger
}

func NewJwtAuthService() *JwtAuthService {
	return &JwtAuthService{}
}

func (j *JwtAuthService) HashPass(pass string) (string, error) {
	const op = "modules.user.service.HashPass"

	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(h), nil
}

func (j *JwtAuthService) IsPassCorrect(hashPass []byte, pass string) (bool, error) {
	const op = "modules.user.service.IsPassCorrect"

	if err := bcrypt.CompareHashAndPassword(hashPass, []byte(pass)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("%s: wrong password", op)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (j *JwtAuthService) GenJwtToken(hashPass []byte, pass string) (bool, error) {
	const op = "modules.user.service.GenJwtToken"

	if err := bcrypt.CompareHashAndPassword(hashPass, []byte(pass)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("%s: wrong password", op)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
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
}

func NewUserService(log *slog.Logger, repo UserRepo, auth AuthService) *SqliteUserService {
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

	if res := FullPasswordValidation(userIn.Pass); res != "" {
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

func Login() string {
	// TODO
	return ""
}
