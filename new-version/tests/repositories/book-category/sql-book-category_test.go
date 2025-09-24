package bookcategory_test

import (
	"context"
	"new-version/internal/storage/repositories/bookcategory"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestBookCategoryRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer db.Close()

	repo := bookcategory.NewBookCatRepo(db)

	args := "fantasy"
	mock.ExpectExec(`INSERT INTO book_categories\(title\) VALUES \(\$1\)`).
		WithArgs(args).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	err = repo.Create(ctx, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookCategoryRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := bookcategory.NewBookCatRepo(db)

	id := 1
	title := "fantasy"
	tn := time.Now()
	rows := mock.NewRows([]string{"id", "title", "created_time"}).AddRow(1, title, tn)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM book_categories WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	ctx := context.Background()
	bc, err := repo.GetById(ctx, id)

	require.NoError(t, err)

	require.Equal(t, id, bc.Id)
	require.Equal(t, title, bc.Title)
	require.Equal(t, tn, bc.CreatedTime)

	require.NoError(t, mock.ExpectationsWereMet())
}
