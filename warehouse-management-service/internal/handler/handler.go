package handler

import (
	"net/http"
	"warehouse-management-service/pkg/database"
	"warehouse-management-service/pkg/log"
)

type handler struct {
	db     database.Service
	logger log.Logger
}

func New(db database.Service, logger log.Logger) http.Handler {
	handler := &handler{
		db:     db,
		logger: logger,
	}
	return handler.router()
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
