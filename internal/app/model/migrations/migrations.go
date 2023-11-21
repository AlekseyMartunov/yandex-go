// Package migrations create migration in db when app starting
package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

// MakeMigration check db and if tables does not exist create it
func MakeMigration(db *sql.DB) error {
	goose.SetDialect("urlpostgres")

	err := goose.UpTo(db, "./internal/app/model/migrations", 1)
	if err != nil {
		return err
	}

	err = goose.UpTo(db, "./internal/app/model/migrations", 2)
	if err != nil {
		return err
	}

	err = goose.UpTo(db, "./internal/app/model/migrations", 3)
	if err != nil {
		return err
	}

	return nil
}
