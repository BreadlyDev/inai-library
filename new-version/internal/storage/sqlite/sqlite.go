package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath+"/db.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, err := os.ReadFile("./docs/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(string(query))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}
