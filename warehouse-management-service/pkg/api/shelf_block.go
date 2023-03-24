package api

import (
	"fmt"
	"strings"
	wms "warehouse-management-service"
)

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

func (c *CreateShelfBlockRequest) IsValid() (error, bool) {
	var errors []string

	if c.Aisle == "" {
		errors = append(errors, "aisle cannot be empty")
	}
	if c.Rack == "" {
		errors = append(errors, "rack cannot be empty")
	}
	if c.StorageType == "" {
		errors = append(errors, "storageType cannot be empty")
	}
	if c.WarehouseId == "" {
		errors = append(errors, "warehouseId cannot be empty")
	}

	if len(errors) != 0 {
		err := fmt.Errorf(strings.Join(errors, ", "))
		return err, false
	}

	return nil, true
}

func (u *UpdateShelfBlockRequest) IsValid() (error, bool) {
	var errors []string

	if u.Id == "" {
		errors = append(errors, "id cannot be empty")
	}
	if u.Aisle == "" {
		errors = append(errors, "aisle cannot be empty")
	}
	if u.Rack == "" {
		errors = append(errors, "rack cannot be empty")
	}
	if u.StorageType == "" {
		errors = append(errors, "storageType cannot be empty")
	}
	if u.WarehouseId == "" {
		errors = append(errors, "warehouseId cannot be empty")
	}

	if len(errors) != 0 {
		err := fmt.Errorf(strings.Join(errors, ", "))
		return err, false
	}

	return nil, true
}
