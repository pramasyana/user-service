// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/shared/model"
	mock "github.com/stretchr/testify/mock"
)

// ClientUsecase is an autogenerated mock type for the ClientUsecase type
type ClientUsecase struct {
	mock.Mock
}

// Logout provides a mock function with given fields: ctx, email
func (_m *ClientUsecase) Logout(ctx context.Context, email string) <-chan model.ResultUseCase {
	ret := _m.Called(ctx, email)

	var r0 <-chan model.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan model.ResultUseCase); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan model.ResultUseCase)
		}
	}

	return r0
}
