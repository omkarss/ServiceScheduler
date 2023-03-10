// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/scheduler/scheduler.go

// Package mock_scheduler is a generated GoMock package.
package mock_scheduler

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSchedulerI is a mock of SchedulerI interface.
type MockSchedulerI struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerIMockRecorder
}

// MockSchedulerIMockRecorder is the mock recorder for MockSchedulerI.
type MockSchedulerIMockRecorder struct {
	mock *MockSchedulerI
}

// NewMockSchedulerI creates a new mock instance.
func NewMockSchedulerI(ctrl *gomock.Controller) *MockSchedulerI {
	mock := &MockSchedulerI{ctrl: ctrl}
	mock.recorder = &MockSchedulerIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedulerI) EXPECT() *MockSchedulerIMockRecorder {
	return m.recorder
}

// GetNextCustomer mocks base method.
func (m *MockSchedulerI) GetNextCustomer(arg0 *context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextCustomer", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNextCustomer indicates an expected call of GetNextCustomer.
func (mr *MockSchedulerIMockRecorder) GetNextCustomer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextCustomer", reflect.TypeOf((*MockSchedulerI)(nil).GetNextCustomer), arg0)
}

// GetNextTicketNumber mocks base method.
func (m *MockSchedulerI) GetNextTicketNumber(arg0 *context.Context) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextTicketNumber", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetNextTicketNumber indicates an expected call of GetNextTicketNumber.
func (mr *MockSchedulerIMockRecorder) GetNextTicketNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextTicketNumber", reflect.TypeOf((*MockSchedulerI)(nil).GetNextTicketNumber), arg0)
}
