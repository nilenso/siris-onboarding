package wms

import (
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

func NewShelfBlock(aisle, rack, storageType, warehouseId string) ShelfBlock {
	return ShelfBlock{
		Id:          generateUUID(),
		Aisle:       aisle,
		Rack:        rack,
		StorageType: storageType,
		WarehouseId: warehouseId,
	}
}
