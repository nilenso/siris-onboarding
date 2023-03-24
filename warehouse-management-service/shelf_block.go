package warehousemanagementservice

import (
	"context"
	"errors"
)

type ShelfBlock struct {
	Id          string
	Aisle       string
	Rack        string
	StorageType string
	WarehouseId string
}

var ShelfBlockDoesNotExist = errors.New("shelf block does not exist")
var InvalidWarehouse = errors.New("invalid warehouse")

// mockgen -source="./shelf_block.go" -destination="./internal/handler/mock/shelf_block.go"
type ShelfService interface {
	GetShelfBlockById(ctx context.Context, id string) (ShelfBlock, error)
	CreateShelfBlock(ctx context.Context, shelfBlock ShelfBlock) error
	UpdateShelfBlock(ctx context.Context, shelfBlock ShelfBlock) error
	DeleteShelfBlockById(ctx context.Context, id string) error
}

func NewShelfBlock(aisle, rack, storageType, warehouseId string) ShelfBlock {
	return ShelfBlock{
		Id:          generateUUID(),
		Aisle:       aisle,
		Rack:        rack,
		StorageType: storageType,
		WarehouseId: warehouseId,
	}
}
