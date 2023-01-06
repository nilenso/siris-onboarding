package handler

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *handler) router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/ping", h.Ping)

	return router
}
