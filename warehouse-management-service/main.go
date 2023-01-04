package main

import (
	"fmt"
	"net/http"
	"os"
	"warehouse-management-service/internal/config"
	"warehouse-management-service/pkg/storage/postgres"

	"github.com/gin-gonic/gin"
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

	err = startServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}
}

func startServer() error {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router.Run(":80")
}
