package bookcategory

import (
	"context"
	"log/slog"
	"net/http"
	srv "new-version/internal/http-server"
	"new-version/utils/json"
	"strconv"
	"strings"
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

func (b *SqliteBookCatHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.CreateCategory"

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	var req BookCatRequest
	var resp srv.Response

	if err := json.ReadRequestBody(r, &req); err != nil {
		resp = srv.NewErrResponse(err.Error(), http.StatusBadRequest)

		json.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	if err := b.repo.Create(ctx, req.Title); err != nil {
		resp = srv.NewErrResponse(err.Error(), http.StatusInternalServerError)

		json.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("created book category", nil, http.StatusCreated)
	json.WriteResponseBody(w, resp, http.StatusCreated)
}

func (b *SqliteBookCatHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "domain.bookcategory.handler.GetCategoryById"

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	var resp srv.Response

	pathParams := strings.Split(r.URL.Path, "/")

	idStr := pathParams[len(pathParams)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp = srv.NewErrResponse("'id' must be int", http.StatusBadRequest)

		json.WriteResponseBody(w, resp, http.StatusBadRequest)
		return
	}

	bc, err := b.repo.GetById(ctx, id)
	if err != nil {
		resp = srv.NewErrResponse(err.Error(), http.StatusInternalServerError)

		json.WriteResponseBody(w, resp, http.StatusInternalServerError)
		return
	}

	resp = srv.NewResponse("fetched book", bc, http.StatusOK)
	json.WriteResponseBody(w, resp, http.StatusOK)
}
