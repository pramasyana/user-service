// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	query "github.com/Bhinneka/user-service/src/auth/v1/query"
	mock "github.com/stretchr/testify/mock"
)

// AuthQueryOA is an autogenerated mock type for the AuthQueryOA type
type AuthQueryOA struct {
	mock.Mock
}

// GetAppleToken provides a mock function with given fields: ctxReq, code, redirectURI, clientID
func (_m *AuthQueryOA) GetAppleToken(ctxReq context.Context, code string, redirectURI string, clientID string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code, redirectURI, clientID)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code, redirectURI, clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetAzureToken provides a mock function with given fields: ctxReq, code, redirectURI
func (_m *AuthQueryOA) GetAzureToken(ctxReq context.Context, code string, redirectURI string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code, redirectURI)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetDetailAzureMember provides a mock function with given fields: ctxReq, token
func (_m *AuthQueryOA) GetDetailAzureMember(ctxReq context.Context, token string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, token)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetDetailFacebookMember provides a mock function with given fields: ctxReq, code
func (_m *AuthQueryOA) GetDetailFacebookMember(ctxReq context.Context, code string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetDetailGoogleMember provides a mock function with given fields: ctxReq, code
func (_m *AuthQueryOA) GetDetailGoogleMember(ctxReq context.Context, code string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetFacebookToken provides a mock function with given fields: ctxReq, code, redirectURI
func (_m *AuthQueryOA) GetFacebookToken(ctxReq context.Context, code string, redirectURI string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code, redirectURI)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetGoogleToken provides a mock function with given fields: ctxReq, code, redirectURI
func (_m *AuthQueryOA) GetGoogleToken(ctxReq context.Context, code string, redirectURI string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, code, redirectURI)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, code, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}

// GetGoogleTokenInfo provides a mock function with given fields: ctxReq, token
func (_m *AuthQueryOA) GetGoogleTokenInfo(ctxReq context.Context, token string) <-chan query.ResultQuery {
	ret := _m.Called(ctxReq, token)

	var r0 <-chan query.ResultQuery
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan query.ResultQuery); ok {
		r0 = rf(ctxReq, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan query.ResultQuery)
		}
	}

	return r0
}
