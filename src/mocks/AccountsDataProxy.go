// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/chff7cb/swissbank/core"
	mock "github.com/stretchr/testify/mock"
)

// AccountsDataProxy is an autogenerated mock type for the AccountsDataProxy type
type AccountsDataProxy struct {
	mock.Mock
}

// CreateAccount provides a mock function with given fields: _a0, _a1
func (_m *AccountsDataProxy) CreateAccount(_a0 context.Context, _a1 *core.Account) (*core.Account, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *core.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.Account) (*core.Account, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *core.Account) *core.Account); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *core.Account) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountByID provides a mock function with given fields: _a0, _a1
func (_m *AccountsDataProxy) GetAccountByID(_a0 context.Context, _a1 string) (*core.Account, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *core.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*core.Account, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.Account); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccountsDataProxy creates a new instance of AccountsDataProxy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountsDataProxy(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountsDataProxy {
	mock := &AccountsDataProxy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}