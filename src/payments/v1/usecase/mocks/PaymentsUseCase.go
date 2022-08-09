// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/payments/v1/model"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/Bhinneka/user-service/src/payments/v1/usecase"
)

// PaymentsUseCase is an autogenerated mock type for the PaymentsUseCase type
type PaymentsUseCase struct {
	mock.Mock
}

// AddUpdatePayments provides a mock function with given fields: ctxReq, data
func (_m *PaymentsUseCase) AddUpdatePayments(ctxReq context.Context, data *model.Payments) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Payments) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CompareHeaderAndBody provides a mock function with given fields: ctxReq, data, basicAuth
func (_m *PaymentsUseCase) CompareHeaderAndBody(ctxReq context.Context, data *model.Payments, basicAuth string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, basicAuth)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Payments, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, basicAuth)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetPaymentDetail provides a mock function with given fields: ctxReq, data
func (_m *PaymentsUseCase) GetPaymentDetail(ctxReq context.Context, data *model.Payments) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Payments) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}
