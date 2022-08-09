// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/auth/v1/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Bhinneka/user-service/src/auth/v1/repo"
)

// AttemptRepository is an autogenerated mock type for the AttemptRepository type
type AttemptRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctxReq, key
func (_m *AttemptRepository) Delete(ctxReq context.Context, key string) <-chan repo.ResultRepository {
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
func (_m *AttemptRepository) Load(ctxReq context.Context, key string) <-chan repo.ResultRepository {
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

// Save provides a mock function with given fields: ctxReq, data
func (_m *AttemptRepository) Save(ctxReq context.Context, data *model.LoginAttempt) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.LoginAttempt) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}
