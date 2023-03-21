package warehousemanagementservice

import (
	"context"
	"errors"
)

type Product struct {
	SKU        string  `json:"sku,omitempty"`
	Name       string  `json:"name,omitempty"`
	MRP        float64 `json:"mrp,omitempty"`
	Variant    string  `json:"variant,omitempty"`
	LengthInCm float64 `json:"lengthInCm,omitempty"`
	WidthInCm  float64 `json:"widthInCm,omitempty"`
	HeightInCm float64 `json:"heightInCm,omitempty"`
	WeightInKg float64 `json:"weightInKg,omitempty"`
	Perishable bool    `json:"perishable,omitempty"`
}

var ProductDoesNotExist = errors.New("shelf does not exist")

type ProductService interface {
	GetProductById(ctx context.Context, id string) (Product, error)
	CreateProduct(ctx context.Context, product Product) error
	UpdateProduct(ctx context.Context, product Product) error
	DeleteProductById(ctx context.Context, id string) error
}

func NewProduct(
	sku string,
	name string,
	mrp float64,
	variant string,
	lengthInCm float64,
	widthInCm float64,
	breadthInCm float64,
	weightInKg float64,
	perishable bool,
) Product {
	return Product{
		SKU:        sku,
		Name:       name,
		MRP:        mrp,
		Variant:    variant,
		LengthInCm: lengthInCm,
		WidthInCm:  widthInCm,
		HeightInCm: breadthInCm,
		WeightInKg: weightInKg,
		Perishable: perishable,
	}
}
