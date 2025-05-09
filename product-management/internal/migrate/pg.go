package migrate

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Do(db *sql.DB) error {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}
