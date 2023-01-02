package main

import (
	"fmt"
	"net/http"
	"os"
	"warehouse-management-service/pkg/storage/postgres"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
	DB_USER     = "wms"
	DB_PASSWORD = "x5t5%h^NsXE3"
	DB_NAME     = "warehouse-management-system"
	DB_SSL_MODE = false
)

func main() {
	db, err := postgres.New(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}
	// Close DB connection
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
