package httpserver

import (
	"log/slog"
	"net/http"
	"new-version/internal/config"
	bookCatHdl "new-version/internal/http/handler/bookcategory"
	userHdl "new-version/internal/http/handler/user"
	bookCatRepo "new-version/internal/repository/bookcategory"
	userRepo "new-version/internal/repository/user"
	authSvc "new-version/internal/service/auth"
	userSvc "new-version/internal/service/user"

	"new-version/internal/storage/postgres"

	"github.com/rs/cors"
	swagger "github.com/swaggo/http-swagger"

	_ "new-version/docs"
)

func New(log *slog.Logger, cfg *config.Config, stg *postgres.Storage) *http.Server {
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

	bcRepo := bookCatRepo.New(stg.DB())
	bcHandler := bookCatHdl.New(log, bcRepo, &cfg.Security)
	bcHandler.RegisterRoutes(mux, log)

	aSvc := authSvc.New(log, &cfg.Security)
	uRepo := userRepo.New(stg.DB())
	uSrv := userSvc.New(log, uRepo, aSvc, &cfg.Security)
	uHandler := userHdl.New(log, uSrv, &cfg.Security)
	uHandler.RegisterRoutes(mux, log)

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
