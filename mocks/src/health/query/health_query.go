// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	query "github.com/Bhinneka/user-service/src/health/query"
	mock "github.com/stretchr/testify/mock"
)

// HealthQuery is an autogenerated mock type for the HealthQuery type
type HealthQuery struct {
	mock.Mock
}

// Ping provides a mock function with given fields:
func (_m *HealthQuery) Ping() <-chan query.ResultQuery {
	ret := _m.Called()

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func() <-chan query.ResultQuery); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}
