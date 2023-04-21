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

// mockgen -destination="./internal/handler/mock/shelf_block.go" warehouse-management-service/internal/handler ShelfBlockService
type ShelfBlockService interface {
	GetShelfBlockById(ctx context.Context, id string) (wms.ShelfBlock, error)
	CreateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	UpdateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	DeleteShelfBlockById(ctx context.Context, id string) error
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

	shelfBlock, err := h.shelfBlockService.GetShelfBlockById(r.Context(), shelfBlockId)
	if err != nil {
		if err == wms.ShelfBlockDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.GetShelfBlockResponse{Error: fmt.Sprintf(
				"failed to get, shelfBlock: %s does not exist",
				shelfBlockId,
			)})
			return
		} else {

			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.GetShelfBlockResponse{Error: "Failed to get shelfBlock"})
			return
		}
	}
	h.response(w, http.StatusOK, api.GetShelfBlockResponse{Response: shelfBlock})
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

	err = h.shelfBlockService.CreateShelfBlock(r.Context(), shelfBlock)
	if err != nil {
		if err == wms.InvalidWarehouse {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				shelfBlock.WarehouseId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to create shelf block"},
			)
			return
		}
	}
	h.response(w, http.StatusOK, api.ShelfBlockResponse{Response: fmt.Sprintf(
		"Successfully created shelf_block: %s",
		shelfBlock.Id,
	)})
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
	err = h.shelfBlockService.UpdateShelfBlock(r.Context(), shelfBlock)
	if err != nil {
		if err == wms.ShelfBlockDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to update, shelf_block: %s does not exist",
				updateShelfBlockRequest.Id,
			)})
			return
		} else if err == wms.InvalidWarehouse {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				shelfBlock.WarehouseId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to update shelf_block"},
			)
			return
		}
	}

	h.response(w, http.StatusOK, api.ShelfBlockResponse{Response: fmt.Sprintf(
		"Successfully updated shelf_block: %s",
		shelfBlock.Id,
	)})
}

func (h *handler) DeleteShelfBlock(w http.ResponseWriter, r *http.Request) {
	shelfBlockId := chi.URLParam(r, "shelfBlockId")

	if shelfBlockId == "" {
		err := fmt.Errorf("%v", api.ShelfBlockResponse{Error: "shelf_block id cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.shelfBlockService.DeleteShelfBlockById(r.Context(), shelfBlockId)
	if err != nil {
		if err == wms.ShelfBlockDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to delete, shelf_block: %s does not exist",
				shelfBlockId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfBlockResponse{Error: "Failed to delete shelf_block"},
			)
			return
		}
	}

	h.response(
		w,
		http.StatusOK,
		api.ShelfBlockResponse{Response: fmt.Sprintf(
			"Successfully deleted shelf_block: %s",
			shelfBlockId,
		)},
	)
}
