package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gopkg.in/validator.v2"
	"net/http"
	wms "warehouse-management-service"
	"warehouse-management-service/pkg/api"
	"warehouse-management-service/pkg/log"
)

// mockgen -source="./warehouse.go" -destination="./internal/handler/mock/warehouse.go"
type WarehouseService interface {
	GetWarehouseById(ctx context.Context, id string) (*wms.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
}

// mockgen -source="./shelf_block.go" -destination="./internal/handler/mock/shelf_block.go"
type ShelfService interface {
	GetShelfBlockById(ctx context.Context, id string) (wms.ShelfBlock, error)
	CreateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	UpdateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	DeleteShelfBlockById(ctx context.Context, id string) error
}

type handler struct {
	warehouseService WarehouseService
	shelfService     ShelfService
	logger           log.Logger
}

func New(
	logger log.Logger,
	warehouseService WarehouseService,
	shelfService ShelfService,
) http.Handler {
	handler := &handler{
		logger:           logger,
		warehouseService: warehouseService,
		shelfService:     shelfService,
	}
	return handler.router()
}

func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		h.logger.Log(log.Error, err)
	}
}

func (h *handler) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseId := chi.URLParam(r, "warehouseId")

	if warehouseId == "" {
		err := fmt.Errorf("%v", api.GetWarehouseResponse{
			Error: "warehouse id cannot be empty",
		})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	warehouse, err := h.warehouseService.GetWarehouseById(r.Context(), warehouseId)
	switch err {
	case wms.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.GetWarehouseResponse{Error: fmt.Sprintf(
				"failed to get, warehouse: %s does not exist",
				warehouseId,
			)})
		}
	case nil:
		{
			h.response(w, http.StatusOK, api.GetWarehouseResponse{Response: *warehouse})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.GetWarehouseResponse{Error: "Failed to get warehouse"})
		}
	}
}

func (h *handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var createWarehouseRequest api.CreateWarehouseRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: err.Error(),
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&createWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: "Failed to parse request",
		})
		return
	}

	err = validator.Validate(createWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	warehouse := wms.NewWarehouse(
		createWarehouseRequest.Name,
		createWarehouseRequest.Latitude,
		createWarehouseRequest.Longitude,
	)

	err = h.warehouseService.CreateWarehouse(r.Context(), warehouse)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusInternalServerError, api.WarehouseResponse{Error: "Failed to create warehouse"})
		return
	}

	h.response(w, http.StatusOK, api.WarehouseResponse{Response: fmt.Sprintf(
		"Successfully created warehouse: %s",
		warehouse.Id,
	)})
}

func (h *handler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	var updateWarehouseRequest api.UpdateWarehouseRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: err.Error(),
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&updateWarehouseRequest); err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: "Failed to parse request"})
		return
	}

	if err := validator.Validate(updateWarehouseRequest); err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err := h.warehouseService.UpdateWarehouse(r.Context(), &wms.Warehouse{
		Id:        updateWarehouseRequest.Id,
		Name:      updateWarehouseRequest.Name,
		Latitude:  updateWarehouseRequest.Latitude,
		Longitude: updateWarehouseRequest.Longitude,
	})
	switch err {
	case wms.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to update, warehouse: %s does not exist",
				updateWarehouseRequest.Id,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				api.WarehouseResponse{
					Response: fmt.Sprintf(
						"Successfully updated warehouse: %s",
						updateWarehouseRequest.Id,
					)},
			)
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.WarehouseResponse{Error: "Failed to update warehouse"},
			)
		}
	}
}

func (h *handler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseId := chi.URLParam(r, "warehouseId")

	if warehouseId == "" {
		err := fmt.Errorf("%v", map[string]string{"error": "warehouse id cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.warehouseService.DeleteWarehouse(r.Context(), warehouseId)
	switch err {
	case wms.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to delete, warehouse: %s does not exist",
				warehouseId,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				api.WarehouseResponse{Response: fmt.Sprintf(
					"Successfully deleted warehouse: %s",
					warehouseId,
				)},
			)
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.WarehouseResponse{Error: "Failed to delete warehouse"},
			)
		}
	}
}

