package api

import (
	warehousemanagementservice "warehouse-management-service"
)

type CreateWarehouseRequest struct {
	Name      string  `json:"name" validate:"nonzero"`
	Latitude  float64 `json:"latitude" validate:"min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"min=-180,max=180"`
}

type GetWarehouseResponse struct {
	Response warehousemanagementservice.Warehouse `json:"response,omitempty"`
	Error    string                               `json:"error,omitempty"`
}

type WarehouseResponse struct {
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

type UpdateWarehouseRequest struct {
	Id        string  `json:"id" validate:"nonzero"`
	Name      string  `json:"name" validate:"nonzero"`
	Latitude  float64 `json:"latitude" validate:"min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"min=-180,max=180"`
}
