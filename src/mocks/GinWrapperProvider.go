// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"

	providers "github.com/chff7cb/swissbank/providers"
)

// GinWrapperProvider is an autogenerated mock type for the GinWrapperProvider type
type GinWrapperProvider struct {
	mock.Mock
}

// Wrap provides a mock function with given fields: ctx
func (_m *GinWrapperProvider) Wrap(ctx *gin.Context) providers.GinWrapper {
	ret := _m.Called(ctx)

	var r0 providers.GinWrapper
	if rf, ok := ret.Get(0).(func(*gin.Context) providers.GinWrapper); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(providers.GinWrapper)
		}
	}

	return r0
}

// NewGinWrapperProvider creates a new instance of GinWrapperProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGinWrapperProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *GinWrapperProvider {
	mock := &GinWrapperProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
