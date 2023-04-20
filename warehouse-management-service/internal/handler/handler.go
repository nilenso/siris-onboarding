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

// mockgen -destination="./internal/handler/mock/shelf_block.go" warehouse-management-service/internal/handler ShelfBlockService
type ShelfBlockService interface {
	GetShelfBlockById(ctx context.Context, id string) (wms.ShelfBlock, error)
	CreateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	UpdateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error
	DeleteShelfBlockById(ctx context.Context, id string) error
}

// mockgen -destination="./internal/handler/mock/shelf.go" warehouse-management-service/internal/handler ShelfService
type ShelfService interface {
	GetShelfById(ctx context.Context, id string) (wms.Shelf, error)
	CreateShelf(ctx context.Context, shelf wms.Shelf) error
	UpdateShelf(ctx context.Context, shelf wms.Shelf) error
	DeleteShelfById(ctx context.Context, id string) error
}

// mockgen -destination="./internal/handler/mock/product.go" warehouse-management-service/internal/handler ProductService
type ProductService interface {
	GetProductById(ctx context.Context, id string) (wms.Product, error)
	CreateProduct(ctx context.Context, product wms.Product) error
	UpdateProduct(ctx context.Context, product wms.Product) error
	DeleteProductById(ctx context.Context, id string) error
}

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

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")

	if sku == "" {
		err := fmt.Errorf("%v", api.WarehouseResponse{
			Error: "sku cannot be empty",
		})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.productService.GetProductById(r.Context(), sku)
	switch err {
	case wms.ProductDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ProductResponse{Error: fmt.Sprintf(
				"failed to get, product: %s does not exist",
				sku,
			)})
		}
	case nil:
		{
			h.response(w, http.StatusOK, api.ProductResponse{Response: product})
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.ProductResponse{Error: "Failed to get product"})
		}
	}
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var createProductRequest wms.Product

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: err.Error(),
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&createProductRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: "Failed to parse request",
		})
		return
	}

	err = validator.Validate(createProductRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err = h.productService.CreateProduct(r.Context(), createProductRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusInternalServerError, api.ProductResponse{Error: "Failed to create product"})
		return
	}

	h.response(w, http.StatusOK, api.WarehouseResponse{Response: fmt.Sprintf(
		"Successfully created product: %s",
		createProductRequest.SKU,
	)})
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var updateProductRequest wms.Product

	if r.Body == nil {
		err := fmt.Errorf("request body cannot be empty")
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: err.Error(),
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&updateProductRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: "Failed to parse request"})
		return
	}

	err = validator.Validate(updateProductRequest)
	if err != nil {
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, api.ProductResponse{
			Error: fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	err = h.productService.UpdateProduct(r.Context(), updateProductRequest)
	switch err {
	case wms.ProductDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ProductResponse{Error: fmt.Sprintf(
				"failed to update, product: %s does not exist",
				updateProductRequest.SKU,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				api.ProductResponse{
					Message: fmt.Sprintf(
						"Successfully updated product: %s",
						updateProductRequest.SKU,
					)},
			)
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ProductResponse{Error: "Failed to update product"},
			)
		}
	}
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")

	if sku == "" {
		err := fmt.Errorf("%v", map[string]string{"error": "sku cannot be empty"})
		h.logger.Log(log.Error, err)
		h.response(w, http.StatusBadRequest, err)
		return
	}

	err := h.productService.DeleteProductById(r.Context(), sku)
	switch err {
	case wms.ProductDoesNotExist:
		{
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ProductResponse{Error: fmt.Sprintf(
				"failed to delete, product: %s does not exist",
				sku,
			)})
		}
	case nil:
		{
			h.response(
				w,
				http.StatusOK,
				api.ProductResponse{Message: fmt.Sprintf(
					"Successfully deleted product: %s",
					sku,
				)},
			)
		}
	default:
		{
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ProductResponse{Error: "Failed to delete product"},
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
