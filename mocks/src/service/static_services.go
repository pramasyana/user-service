// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/service/model"
	mock "github.com/stretchr/testify/mock"
)

// StaticServices is an autogenerated mock type for the StaticServices type
type StaticServices struct {
	mock.Mock
}

// FindStaticsByID provides a mock function with given fields: ctxReq, id
func (_m *StaticServices) FindStaticsByID(ctxReq context.Context, id string) <-chan model.ServiceResult {
	ret := _m.Called(ctxReq, id)

	var r0 <-chan model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan model.ServiceResult); ok {
		r0 = rf(ctxReq, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan model.ServiceResult)
		}
	}

	return r0
}

// FindStaticsGwsByID provides a mock function with given fields: ctxReq, id
func (_m *StaticServices) FindStaticsGwsByID(ctxReq context.Context, id string) <-chan model.ServiceResult {
	ret := _m.Called(ctxReq, id)

	var r0 <-chan model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan model.ServiceResult); ok {
		r0 = rf(ctxReq, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan model.ServiceResult)
		}
	}

	return r0
}