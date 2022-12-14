// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/auth/v1/model"
	mock "github.com/stretchr/testify/mock"
)

// LDAPService is an autogenerated mock type for the LDAPService type
type LDAPService struct {
	mock.Mock
}

// Auth provides a mock function with given fields: _a0, _a1, _a2
func (_m *LDAPService) Auth(_a0 context.Context, _a1 string, _a2 string) (*model.LDAPProfile, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.LDAPProfile
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.LDAPProfile); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LDAPProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
