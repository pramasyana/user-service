// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/auth/v1/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Bhinneka/user-service/src/auth/v1/repo"
)

// RefreshTokenRepository is an autogenerated mock type for the RefreshTokenRepository type
type RefreshTokenRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctxReq, key
func (_m *RefreshTokenRepository) Delete(ctxReq context.Context, key string) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, key)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Load provides a mock function with given fields: ctxReq, key
func (_m *RefreshTokenRepository) Load(ctxReq context.Context, key string) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, key)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Save provides a mock function with given fields: ctxReq, refreshToken
func (_m *RefreshTokenRepository) Save(ctxReq context.Context, refreshToken *model.RefreshToken) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, refreshToken)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.RefreshToken) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, refreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}