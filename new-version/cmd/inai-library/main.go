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

	bookCatRepo := bookCat.NewBookCatRepo(storage.DB)

	book, err := bookCatRepo.GetById(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d, %s, %s", book.Id, book.Title, book.CreatedTime.UTC().Format("2006-01-02 15:04:05"))
}
