package user

import (
	"context"
	"fmt"
	"log/slog"
	"new-version/internal/config"

	"github.com/google/uuid"
)

type UserService interface {
	Login(ctx context.Context, userIn UserLogin) (string, error)
	Register(ctx context.Context, userIn UserCreate) (uuid.UUID, error)
}

type UserServiceImpl struct {
	log  *slog.Logger
	repo UserRepo
	auth AuthService
	cfg  *config.Security
}

func NewUserService(log *slog.Logger, repo UserRepo, auth AuthService, cfg *config.Security) *UserServiceImpl {
	return &UserServiceImpl{
		log:  log,
		repo: repo,
		auth: auth,
		cfg:  cfg,
	}
}

func (u *UserServiceImpl) Register(ctx context.Context, userIn UserCreate) (uuid.UUID, error) {
	const op = "modules.user.service.Register"

	if !ValidateEmailFormat(userIn.Email) {
		return uuid.UUID{}, fmt.Errorf("%s", WrongEmailFormat(userIn.Email))
	}

	if res := ValidatePassword(userIn.Pass, u.cfg.PasswordMinLen); res != "" {
		return uuid.UUID{}, fmt.Errorf("%s", res)
	}

	pass, err := u.auth.HashPass(userIn.Pass)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	userIn.Pass = pass

	id, err := u.repo.Create(ctx, userIn)
	if err != nil {
		u.log.Error("%s: %w", op, err)
		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (u *UserServiceImpl) Login(ctx context.Context, userIn UserLogin) (string, error) {
	const op = "modules.user.service.Login"

	user, err := u.repo.GetInfoByEmail(ctx, userIn.Email)
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

	token, err := u.auth.GenJwtToken(UserLoginResp{
		Email:       user.Email,
		Pass:        user.Pass,
		AccessLevel: user.AccessLevel,
	})

	return token, nil
}
