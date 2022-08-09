// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/corporate/v2/model"
	mock "github.com/stretchr/testify/mock"
)

// AccountTemporaryRepository is an autogenerated mock type for the AccountTemporaryRepository type
type AccountTemporaryRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *AccountTemporaryRepository) Delete(_a0 context.Context, _a1 model.B2BAccountTemporary) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BAccountTemporary) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *AccountTemporaryRepository) Save(_a0 context.Context, _a1 model.B2BAccountTemporary) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BAccountTemporary) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *AccountTemporaryRepository) Update(_a0 context.Context, _a1 model.B2BAccountTemporary) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2BAccountTemporary) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}