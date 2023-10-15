// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/chff7cb/swissbank/core"
	mock "github.com/stretchr/testify/mock"
)

// TransactionsDataProxy is an autogenerated mock type for the TransactionsDataProxy type
type TransactionsDataProxy struct {
	mock.Mock
}

// CreateTransaction provides a mock function with given fields: _a0, _a1
func (_m *TransactionsDataProxy) CreateTransaction(_a0 context.Context, _a1 *core.Transaction) (*core.Transaction, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *core.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.Transaction) (*core.Transaction, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *core.Transaction) *core.Transaction); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *core.Transaction) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionsDataProxy creates a new instance of TransactionsDataProxy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionsDataProxy(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionsDataProxy {
	mock := &TransactionsDataProxy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}