func (h *handler) GetShelfBlock(w http.ResponseWriter, r *http.Request) {
	shelfBlockId := chi.URLParam(r, "shelfBlockId")

	if shelfBlockId == "" {
		err := fmt.Errorf("%v", api.GetShelfBlockResponse{
			Error: "shelf block id cannot be empty",
		})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	shelfBlock, err := h.shelfService.GetShelfBlockById(r.Context(), shelfBlockId)
	switch err {
	case wms.ShelfBlockDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.GetShelfBlockResponse{Error: fmt.Sprintf(
				"failed to get, shelfBlock: %s does not exist",
				shelfBlockId,
			)})
		}
	case nil:
		{
			h.response(w, http.StatusOK, api.GetShelfBlockResponse{Response: shelfBlock})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.GetShelfBlockResponse{Error: "Failed to get shelfBlock"})
		}
	}
}

func (h *handler) CreateShelfBlock(w http.ResponseWriter, r *http.Request) {
	var createShelfBlockRequest api.CreateShelfBlockRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{
			Error: err.Error()})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createShelfBlockRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{
			Error: "Failed to parse request"})
		return
	}

	err = validator.Validate(createShelfBlockRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	shelfBlock := wms.NewShelfBlock(
		createShelfBlockRequest.Aisle,
		createShelfBlockRequest.Rack,
		createShelfBlockRequest.StorageType,
		createShelfBlockRequest.WarehouseId)

	err = h.shelfService.CreateShelfBlock(r.Context(), shelfBlock)
	switch err {
	case nil:
		{
			h.response(w, http.StatusOK, api.ShelfBlockResponse{Response: fmt.Sprintf(
				"Successfully created shelf_block: %s",
				shelfBlock.Id,
			)})
		}
	case wms.InvalidWarehouse:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				shelfBlock.WarehouseId,
			)})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to create shelf block"},
			)
		}
	}
}

func (h *handler) UpdateShelfBlock(w http.ResponseWriter, r *http.Request) {
	var updateShelfBlockRequest api.UpdateShelfBlockRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(
			w,
			http.StatusBadRequest,
			api.ShelfBlockResponse{Error: err.Error()},
		)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&updateShelfBlockRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{
			Error: "Failed to parse request",
		})
		return
	}

	err = validator.Validate(updateShelfBlockRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	shelfBlock := wms.ShelfBlock{
		Id:          updateShelfBlockRequest.Id,
		Aisle:       updateShelfBlockRequest.Aisle,
		Rack:        updateShelfBlockRequest.Rack,
		StorageType: updateShelfBlockRequest.StorageType,
		WarehouseId: updateShelfBlockRequest.WarehouseId,
	}
	err = h.shelfService.UpdateShelfBlock(r.Context(), shelfBlock)
	switch err {
	case nil:
		{
			h.response(w, http.StatusOK, api.ShelfBlockResponse{Response: fmt.Sprintf(
				"Successfully updated shelf_block: %s",
				shelfBlock.Id,
			)})
		}
	case wms.ShelfBlockDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to update, shelf_block: %s does not exist",
				updateShelfBlockRequest.Id,
			)})
		}
	case wms.InvalidWarehouse:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				shelfBlock.WarehouseId,
			)})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to update shelf_block"},
			)
		}
	}
}

func (h *handler) DeleteShelfBlock(w http.ResponseWriter, r *http.Request) {
	shelfBlockId := chi.URLParam(r, "shelfBlockId")

	if shelfBlockId == "" {
		err := fmt.Errorf("%v", api.ShelfBlockResponse{Error: "shelf_block id cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.shelfService.DeleteShelfBlockById(r.Context(), shelfBlockId)
	switch err {
	case wms.ShelfBlockDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to delete, shelf_block: %s does not exist",
				shelfBlockId,
			)})
			return
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				api.ShelfBlockResponse{Response: fmt.Sprintf(
					"Successfully deleted shelf_block: %s",
					shelfBlockId,
				)},
			)
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to delete shelf_block"},
			)
			return
		}
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
