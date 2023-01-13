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
	_, err := w.Write([]byte("pong"))
	if err != nil {
		h.logger.Log(log.Error, err)
	}
}
