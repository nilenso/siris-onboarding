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

func main() {
	runDBMigrations := *flag.Bool("migrate", false, "true or false, specifies if database migrations should be run")

	appConfig, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config %v", err))
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
		err := pg.RunMigration(appConfig.DBMigrationSourcePath)
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
	shelfBlockService := postgres.NewShelfBlockService(db)
	shelfService := postgres.NewShelfService(db)

	h := handler.New(logger, warehouseService, shelfBlockService, shelfService)

	err = http.ListenAndServe(":80", h)
	switch err {
	case http.ErrServerClosed:
		logger.Log(log.Info, "server shut down successfully")
	default:
		logger.Log(log.Error, fmt.Sprintf("error starting up server: %v", err))
		os.Exit(1)
	}
}
