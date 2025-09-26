package main

import (
	"log"
	"net/http"
	"new-version/internal/config"
	bc "new-version/internal/modules/book-category"
	"new-version/internal/storage/sqlite"
	"new-version/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	_ = storage

	log := logger.SetupLogger(cfg.Env)
	bcRepo := bc.NewBookCatRepo(storage.DB)
	bcSrv := bc.NewBookCatHandler(log, bcRepo)

	mux := http.NewServeMux()
	bcSrv.RegisterRoutes(mux)

	http.ListenAndServe(":8080", mux)
}
