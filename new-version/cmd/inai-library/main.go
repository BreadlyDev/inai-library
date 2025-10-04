package main

import (
	"context"
	"log"
	"new-version/internal/config"
	httpserver "new-version/internal/http-server"
	"new-version/internal/storage/sqlite"
	"new-version/pkg/logger"
	"os"
	"time"
)

// @title INAI Library API
// @version 2.0
// @description API Server for library application of the university INAI

// @host localhost:8080
// @BasePath

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	log := logger.SetupLogger(cfg.Env)

	done := make(chan os.Signal, 1)

	srv := httpserver.NewServer(log, cfg, storage)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done

	log.Info("stop server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server: ", err)

		return
	}

	defer storage.DB.Close()

	log.Info("server stopped")

	// TODO: Refactor project
	// TODO: Add other entities
	// TODO: Add business logic
	// TODO: Move to Postgres
	// TODO: Add dockerfile & deploy on ...
}
