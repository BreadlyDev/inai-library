package postgres

import (
	"database/sql"
	"fmt"
	"new-version/internal/config"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Database) (*Storage, error) {
	const op = "storage.postgres.New"

	dbUrl := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Host,
		cfg.Port,
	)

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) DB() *sql.DB {
	return s.db
}
