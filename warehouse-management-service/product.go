package warehousemanagementservice

import "context"

type Product struct {
	SKU         string  `json:"sku,omitempty"`
	Name        string  `json:"name,omitempty"`
	MRP         float64 `json:"mrp,omitempty"`
	Variant     string  `json:"variant,omitempty"`
	LengthInCm  float64 `json:"lengthInCm,omitempty"`
	WidthInCm   float64 `json:"widthInCm,omitempty"`
	BreadthInCm float64 `json:"breadthInCm,omitempty"`
	WeightInKg  float64 `json:"weightInKg,omitempty"`
	Perishable  bool    `json:"perishable,omitempty"`
}

type ProductService interface {
	GetProductById(ctx context.Context, id string)
	CreateProduct(ctx context.Context, product Product)
	UpdateProduct(ctx context.Context, product Product)
	DeleteProductById(ctx context.Context, id string)
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
		SKU:         sku,
		Name:        name,
		MRP:         mrp,
		Variant:     variant,
		LengthInCm:  lengthInCm,
		WidthInCm:   widthInCm,
		BreadthInCm: breadthInCm,
		WeightInKg:  weightInKg,
		Perishable:  perishable,
	}
}
