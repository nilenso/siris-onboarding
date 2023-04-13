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

	router.Get("/warehouse/{warehouseId}", h.GetWarehouse)
	router.Post("/warehouse", h.CreateWarehouse)
	router.Put("/warehouse", h.UpdateWarehouse)
	router.Delete("/warehouse/{warehouseId}", h.DeleteWarehouse)

	router.Get("/shelf_block/{shelfBlockId}", h.GetShelfBlock)
	router.Post("/shelf_block", h.CreateShelfBlock)
	router.Put("/shelf_block", h.UpdateShelfBlock)
	router.Delete("/shelf_block/{shelfBlockId}", h.DeleteShelfBlock)

	router.Get("/shelf/{shelfId}", h.GetShelf)
	router.Post("/shelf", h.CreateShelf)
	router.Put("/shelf", h.UpdateShelf)
	router.Delete("/shelf/{shelfId}", h.DeleteShelf)

	return router
}
