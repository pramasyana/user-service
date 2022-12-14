// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	model "github.com/Bhinneka/user-service/src/auth/v1/model"
	repo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	mock "github.com/stretchr/testify/mock"
)

// ClientAppRepository is an autogenerated mock type for the ClientAppRepository type
type ClientAppRepository struct {
	mock.Mock
}

// FindByClientID provides a mock function with given fields: _a0
func (_m *ClientAppRepository) FindByClientID(_a0 string) <-chan repo.ResultRepository {
	ret := _m.Called(_a0)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(string) <-chan repo.ResultRepository); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Load provides a mock function with given fields: _a0
func (_m *ClientAppRepository) Load(_a0 int) <-chan repo.ResultRepository {
	ret := _m.Called(_a0)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(int) <-chan repo.ResultRepository); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *ClientAppRepository) Save(_a0 *model.ClientApp) <-chan repo.ResultRepository {
	ret := _m.Called(_a0)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(*model.ClientApp) <-chan repo.ResultRepository); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}
