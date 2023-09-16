// Code generated by MockGen. DO NOT EDIT.
// Source: pthd-notifications/pkg/domain (interfaces: IConnector)

// Package mocks is a generated GoMock package.
package mocks

import (
	model "pthd-notifications/pkg/domain/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIConnector is a mock of IConnector interface.
type MockIConnector struct {
	ctrl     *gomock.Controller
	recorder *MockIConnectorMockRecorder
}

// MockIConnectorMockRecorder is the mock recorder for MockIConnector.
type MockIConnectorMockRecorder struct {
	mock *MockIConnector
}

// NewMockIConnector creates a new mock instance.
func NewMockIConnector(ctrl *gomock.Controller) *MockIConnector {
	mock := &MockIConnector{ctrl: ctrl}
	mock.recorder = &MockIConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConnector) EXPECT() *MockIConnectorMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockIConnector) Send(arg0 *model.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockIConnectorMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockIConnector)(nil).Send), arg0)
}
