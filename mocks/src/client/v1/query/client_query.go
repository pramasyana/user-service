// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ClientQuery is an autogenerated mock type for the ClientQuery type
type ClientQuery struct {
	mock.Mock
}

// Validate provides a mock function with given fields: ctxReq, clientID, clientSecret
func (_m *ClientQuery) Validate(ctxReq context.Context, clientID string, clientSecret string) (bool, error) {
	ret := _m.Called(ctxReq, clientID, clientSecret)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctxReq, clientID, clientSecret)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctxReq, clientID, clientSecret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}