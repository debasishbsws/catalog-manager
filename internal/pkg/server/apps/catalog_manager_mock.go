// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/napptive/catalog-manager/internal/pkg/server/catalog-manager (interfaces: Manager)

// Package apps is a generated GoMock package.
package apps

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/napptive/catalog-manager/internal/pkg/entities"
)

// MockCatalogManager is a mock of Manager interface.
type MockCatalogManager struct {
	ctrl     *gomock.Controller
	recorder *MockCatalogManagerMockRecorder
}

// MockCatalogManagerMockRecorder is the mock recorder for MockCatalogManager.
type MockCatalogManagerMockRecorder struct {
	mock *MockCatalogManager
}

// NewMockCatalogManager creates a new mock instance.
func NewMockCatalogManager(ctrl *gomock.Controller) *MockCatalogManager {
	mock := &MockCatalogManager{ctrl: ctrl}
	mock.recorder = &MockCatalogManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCatalogManager) EXPECT() *MockCatalogManagerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCatalogManager) Add(arg0 string, arg1 []*entities.FileInfo, arg2 bool, arg3 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockCatalogManagerMockRecorder) Add(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCatalogManager)(nil).Add), arg0, arg1, arg2, arg3)
}

// Download mocks base method.
func (m *MockCatalogManager) Download(arg0 string, arg1 bool, arg2 string) ([]*entities.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*entities.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockCatalogManagerMockRecorder) Download(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockCatalogManager)(nil).Download), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockCatalogManager) Get(arg0, arg1 string) (*entities.ExtendedApplicationMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*entities.ExtendedApplicationMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCatalogManagerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCatalogManager)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockCatalogManager) List(arg0, arg1 string) ([]*entities.AppSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*entities.AppSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCatalogManagerMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCatalogManager)(nil).List), arg0, arg1)
}

// Remove mocks base method.
func (m *MockCatalogManager) Remove(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockCatalogManagerMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockCatalogManager)(nil).Remove), arg0)
}

// Summary mocks base method.
func (m *MockCatalogManager) Summary() (*entities.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Summary")
	ret0, _ := ret[0].(*entities.Summary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Summary indicates an expected call of Summary.
func (mr *MockCatalogManagerMockRecorder) Summary() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Summary", reflect.TypeOf((*MockCatalogManager)(nil).Summary))
}

// UpdateApplicationVisibility mocks base method.
func (m *MockCatalogManager) UpdateApplicationVisibility(arg0, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateApplicationVisibility", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateApplicationVisibility indicates an expected call of UpdateApplicationVisibility.
func (mr *MockCatalogManagerMockRecorder) UpdateApplicationVisibility(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApplicationVisibility", reflect.TypeOf((*MockCatalogManager)(nil).UpdateApplicationVisibility), arg0, arg1, arg2)
}
