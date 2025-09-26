package bookcategory

import (
	"context"
	"log/slog"
	"net/http"
	help "new-version/pkg/http-helpers"
	"new-version/pkg/json"
	"time"
)

type BookCatHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryById(w http.ResponseWriter, r *http.Request)
	DeleteCategoryById(w http.ResponseWriter, r *http.Request)
	UpdateCategoryById(w http.ResponseWriter, r *http.Request)
	ListCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryByTitle(w http.ResponseWriter, r *http.Request)
}

type SqliteBookCatHandler struct {
	log  *slog.Logger
	repo BookCatRepo
}

func NewBookCatHandler(log *slog.Logger, repo BookCatRepo) *SqliteBookCatHandler {
	return &SqliteBookCatHandler{
		log:  log,
		repo: repo,
	}
}

func (b *SqliteBookCatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /book-category/", b.CreateCategory)
	mux.HandleFunc("GET /book-category/{id}", b.GetCategoryById)
	mux.HandleFunc("PATCH /book-category/{id}", b.UpdateCategoryById)
	mux.HandleFunc("DELETE /book-category/{id}", b.DeleteCategoryById)
	mux.HandleFunc("GET /book-category/title", b.GetCategoryByTitle)
	mux.HandleFunc("GET /book-category/", b.ListCategories)
}

func (b *SqliteBookCatHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.CreateCategory"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	var req BookCatRequest

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.repo.Create(ctx, req.Title); err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "created book category", nil, http.StatusOK)
}

func (b *SqliteBookCatHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.GetCategoryById"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	id, err := help.ParseIdFromPath(r)

	if err != nil {
		// id must be int
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	bc, err := b.repo.GetById(ctx, id)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "fetched book category", bc, http.StatusOK)
}

func (b *SqliteBookCatHandler) GetCategoryByTitle(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.GetCategoryByTitle"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	title := r.URL.Query().Get("title")

	if title == "" {
		// id must be int
		json.WriteError(w, "title must not be empty", http.StatusBadRequest)
		return
	}

	bc, err := b.repo.GetByTitle(ctx, title)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "fetched book category", bc, http.StatusOK)
}

func (b *SqliteBookCatHandler) UpdateCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.UpdateCategoryById"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	id, err := help.ParseIdFromPath(r)

	var req BookCatRequest

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		// id must be int
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = b.repo.UpdateById(ctx, req.Title, id)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "updated book category", nil, http.StatusOK)
}

func (b *SqliteBookCatHandler) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.DeleteCategoryById"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	id, err := help.ParseIdFromPath(r)

	if err != nil {
		// id must be int
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = b.repo.DeleteById(ctx, id)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "deleted book category", nil, http.StatusOK)
}

func (b *SqliteBookCatHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.ListCategories"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	bcList, err := b.repo.GetList(ctx)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "fetched book categories", bcList, http.StatusOK)
}
