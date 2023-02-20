package postgres

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"testing"
	wms "warehouse-management-service"
)

func TestGetShelfBlockByIdTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "get_test",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	// insert warehouse and shelf_block records
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

	shelfBlockFromDB, err := shelfBlockService.queries.getShelfBlockByIdTx(
		context.Background(),
		tx,
		shelfBlock.Id,
	)

	if shelfBlockFromDB != shelfBlock {
		t.Errorf("expected: %v, got: %v", shelfBlock, shelfBlockFromDB)
	}
}

func TestGetShelfBlockByIdTxError(t *testing.T) {
	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	_, err = shelfBlockService.queries.getShelfBlockByIdTx(context.Background(), tx, "bad_id")
	if err == nil {
		t.Errorf("expected: %v, got: %v", "error", err)
	}
}

func TestCreateShelfBlockTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "get_test",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	// insert warehouse
	warehouseQuery := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	shelfBlockQuery := "SELECT id, aisle, rack, storage_type, warehouse_id FROM shelf_block WHERE id=$1"
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

	err = shelfBlockService.queries.createShelfBlockTx(context.Background(), tx, shelfBlock)
	if err != nil {
		t.Error(err)
	}

	// check shelf_block in db
	row := tx.QueryRowContext(context.Background(), shelfBlockQuery, shelfBlock.Id)

	var shelfBlockFromDB wms.ShelfBlock
	err = row.Scan(
		&shelfBlockFromDB.Id,
		&shelfBlockFromDB.Aisle,
		&shelfBlockFromDB.Rack,
		&shelfBlockFromDB.StorageType,
		&shelfBlockFromDB.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if shelfBlock != shelfBlockFromDB {
		t.Errorf("expected: %v, got: %v", shelfBlock, shelfBlockFromDB)
	}
}

