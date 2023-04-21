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

// mockgen -destination="./internal/handler/mock/warehouse.go" warehouse-management-service/internal/handler WarehouseService
type WarehouseService interface {
	GetWarehouseById(ctx context.Context, id string) (*wms.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
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
	if err != nil {
		if err == wms.WarehouseDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.GetWarehouseResponse{Error: fmt.Sprintf(
				"failed to get, warehouse: %s does not exist",
				warehouseId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.GetWarehouseResponse{Error: "Failed to get warehouse"})
			return
		}
	}
	h.response(w, http.StatusOK, api.GetWarehouseResponse{Response: *warehouse})
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

	err := decoder.Decode(&updateWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: "Failed to parse request"})
		return
	}

	err = validator.Validate(updateWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.WarehouseResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err = h.warehouseService.UpdateWarehouse(r.Context(), &wms.Warehouse{
		Id:        updateWarehouseRequest.Id,
		Name:      updateWarehouseRequest.Name,
		Latitude:  updateWarehouseRequest.Latitude,
		Longitude: updateWarehouseRequest.Longitude,
	})
	if err != nil {
		if err == wms.WarehouseDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to update, warehouse: %s does not exist",
				updateWarehouseRequest.Id,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.WarehouseResponse{Error: "Failed to update warehouse"},
			)
		}
		return
	}

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

func (h *handler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseId := chi.URLParam(r, "warehouseId")

	if warehouseId == "" {
		err := fmt.Errorf("%v", map[string]string{"error": "warehouse id cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.warehouseService.DeleteWarehouse(r.Context(), warehouseId)
	if err != nil {
		if err == wms.WarehouseDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to delete, warehouse: %s does not exist",
				warehouseId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.WarehouseResponse{Error: "Failed to delete warehouse"},
			)
			return
		}
	}

	h.response(
		w,
		http.StatusOK,
		api.WarehouseResponse{Response: fmt.Sprintf(
			"Successfully deleted warehouse: %s",
			warehouseId,
		)},
	)
}
