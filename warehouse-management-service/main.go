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

const EnvConfigFilePath = "CONFIG_FILE_PATH"

func main() {
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

	db, err := postgres.New(config.Postgres)
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("App startup error: %v", err))
		os.Exit(1)
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
