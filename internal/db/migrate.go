package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) RunMigrations() error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %v", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations", // file source
		"postgres",                      // database type
		driver,                          // database instance
	)
	if err != nil {
		return fmt.Errorf("could not create migrations: %v", err)
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	return nil
}
