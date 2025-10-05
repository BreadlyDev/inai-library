package user

import (
	"context"
	"fmt"
	"log/slog"
	"new-version/internal/config"
	userDto "new-version/internal/contract/user"
	userRepo "new-version/internal/repository/user"
	"new-version/internal/service/auth"
	userVal "new-version/internal/validator/user"
)

type Service interface {
	Login(ctx context.Context, userReq userDto.Request) (string, error)
	Register(ctx context.Context, userReq userDto.Request) error
}

type DefaultService struct {
	log  *slog.Logger
	repo userRepo.Repository
	auth auth.Service
	cfg  *config.Security
}

func New(
	log *slog.Logger,
	repo userRepo.Repository,
	auth auth.Service,
	cfg *config.Security,
) *DefaultService {
	return &DefaultService{
		log:  log,
		repo: repo,
		auth: auth,
		cfg:  cfg,
	}
}

func (u *DefaultService) Register(ctx context.Context, userReq userDto.Request) error {
	const op = "service.user.Register"

	if !userVal.RightEmailFormat(userReq.Email) {
		return fmt.Errorf("%s", userVal.WrongEmailFormat(userReq.Email))
	}

	if res := userVal.ValidatePassword(userReq.Password, u.cfg.PasswordMinLen); res != "" {
		return fmt.Errorf("%s", res)
	}

	pass, err := u.auth.HashPassword(userReq.Password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userReq.Password = pass

	err = u.repo.Create(ctx, userReq)
	if err != nil {
		u.log.Error("%s: %w", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *DefaultService) Login(ctx context.Context, userReq userDto.Request) (string, error) {
	const op = "service.user.Login"

	pass, err := u.repo.GetPasswordByEmail(ctx, userReq.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	valid, err := u.auth.ComparePassword(pass, userReq.Password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !valid {
		return "", fmt.Errorf("wrong password")
	}

	userInfo, err := u.repo.GetInfoByEmail(ctx, userReq.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.auth.GenerateJwtToken(userDto.Model{
		Email:       userInfo.Email,
		Password:    pass,
		AccessLevel: userInfo.AccessLevel,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
