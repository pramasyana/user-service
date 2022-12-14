// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import token "github.com/Bhinneka/user-service/src/auth/v1/token"

// AccessTokenGenerator is an autogenerated mock type for the AccessTokenGenerator type
type AccessTokenGenerator struct {
	mock.Mock
}

// GenerateAccessToken provides a mock function with given fields: cl
func (_m *AccessTokenGenerator) GenerateAccessToken(cl token.Claim) <-chan token.AccessTokenResponse {
	ret := _m.Called(cl)

	var r0 <-chan token.AccessTokenResponse
	if rf, ok := ret.Get(0).(func(token.Claim) <-chan token.AccessTokenResponse); ok {
		r0 = rf(cl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan token.AccessTokenResponse)
		}
	}

	return r0
}

// GenerateAnonymous provides a mock function with given fields: ctxReq
func (_m *AccessTokenGenerator) GenerateAnonymous(ctxReq context.Context) <-chan token.AccessTokenResponse {
	ret := _m.Called(ctxReq)

	var r0 <-chan token.AccessTokenResponse
	if rf, ok := ret.Get(0).(func(context.Context) <-chan token.AccessTokenResponse); ok {
		r0 = rf(ctxReq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan token.AccessTokenResponse)
		}
	}

	return r0
}
