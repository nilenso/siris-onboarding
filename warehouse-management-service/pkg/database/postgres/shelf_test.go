package postgres

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"testing"
	wms "warehouse-management-service"
)

func TestGetShelfByIdTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "863e835b-a05b-4554-b0af-a45389ebbb78",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelf := wms.Shelf{
		Id:           "get_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"
	shelfQuery := "INSERT INTO shelf(id, label, section, level, shelf_block) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(
		context.Background(),
		warehouseQuery,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = tx.ExecContext(
		context.Background(),
		shelfBlockQuery,
		shelfBlock.Id,
		shelfBlock.Aisle,
		shelfBlock.Rack,
		shelfBlock.StorageType,
		shelfBlock.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(context.Background(),
		shelfQuery,
		shelf.Id,
		shelf.Label,
		shelf.Section,
		shelf.Level,
		shelf.ShelfBlockId)
	if err != nil {
		t.Error(err)
		return
	}

	shelfFromDB, err := shelfService.queries.getShelfByIdTx(
		context.Background(),
		tx,
		shelf.Id,
	)
	if err != nil {
		t.Error(err)
	}

	if shelfFromDB != shelf {
		t.Errorf("expected: %v, got: %v", shelf, shelfFromDB)
	}
}

func TestGetShelfByIdTxError(t *testing.T) {
	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = shelfService.queries.getShelfByIdTx(context.Background(), tx, "bad_id")
	if err == nil {
		t.Errorf("expected: %v, got: %v", "error", err)
	}
}

func TestCreateShelfTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "863e835b-a05b-4554-b0af-a45389ebbb78",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelf := wms.Shelf{
		Id:           "get_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	// insert warehouse and shelf_block
	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"
	shelfQuery := "SELECT id, label, section, level, shelf_block FROM shelf WHERE id=$1"
	_, err = tx.ExecContext(
		context.Background(),
		warehouseQuery,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(
		context.Background(),
		shelfBlockQuery,
		shelfBlock.Id,
		shelfBlock.Aisle,
		shelfBlock.Rack,
		shelfBlock.StorageType,
		shelfBlock.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = shelfService.queries.createShelfTx(context.Background(), tx, shelf)
	if err != nil {
		t.Error(err)
	}

	// check shelf in db
	row := tx.QueryRowContext(context.Background(), shelfQuery, shelf.Id)

	var shelfFromDB wms.Shelf
	err = row.Scan(
		&shelfFromDB.Id,
		&shelfFromDB.Label,
		&shelfFromDB.Section,
		&shelfFromDB.Level,
		&shelfFromDB.ShelfBlockId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if shelf != shelfFromDB {
		t.Errorf("expected: %v, got: %v", shelf, shelfFromDB)
	}
}

func TestCreateShelfTxError(t *testing.T) {
	shelf := wms.Shelf{
		Id:           "get_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	err = shelfService.queries.createShelfTx(context.Background(), tx, shelf)
	if err == nil {
		t.Errorf("want: %v, got: %v", "error", err)
	}
}

func TestUpdateShelfTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "863e835b-a05b-4554-b0af-a45389ebbb78",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelf := wms.Shelf{
		Id:           "get_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	shelfUpdated := wms.Shelf{
		Id:           "get_test",
		Label:        "11G",
		Section:      "G",
		Level:        "11",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	// insert warehouse and shelf_block
	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"
	shelfQuery := "INSERT INTO shelf(id, label, section, level, shelf_block) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(
		context.Background(),
		warehouseQuery,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(
		context.Background(),
		shelfBlockQuery,
		shelfBlock.Id,
		shelfBlock.Aisle,
		shelfBlock.Rack,
		shelfBlock.StorageType,
		shelfBlock.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(context.Background(),
		shelfQuery,
		shelf.Id,
		shelf.Label,
		shelf.Section,
		shelf.Level,
		shelf.ShelfBlockId)
	if err != nil {
		t.Error(err)
		return
	}

	err = shelfService.queries.updateShelfTx(context.Background(), tx, shelfUpdated)
	if err != nil {
		t.Error(err)
	}

	// check shelf in db
	shelfUpdatedQuery := "SELECT id, label, section, level, shelf_block FROM shelf WHERE id=$1"
	row := tx.QueryRowContext(context.Background(), shelfUpdatedQuery, shelf.Id)

	var shelfFromDB wms.Shelf
	err = row.Scan(
		&shelfFromDB.Id,
		&shelfFromDB.Label,
		&shelfFromDB.Section,
		&shelfFromDB.Level,
		&shelfFromDB.ShelfBlockId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if shelfFromDB != shelfUpdated {
		t.Errorf("want: %v, got: %v", shelfUpdated, shelfFromDB)
	}
}

func TestUpdateShelfTxError(t *testing.T) {
	shelfUpdated := wms.Shelf{
		Id:           "get_test",
		Label:        "11G",
		Section:      "G",
		Level:        "11",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	err = shelfService.queries.updateShelfTx(context.Background(), tx, shelfUpdated)
	if err == nil {
		t.Errorf("want: %v, got: %v", "error", err)
	}
}

func TestUpdateShelfTxNoRowsError(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "863e835b-a05b-4554-b0af-a45389ebbb78",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelfUpdated := wms.Shelf{
		Id:           "get_test",
		Label:        "11G",
		Section:      "G",
		Level:        "11",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	// insert warehouse and shelf_block
	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(
		context.Background(),
		warehouseQuery,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(
		context.Background(),
		shelfBlockQuery,
		shelfBlock.Id,
		shelfBlock.Aisle,
		shelfBlock.Rack,
		shelfBlock.StorageType,
		shelfBlock.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	err = shelfService.queries.updateShelfTx(context.Background(), tx, shelfUpdated)
	if err != RowDoesNotExist {
		t.Errorf("want: %v, got: %v", RowDoesNotExist, err)
	}
}

func TestDeleteShelfByIdTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "863e835b-a05b-4554-b0af-a45389ebbb78",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelf := wms.Shelf{
		Id:           "delete_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "INSERT INTO shelf_block(id, aisle, rack, storage_type, warehouse_id) VALUES ($1, $2, $3, $4, $5)"
	shelfQuery := "INSERT INTO shelf(id, label, section, level, shelf_block) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(
		context.Background(),
		warehouseQuery,
		warehouse.Id,
		warehouse.Name,
		warehouse.Longitude,
		warehouse.Latitude)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = tx.ExecContext(
		context.Background(),
		shelfBlockQuery,
		shelfBlock.Id,
		shelfBlock.Aisle,
		shelfBlock.Rack,
		shelfBlock.StorageType,
		shelfBlock.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = tx.ExecContext(context.Background(),
		shelfQuery,
		shelf.Id,
		shelf.Label,
		shelf.Section,
		shelf.Level,
		shelf.ShelfBlockId)
	if err != nil {
		t.Error(err)
		return
	}

	err = shelfService.queries.deleteShelfTx(
		context.Background(),
		tx,
		shelf.Id,
	)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteShelfByIdTxError(t *testing.T) {
	tx, err := shelfService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Rollback() }()

	err = shelfService.queries.deleteShelfTx(context.Background(), tx, "non-existent-id")
	if err != RowDoesNotExist {
		t.Error(err)
	}
}

func TestGetShelfById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfQueries(mockCtrl)

	ctx := context.Background()

	tests := []struct {
		getShelfByIdRequest    string
		getShelfByIdTxResponse wms.Shelf
		getShelfByIdTxErr      error
		wantResponse           wms.Shelf
		wantErr                error
	}{
		{
			getShelfByIdRequest: "test_get_by_id",
			getShelfByIdTxResponse: wms.Shelf{
				Id:           "test_get_by_id",
				Label:        "12A",
				Section:      "A",
				Level:        "12",
				ShelfBlockId: "foo",
			},
			getShelfByIdTxErr: nil,
			wantResponse: wms.Shelf{
				Id:           "test_get_by_id",
				Label:        "12A",
				Section:      "A",
				Level:        "12",
				ShelfBlockId: "foo",
			},
			wantErr: nil,
		},
		{
			getShelfByIdRequest:    "test_get_by_id",
			getShelfByIdTxResponse: wms.Shelf{},
			getShelfByIdTxErr:      sql.ErrConnDone,
			wantResponse:           wms.Shelf{},
			wantErr:                sql.ErrConnDone,
		},
		{
			getShelfByIdRequest:    "test_get_by_id",
			getShelfByIdTxResponse: wms.Shelf{},
			getShelfByIdTxErr:      sql.ErrNoRows,
			wantResponse:           wms.Shelf{},
			wantErr:                wms.ShelfDoesNotExist,
		},
		{
			getShelfByIdRequest:    "test_get_by_id",
			getShelfByIdTxResponse: wms.Shelf{},
			getShelfByIdTxErr:      sql.ErrTxDone,
			wantResponse:           wms.Shelf{},
			wantErr:                sql.ErrTxDone,
		},
		{
			getShelfByIdRequest:    "test_get_by_id",
			getShelfByIdTxResponse: wms.Shelf{},
			getShelfByIdTxErr:      context.Canceled,
			wantResponse:           wms.Shelf{},
			wantErr:                context.Canceled,
		},
	}

	for _, test := range tests {
		mockObj.EXPECT().getShelfByIdTx(
			ctx,
			gomock.Any(),
			test.getShelfByIdRequest,
		).Return(test.getShelfByIdTxResponse, test.getShelfByIdTxErr)

		mockShelfService := &ShelfService{
			queries: mockObj,
			db:      shelfService.db,
		}

		response, err := mockShelfService.GetShelfById(ctx, test.getShelfByIdRequest)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}

		if response != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, response)
		}
	}
}

func TestCreateShelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfQueries(mockCtrl)

	ctx := context.Background()

	request := wms.Shelf{
		Id:           "test_create",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "foo",
	}

	tests := []struct {
		createShelfError error
		wantErr          error
	}{
		{
			createShelfError: nil,
			wantErr:          nil,
		},
		{
			createShelfError: sql.ErrConnDone,
			wantErr:          sql.ErrConnDone,
		},
		{
			createShelfError: InvalidShelfBlock,
			wantErr:          wms.InvalidShelfBlock,
		},
		{
			createShelfError: sql.ErrTxDone,
			wantErr:          sql.ErrTxDone,
		},
		{
			createShelfError: context.Canceled,
			wantErr:          context.Canceled,
		},
	}
	for _, test := range tests {
		mockObj.EXPECT().createShelfTx(
			ctx,
			gomock.Any(),
			request,
		).Return(test.createShelfError)

		mockShelfService := &ShelfService{
			queries: mockObj,
			db:      shelfService.db,
		}

		err := mockShelfService.CreateShelf(ctx, request)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}
	}
}

func TestUpdateShelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfQueries(mockCtrl)

	ctx := context.Background()

	request := wms.Shelf{
		Id:           "test_create",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "foo",
	}

	tests := []struct {
		updateShelfError error
		wantErr          error
	}{
		{
			updateShelfError: nil,
			wantErr:          nil,
		},
		{
			updateShelfError: sql.ErrConnDone,
			wantErr:          sql.ErrConnDone,
		},
		{
			updateShelfError: InvalidShelfBlock,
			wantErr:          wms.InvalidShelfBlock,
		},
		{
			updateShelfError: RowDoesNotExist,
			wantErr:          wms.ShelfDoesNotExist,
		},
		{
			updateShelfError: sql.ErrTxDone,
			wantErr:          sql.ErrTxDone,
		},
		{
			updateShelfError: context.Canceled,
			wantErr:          context.Canceled,
		},
	}
	for _, test := range tests {
		mockObj.EXPECT().updateShelfTx(
			ctx,
			gomock.Any(),
			request,
		).Return(test.updateShelfError)

		mockShelfService := &ShelfService{
			queries: mockObj,
			db:      shelfService.db,
		}

		err := mockShelfService.UpdateShelf(ctx, request)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}
	}
}

func TestDeleteShelfById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfQueries(mockCtrl)

	ctx := context.Background()

	request := "test_delete"
	tests := []struct {
		deleteShelfErr error
		wantErr        error
	}{
		{
			deleteShelfErr: nil,
			wantErr:        nil,
		},
		{
			deleteShelfErr: sql.ErrConnDone,
			wantErr:        sql.ErrConnDone,
		},
		{
			deleteShelfErr: RowDoesNotExist,
			wantErr:        wms.ShelfDoesNotExist,
		},
		{
			deleteShelfErr: sql.ErrTxDone,
			wantErr:        sql.ErrTxDone,
		},
		{
			deleteShelfErr: context.Canceled,
			wantErr:        context.Canceled,
		},
	}
	for _, test := range tests {
		mockObj.EXPECT().deleteShelfTx(
			ctx,
			gomock.Any(),
			request,
		).Return(test.deleteShelfErr)

		mockShelfService := &ShelfService{
			queries: mockObj,
			db:      shelfService.db,
		}

		err := mockShelfService.DeleteShelfById(ctx, request)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}
	}
}
