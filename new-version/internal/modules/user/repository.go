package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (User, error)
	GetInfoById(ctx context.Context, id uuid.UUID) (UserInfo, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetInfoByEmail(ctx context.Context, email string) (UserInfo, error)
	GetPassByEmail(ctx context.Context, email string) (string, error)
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

// Not used
func (u *SqliteUserRepo) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	const op = "modules.user.repository.GetById"

	var user User

	row := u.db.QueryRowContext(ctx, `SELECT id, email FROM users WHERE id = $1`, id)
	if err := row.Scan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("%s: user with this id does not exist", op)
		}

		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Not used
func (u *SqliteUserRepo) GetInfoById(ctx context.Context, id uuid.UUID) (UserInfo, error) {
	const op = "modules.user.repository.GetInfoById"

	var user UserInfo

	row := u.db.QueryRowContext(ctx,
		`SELECT id, email, joined_at, access_level FROM users WHERE id = $1`, id)
	if err := row.Scan(&user.Id, &user.Email, &user.JoinedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserInfo{}, fmt.Errorf("%s: user with this id does not exist", op)
		}

		return UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetByEmail(ctx context.Context, email string) (User, error) {
	const op = "modules.user.repository.GetByEmail"

	var user User

	row := u.db.QueryRowContext(ctx, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.Id, &user.Email, &user.Pass); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("%s: user with this email does not exist", op)
		}

		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetInfoByEmail(ctx context.Context, email string) (UserInfo, error) {
	const op = "modules.user.repository.GetInfoByEmail"

	var user UserInfo

	row := u.db.QueryRowContext(ctx,
		`SELECT id, email, joined_at, access_level, pass_hash FROM users WHERE email = $1`, email)
	if err := row.Scan(
		&user.Id, &user.Email, &user.JoinedAt, &user.AccessLevel, &user.Pass,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserInfo{}, fmt.Errorf("%s: user with this email does not exist", op)
		}

		return UserInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SqliteUserRepo) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
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

func (u *SqliteUserRepo) Create(ctx context.Context, userIn UserCreate) error {
	const op = "modules.user.repository.Create"

	_, err := u.db.ExecContext(ctx,
		`INSERT INTO users(id, email, pass_hash) VALUES($1, $2, $3)`, uuid.New(), userIn.Email, userIn.Pass)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return fmt.Errorf("%s: user with this email '%s' already exists", op, userIn.Email)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
