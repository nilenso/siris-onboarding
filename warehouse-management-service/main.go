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

// EnvConfigFilePath is the env variable that specifies absolute path of the config file
const EnvConfigFilePath = "CONFIG_FILE_PATH"

func main() {
	runDBMigrations := *flag.Bool("migrate", false, "true or false, specifies if database migrations should be run")

	logger := log.New(log.Warning)

	configFilePath, ok := os.LookupEnv(EnvConfigFilePath)
	if !ok {
		logger.Log(log.Fatal, fmt.Sprintf("%s env variable not set", EnvConfigFilePath))
	}

	config, err := config.FromFile(configFilePath)
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
	}

	connectionURL := postgres.Connection(config.Postgres)
	db, err := postgres.New(connectionURL)
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
	}

	if runDBMigrations {
		err = RunMigration(config.DBMigration.SourcePath, connectionURL)
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

	handler := handler.New(db, logger)

	err = http.ListenAndServe(":80", handler)
	switch err {
	case http.ErrServerClosed:
		logger.Log(log.Info, "server shut down successfully")
	default:
		logger.Log(log.Error, fmt.Sprintf("error starting up server: %v", err))
		os.Exit(1)
	}
}
