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

// mockgen -destination="./internal/handler/mock/shelf.go" warehouse-management-service/internal/handler ShelfService
type ShelfService interface {
	GetShelfById(ctx context.Context, id string) (wms.Shelf, error)
	CreateShelf(ctx context.Context, shelf wms.Shelf) error
	UpdateShelf(ctx context.Context, shelf wms.Shelf) error
	DeleteShelfById(ctx context.Context, id string) error
}

func (h *handler) GetShelf(w http.ResponseWriter, r *http.Request) {
	shelfId := chi.URLParam(r, "shelfId")

	if shelfId == "" {
		err := fmt.Errorf("%v", api.ShelfResponse{
			Error: "shelf id cannot be empty",
		})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	shelf, err := h.shelfService.GetShelfById(r.Context(), shelfId)
	if err != nil {
		if err == wms.ShelfDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfResponse{Error: fmt.Sprintf(
				"failed to get, shelf: %s does not exist",
				shelfId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.ShelfResponse{Error: "Failed to get shelf"})
			return
		}
	}
	h.response(w, http.StatusOK, api.ShelfResponse{Response: shelf})
}

func (h *handler) CreateShelf(w http.ResponseWriter, r *http.Request) {
	var createShelfRequest wms.Shelf

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfResponse{
			Error: err.Error()})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createShelfRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfResponse{
			Error: "Failed to parse request"})
		return
	}

	err = validator.Validate(createShelfRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	shelf := wms.NewShelf(
		createShelfRequest.Label,
		createShelfRequest.Section,
		createShelfRequest.Level,
		createShelfRequest.ShelfBlockId)

	err = h.shelfService.CreateShelf(r.Context(), shelf)
	if err != nil {
		if err == wms.InvalidShelfBlock {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				shelf.ShelfBlockId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfResponse{Error: "Failed to create shelf"},
			)
			return
		}
	}

	h.response(w, http.StatusOK, api.ShelfResponse{Message: fmt.Sprintf(
		"Successfully created shelf: %s",
		shelf.Id,
	)})
}

func (h *handler) UpdateShelf(w http.ResponseWriter, r *http.Request) {
	var updateShelfRequest wms.Shelf

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(
			w,
			http.StatusBadRequest,
			api.ShelfResponse{Error: err.Error()})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&updateShelfRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfResponse{
			Error: "Failed to parse request",
		})
		return
	}

	err = validator.Validate(updateShelfRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ShelfResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err = h.shelfService.UpdateShelf(r.Context(), updateShelfRequest)
	if err != nil {
		if err == wms.ShelfDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfResponse{Error: fmt.Sprintf(
				"failed to update, shelf: %s does not exist",
				updateShelfRequest.Id,
			)})
			return
		} else if err == wms.InvalidShelfBlock {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusBadRequest, api.ShelfResponse{Error: fmt.Sprintf("%s: %s",
				err.Error(),
				updateShelfRequest.ShelfBlockId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfResponse{Error: "Failed to update shelf"},
			)
			return
		}
	}

	h.response(w, http.StatusOK, api.ShelfResponse{Message: fmt.Sprintf(
		"Successfully updated shelf: %s",
		updateShelfRequest.Id,
	)})
}

func (h *handler) DeleteShelf(w http.ResponseWriter, r *http.Request) {
	shelfId := chi.URLParam(r, "shelfId")

	if shelfId == "" {
		err := fmt.Errorf("%v", api.ShelfResponse{Error: "shelf id cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.shelfService.DeleteShelfById(r.Context(), shelfId)
	if err != nil {
		if err == wms.ShelfDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ShelfResponse{Error: fmt.Sprintf(
				"failed to delete, shelf: %s does not exist",
				shelfId,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ShelfResponse{Error: "Failed to delete shelf"},
			)
			return
		}
	}

	h.response(
		w,
		http.StatusOK,
		api.ShelfResponse{Message: fmt.Sprintf(
			"Successfully deleted shelf: %s",
			shelfId,
		)},
	)
}
