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

// mockgen -destination="./internal/handler/mock/product.go" warehouse-management-service/internal/handler ProductService
type ProductService interface {
	GetProductById(ctx context.Context, id string) (wms.Product, error)
	CreateProduct(ctx context.Context, product wms.Product) error
	UpdateProduct(ctx context.Context, product wms.Product) error
	DeleteProductById(ctx context.Context, id string) error
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
	if err != nil {
		if err == wms.ProductDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ProductResponse{Error: fmt.Sprintf(
				"failed to get, product: %s does not exist",
				sku,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusInternalServerError, api.ProductResponse{Error: "Failed to get product"})
			return
		}
	}
	h.response(w, http.StatusOK, api.ProductResponse{Response: product})
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
	if err != nil {
		if err == wms.ProductDoesNotExist {
			h.logger.Log(log.Error, err)
			h.response(w, http.StatusNotFound, api.ProductResponse{Error: fmt.Sprintf(
				"failed to update, product: %s does not exist",
				updateProductRequest.SKU,
			)})
			return
		} else {
			h.logger.Log(log.Error, err)
			h.response(
				w,
				http.StatusInternalServerError,
				api.ProductResponse{Error: "Failed to update product"},
			)
			return
		}
	}

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
