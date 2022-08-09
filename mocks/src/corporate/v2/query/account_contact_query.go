// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	query "github.com/Bhinneka/user-service/src/corporate/v2/query"
	mock "github.com/stretchr/testify/mock"
)

// AccountContactQuery is an autogenerated mock type for the AccountContactQuery type
type AccountContactQuery struct {
	mock.Mock
}

// FindAccountMicrositeByContactID provides a mock function with given fields: id
func (_m *AccountContactQuery) FindAccountMicrositeByContactID(id int) <-chan query.ResultQuery {
	ret := _m.Called(id)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(int) <-chan query.ResultQuery); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// FindByAccountContactID provides a mock function with given fields: id
func (_m *AccountContactQuery) FindByAccountContactID(id int) <-chan query.ResultQuery {
	ret := _m.Called(id)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(int) <-chan query.ResultQuery); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// FindByAccountMicrositeContactID provides a mock function with given fields: id
func (_m *AccountContactQuery) FindByAccountMicrositeContactID(id int) <-chan query.ResultQuery {
	ret := _m.Called(id)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(int) <-chan query.ResultQuery); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}
