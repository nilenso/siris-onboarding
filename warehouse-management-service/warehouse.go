package warehousemanagementservice

import (
	"context"
	"errors"
)

type Warehouse struct {
	Id        string
	Name      string
	Latitude  float64
	Longitude float64
}

var WarehouseDoesNotExist = errors.New("warehouse does not exist")

// mockgen -source="./warehouse.go" -destination="./internal/handler/mock/warehouse.go"
type WarehouseService interface {
	GetWarehouseById(ctx context.Context, id string) (*Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
}

func NewWarehouse(name string, latitude float64, longitude float64) *Warehouse {
	return &Warehouse{
		Id:        generateUUID(),
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
}
