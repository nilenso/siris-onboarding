package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"warehouse-management-service/pkg/log"
)

type handler struct {
	warehouseService  WarehouseService
	shelfBlockService ShelfBlockService
	shelfService      ShelfService
	productService    ProductService
	logger            log.Logger
}

func New(
	logger log.Logger,
	warehouseService WarehouseService,
	shelfBlockService ShelfBlockService,
	shelfService ShelfService,
	productService ProductService,
) http.Handler {
	handler := &handler{
		logger:            logger,
		warehouseService:  warehouseService,
		shelfBlockService: shelfBlockService,
		shelfService:      shelfService,
		productService:    productService,
	}
	return handler.router()
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		h.logger.Log(log.Error, err)
	}
}

func (h *handler) response(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")

	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		h.logger.Log(log.Error, err)

		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := fmt.Sprintf("%v", map[string]string{"error": "Failed to parse response"})
		_, err = w.Write([]byte(errorResponse))
		if err != nil {
			h.logger.Log(log.Error, err)
		}
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(marshalledResponse)
	if err != nil {
		h.logger.Log(log.Error, err)
	}
}
