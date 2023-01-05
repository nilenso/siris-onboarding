package server

import (
	"warehouse-management-service/pkg/storage"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	storage storage.Service
	Router  *chi.Mux
}

// Creates and returns a chi Router and configures the server routes
func New(storage storage.Service) *Server {
	server := &Server{
		Router:  chi.NewRouter(),
		storage: storage,
	}
	server.routes()
	return server
}
