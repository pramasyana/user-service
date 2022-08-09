// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/shared/model"
	mock "github.com/stretchr/testify/mock"
)

// ContactRepository is an autogenerated mock type for the ContactRepository type
type ContactRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *ContactRepository) Delete(_a0 context.Context, _a1 model.B2BContact) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BContact) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *ContactRepository) Save(_a0 context.Context, _a1 model.B2BContact) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BContact) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *ContactRepository) Update(_a0 context.Context, _a1 model.B2BContact) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BContact) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
