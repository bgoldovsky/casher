// Code generated by MockGen. DO NOT EDIT.
// Source: operations.go

// Package operations is a generated GoMock package.
package operations

import (
	reflect "reflect"

	models "github.com/bgoldovsky/casher/app/models"
	gomock "github.com/golang/mock/gomock"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockrepository) Create(operation *models.Operation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", operation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockrepositoryMockRecorder) Create(operation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockrepository)(nil).Create), operation)
}

// Get mocks base method.
func (m *Mockrepository) Get(userID, page, size int64) (*models.OperationPaginator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", userID, page, size)
	ret0, _ := ret[0].(*models.OperationPaginator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockrepositoryMockRecorder) Get(userID, page, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*Mockrepository)(nil).Get), userID, page, size)
}

// Remove mocks base method.
func (m *Mockrepository) Remove(operationID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", operationID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockrepositoryMockRecorder) Remove(operationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*Mockrepository)(nil).Remove), operationID)
}
