package handler

import (
	"net/http"
	"warehouse-management-service/pkg/database"
)

type handler struct {
	db database.Service
}

func New(db database.Service) http.Handler {
	handler := &handler{
		db: db,
	}
	return handler.router()
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
