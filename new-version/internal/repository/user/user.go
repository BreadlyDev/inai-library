package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"new-version/internal/contract/user"
)

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (user.Response, error)
	GetInfoById(ctx context.Context, id uuid.UUID) (user.InfoResponse, error)
	GetByEmail(ctx context.Context, email string) (user.Response, error)
	GetInfoByEmail(ctx context.Context, email string) (user.InfoResponse, error)
	GetPasswordByEmail(ctx context.Context, email string) (string, error)
	Create(ctx context.Context, userReq user.Request) error
}

type DefaultRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *DefaultRepository {
	return &DefaultRepository{
		db: db,
	}
}

// Not used
func (u *DefaultRepository) GetById(ctx context.Context, id uuid.UUID) (user.Response, error) {
	const op = "modules.user.repository.GetById"

	var resp user.Response

	row := u.db.QueryRowContext(ctx, `SELECT id, email FROM users WHERE id = $1`, id)
	if err := row.Scan(&resp.Id, &resp.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.Response{}, fmt.Errorf("%s: user with this id does not exist", op)
		}

		return user.Response{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

// Not used
func (u *DefaultRepository) GetInfoById(ctx context.Context, id uuid.UUID) (user.InfoResponse, error) {
	const op = "modules.user.repository.GetInfoById"

	var resp user.InfoResponse

	row := u.db.QueryRowContext(ctx,
		`SELECT id, email, joined_at, access_level FROM users WHERE id = $1`, id)
	if err := row.Scan(&resp.Id, &resp.Email, &resp.JoinedAt, &resp.AccessLevel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.InfoResponse{}, fmt.Errorf("%s: user with this id does not exist", op)
		}

		return user.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (u *DefaultRepository) GetByEmail(ctx context.Context, email string) (user.Response, error) {
	const op = "modules.user.repository.GetByEmail"

	var resp user.Response

	row := u.db.QueryRowContext(ctx, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&resp.Id, &resp.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.Response{}, fmt.Errorf("%s: user with this email does not exist", op)
		}

		return user.Response{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (u *DefaultRepository) GetInfoByEmail(ctx context.Context, email string) (user.InfoResponse, error) {
	const op = "modules.user.repository.GetInfoByEmail"

	var resp user.InfoResponse

	row := u.db.QueryRowContext(ctx,
		`SELECT id, email, joined_at, access_level FROM users WHERE email = $1`, email)
	if err := row.Scan(
		&resp.Id, &resp.Email, &resp.JoinedAt, &resp.AccessLevel,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.InfoResponse{}, fmt.Errorf("%s: user with this email does not exist", op)
		}

		return user.InfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (u *DefaultRepository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	const op = "modules.user.repository.GetPasswordByEmail"

	var pass string

	row := u.db.QueryRowContext(ctx,
		`SELECT pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&pass); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: user with this email does not exist", op)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return pass, nil
}

func (u *DefaultRepository) Create(ctx context.Context, userReq user.Request) error {
	const op = "modules.user.repository.Create"

	_, err := u.db.ExecContext(ctx,
		`INSERT INTO users(id, email, pass_hash) VALUES($1, $2, $3)`, uuid.New(), userReq.Email, userReq.Password)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return fmt.Errorf("%s: user with this email '%s' already exists", op, userReq.Email)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
