package warehousemanagementservice

import (
	"context"
	"errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Warehouse struct {
	Id        string
	Name      string
	Latitude  float64
	Longitude float64
}

type WarehouseUpdate struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var WarehouseDoesNotExist = errors.New("warehouse does not exist")

// mockgen -source="./warehouse.go" -destination="./internal/handler/mock/warehouse.go"
type WarehouseService interface {
	GetWarehouseById(ctx context.Context, id string) (*Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *Warehouse) error
	DeleteWarehouse(ctx context.Context, id string) error
}

func New(name string, latitude float64, longitude float64) *Warehouse {
	return &Warehouse{
		Id:        generateUUID(),
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (w *Warehouse) SetID() {
	w.Id = generateUUID()
}

func generateUUID() string {
	return uuid.New().String()
}
