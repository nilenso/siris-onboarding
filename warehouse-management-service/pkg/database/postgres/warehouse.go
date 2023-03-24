package postgres

import (
	"context"
	"database/sql"
	"errors"
	wms "warehouse-management-service"

	_ "github.com/lib/pq"
)

// mockgen -source="./pkg/database/postgres/warehouse.go" -destination="./pkg/database/postgres/warehouse_mock.go"
type warehouseQueries interface {
	createWarehouseTx(ctx context.Context, tx *sql.Tx, warehouse *wms.Warehouse) error
	getWarehouseByIdTx(ctx context.Context, tx *sql.Tx, id string) (*wms.Warehouse, error)
	updateWarehouseTx(ctx context.Context, tx *sql.Tx, warehouse *wms.Warehouse) error
	deleteWarehouseTx(ctx context.Context, tx *sql.Tx, id string) error
}

type queriesImpl struct{}

type WarehouseService struct {
	queries warehouseQueries
	db      *sql.DB
}

var RowDoesNotExist = errors.New("postgres: queried row does not exist")

func NewWarehouseService(db *sql.DB) *WarehouseService {
	return &WarehouseService{
		queries: new(queriesImpl),
		db:      db,
	}
}

func (w *WarehouseService) GetWarehouseById(ctx context.Context, id string) (*wms.Warehouse, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	warehouse, err := w.queries.getWarehouseByIdTx(ctx, tx, id)
	switch err {
	case nil:
		return warehouse, tx.Commit()
	case sql.ErrNoRows:
		return nil, wms.WarehouseDoesNotExist
	default:
		return nil, err
	}
}

func (w *WarehouseService) CreateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = w.queries.createWarehouseTx(ctx, tx, warehouse)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (w *WarehouseService) UpdateWarehouse(ctx context.Context, warehouse *wms.Warehouse) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = w.queries.updateWarehouseTx(ctx, tx, warehouse)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.WarehouseDoesNotExist
	default:
		return err
	}
}

func (w *WarehouseService) DeleteWarehouse(ctx context.Context, id string) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = w.queries.deleteWarehouseTx(ctx, tx, id)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.WarehouseDoesNotExist
	default:
		return err
	}
}

func (q *queriesImpl) getWarehouseByIdTx(ctx context.Context, tx *sql.Tx, id string) (*wms.Warehouse, error) {
	row := tx.QueryRowContext(ctx, `SELECT id, name, geolocation[0], geolocation[1] FROM warehouse WHERE id=$1`, id)

	var warehouse wms.Warehouse

	err := row.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Longitude, &warehouse.Latitude)
	if err != nil {
		return nil, err
	}

	return &warehouse, nil
}

func (q *queriesImpl) createWarehouseTx(ctx context.Context, tx *sql.Tx, warehouse *wms.Warehouse) error {
	query := `INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4));`

	_, err := tx.ExecContext(ctx,
		query,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude)
	if err != nil {
		return err
	}

	return nil
}

func (q *queriesImpl) updateWarehouseTx(ctx context.Context, tx *sql.Tx, warehouse *wms.Warehouse) error {
	query := `UPDATE warehouse SET name = $1, geolocation = point($2, $3) WHERE id = $4;`

	result, err := tx.ExecContext(ctx,
		query,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude,
		warehouse.Id)
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

func (q *queriesImpl) deleteWarehouseTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM warehouse WHERE id=$1`

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
