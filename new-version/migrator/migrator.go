package migrator

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigrations(migrUrl string, dbUrl string) {
	if migrUrl == "" {
		log.Fatal("migration path required")
	}

	m, err := migrate.New(migrUrl, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")

			return
		}

		log.Fatal(err)
	}

	log.Println("migrations applied")
}
