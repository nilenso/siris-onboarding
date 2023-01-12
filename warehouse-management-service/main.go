package main

import (
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
	CONFIG_FILE_PATH_ENV_VARIABLE = "CONFIG_FILE_PATH"
)

func main() {
	logger := log.New(log.Warning)

	// Get absolute path of config file from env variable
	configFilePath, ok := os.LookupEnv(CONFIG_FILE_PATH_ENV_VARIABLE)
	if !ok {
		logger.Log(log.Fatal, fmt.Sprintf("%s env variable not set", CONFIG_FILE_PATH_ENV_VARIABLE))
	}

	// Get config
	config, err := config.FromFile(configFilePath)
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
	}

	// Get Postgres DB service
	db, err := postgres.New(config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.SSLMode)
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
	}

	// Close DB connection at shutdown
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Log(log.Error, err)
		}
	}()

	handler := handler.New(db)

	err = http.ListenAndServe(":80", handler)
	switch err {
	case http.ErrServerClosed:
		logger.Log(log.Info, "server shut down successfully")
	default:
		logger.Log(log.Error, fmt.Sprintf("error starting up server: %v", err))
		os.Exit(1)
	}
}
