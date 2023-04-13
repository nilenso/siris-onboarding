package postgres

import (
	"context"
	"database/sql"
	"errors"
	wms "warehouse-management-service"
)

type shelfQueries interface {
	createShelfTx(ctx context.Context, tx *sql.Tx, shelf wms.Shelf) error
	getShelfByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.Shelf, error)
	updateShelfTx(ctx context.Context, tx *sql.Tx, shelf wms.Shelf) error
	deleteShelfTx(ctx context.Context, tx *sql.Tx, id string) error
	shelfBlockExistsTx(ctx context.Context, tx *sql.Tx, id string) (bool, error)
}

type shelfQueriesImpl struct{}

type ShelfService struct {
	queries shelfQueries
	db      *sql.DB
}

func NewShelfService(db *sql.DB) *ShelfService {
	return &ShelfService{
		queries: new(shelfQueriesImpl),
		db:      db,
	}
}

var InvalidShelfBlock = errors.New("invalid shelfBlockId")

func (s *ShelfService) GetShelfById(ctx context.Context, id string) (wms.Shelf, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return wms.Shelf{}, err
	}

	shelf, err := s.queries.getShelfByIdTx(ctx, tx, id)
	switch err {
	case nil:
		return shelf, tx.Commit()
	case sql.ErrNoRows:
		return wms.Shelf{}, wms.ShelfDoesNotExist
	default:
		return wms.Shelf{}, err
	}
}

func (s *ShelfService) CreateShelf(ctx context.Context, shelf wms.Shelf) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.createShelfTx(ctx, tx, shelf)
	switch err {
	case InvalidShelfBlock:
		return wms.InvalidShelfBlock
	case nil:
		return tx.Commit()
	default:
		return err
	}
}

func (s *ShelfService) UpdateShelf(ctx context.Context, shelf wms.Shelf) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.updateShelfTx(ctx, tx, shelf)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.ShelfDoesNotExist
	case InvalidShelfBlock:
		return wms.InvalidShelfBlock
	default:
		return err
	}
}

func (s *ShelfService) DeleteShelfById(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.queries.deleteShelfTx(ctx, tx, id)
	switch err {
	case nil:
		return tx.Commit()
	case RowDoesNotExist:
		return wms.ShelfDoesNotExist
	default:
		return err
	}
}

func (s *shelfQueriesImpl) createShelfTx(ctx context.Context, tx *sql.Tx, shelf wms.Shelf) error {
	if shelfBlockExists, err := s.shelfBlockExistsTx(ctx, tx, shelf.ShelfBlockId); err != nil {
		return err
	} else if !shelfBlockExists {
		return InvalidShelfBlock
	}

	query := "INSERT INTO shelf(id, label, section, level, shelf_block) VALUES ($1, $2, $3, $4, $5)"

	_, err := tx.ExecContext(
		ctx,
		query,
		shelf.Id,
		shelf.Label,
		shelf.Section,
		shelf.Level,
		shelf.ShelfBlockId,
	)
	return err
}

func (s *shelfQueriesImpl) getShelfByIdTx(ctx context.Context, tx *sql.Tx, id string) (wms.Shelf, error) {
	row := tx.QueryRowContext(ctx, `SELECT id, label, section, level, shelf_block FROM shelf WHERE id=$1`, id)

	var shelf wms.Shelf

	err := row.Scan(&shelf.Id, &shelf.Label, &shelf.Section, &shelf.Level, &shelf.ShelfBlockId)
	if err != nil {
		return wms.Shelf{}, err
	}
	return shelf, nil
}

func (s *shelfQueriesImpl) updateShelfTx(ctx context.Context, tx *sql.Tx, shelf wms.Shelf) error {
	if shelfBlockExists, err := s.shelfBlockExistsTx(ctx, tx, shelf.ShelfBlockId); err != nil {
		return err
	} else if !shelfBlockExists {
		return InvalidShelfBlock
	}

	query := `UPDATE shelf SET label = $1, section = $2, level = $3, shelf_block = $4 where id = $5`

	result, err := tx.ExecContext(
		ctx,
		query,
		shelf.Label,
		shelf.Section,
		shelf.Level,
		shelf.ShelfBlockId,
		shelf.Id)
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

func (s *shelfQueriesImpl) deleteShelfTx(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM shelf WHERE id=$1`

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

func (s *shelfQueriesImpl) shelfBlockExistsTx(ctx context.Context, tx *sql.Tx, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shelf_block WHERE id = $1)`

	var exists bool
	row := tx.QueryRowContext(ctx, query, id)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
