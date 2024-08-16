// Code generated by MockGen. DO NOT EDIT.
// Source: init.go

// Package mock_endpoint is a generated GoMock package.
package mock_endpoint

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/reversersed/AuthService/internal/service"
)

// Mockservice is a mock of service interface.
type Mockservice struct {
	ctrl     *gomock.Controller
	recorder *MockserviceMockRecorder
}

// MockserviceMockRecorder is the mock recorder for Mockservice.
type MockserviceMockRecorder struct {
	mock *Mockservice
}

// NewMockservice creates a new mock instance.
func NewMockservice(ctrl *gomock.Controller) *Mockservice {
	mock := &Mockservice{ctrl: ctrl}
	mock.recorder = &MockserviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockservice) EXPECT() *MockserviceMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *Mockservice) GenerateAccessToken(arg0 context.Context, arg1, arg2 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockserviceMockRecorder) GenerateAccessToken(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*Mockservice)(nil).GenerateAccessToken), arg0, arg1, arg2)
}

// ValidateUserToken mocks base method.
func (m *Mockservice) ValidateUserToken(arg0 context.Context, arg1, arg2, arg3 string) (*service.Claims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUserToken", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*service.Claims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUserToken indicates an expected call of ValidateUserToken.
func (mr *MockserviceMockRecorder) ValidateUserToken(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUserToken", reflect.TypeOf((*Mockservice)(nil).ValidateUserToken), arg0, arg1, arg2, arg3)
}

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Info mocks base method.
func (m *Mocklogger) Info(arg0 ...any) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), arg0...)
}

// Mockvalidator is a mock of validator interface.
type Mockvalidator struct {
	ctrl     *gomock.Controller
	recorder *MockvalidatorMockRecorder
}

// MockvalidatorMockRecorder is the mock recorder for Mockvalidator.
type MockvalidatorMockRecorder struct {
	mock *Mockvalidator
}

// NewMockvalidator creates a new mock instance.
func NewMockvalidator(ctrl *gomock.Controller) *Mockvalidator {
	mock := &Mockvalidator{ctrl: ctrl}
	mock.recorder = &MockvalidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockvalidator) EXPECT() *MockvalidatorMockRecorder {
	return m.recorder
}

// StructValidation mocks base method.
func (m *Mockvalidator) StructValidation(data any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StructValidation", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// StructValidation indicates an expected call of StructValidation.
func (mr *MockvalidatorMockRecorder) StructValidation(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StructValidation", reflect.TypeOf((*Mockvalidator)(nil).StructValidation), data)
}
