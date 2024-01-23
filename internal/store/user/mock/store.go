// Code generated by MockGen. DO NOT EDIT.
// Source: internal/store/user/store.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/jaysm12/multifinance-apps/models"
)

// MockUserStoreMethod is a mock of UserStoreMethod interface.
type MockUserStoreMethod struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMethodMockRecorder
}

// MockUserStoreMethodMockRecorder is the mock recorder for MockUserStoreMethod.
type MockUserStoreMethodMockRecorder struct {
	mock *MockUserStoreMethod
}

// NewMockUserStoreMethod creates a new mock instance.
func NewMockUserStoreMethod(ctrl *gomock.Controller) *MockUserStoreMethod {
	mock := &MockUserStoreMethod{ctrl: ctrl}
	mock.recorder = &MockUserStoreMethodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStoreMethod) EXPECT() *MockUserStoreMethodMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockUserStoreMethod) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockUserStoreMethodMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockUserStoreMethod)(nil).Count))
}

// CreateUser mocks base method.
func (m *MockUserStoreMethod) CreateUser(userinfo models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", userinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserStoreMethodMockRecorder) CreateUser(userinfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserStoreMethod)(nil).CreateUser), userinfo)
}

// DeleteUser mocks base method.
func (m *MockUserStoreMethod) DeleteUser(userid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserStoreMethodMockRecorder) DeleteUser(userid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserStoreMethod)(nil).DeleteUser), userid)
}

// GetUserInfoByID mocks base method.
func (m *MockUserStoreMethod) GetUserInfoByID(userid uint) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoByID", userid)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoByID indicates an expected call of GetUserInfoByID.
func (mr *MockUserStoreMethodMockRecorder) GetUserInfoByID(userid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoByID", reflect.TypeOf((*MockUserStoreMethod)(nil).GetUserInfoByID), userid)
}

// GetUserInfoByUsername mocks base method.
func (m *MockUserStoreMethod) GetUserInfoByUsername(username string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfoByUsername", username)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfoByUsername indicates an expected call of GetUserInfoByUsername.
func (mr *MockUserStoreMethodMockRecorder) GetUserInfoByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfoByUsername", reflect.TypeOf((*MockUserStoreMethod)(nil).GetUserInfoByUsername), username)
}

// UpdateUser mocks base method.
func (m *MockUserStoreMethod) UpdateUser(userinfo models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", userinfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserStoreMethodMockRecorder) UpdateUser(userinfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserStoreMethod)(nil).UpdateUser), userinfo)
}