func TestCreateShelfBlockTxError(t *testing.T) {
	shelfBlock := wms.ShelfBlock{
		Id:          "16acfb14-2486-4f51-8529-aa721ead07cb",
		Aisle:       "2",
		Rack:        "4",
		StorageType: "regular",
		WarehouseId: "c69a7de2-25dc-456d-af3a-c9c4674174b2",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	err = shelfBlockService.queries.createShelfBlockTx(context.Background(), tx, shelfBlock)
	if err == nil {
		t.Errorf("want: %v, got: %v", "error", err)
	}
}

func TestUpdateShelfBlockTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "get_test",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	shelfBlockUpdated := wms.ShelfBlock{
		Id:          "get_test",
		Aisle:       "2",
		Rack:        "3",
		StorageType: "refrigerated",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	// insert warehouse
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

	err = shelfBlockService.queries.updateShelfBlockTx(context.Background(), tx, shelfBlockUpdated)
	if err != nil {
		t.Error(err)
	}

	// check shelf_block in db
	shelfBlockUpdatedQuery := "SELECT id, aisle, rack, storage_type, warehouse_id FROM shelf_block WHERE id=$1"
	row := tx.QueryRowContext(context.Background(), shelfBlockUpdatedQuery, shelfBlock.Id)

	var shelfBlockFromDB wms.ShelfBlock
	err = row.Scan(
		&shelfBlockFromDB.Id,
		&shelfBlockFromDB.Aisle,
		&shelfBlockFromDB.Rack,
		&shelfBlockFromDB.StorageType,
		&shelfBlockFromDB.WarehouseId,
	)
	if err != nil {
		t.Error(err)
		return
	}

	if shelfBlockFromDB != shelfBlockUpdated {
		t.Errorf("want: %v, got: %v", shelfBlockUpdated, shelfBlockFromDB)
	}
}

func TestUpdateShelfBlockTxError(t *testing.T) {
	shelfBlockUpdated := wms.ShelfBlock{
		Id:          "get_test",
		Aisle:       "2",
		Rack:        "3",
		StorageType: "refrigerated",
		WarehouseId: "xx",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	err = shelfBlockService.queries.updateShelfBlockTx(context.Background(), tx, shelfBlockUpdated)
	if err == nil {
		t.Errorf("want: %v, got: %v", "error", err)
	}
}

func TestDeleteShelfBlockTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_get_shelf_block_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	shelfBlock := wms.ShelfBlock{
		Id:          "bb3baaa2-1278-44bf-998f-5e3760b7eb82",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}

	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	// insert warehouse and shelf_block records
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

	err = shelfBlockService.queries.deleteShelfBlockTx(
		context.Background(),
		tx,
		shelfBlock.Id,
	)

	if err != nil {
		t.Error(err)
		return
	}

	// check shelf_block in db
	getShelfBlockQuery := "SELECT id, aisle, rack, storage_type, warehouse_id FROM shelf_block WHERE id=$1"
	err = tx.QueryRowContext(context.Background(), getShelfBlockQuery, shelfBlock.Id).Scan()
	if err != sql.ErrNoRows {
		t.Errorf("expected: %v, got: %v", sql.ErrNoRows, err)
	}
}

func TestDeleteShelfBlockTxError(t *testing.T) {
	tx, err := shelfBlockService.db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	err = shelfBlockService.queries.deleteShelfBlockTx(context.Background(), tx, "non-existent-id")
	if err != RowDoesNotExist {
		t.Error(err)
	}
}

var getShelfBlockTests = []struct {
	getShelfBlockByIdRequest    string
	getShelfBlockByIdTxResponse wms.ShelfBlock
	getShelfBlockByIdTxErr      error
	wantResponse                wms.ShelfBlock
	wantErr                     error
}{
	{
		getShelfBlockByIdRequest: "test_get_by_id",
		getShelfBlockByIdTxResponse: wms.ShelfBlock{
			Id:          "test_get_by_id",
			Aisle:       "2",
			Rack:        "3",
			StorageType: "regular",
			WarehouseId: "foo",
		},
		getShelfBlockByIdTxErr: nil,
		wantResponse: wms.ShelfBlock{
			Id:          "test_get_by_id",
			Aisle:       "2",
			Rack:        "3",
			StorageType: "regular",
			WarehouseId: "foo",
		},
		wantErr: nil,
	},
	{
		getShelfBlockByIdRequest:    "test_get_by_id",
		getShelfBlockByIdTxResponse: wms.ShelfBlock{},
		getShelfBlockByIdTxErr:      sql.ErrConnDone,
		wantResponse:                wms.ShelfBlock{},
		wantErr:                     sql.ErrConnDone,
	},
	{
		getShelfBlockByIdRequest:    "test_get_by_id",
		getShelfBlockByIdTxResponse: wms.ShelfBlock{},
		getShelfBlockByIdTxErr:      sql.ErrNoRows,
		wantResponse:                wms.ShelfBlock{},
		wantErr:                     wms.ShelfBlockDoesNotExist,
	},
	{
		getShelfBlockByIdRequest:    "test_get_by_id",
		getShelfBlockByIdTxResponse: wms.ShelfBlock{},
		getShelfBlockByIdTxErr:      sql.ErrTxDone,
		wantResponse:                wms.ShelfBlock{},
		wantErr:                     sql.ErrTxDone,
	},
	{
		getShelfBlockByIdRequest:    "test_get_by_id",
		getShelfBlockByIdTxResponse: wms.ShelfBlock{},
		getShelfBlockByIdTxErr:      context.Canceled,
		wantResponse:                wms.ShelfBlock{},
		wantErr:                     context.Canceled,
	},
}

func TestGetShelfBlockById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfBlockQueries(mockCtrl)

	ctx := context.Background()

	for _, test := range getShelfBlockTests {
		mockObj.EXPECT().getShelfBlockByIdTx(
			ctx,
			gomock.Any(),
			test.getShelfBlockByIdRequest,
		).Return(test.getShelfBlockByIdTxResponse, test.getShelfBlockByIdTxErr)

		mockShelfBlockService := &ShelfBlockService{
			queries: mockObj,
			db:      shelfBlockService.db,
		}

		response, err := mockShelfBlockService.GetShelfBlockById(ctx, test.getShelfBlockByIdRequest)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}

		if response != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, response)
		}
	}
}

func TestCreateShelfBlockById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := NewMockshelfBlockQueries(mockCtrl)

	ctx := context.Background()

	for _, test := range getShelfBlockTests {
		mockObj.EXPECT().createShelfBlockTx(
			ctx,
			gomock.Any(),
			test.getShelfBlockByIdRequest,
		).Return(test.getShelfBlockByIdTxResponse, test.getShelfBlockByIdTxErr)

		mockShelfBlockService := &ShelfBlockService{
			queries: mockObj,
			db:      shelfBlockService.db,
		}

		response, err := mockShelfBlockService.GetShelfBlockById(ctx, test.getShelfBlockByIdRequest)

		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}

		if response != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, response)
		}
	}
}
