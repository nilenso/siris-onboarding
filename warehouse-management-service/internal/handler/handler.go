package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	warehousemanagementservice "warehouse-management-service"
	"warehouse-management-service/pkg/log"
	"warehouse-management-service/pkg/wms"
)

type handler struct {
	warehouseService warehousemanagementservice.WarehouseService
	logger           log.Logger
}

func New(logger log.Logger, warehouseService warehousemanagementservice.WarehouseService) http.Handler {
	handler := &handler{
		logger:           logger,
		warehouseService: warehouseService,
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
		err := fmt.Errorf("%v", wms.GetWarehouseResponse{
			Error: "warehouse id cannot be empty",
		})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	warehouse, err := h.warehouseService.GetWarehouseById(r.Context(), warehouseId)
	switch err {
	case warehousemanagementservice.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, wms.GetWarehouseResponse{Error: fmt.Sprintf(
				"failed to get, warehouse: %s does not exist",
				warehouseId,
			)})
		}
	case nil:
		{
			h.response(w, http.StatusOK, wms.GetWarehouseResponse{Response: *warehouse})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, wms.GetWarehouseResponse{Error: "Failed to get warehouse"})
		}
	}
}

func (h *handler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var createWarehouseRequest wms.CreateWarehouseRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: err.Error(),
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&createWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: "Failed to parse request",
		})
		return
	}

	if err, ok := createWarehouseRequest.IsValid(); !ok {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	warehouse := warehousemanagementservice.New(
		createWarehouseRequest.Name,
		createWarehouseRequest.Latitude,
		createWarehouseRequest.Longitude,
	)

	err = h.warehouseService.CreateWarehouse(r.Context(), warehouse)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusInternalServerError, wms.WarehouseResponse{Error: "Failed to create warehouse"})
		return
	}

	h.response(w, http.StatusOK, wms.WarehouseResponse{Response: fmt.Sprintf(
		"Successfully created warehouse: %s",
		warehouse.Id,
	)})
}

func (h *handler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	var updateWarehouseRequest wms.UpdateWarehouseRequest

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: err.Error(),
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&updateWarehouseRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: "Failed to parse request"})
		return
	}

	if err, ok := updateWarehouseRequest.IsValid(); !ok {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, wms.WarehouseResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err = h.warehouseService.UpdateWarehouse(r.Context(), &warehousemanagementservice.Warehouse{
		Id:        updateWarehouseRequest.Id,
		Name:      updateWarehouseRequest.Name,
		Latitude:  updateWarehouseRequest.Latitude,
		Longitude: updateWarehouseRequest.Longitude,
	})
	switch err {
	case warehousemanagementservice.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, wms.WarehouseResponse{Error: fmt.Sprintf(
				"failed to update, warehouse: %s does not exist",
				updateWarehouseRequest.Id,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				wms.WarehouseResponse{
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
				wms.WarehouseResponse{Error: "Failed to update warehouse"},
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
	case warehousemanagementservice.WarehouseDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, wms.WarehouseResponse{Error: fmt.Sprintf(
				"failed to delete, warehouse: %s does not exist",
				warehouseId,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				wms.WarehouseResponse{Response: fmt.Sprintf(
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
				wms.WarehouseResponse{Error: "Failed to delete warehouse"},
			)
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
