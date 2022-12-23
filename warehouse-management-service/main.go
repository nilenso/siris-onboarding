package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	err := initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "App startup error: %v", err)
		os.Exit(1)
	}
}

func initialize() error {
	err := setupDatabase()
	if err != nil {
		return err
	}

	err = startServer()
	if err != nil {
		return err
	}

	return nil
}

func setupDatabase() error {
	// TODO
	return nil
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
