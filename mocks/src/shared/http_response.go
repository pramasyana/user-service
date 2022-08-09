// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo"
	mock "github.com/stretchr/testify/mock"
)

// HTTPResponse is an autogenerated mock type for the HTTPResponse type
type HTTPResponse struct {
	mock.Mock
}

// JSON provides a mock function with given fields: c
func (_m *HTTPResponse) JSON(c echo.Context) error {
	ret := _m.Called(c)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetSuccess provides a mock function with given fields: isSuccess
func (_m *HTTPResponse) SetSuccess(isSuccess bool) {
	_m.Called(isSuccess)
}
