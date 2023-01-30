package wms

import (
	"fmt"
	"strings"
	warehousemanagementservice "warehouse-management-service"
)

type CreateWarehouseRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *CreateWarehouseRequest) IsValid() (error, bool) {
	var errors []string

	if !validateName(c.Name) {
		errors = append(errors, "name cannot be empty")
	}
	if !validateLatitude(c.Latitude) {
		errors = append(errors, "latitude has to be in the range [-90, 90]")
	}
	if !validateLongitude(c.Longitude) {
		errors = append(errors, "longitude has to be in the range [-180, 180]")
	}

	if len(errors) != 0 {
		err := fmt.Errorf(strings.Join(errors, ", "))
		return err, false
	}

	return nil, true
}

func (u *UpdateWarehouseRequest) IsValid() (error, bool) {
	var errors []string

	if u.Id == "" {
		errors = append(errors, "id cannot be empty")
	}
	if !validateName(u.Name) {
		errors = append(errors, "name cannot be empty")
	}
	if !validateLatitude(u.Latitude) {
		errors = append(errors, "latitude has to be in the range [-90, 90]")
	}
	if !validateLongitude(u.Longitude) {
		errors = append(errors, "longitude has to be in the range [-180, 180]")
	}

	if len(errors) != 0 {
		err := fmt.Errorf(strings.Join(errors, ", "))
		return err, false
	}

	return nil, true
}

func validateName(name string) bool {
	return name != ""
}

func validateLatitude(latitude float64) bool {
	return latitude > -90 && latitude < 90
}

func validateLongitude(longitude float64) bool {
	return longitude > -180 && longitude < 180
}
