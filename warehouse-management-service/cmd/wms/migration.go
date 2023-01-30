package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(sourceUrl string, databaseUrl string) error {
	migrateInstance, err := migrate.New(sourceUrl, databaseUrl)
	if err != nil {
		return err
	}
	if err := migrateInstance.Up(); err != nil {
		return err
	}
	return nil
}
