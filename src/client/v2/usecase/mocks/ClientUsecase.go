// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import sharedmodel "github.com/Bhinneka/user-service/src/shared/model"

// ClientUsecase is an autogenerated mock type for the ClientUsecase type
type ClientUsecase struct {
	mock.Mock
}

// Logout provides a mock function with given fields: ctx, email
func (_m *ClientUsecase) Logout(ctx context.Context, email string) <-chan sharedmodel.ResultUseCase {
	ret := _m.Called(ctx, email)

	var r0 <-chan sharedmodel.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan sharedmodel.ResultUseCase); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan sharedmodel.ResultUseCase)
		}
	}

	return r0
}
