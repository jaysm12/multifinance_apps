// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/authentication/service.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	authentication "github.com/jaysm12/multifinance-apps/internal/service/authentication"
)

// MockAuthenticationServiceMethod is a mock of AuthenticationServiceMethod interface.
type MockAuthenticationServiceMethod struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationServiceMethodMockRecorder
}

// MockAuthenticationServiceMethodMockRecorder is the mock recorder for MockAuthenticationServiceMethod.
type MockAuthenticationServiceMethodMockRecorder struct {
	mock *MockAuthenticationServiceMethod
}

// NewMockAuthenticationServiceMethod creates a new mock instance.
func NewMockAuthenticationServiceMethod(ctrl *gomock.Controller) *MockAuthenticationServiceMethod {
	mock := &MockAuthenticationServiceMethod{ctrl: ctrl}
	mock.recorder = &MockAuthenticationServiceMethodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticationServiceMethod) EXPECT() *MockAuthenticationServiceMethodMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthenticationServiceMethod) Login(arg0 authentication.LoginServiceRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthenticationServiceMethodMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthenticationServiceMethod)(nil).Login), arg0)
}

// Register mocks base method.
func (m *MockAuthenticationServiceMethod) Register(arg0 authentication.RegisterServiceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockAuthenticationServiceMethodMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthenticationServiceMethod)(nil).Register), arg0)
}
