package bookcategory

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BookCatRepo interface {
	GetById(ctx context.Context, id int) (BookCat, error)
	Create(ctx context.Context, title string) error
	DeleteById(ctx context.Context, id int) error
	UpdateById(ctx context.Context, newTitle string, id int) error
	GetByTitle(ctx context.Context, title string) (BookCat, error)
	GetList(ctx context.Context) ([]BookCat, error)
}

type BookCat struct {
	Id          int
	Title       string
	CreatedTime time.Time
}

type SqliteBookCatRepo struct {
	db *sql.DB
}

func NewBookCatRepo(db *sql.DB) *SqliteBookCatRepo {
	return &SqliteBookCatRepo{db: db}
}

func (b *SqliteBookCatRepo) GetById(ctx context.Context, id int) (BookCat, error) {
	const op = "domain.bookcategory.repository.GetById"

	row := b.db.QueryRowContext(ctx, `SELECT * FROM book_categories WHERE id = $1`, id)

	var bookCat BookCat

	err := row.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
	if err != nil {
		return BookCat{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookCat, nil
}

func (b *SqliteBookCatRepo) GetByTitle(ctx context.Context, title string) (BookCat, error) {
	const op = "domain.bookcategory.repository.GetByTitle"

	row := b.db.QueryRowContext(ctx, `SELECT * FROM book_categories WHERE title = $1`, title)

	var bookCat BookCat

	err := row.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime)
	if err != nil {
		return BookCat{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookCat, nil
}

func (b *SqliteBookCatRepo) Create(ctx context.Context, title string) error {
	const op = "domain.bookcategory.repository.Create"

	_, err := b.db.ExecContext(ctx, `INSERT INTO book_categories(title) VALUES ($1)`, title)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (b *SqliteBookCatRepo) UpdateById(ctx context.Context, newTitle string, id int) error {
	const op = "domain.bookcategory.repository.Update"

	_, err := b.db.ExecContext(ctx, `UPDATE book_categories SET title = $1 WHERE id = $2`, newTitle, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (b *SqliteBookCatRepo) GetList(ctx context.Context) ([]BookCat, error) {
	const op = "domain.bookcategory.repository.GetList"

	rows, err := b.db.QueryContext(ctx, `SELECT * FROM book_categories`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var bookCatList []BookCat

	for rows.Next() {
		var bookCat BookCat
		if err := rows.Scan(&bookCat.Id, &bookCat.Title, &bookCat.CreatedTime); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		bookCatList = append(bookCatList, bookCat)
	}

	return bookCatList, nil
}

func (b *SqliteBookCatRepo) DeleteById(ctx context.Context, id int) error {
	const op = "domain.bookcategory.repository.Delete"

	_, err := b.db.ExecContext(ctx, `DELETE FROM book_categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
