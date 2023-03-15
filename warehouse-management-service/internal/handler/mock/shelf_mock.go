// Code generated by MockGen. DO NOT EDIT.
// Source: ./shelf.go

// Package mock_warehousemanagementservice is a generated GoMock package.
package mock_warehousemanagementservice

import (
	context "context"
	reflect "reflect"
	warehousemanagementservice "warehouse-management-service"

	gomock "github.com/golang/mock/gomock"
)

// MockShelfService is a mock of ShelfService interface.
type MockShelfService struct {
	ctrl     *gomock.Controller
	recorder *MockShelfServiceMockRecorder
}

// MockShelfServiceMockRecorder is the mock recorder for MockShelfService.
type MockShelfServiceMockRecorder struct {
	mock *MockShelfService
}

// NewMockShelfService creates a new mock instance.
func NewMockShelfService(ctrl *gomock.Controller) *MockShelfService {
	mock := &MockShelfService{ctrl: ctrl}
	mock.recorder = &MockShelfServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShelfService) EXPECT() *MockShelfServiceMockRecorder {
	return m.recorder
}

// CreateShelf mocks base method.
func (m *MockShelfService) CreateShelf(ctx context.Context, shelf warehousemanagementservice.Shelf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShelf", ctx, shelf)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateShelf indicates an expected call of CreateShelf.
func (mr *MockShelfServiceMockRecorder) CreateShelf(ctx, shelf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShelf", reflect.TypeOf((*MockShelfService)(nil).CreateShelf), ctx, shelf)
}

// DeleteShelfById mocks base method.
func (m *MockShelfService) DeleteShelfById(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShelfById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShelfById indicates an expected call of DeleteShelfById.
func (mr *MockShelfServiceMockRecorder) DeleteShelfById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShelfById", reflect.TypeOf((*MockShelfService)(nil).DeleteShelfById), ctx, id)
}

// GetShelfById mocks base method.
func (m *MockShelfService) GetShelfById(ctx context.Context, id string) (warehousemanagementservice.Shelf, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShelfById", ctx, id)
	ret0, _ := ret[0].(warehousemanagementservice.Shelf)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShelfById indicates an expected call of GetShelfById.
func (mr *MockShelfServiceMockRecorder) GetShelfById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShelfById", reflect.TypeOf((*MockShelfService)(nil).GetShelfById), ctx, id)
}

// UpdateShelf mocks base method.
func (m *MockShelfService) UpdateShelf(ctx context.Context, shelf warehousemanagementservice.Shelf) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShelf", ctx, shelf)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShelf indicates an expected call of UpdateShelf.
func (mr *MockShelfServiceMockRecorder) UpdateShelf(ctx, shelf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShelf", reflect.TypeOf((*MockShelfService)(nil).UpdateShelf), ctx, shelf)
}
