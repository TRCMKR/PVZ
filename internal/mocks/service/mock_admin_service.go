// Code generated by MockGen. DO NOT EDIT.
// Source: admin_handler.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// MockadminService is a mock of adminService interface.
type MockadminService struct {
	ctrl     *gomock.Controller
	recorder *MockadminServiceMockRecorder
}

// MockadminServiceMockRecorder is the mock recorder for MockadminService.
type MockadminServiceMockRecorder struct {
	mock *MockadminService
}

// NewMockadminService creates a new mock instance.
func NewMockadminService(ctrl *gomock.Controller) *MockadminService {
	mock := &MockadminService{ctrl: ctrl}
	mock.recorder = &MockadminServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockadminService) EXPECT() *MockadminServiceMockRecorder {
	return m.recorder
}

// ContainsID mocks base method.
func (m *MockadminService) ContainsID(arg0 context.Context, arg1 int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainsID", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainsID indicates an expected call of ContainsID.
func (mr *MockadminServiceMockRecorder) ContainsID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainsID", reflect.TypeOf((*MockadminService)(nil).ContainsID), arg0, arg1)
}

// ContainsUsername mocks base method.
func (m *MockadminService) ContainsUsername(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainsUsername", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainsUsername indicates an expected call of ContainsUsername.
func (mr *MockadminServiceMockRecorder) ContainsUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainsUsername", reflect.TypeOf((*MockadminService)(nil).ContainsUsername), arg0, arg1)
}

// CreateAdmin mocks base method.
func (m *MockadminService) CreateAdmin(arg0 context.Context, arg1 models.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdmin", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAdmin indicates an expected call of CreateAdmin.
func (mr *MockadminServiceMockRecorder) CreateAdmin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdmin", reflect.TypeOf((*MockadminService)(nil).CreateAdmin), arg0, arg1)
}

// DeleteAdmin mocks base method.
func (m *MockadminService) DeleteAdmin(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdmin", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdmin indicates an expected call of DeleteAdmin.
func (mr *MockadminServiceMockRecorder) DeleteAdmin(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdmin", reflect.TypeOf((*MockadminService)(nil).DeleteAdmin), arg0, arg1, arg2)
}

// GetAdminByUsername mocks base method.
func (m *MockadminService) GetAdminByUsername(arg0 context.Context, arg1 string) (models.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminByUsername", arg0, arg1)
	ret0, _ := ret[0].(models.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminByUsername indicates an expected call of GetAdminByUsername.
func (mr *MockadminServiceMockRecorder) GetAdminByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminByUsername", reflect.TypeOf((*MockadminService)(nil).GetAdminByUsername), arg0, arg1)
}

// UpdateAdmin mocks base method.
func (m *MockadminService) UpdateAdmin(arg0 context.Context, arg1, arg2 string, arg3 models.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdmin", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAdmin indicates an expected call of UpdateAdmin.
func (mr *MockadminServiceMockRecorder) UpdateAdmin(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdmin", reflect.TypeOf((*MockadminService)(nil).UpdateAdmin), arg0, arg1, arg2, arg3)
}
