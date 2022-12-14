// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/service/model"
	mock "github.com/stretchr/testify/mock"

	v1model "github.com/Bhinneka/user-service/src/member/v1/model"

	v2model "github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// MerchantServices is an autogenerated mock type for the MerchantServices type
type MerchantServices struct {
	mock.Mock
}

// FindMerchantServiceByID provides a mock function with given fields: ctxReq, id, token, merchantID
func (_m *MerchantServices) FindMerchantServiceByID(ctxReq context.Context, id string, token string, merchantID string) <-chan model.ServiceResult {
	ret := _m.Called(ctxReq, id, token, merchantID)

	var r0 <-chan model.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan model.ServiceResult); ok {
		r0 = rf(ctxReq, id, token, merchantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan model.ServiceResult)
		}
	}

	return r0
}

// InsertLogMerchant provides a mock function with given fields: ctxReq, oldData, newData, action, module
func (_m *MerchantServices) InsertLogMerchant(ctxReq context.Context, oldData v2model.B2CMerchantDataV2, newData v2model.B2CMerchantDataV2, action string, module string) error {
	ret := _m.Called(ctxReq, oldData, newData, action, module)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, v2model.B2CMerchantDataV2, v2model.B2CMerchantDataV2, string, string) error); ok {
		r0 = rf(ctxReq, oldData, newData, action, module)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertLogMerchantPIC provides a mock function with given fields: ctxReq, oldData, newData, action, module, member
func (_m *MerchantServices) InsertLogMerchantPIC(ctxReq context.Context, oldData v2model.B2CMerchantDataV2, newData v2model.B2CMerchantDataV2, action string, module string, member v1model.Member) error {
	ret := _m.Called(ctxReq, oldData, newData, action, module, member)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, v2model.B2CMerchantDataV2, v2model.B2CMerchantDataV2, string, string, v1model.Member) error); ok {
		r0 = rf(ctxReq, oldData, newData, action, module, member)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishToKafkaUserMerchant provides a mock function with given fields: ctxReq, data, eventType, producer
func (_m *MerchantServices) PublishToKafkaUserMerchant(ctxReq context.Context, data *v2model.B2CMerchantDataV2, eventType string, producer string) error {
	ret := _m.Called(ctxReq, data, eventType, producer)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2model.B2CMerchantDataV2, string, string) error); ok {
		r0 = rf(ctxReq, data, eventType, producer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
