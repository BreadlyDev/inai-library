package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetInfoById(ctx context.Context, id uuid.UUID) (UserInfo, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetInfoByEmail(ctx context.Context, email string) (UserInfo, error)
	Create(ctx context.Context, userIn UserCreate) (uuid.UUID, error)
}

type SqliteUserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *SqliteUserRepo {
	return &SqliteUserRepo{
		db: db,
	}
}

func (u *SqliteUserRepo) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	const op = "modules.user.repository.GetById"

	var user User

	row := u.db.QueryRowContext(ctx, `SELECT id, email FROM users WHERE id = $1`, id)
	if err := row.Scan(&user); err != nil {
		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetInfoById(ctx context.Context, id uuid.UUID) (UserInfo, error) {
	const op = "modules.user.repository.GetInfoById"

	var user UserInfo

	row := u.db.QueryRowContext(ctx, `SELECT id, email, joined_at, access_level FROM users WHERE id = $1`, id)
	if err := row.Scan(&user); err != nil {
		return UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	const op = "modules.user.repository.GetByEmail"

	var user User

	row := u.db.QueryRowContext(ctx, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.Id, &user.Email, &user.Pass); err != nil {
		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetInfoByEmail(ctx context.Context, email string) (UserInfo, error) {
	const op = "modules.user.repository.GetInfoByEmail"

	var user UserInfo

	row := u.db.QueryRowContext(ctx,
		`SELECT id, email, joined_at, access_level, pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.Id, &user.Email, &user.JoinedIn, &user.AccessLevel, &user.Pass); err != nil {
		return UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) Create(ctx context.Context, userIn UserCreate) (uuid.UUID, error) {
	const op = "modules.user.repository.Create"

	var id uuid.UUID

	row := u.db.QueryRowContext(ctx,
		`INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`, userIn.Email, userIn.Pass)
	if err := row.Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
