package auth

import "context"

type AuthRepository interface {
	Login(ctx context.Context)
	Logout(ctx context.Context)
	Register(ctx context.Context)
}
