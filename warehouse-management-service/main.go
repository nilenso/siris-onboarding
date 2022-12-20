package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	CONNECTION_STRING = "postgres://siripr@localhost/warehouse-management-system?sslmode=disable"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}
}

func initialize() error {
	db, err := setupDatabase(CONNECTION_STRING)
	if err != nil {
		return err
	}

	fmt.Println(db.Stats())

	err = startServer()
	if err != nil {
		return err
	}

	return nil
}

func setupDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
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
