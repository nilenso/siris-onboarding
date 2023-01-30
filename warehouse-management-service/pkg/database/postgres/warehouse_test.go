package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	wms "warehouse-management-service"
)

var genericError = errors.New("generic error")
var warehouse = wms.Warehouse{
	Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	Name:      "test_find_by_id",
	Latitude:  12.9716,
	Longitude: 77.5946,
}

var getWarehouseByIdTests = []struct {
	getWarehouseByIdTxReturns *wms.Warehouse
	getWarehouseByIdTxErr     error
	want                      *wms.Warehouse
	wantErr                   error
}{
	{getWarehouseByIdTxReturns: &warehouse, getWarehouseByIdTxErr: nil, want: &warehouse, wantErr: nil},
	{getWarehouseByIdTxReturns: &warehouse, getWarehouseByIdTxErr: genericError, want: nil, wantErr: genericError},
	{getWarehouseByIdTxReturns: &warehouse, getWarehouseByIdTxErr: sql.ErrConnDone, want: nil, wantErr: sql.ErrConnDone},
	{getWarehouseByIdTxReturns: &warehouse, getWarehouseByIdTxErr: sql.ErrTxDone, want: nil, wantErr: sql.ErrTxDone},
	{getWarehouseByIdTxReturns: &warehouse, getWarehouseByIdTxErr: context.Canceled, want: nil, wantErr: context.Canceled},
}

var updateWarehouseTests = []struct {
	updateWarehouseTxReturns error
	want                     error
}{
	{updateWarehouseTxReturns: nil, want: nil},
	{updateWarehouseTxReturns: genericError, want: genericError},
	{updateWarehouseTxReturns: RowDoesNotExist, want: wms.WarehouseDoesNotExist},
	{updateWarehouseTxReturns: sql.ErrConnDone, want: sql.ErrConnDone},
	{updateWarehouseTxReturns: sql.ErrTxDone, want: sql.ErrTxDone},
	{updateWarehouseTxReturns: context.Canceled, want: context.Canceled},
}

var deleteWarehouseTests = []struct {
	deleteWarehouseTxReturns error
	want                     error
}{
	{deleteWarehouseTxReturns: nil, want: nil},
	{deleteWarehouseTxReturns: genericError, want: genericError},
	{deleteWarehouseTxReturns: RowDoesNotExist, want: wms.WarehouseDoesNotExist},
	{deleteWarehouseTxReturns: sql.ErrConnDone, want: sql.ErrConnDone},
	{deleteWarehouseTxReturns: sql.ErrTxDone, want: sql.ErrTxDone},
	{deleteWarehouseTxReturns: context.Canceled, want: context.Canceled},
}

func TestGetWarehouseById(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := NewMockqueries(mockCtrl)

	for _, test := range getWarehouseByIdTests {
		mockObj.EXPECT().getWarehouseByIdTx(ctx, gomock.Any(), warehouse.Id).Return(test.getWarehouseByIdTxReturns, test.getWarehouseByIdTxErr)

		ws := WarehouseService{
			queries: mockObj,
			db:      warehouseService.db,
		}
		got, err := ws.GetWarehouseById(ctx, warehouse.Id)
		if got != test.want {
			t.Errorf("want: %v, got: %v", test.want, got)
		}
		if err != test.wantErr {
			t.Errorf("want: %v, got: %v", test.wantErr, err)
		}
	}
}

func TestCreateWarehouse(t *testing.T) {
	ctx := context.Background()
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_find_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := NewMockqueries(mockCtrl)

	mockObj.EXPECT().createWarehouseTx(ctx, gomock.Any(), &warehouse).Return(nil)

	ws := WarehouseService{db: warehouseService.db, queries: mockObj}
	err := ws.CreateWarehouse(ctx, &warehouse)

	if err != nil {
		t.Error(err)
	}
}

func TestCreateWarehouseError(t *testing.T) {
	ctx := context.Background()
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_find_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}
	errs := []error{
		errors.New("generic error"),
		sql.ErrConnDone,
		sql.ErrTxDone,
		context.Canceled,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := NewMockqueries(mockCtrl)

	for _, testErr := range errs {
		mockObj.EXPECT().createWarehouseTx(ctx, gomock.Any(), &warehouse).Return(testErr)

		ws := WarehouseService{db: warehouseService.db, queries: mockObj}
		err := ws.CreateWarehouse(ctx, &warehouse)

		if err != testErr {
			t.Error(err)
		}
	}
}

func TestUpdateWarehouse(t *testing.T) {
	ctx := context.Background()
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_find_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := NewMockqueries(mockCtrl)

	for _, test := range updateWarehouseTests {
		mockObj.EXPECT().updateWarehouseTx(ctx, gomock.Any(), &warehouse).Return(test.updateWarehouseTxReturns)

		ws := WarehouseService{db: warehouseService.db, queries: mockObj}
		err := ws.UpdateWarehouse(ctx, &warehouse)

		if err != test.want {
			t.Error(err)
		}
	}
}

func TestDeleteWarehouse(t *testing.T) {
	ctx := context.Background()
	id := "85bd3b85-ad4d-4224-b589-fb2a80a6ce45"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := NewMockqueries(mockCtrl)

	for _, test := range deleteWarehouseTests {
		mockObj.EXPECT().deleteWarehouseTx(ctx, gomock.Any(), id).Return(test.deleteWarehouseTxReturns)

		ws := WarehouseService{db: warehouseService.db, queries: mockObj}
		err := ws.DeleteWarehouse(ctx, id)

		if err != test.want {
			t.Error(err)
		}
	}
}

func TestGetWarehouseByIDTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		Name:      "test_find_by_id",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	// insert row
	query := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	_, err = tx.ExecContext(context.Background(), query, warehouse.Id, warehouse.Name, warehouse.Longitude, warehouse.Latitude)
	if err != nil {
		t.Error(err)
	}

	warehouseFromDB, err := warehouseService.queries.getWarehouseByIdTx(context.Background(), tx, warehouse.Id)
	if err != nil {
		t.Error(err)
		return
	}

	if *warehouseFromDB != warehouse {
		t.Errorf("expected: %v, got: %v", warehouse, *warehouseFromDB)
	}
}

func TestGetWarehouseByIdTxError(t *testing.T) {
	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()
	_, err = warehouseService.queries.getWarehouseByIdTx(context.Background(), tx, "non-existent-id")
	if err == nil {
		t.Errorf("expected: %v, got: %v", "error", err)
	}
}

func TestCreateWarehouseTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "f405e23a-927f-4c58-9772-4497db2f62a0",
		Name:      "test_create",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	tx, err := warehouseService.db.BeginTx(context.Background(), nil)
	defer tx.Rollback()

	err = warehouseService.queries.createWarehouseTx(context.Background(), tx, &warehouse)
	if err != nil {
		t.Error(err)
	}

	query := "SELECT id, name, geolocation[0], geolocation[1] from warehouse where id=$1;"
	row := tx.QueryRowContext(context.Background(), query, warehouse.Id)

	var warehouseFromDB wms.Warehouse
	err = row.Scan(&warehouseFromDB.Id, &warehouseFromDB.Name, &warehouseFromDB.Longitude, &warehouseFromDB.Latitude)
	if err != nil {
		t.Error(err)
	}

	if warehouse != warehouseFromDB {
		t.Errorf("expected: %v, got: %v", warehouse, warehouseFromDB)
	}
}

func TestCreateWarehouseQueryTxError(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "6ef20892-d417-4e1b-b5c8-735509b2938b",
		Name:      "test_create_error",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = warehouseService.queries.createWarehouseTx(context.Background(), tx, &warehouse)
	if err != nil {
		t.Error(err)
	}

	// Create duplicate warehouse, expect error
	err = warehouseService.queries.createWarehouseTx(context.Background(), tx, &warehouse)
	if err == nil {
		t.Errorf("expected: %v, got: %v", "error", err)
	}

}

func TestCreateWarehouseQueryError(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "6ef20892-d417-4e1b-b5c8-735509b2938b",
		Name:      "test_create_error",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	err := warehouseService.CreateWarehouse(context.Background(), &warehouse)
	if err != nil {
		t.Error(err)
	}

	// Create duplicate warehouse, expect error
	err = warehouseService.CreateWarehouse(context.Background(), &warehouse)
	if err == nil {
		t.Errorf("expected: %v, got: %v", "error", err)
	}
}

func TestUpdateWarehouseTx(t *testing.T) {

	warehouse := wms.Warehouse{
		Id:        "36935050-9aa8-4425-ba05-48be9557f580",
		Name:      "test_update",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	warehouseUpdated := wms.Warehouse{
		Id:        "36935050-9aa8-4425-ba05-48be9557f580",
		Name:      "test_update_post_exec",
		Latitude:  12.9822,
		Longitude: 77.5898,
	}

	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	// insert row
	query := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	_, err = tx.ExecContext(context.Background(), query, warehouse.Id, warehouse.Name, warehouse.Longitude, warehouse.Latitude)
	if err != nil {
		t.Error(err)
	}

	// update row
	err = warehouseService.queries.updateWarehouseTx(context.Background(), tx, &warehouseUpdated)
	if err != nil {
		t.Error(err)
	}

	// select row
	query = "SELECT id, name, geolocation[0], geolocation[1] from warehouse where id=$1;"
	row := tx.QueryRowContext(context.Background(), query, warehouse.Id)

	var warehouseFromDB wms.Warehouse
	err = row.Scan(&warehouseFromDB.Id, &warehouseFromDB.Name, &warehouseFromDB.Longitude, &warehouseFromDB.Latitude)
	if err != nil {
		t.Error(err)
	}

	if warehouseUpdated != warehouseFromDB {
		t.Errorf("expected: %v, got: %v", warehouse, warehouseFromDB)
	}
}

func TestUpdateWarehouseTxError(t *testing.T) {
	warehouseUpdated := wms.Warehouse{
		Id:        "36935050-9aa8-4425-ba05-48be9557f580",
		Name:      "test_update_post_exec",
		Latitude:  12.9822,
		Longitude: 77.5898,
	}

	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	// update row
	err = warehouseService.queries.updateWarehouseTx(context.Background(), tx, &warehouseUpdated)
	if err == nil {
		t.Errorf(
			"expected: %v, got: %v",
			RowDoesNotExist,
			err,
		)
	}
}

func TestDeleteWarehouseTx(t *testing.T) {
	warehouse := wms.Warehouse{
		Id:        "6091629f-28fc-4dfa-a115-b94651652008",
		Name:      "test_delete",
		Latitude:  12.9716,
		Longitude: 77.5946,
	}

	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	// insert row
	query := "INSERT INTO warehouse (id, name, geolocation) VALUES ($1, $2, point($3, $4))"
	_, err = tx.ExecContext(context.Background(), query, warehouse.Id, warehouse.Name, warehouse.Longitude, warehouse.Latitude)
	if err != nil {
		t.Error(err)
		return
	}

	err = warehouseService.queries.deleteWarehouseTx(context.Background(), tx, warehouse.Id)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteWarehouseTxError(t *testing.T) {
	tx, err := warehouseService.db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	err = warehouseService.queries.deleteWarehouseTx(context.Background(), tx, "non-existent-id")
	if err == nil {
		t.Errorf(
			"expected: %v, got: %v",
			RowDoesNotExist,
			err,
		)
	}
}
