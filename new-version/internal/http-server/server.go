package httpserver

import (
	"log/slog"
	"net/http"
	"new-version/internal/config"
	bc "new-version/internal/modules/book-category"
	"new-version/internal/storage/sqlite"

	swagger "github.com/swaggo/http-swagger"

	_ "new-version/docs"
)

func NewServer(log *slog.Logger, cfg *config.HTTPServer, stg *sqlite.Storage) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", swagger.WrapHandler)

	bcRepo := bc.NewBookCatRepo(stg.DB)
	bcHandler := bc.NewBookCatHandler(bcRepo)

	bcHandler.RegisterRoutes(mux, log)

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
