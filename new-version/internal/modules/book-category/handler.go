package bookcategory

import (
	"context"
	"log/slog"
	"net/http"
	"new-version/internal/config"
	mw "new-version/internal/http-server/middleware"
	"new-version/internal/modules/common"
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

type BookCatHandlerImpl struct {
	log  *slog.Logger
	repo BookCatRepo
	cfg  *config.Security
}

func NewBookCatHandler(log *slog.Logger, repo BookCatRepo, cfg *config.Security) *BookCatHandlerImpl {
	return &BookCatHandlerImpl{
		log:  log,
		repo: repo,
		cfg:  cfg,
	}
}

func (b *BookCatHandlerImpl) RegisterRoutes(mux *http.ServeMux, log *slog.Logger) {
	logMw := mw.LoggerMiddleware(log)
	authMw := mw.AuthMiddleware(b.cfg)

	mux.Handle("POST /book-category/", authMw(logMw(http.HandlerFunc(b.CreateCategory)), common.USER_ACCESS_LEVEL))
	mux.Handle("GET /book-category/{id}", authMw(logMw(http.HandlerFunc(b.GetCategoryById)), common.USER_ACCESS_LEVEL))
	mux.Handle("PATCH /book-category/{id}", authMw(logMw(http.HandlerFunc(b.UpdateCategoryById)), common.USER_ACCESS_LEVEL))
	mux.Handle("DELETE /book-category/{id}", authMw(logMw(http.HandlerFunc(b.DeleteCategoryById)), common.ADMIN_ACCESS_LEVEL))
	mux.Handle("GET /book-category/title", logMw(http.HandlerFunc(b.GetCategoryByTitle)))
	mux.Handle("GET /book-category/", logMw(http.HandlerFunc(b.ListCategories)))
}

// CreateCategory adds a new book category to library.
// @ID createBookCategory
// @Summary CreateCategory
// @Tags book-category
// @Description create book category
// @Accept json
// @Produce json
// @Param req body BookCatRequest true "CatRequest"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/ [post]
func (b *BookCatHandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.CreateCategory"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	var req BookCatRequest
	if err := json.ReadRequestBody(r, &req); err != nil {

		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if common.ValidateFieldNotEmpty(req.Title) {
		json.WriteError(w, common.ReqFieldMessage(req.Title), http.StatusBadRequest)
		return
	}

	id, err := b.repo.Create(ctx, req.Title)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "created book category", map[string]any{"id": id}, http.StatusOK)
}

// GetCategoryById gets a book category by id from library.
// @ID getBookCategoryById
// @Summary GetCategoryById
// @Tags book-category
// @Description get book category by id
// @Accept json
// @Produce json
// @Param id path int true "Category Id"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/{id} [get]
func (b *BookCatHandlerImpl) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.GetCategoryById"

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

// GetCategoryByTitle gets a book category by title from library.
// @ID getBookCategoryTitle
// @Summary GetCategoryByTitle
// @Tags book-category
// @Description get book category by title
// @Accept json
// @Produce json
// @Param title query string true "Category Title"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/title [get]
func (b *BookCatHandlerImpl) GetCategoryByTitle(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.GetCategoryByTitle"

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()
	defer r.Body.Close()

	title := r.URL.Query().Get("title")

	if common.ValidateFieldNotEmpty(title) {
		json.WriteError(w, common.ReqFieldMessage(title), http.StatusBadRequest)
		return
	}

	bc, err := b.repo.GetByTitle(ctx, title)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "fetched book category", bc, http.StatusOK)
}

// UpdateCategoryById updates a book category by id.
// @ID updateBookCategoryById
// @Summary UpdateCategoryById
// @Tags book-category
// @Description update book category by id
// @Accept json
// @Produce json
// @Param inpur body BookCatRequest true "CatRequest"
// @Param id path int true "Category Id"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/{id} [patch]
func (b *BookCatHandlerImpl) UpdateCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.UpdateCategoryById"

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

	json.WriteSuccess(w, "updated book category", map[string]any{"id": id}, http.StatusOK)
}

// DeleteCategoryById deletes a book category by id from library.
// @ID deleteBookCategoryById
// @Summary DeleteCategoryById
// @Tags book-category
// @Description delete book category by id
// @Accept json
// @Produce json
// @Param id path int true "Category Id"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/{id} [delete]
func (b *BookCatHandlerImpl) DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.DeleteCategoryById"

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

	json.WriteSuccess(w, "deleted book category", map[string]any{"id": id}, http.StatusOK)
}

// ListCategories gets a serie of book categories from library.
// @ID listBookCategories
// @Summary ListCategories
// @Tags book-category
// @Description get list of book categories
// @Accept json
// @Produce json
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /book-category/ [get]
func (b *BookCatHandlerImpl) ListCategories(w http.ResponseWriter, r *http.Request) {
	const op = "modules.bookcategory.handler.ListCategories"

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
