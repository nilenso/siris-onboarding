package postgres

import (
	"fmt"
	"os"
	"testing"
	"warehouse-management-service/internal/config"

	_ "github.com/lib/pq"
)

const (
	EnvConfigFilePath = "CONFIG_FILE_PATH"
)

var warehouseService *WarehouseService
var shelfService *ShelfBlockService
var postgres *Postgres

func TestMain(m *testing.M) {
	configFilePath, ok := os.LookupEnv(EnvConfigFilePath)
	if !ok {
		panic("Failed to read environment variable")
	}

	cfg, err := config.FromFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}

	postgres = New(cfg.Postgres)

	db, err := postgres.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	err = postgres.RunMigration(cfg.DBMigration.SourcePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to run database migrations: %v", err))
	}

	warehouseService = NewWarehouseService(db)
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

	err = postgres.MigrateDown(cfg.DBMigration.SourcePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate down: %v", err))
	}
}
