package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (p *Postgres) RunMigration(sourceUrl string) error {
	databaseUrl := buildConnectionURL(p.PostgresConfig)
	migrateInstance, err := migrate.New(sourceUrl, databaseUrl)
	if err != nil {
		return err
	}
	defer migrateInstance.Close()
	if err := migrateInstance.Up(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) MigrateDown(sourceUrl string) error {
	databaseUrl := buildConnectionURL(p.PostgresConfig)
	migrateInstance, err := migrate.New(sourceUrl, databaseUrl)
	if err != nil {
		return err
	}
	defer migrateInstance.Close()
	if err := migrateInstance.Down(); err != nil {
		return err
	}
	return nil
}
