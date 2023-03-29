package postgres

import (
	"context"
	"database/sql"
	wms "warehouse-management-service"
)

type productQueries interface {
	createProductTx(ctx context.Context, tx *sql.Tx, Product wms.Product) error
	getProductByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.Product, error)
	updateProductTx(ctx context.Context, tx *sql.Tx, Product wms.Product) error
	deleteProductTx(ctx context.Context, tx *sql.Tx, id string) error
}

type productQueriesImpl struct{}

type ProductService struct {
	queries productQueries
	db      *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{
		queries: new(productQueriesImpl),
		db:      db,
	}
}

func (p *ProductService) GetProductById(ctx context.Context, id string) (wms.Product, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return wms.Product{}, err
	}
	defer func() { _ = tx.Rollback() }()

	product, err := p.queries.getProductByIdTx(ctx, tx, id)
	switch err {
	case nil:
		return product, tx.Commit()
	case sql.ErrNoRows:
		return wms.Product{}, wms.ProductDoesNotExist
	default:
		return wms.Product{}, err
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, product wms.Product) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = p.queries.createProductTx(ctx, tx, product)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (p *ProductService) UpdateProduct(ctx context.Context, product wms.Product) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = p.queries.updateProductTx(ctx, tx, product)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (p *ProductService) DeleteProductById(ctx context.Context, id string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = p.queries.deleteProductTx(ctx, tx, id)
	switch err {
	case nil:
		return tx.Commit()
	case sql.ErrNoRows:
		return wms.ProductDoesNotExist
	default:
		return err
	}
}

func (p *productQueriesImpl) createProductTx(ctx context.Context, tx *sql.Tx, product wms.Product) error {
	query := `INSERT INTO product(sku, name, mrp, variant, length_in_cm, 
		width_in_cm, height_in_cm, weight_in_kg, perishable)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := tx.ExecContext(
		ctx,
		query,
		product.SKU,
		product.Name,
		product.MRP,
		product.Variant,
		product.LengthInCm,
		product.WidthInCm,
		product.HeightInCm,
		product.WeightInKg,
		product.Perishable,
	)
	return err
}

func (p *productQueriesImpl) getProductByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.Product, error) {
	row := tx.QueryRowContext(ctx, `SELECT sku, name, mrp, variant, length_in_cm, 
       width_in_cm, height_in_cm, weight_in_kg, perishable FROM product WHERE sku=$1`, id)

	var product wms.Product

	err := row.Scan(&product.SKU,
		&product.Name,
		&product.MRP,
		&product.Variant,
		&product.LengthInCm,
		&product.WidthInCm,
		&product.HeightInCm,
		&product.WeightInKg,
		&product.Perishable)
	if err != nil {
		return wms.Product{}, err
	}
	return product, nil
}

func (p *productQueriesImpl) updateProductTx(ctx context.Context, tx *sql.Tx, product wms.Product) error {
	query := `UPDATE product set name = $2, mrp = $3, variant = $4, length_in_cm = $5, width_in_cm = $6, height_in_cm = $7, weight_in_kg = $8, perishable = $9 where sku = $1`

	result, err := tx.ExecContext(
		ctx,
		query,
		product.SKU,
		product.Name,
		product.MRP,
		product.Variant,
		product.LengthInCm,
		product.WidthInCm,
		product.HeightInCm,
		product.WeightInKg,
		product.Perishable,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return RowDoesNotExist
	}
	return nil
}

func (p *productQueriesImpl) deleteProductTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM product WHERE id=$1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return RowDoesNotExist
	}

	return nil
}
