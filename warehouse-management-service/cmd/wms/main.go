package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"warehouse-management-service/internal/config"
	"warehouse-management-service/internal/handler"
	"warehouse-management-service/pkg/database/postgres"
	"warehouse-management-service/pkg/log"

	_ "github.com/lib/pq"
)

func main() {
	runDBMigrations := flag.Bool("migrate", false, "true or false, specifies if database migrations should be run")
	flag.Parse()

	appConfig, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config %v", err))
	}

	logger := log.New()
	logger.SetLevel(appConfig.LogLevel)

	pg := postgres.New(appConfig.Postgres)

	db, err := pg.Open()
	if err != nil {
		logger.Log(log.Fatal, fmt.Sprintf("Failed to connect to database: %v", err))
		return
	}

	// close db gracefully
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Log(log.Error, err)
		}
	}()

	if *runDBMigrations {
		logger.Log(log.Info, "Running database migrations")
		err := pg.RunMigration(appConfig.DBMigrationSourcePath)
		if err != nil {
			logger.Log(log.Fatal, fmt.Sprintf("Failed to run database migrations: %v", err))
			return
		}
	}

	// Instantiate Postgres-backed services
	warehouseService := postgres.NewWarehouseService(db)
	shelfBlockService := postgres.NewShelfBlockService(db)
	shelfService := postgres.NewShelfService(db)
	productService := postgres.NewProductService(db)

	h := handler.New(logger, warehouseService, shelfBlockService, shelfService, productService)

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	server := &http.Server{Addr: ":80", Handler: h}
	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logger.Log(log.Info, "server shut down successfully")
			} else {
				logger.Log(log.Error, fmt.Sprintf("error starting up server: %v", err))
			}
		}
	}()
	logger.Log(log.Info, "Server listening on port 80")

	// listen for exit signals
	<-exitChan

	// shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log(log.Error, fmt.Sprintf("Failed to shutdown the server %v", err))
	}

	// wait for server shutdown
	wg.Wait()
}
