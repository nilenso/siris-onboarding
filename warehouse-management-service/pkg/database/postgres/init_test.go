package postgres

import (
	"fmt"
	"testing"
	"warehouse-management-service/internal/config"

	_ "github.com/lib/pq"
)

var warehouseService *WarehouseService
var shelfBlockService *ShelfBlockService
var shelfService *ShelfService
var postgres *Postgres

func TestMain(m *testing.M) {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config: %v", err))
	}

	postgres = New(cfg.Postgres)

	db, err := postgres.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	err = postgres.RunMigration(cfg.DBMigrationSourcePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to run database migrations: %v", err))
	}

	warehouseService = NewWarehouseService(db)
	shelfBlockService = NewShelfBlockService(db)
	shelfService = NewShelfService(db)

	mockDB, err := postgres.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	err = mockDB.Close()
	if err != nil {
		panic(fmt.Sprintf("Failed to close database connection: %v", err))
	}

	m.Run()

	err = postgres.MigrateDown(cfg.DBMigrationSourcePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate down: %v", err))
	}
}
