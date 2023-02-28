package postgres

import (
	"context"
	"database/sql"
	"errors"
	wms "warehouse-management-service"
)

type shelfBlockQueries interface {
	createShelfBlockTx(ctx context.Context, tx *sql.Tx, block wms.ShelfBlock) error
	getShelfBlockByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.ShelfBlock, error)
	updateShelfBlockTx(ctx context.Context, tx *sql.Tx, block wms.ShelfBlock) error
	deleteShelfBlockTx(ctx context.Context, tx *sql.Tx, id string) error
	warehouseExistsTx(ctx context.Context, tx *sql.Tx, id string) (bool, error)
}

type shelfBlockQueriesImpl struct{}

type ShelfBlockService struct {
	queries shelfBlockQueries
	db      *sql.DB
}

var InvalidWarehouse = errors.New("invalid warehouseId")

func NewShelfService(db *sql.DB) *ShelfBlockService {
	return &ShelfBlockService{
		queries: new(shelfBlockQueriesImpl),
		db:      db,
	}
}

func (s *ShelfBlockService) GetShelfBlockById(ctx context.Context, id string) (wms.ShelfBlock, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return wms.ShelfBlock{}, err
	}

	shelfBlock, err := s.queries.getShelfBlockByIdTx(ctx, tx, id)
	switch err {
	case nil:
		return shelfBlock, tx.Commit()
	case sql.ErrNoRows:
		return wms.ShelfBlock{}, wms.ShelfBlockDoesNotExist
	default:
		return wms.ShelfBlock{}, err
	}
}

func (s *ShelfBlockService) CreateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.createShelfBlockTx(ctx, tx, shelfBlock)
	switch err {
	case InvalidWarehouse:
		return wms.InvalidWarehouse
	case nil:
		return tx.Commit()
	default:
		return err
	}
}

func (s *ShelfBlockService) UpdateShelfBlock(ctx context.Context, shelfBlock wms.ShelfBlock) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.updateShelfBlockTx(ctx, tx, shelfBlock)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.ShelfBlockDoesNotExist
	case InvalidWarehouse:
		return wms.InvalidWarehouse
	default:
		return err
	}
}

func (s *ShelfBlockService) DeleteShelfBlockById(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.deleteShelfBlockTx(ctx, tx, id)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.ShelfBlockDoesNotExist
	default:
		return err
	}
}

func (s *shelfBlockQueriesImpl) createShelfBlockTx(ctx context.Context, tx *sql.Tx, block wms.ShelfBlock) error {
	if warehouseExists, err := s.warehouseExistsTx(ctx, tx, block.WarehouseId);
		err != nil {
		return err
	} else if !warehouseExists {
		return InvalidWarehouse
	}

	query := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"

	_, err := tx.ExecContext(
		ctx,
		query,
		block.Id,
		block.Aisle,
		block.Rack,
		block.StorageType,
		block.WarehouseId,
	)
	return err
}

func (s *shelfBlockQueriesImpl) getShelfBlockByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.ShelfBlock, error) {
	row := tx.QueryRowContext(ctx, `SELECT id, aisle, rack, storage_type, warehouse_id FROM shelf_block WHERE id=$1`, id)

	var shelfBlock wms.ShelfBlock

	err := row.Scan(&shelfBlock.Id, &shelfBlock.Aisle, &shelfBlock.Rack, &shelfBlock.StorageType, &shelfBlock.WarehouseId)
	if err != nil {
		return wms.ShelfBlock{}, err
	}
	return shelfBlock, nil
}

func (s *shelfBlockQueriesImpl) updateShelfBlockTx(ctx context.Context, tx *sql.Tx, block wms.ShelfBlock) error {
	if warehouseExists, err := s.warehouseExistsTx(ctx, tx, block.WarehouseId);
		err != nil {
		return err
	} else if !warehouseExists {
		return InvalidWarehouse
	}

	query := `UPDATE shelf_block SET aisle = $1, rack = $2, storage_type=$3, warehouse_id=$4 WHERE id = $5;`

	result, err := tx.ExecContext(ctx,
		query,
		block.Aisle,
		block.Rack,
		block.StorageType,
		block.WarehouseId,
		block.Id)
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

func (s *shelfBlockQueriesImpl) deleteShelfBlockTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM shelf_block WHERE id=$1`

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

func (s *shelfBlockQueriesImpl) warehouseExistsTx(ctx context.Context, tx *sql.Tx, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM warehouse WHERE id = $1)`

	var exists bool
	row := tx.QueryRowContext(ctx, query, id)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
