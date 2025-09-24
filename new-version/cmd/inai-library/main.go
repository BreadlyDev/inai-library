package main

import (
	"context"
	"fmt"
	"log"
	"new-version/internal/config"
	bookCat "new-version/internal/storage/repositories/bookcategory"
	"new-version/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	_ = storage
}
