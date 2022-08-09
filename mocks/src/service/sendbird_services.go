// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/service/model"
	mock "github.com/stretchr/testify/mock"
)

// SendbirdServices is an autogenerated mock type for the SendbirdServices type
type SendbirdServices struct {
	mock.Mock
}

// CheckUserSenbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CheckUserSenbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// CheckUserSenbirdV4 provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CheckUserSenbirdV4(ctxReq context.Context, data *model.SendbirdRequestV4) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequestV4) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// CreateTokenUserSendbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CreateTokenUserSendbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// CreateTokenUserSendbirdV4 provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CreateTokenUserSendbirdV4(ctxReq context.Context, data *model.SendbirdRequestV4) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequestV4) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// CreateUserSendbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CreateUserSendbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// CreateUserSendbirdV4 provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) CreateUserSendbirdV4(ctxReq context.Context, data *model.SendbirdRequestV4) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequestV4) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// GetTokenUserSendbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) GetTokenUserSendbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// GetUserSendbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) GetUserSendbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// GetUserSendbirdV4 provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) GetUserSendbirdV4(ctxReq context.Context, data *model.SendbirdRequestV4) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequestV4) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// UpdateUserSendbird provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) UpdateUserSendbird(ctxReq context.Context, data *model.SendbirdRequest) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequest) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}

// UpdateUserSendbirdV4 provides a mock function with given fields: ctxReq, data
func (_m *SendbirdServices) UpdateUserSendbirdV4(ctxReq context.Context, data *model.SendbirdRequestV4) model.ServiceResult {
	ret := _m.Called(ctxReq, data)

	var r0 model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *model.SendbirdRequestV4) model.ServiceResult); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Get(0).(model.ServiceResult)
	}

	return r0
}
