package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"warehouse-management-service/internal/config"
	"warehouse-management-service/internal/handler"
	"warehouse-management-service/pkg/database/postgres"
	"warehouse-management-service/pkg/log"

	_ "github.com/lib/pq"
)

const (
	EnvConfigFilePath = "CONFIG_FILE_PATH"
)

func main() {
	runDBMigrations := *flag.Bool("migrate", false, "true or false, specifies if database migrations should be run")

	configFilePath, ok := os.LookupEnv(EnvConfigFilePath)
	if !ok {
		panic("Failed to read environment variable")
	}

	appConfig, err := config.FromFile(configFilePath)
	if err != nil {
		panic("Failed to read appConfig file")
	}

	logger := log.New()
	logger.SetLevel(appConfig.LogLevel)

	pg := postgres.New(appConfig.Postgres)

	db, err := pg.Open()
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
	}

	if runDBMigrations {
		err := pg.RunMigration(appConfig.DBMigration.SourcePath)
		if err != nil {
			logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
			os.Exit(1)
		}
	}

	defer func() {
		err = db.Close()
		if err != nil {
			logger.Log(log.Error, err)
		}
	}()

	// Instantiate Postgres-backed services
	warehouseService := postgres.NewWarehouseService(db)
	shelfService := postgres.NewShelfService(db)

	h := handler.New(logger, warehouseService, shelfService)

	err = http.ListenAndServe(":80", h)
	switch err {
	case http.ErrServerClosed:
		logger.Log(log.Info, "server shut down successfully")
	default:
		logger.Log(log.Error, fmt.Sprintf("error starting up server: %v", err))
		os.Exit(1)
	}
}
