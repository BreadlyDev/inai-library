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

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	log := logger.SetupLogger(cfg.Env)

	done := make(chan os.Signal, 1)

	srv := httpserver.NewServer(log, &cfg.HTTPServer, storage)

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

	// TODO: Add migrator
	// TODO: Add swagger
	// TODO: Refactor project
	// TODO: Add CORS
	// TODO: Add other entities
	// TODO: Add business logic
	// TODO: Move to Postgres
	// TODO: Add dockerfile & deploy on ...
}
