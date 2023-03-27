package wms

import (
	"errors"
)

type Warehouse struct {
	Id        string
	Name      string
	Latitude  float64
	Longitude float64
}

var WarehouseDoesNotExist = errors.New("warehouse does not exist")

func NewWarehouse(name string, latitude float64, longitude float64) *Warehouse {
	return &Warehouse{
		Id:        generateUUID(),
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
}
