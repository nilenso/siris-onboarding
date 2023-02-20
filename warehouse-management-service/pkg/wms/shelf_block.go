package wms

import wms "warehouse-management-service"

type CreateShelfBlockRequest struct {
	Aisle       string `json:"aisle"`
	Rack        string `json:"rack"`
	StorageType string `json:"storageType"`
	WarehouseId string `json:"warehouseId"`
}

type UpdateShelfBlockRequest struct {
	Id          string `json:"id"`
	Aisle       string `json:"aisle"`
	Rack        string `json:"rack"`
	StorageType string `json:"storageType"`
	WarehouseId string `json:"warehouseId"`
}

type GetShelfBlockResponse struct {
	Response wms.ShelfBlock `json:"response,omitempty"`
	Error    string         `json:"error,omitempty"`
}

type ShelfBlockResponse struct {
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}
