// Code generated by MockGen. DO NOT EDIT.
// Source: ../root.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	todo "elastic-search/pkg/todo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStoreIface is a mock of StoreIface interface
type MockStoreIface struct {
	ctrl     *gomock.Controller
	recorder *MockStoreIfaceMockRecorder
}

// MockStoreIfaceMockRecorder is the mock recorder for MockStoreIface
type MockStoreIfaceMockRecorder struct {
	mock *MockStoreIface
}

// NewMockStoreIface creates a new mock instance
func NewMockStoreIface(ctrl *gomock.Controller) *MockStoreIface {
	mock := &MockStoreIface{ctrl: ctrl}
	mock.recorder = &MockStoreIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStoreIface) EXPECT() *MockStoreIfaceMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockStoreIface) Save(todo todo.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", todo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockStoreIfaceMockRecorder) Save(todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStoreIface)(nil).Save), todo)
}

// GetByID mocks base method
func (m *MockStoreIface) GetByID(ID string) (todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockStoreIfaceMockRecorder) GetByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockStoreIface)(nil).GetByID), ID)
}

// Remove mocks base method
func (m *MockStoreIface) Remove(ID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockStoreIfaceMockRecorder) Remove(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockStoreIface)(nil).Remove), ID)
}

// MockElasticSearchServiceIface is a mock of ElasticSearchServiceIface interface
type MockElasticSearchServiceIface struct {
	ctrl     *gomock.Controller
	recorder *MockElasticSearchServiceIfaceMockRecorder
}

// MockElasticSearchServiceIfaceMockRecorder is the mock recorder for MockElasticSearchServiceIface
type MockElasticSearchServiceIfaceMockRecorder struct {
	mock *MockElasticSearchServiceIface
}

// NewMockElasticSearchServiceIface creates a new mock instance
func NewMockElasticSearchServiceIface(ctrl *gomock.Controller) *MockElasticSearchServiceIface {
	mock := &MockElasticSearchServiceIface{ctrl: ctrl}
	mock.recorder = &MockElasticSearchServiceIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockElasticSearchServiceIface) EXPECT() *MockElasticSearchServiceIfaceMockRecorder {
	return m.recorder
}

// Index mocks base method
func (m *MockElasticSearchServiceIface) Index(ctx context.Context, td todo.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", ctx, td)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockElasticSearchServiceIfaceMockRecorder) Index(ctx, td interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockElasticSearchServiceIface)(nil).Index), ctx, td)
}

// Search mocks base method
func (m *MockElasticSearchServiceIface) Search(ctx context.Context, query string) ([]todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].([]todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockElasticSearchServiceIfaceMockRecorder) Search(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockElasticSearchServiceIface)(nil).Search), ctx, query)
}
