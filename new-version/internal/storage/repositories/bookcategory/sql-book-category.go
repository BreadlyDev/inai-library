package bookcategory

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BookCat struct {
	Id          int
	Title       string
	CreatedTime time.Time
}

type BookCatRepo interface {
	GetById(ctx context.Context, id int) (BookCat, error)
	Create(ctx context.Context, title string) error
}

type SqlBookCatRepo struct {
	db *sql.DB
}

func NewBookCatRepo(db *sql.DB) *SqlBookCatRepo {
	return &SqlBookCatRepo{db: db}
}

func (b *SqlBookCatRepo) GetById(ctx context.Context, id int) (BookCat, error) {
	const op = "storage.repositories.book.GetById"

	row := b.db.QueryRowContext(ctx, `SELECT * FROM book_categories WHERE id = $1`, id)

	var bookCat BookCat

	err := row.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
	if err != nil {
		return BookCat{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookCat, nil
}

func (b *SqlBookCatRepo) Create(ctx context.Context, title string) error {
	const op = "storage.repositories.book.Create"

	_, err := b.db.ExecContext(ctx, `INSERT INTO book_categories(title) VALUES ($1)`, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
