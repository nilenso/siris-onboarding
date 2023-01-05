package main

import (
	"fmt"
	"net/http"
	"os"
	"warehouse-management-service/internal/config"
	"warehouse-management-service/internal/server"
	"warehouse-management-service/pkg/storage/postgres"

	_ "github.com/lib/pq"
)

const (
	CONFIG_FILE_PATH_ENV_VARIABLE = "CONFIG_FILE_PATH"
)

func main() {
	// Get absolute path of config file from env variable
	configFilePath, ok := os.LookupEnv(CONFIG_FILE_PATH_ENV_VARIABLE)
	if !ok {
		fmt.Fprintf(os.Stderr, "%s env variable not set", CONFIG_FILE_PATH_ENV_VARIABLE)
	}

	// Get config
	config, err := config.FromFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
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
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}

	// Close DB connection at shutdown
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	server := server.New(db)

	http.ListenAndServe(":80", server.Router)

	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}
}
