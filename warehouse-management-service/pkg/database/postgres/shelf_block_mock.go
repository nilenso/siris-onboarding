// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/database/postgres/shelf_block.go

// Package mock_postgres is a generated GoMock package.
package postgres

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	warehousemanagementservice "warehouse-management-service"

	gomock "github.com/golang/mock/gomock"
)

// MockshelfBlockQueries is a mock of shelfBlockQueries interface.
type MockshelfBlockQueries struct {
	ctrl     *gomock.Controller
	recorder *MockshelfBlockQueriesMockRecorder
}

// MockshelfBlockQueriesMockRecorder is the mock recorder for MockshelfBlockQueries.
type MockshelfBlockQueriesMockRecorder struct {
	mock *MockshelfBlockQueries
}

// NewMockshelfBlockQueries creates a new mock instance.
func NewMockshelfBlockQueries(ctrl *gomock.Controller) *MockshelfBlockQueries {
	mock := &MockshelfBlockQueries{ctrl: ctrl}
	mock.recorder = &MockshelfBlockQueriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockshelfBlockQueries) EXPECT() *MockshelfBlockQueriesMockRecorder {
	return m.recorder
}

// createShelfBlockTx mocks base method.
func (m *MockshelfBlockQueries) createShelfBlockTx(ctx context.Context, tx *sql.Tx, block warehousemanagementservice.ShelfBlock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createShelfBlockTx", ctx, tx, block)
	ret0, _ := ret[0].(error)
	return ret0
}

// createShelfBlockTx indicates an expected call of createShelfBlockTx.
func (mr *MockshelfBlockQueriesMockRecorder) createShelfBlockTx(ctx, tx, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createShelfBlockTx", reflect.TypeOf((*MockshelfBlockQueries)(nil).createShelfBlockTx), ctx, tx, block)
}

// deleteShelfBlockTx mocks base method.
func (m *MockshelfBlockQueries) deleteShelfBlockTx(ctx context.Context, tx *sql.Tx, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteShelfBlockTx", ctx, tx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteShelfBlockTx indicates an expected call of deleteShelfBlockTx.
func (mr *MockshelfBlockQueriesMockRecorder) deleteShelfBlockTx(ctx, tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteShelfBlockTx", reflect.TypeOf((*MockshelfBlockQueries)(nil).deleteShelfBlockTx), ctx, tx, id)
}

// getShelfBlockByIdTx mocks base method.
func (m *MockshelfBlockQueries) getShelfBlockByIdTx(ctx context.Context, tx *sql.Tx, id string) (warehousemanagementservice.ShelfBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getShelfBlockByIdTx", ctx, tx, id)
	ret0, _ := ret[0].(warehousemanagementservice.ShelfBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getShelfBlockByIdTx indicates an expected call of getShelfBlockByIdTx.
func (mr *MockshelfBlockQueriesMockRecorder) getShelfBlockByIdTx(ctx, tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getShelfBlockByIdTx", reflect.TypeOf((*MockshelfBlockQueries)(nil).getShelfBlockByIdTx), ctx, tx, id)
}

// updateShelfBlockTx mocks base method.
func (m *MockshelfBlockQueries) updateShelfBlockTx(ctx context.Context, tx *sql.Tx, block warehousemanagementservice.ShelfBlock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "updateShelfBlockTx", ctx, tx, block)
	ret0, _ := ret[0].(error)
	return ret0
}

// updateShelfBlockTx indicates an expected call of updateShelfBlockTx.
func (mr *MockshelfBlockQueriesMockRecorder) updateShelfBlockTx(ctx, tx, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "updateShelfBlockTx", reflect.TypeOf((*MockshelfBlockQueries)(nil).updateShelfBlockTx), ctx, tx, block)
}

// warehouseExistsTx mocks base method.
func (m *MockshelfBlockQueries) warehouseExistsTx(ctx context.Context, tx *sql.Tx, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "warehouseExistsTx", ctx, tx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// warehouseExistsTx indicates an expected call of warehouseExistsTx.
func (mr *MockshelfBlockQueriesMockRecorder) warehouseExistsTx(ctx, tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "warehouseExistsTx", reflect.TypeOf((*MockshelfBlockQueries)(nil).warehouseExistsTx), ctx, tx, id)
}
