package httpserver

import (
	"log/slog"
	"net/http"
	"new-version/internal/config"
	bc "new-version/internal/modules/book-category"
	"new-version/internal/modules/user"
	"new-version/internal/storage/sqlite"

	"github.com/rs/cors"
	swagger "github.com/swaggo/http-swagger"

	_ "new-version/docs"
)

func NewServer(log *slog.Logger, cfg *config.Config, stg *sqlite.Storage) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", swagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	bcRepo := bc.NewBookCatRepo(stg.DB)
	bcHandler := bc.NewBookCatHandler(log, bcRepo, &cfg.Security)
	bcHandler.RegisterRoutes(mux, log)

	aSrv := user.NewJwtAuthService(log, &cfg.Security)
	uRepo := user.NewUserRepo(stg.DB)
	uSrv := user.NewUserService(log, uRepo, aSrv, &cfg.Security)
	uHandler := user.NewUserHandler(log, uSrv, &cfg.Security)
	uHandler.RegisterRoutes(mux, log)

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
