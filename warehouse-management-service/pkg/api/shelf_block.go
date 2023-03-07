package api

import (
	wms "warehouse-management-service"
)

type CreateShelfBlockRequest struct {
	Aisle       string `json:"aisle" validate:"nonzero"`
	Rack        string `json:"rack" validate:"nonzero"`
	StorageType string `json:"storageType" validate:"nonzero"`
	WarehouseId string `json:"warehouseId" validate:"nonzero"`
}

type UpdateShelfBlockRequest struct {
	Id          string `json:"id" validate:"nonzero"`
	Aisle       string `json:"aisle" validate:"nonzero"`
	Rack        string `json:"rack" validate:"nonzero"`
	StorageType string `json:"storageType" validate:"nonzero"`
	WarehouseId string `json:"warehouseId" validate:"nonzero"`
}

type GetShelfBlockResponse struct {
	Response wms.ShelfBlock `json:"response,omitempty"`
	Error    string         `json:"error,omitempty"`
}

type ShelfBlockResponse struct {
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}
