package bookcategory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"new-version/internal/contract/bookcategory"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Repository interface {
	GetById(ctx context.Context, id int) (bookcategory.Response, error)
	Create(ctx context.Context, bookCat bookcategory.Request) (int, error)
	DeleteById(ctx context.Context, id int) error
	UpdateById(ctx context.Context, bookCat bookcategory.Request, id int) error
	GetByTitle(ctx context.Context, title string) (bookcategory.Response, error)
	GetList(ctx context.Context) ([]bookcategory.Response, error)
}

type DefaultRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *DefaultRepository {
	return &DefaultRepository{db: db}
}

func (b *DefaultRepository) GetById(ctx context.Context, id int) (bookcategory.Response, error) {
	const op = "modules.bookcategory.repository.GetById"

	row := b.db.QueryRowContext(ctx, `SELECT * FROM book_categories WHERE id = $1`, id)

	var bookCat bookcategory.Response

	err := row.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bookcategory.Response{}, fmt.Errorf("%s: no book category category with id = %d", op, id)
		}

		return bookcategory.Response{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookCat, nil
}

func (b *DefaultRepository) GetByTitle(ctx context.Context, title string) (bookcategory.Response, error) {
	const op = "modules.bookcategory.repository.GetByTitle"

	row := b.db.QueryRowContext(ctx, `SELECT * FROM book_categories WHERE title = $1`, title)

	var bookCat bookcategory.Response

	err := row.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bookcategory.Response{}, fmt.Errorf("%s: no book category category with title = %s", op, title)
		}

		return bookcategory.Response{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookCat, nil
}

func (b *DefaultRepository) Create(ctx context.Context, bookCat bookcategory.Request) (int, error) {
	const op = "modules.bookcategory.repository.Create"

	var id int

	err := b.db.QueryRowContext(
		ctx,
		`INSERT INTO book_categories(title) VALUES ($1) RETURNING id`,
		bookCat.Title,
	).Scan(&id)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return 0, fmt.Errorf("%s: book category with title '%s' already exists", op, bookCat.Title)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (b *DefaultRepository) UpdateById(ctx context.Context, bookCat bookcategory.Request, id int) error {
	const op = "modules.bookcategory.repository.Update"

	_, err := b.db.ExecContext(ctx, `UPDATE book_categories SET title = $1 WHERE id = $2`, bookCat.Title, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: no book category category with id = %d", op, id)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (b *DefaultRepository) GetList(ctx context.Context) ([]bookcategory.Response, error) {
	const op = "modules.bookcategory.repository.GetList"

	rows, err := b.db.QueryContext(ctx, `SELECT * FROM book_categories`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var bookCatList []bookcategory.Response

	for rows.Next() {
		var bookCat bookcategory.Response
		err := rows.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		bookCatList = append(bookCatList, bookCat)
	}

	return bookCatList, nil
}

func (b *DefaultRepository) DeleteById(ctx context.Context, id int) error {
	const op = "modules.bookcategory.repository.Delete"

	_, err := b.db.ExecContext(ctx, `DELETE FROM book_categories WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: no book category category found with id = %d", op, id)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
