package bookcategory_test

import (
	"context"
	bookCatDto "new-version/internal/contract/bookcategory"
	bookCatRepo "new-version/internal/repository/bookcategory"
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
	repo := bookCatRepo.New(db)

	req := bookCatDto.Request{
		Title: "fantasy",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO book_categories(title) VALUES ($1) RETURNING id`)).
		WithArgs(req.Title).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ctx := context.Background()

	_, err = repo.Create(ctx, req)
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

	repo := bookCatRepo.New(db)

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

func TestBookCategoryRepository_GetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := bookCatRepo.New(db)

	id := 1
	title := "fantasy"
	tn := time.Now()
	rows := mock.NewRows([]string{"id", "title", "created_time"}).AddRow(1, title, tn)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM book_categories WHERE title = $1`)).
		WithArgs(title).
		WillReturnRows(rows)

	ctx := context.Background()
	bc, err := repo.GetByTitle(ctx, title)

	require.NoError(t, err)

	require.Equal(t, id, bc.Id)
	require.Equal(t, title, bc.Title)
	require.Equal(t, tn, bc.CreatedTime)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestBookCategoryRepository_DeleteById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	id := 1
	title := "fantasy"
	tn := time.Now()
	mock.NewRows([]string{"id", "title", "created_time"}).AddRow(id, title, tn)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM book_categories WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := bookCatRepo.New(db)
	ctx := context.Background()

	err = repo.DeleteById(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookCategoryRepository_UpdateById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	id := 1
	title := "fantasy"
	tn := time.Now()

	mock.NewRows([]string{"id", "title", "created_time"}).AddRow(id, title, tn)

	req := bookCatDto.Request{
		Title: "mistery",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE book_categories SET title = $1 WHERE id = $2`)).
		WithArgs(req.Title, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := bookCatRepo.New(db)

	ctx := context.Background()

	err = repo.UpdateById(ctx, req, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestBookCategoryRepository_GetList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	tn := time.Now()
	bookCatList := []bookCatDto.Response{
		{Id: 1, Title: "fantasy", CreatedTime: tn},
		{Id: 2, Title: "mistery", CreatedTime: tn.Add(10 * time.Second)},
		{Id: 3, Title: "fiction", CreatedTime: tn.Add(20 * time.Second)},
		{Id: 4, Title: "science", CreatedTime: tn.Add(30 * time.Second)},
		{Id: 5, Title: "romance", CreatedTime: tn.Add(40 * time.Second)},
	}
	rows := mock.NewRows([]string{"id", "title", "created_time"})

	for _, b := range bookCatList {
		rows.AddRow(b.Id, b.Title, b.CreatedTime)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM book_categories`)).WillReturnRows(rows)

	repo := bookCatRepo.New(db)

	ctx := context.Background()
	bookCats, err := repo.GetList(ctx)

	require.NoError(t, err)

	for i, bookCat := range bookCats {
		require.Equal(t, bookCat, bookCatList[i])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
