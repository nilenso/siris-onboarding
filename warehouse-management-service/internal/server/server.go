package server

import (
	"warehouse-management-service/pkg/database"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	db     database.Service
	Router *chi.Mux
}

// Creates and returns a chi Router and configures the server routes
func New(db database.Service) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	server := &Server{
		Router: router,
		db:     db,
	}
	server.routes()
	return server
}
