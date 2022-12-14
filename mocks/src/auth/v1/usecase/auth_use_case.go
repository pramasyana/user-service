// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/auth/v1/model"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/Bhinneka/user-service/src/auth/v1/usecase"
)

// AuthUseCase is an autogenerated mock type for the AuthUseCase type
type AuthUseCase struct {
	mock.Mock
}

// CheckEmail provides a mock function with given fields: ctxReq, email
func (_m *AuthUseCase) CheckEmail(ctxReq context.Context, email string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, email)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CheckEmailV3 provides a mock function with given fields: ctxReq, data
func (_m *AuthUseCase) CheckEmailV3(ctxReq context.Context, data model.CheckEmailPayload) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.CheckEmailPayload) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CreateClientApp provides a mock function with given fields: name
func (_m *AuthUseCase) CreateClientApp(name string) <-chan usecase.ResultUseCase {
	ret := _m.Called(name)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GenerateToken provides a mock function with given fields: ctxReq, mode, data
func (_m *AuthUseCase) GenerateToken(ctxReq context.Context, mode string, data model.RequestToken) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, mode, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, model.RequestToken) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, mode, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GenerateTokenB2b provides a mock function with given fields: ctxReq, mode, data
func (_m *AuthUseCase) GenerateTokenB2b(ctxReq context.Context, mode string, data model.RequestToken) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, mode, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, model.RequestToken) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, mode, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GenerateTokenFromUserID provides a mock function with given fields: ctxReq, data
func (_m *AuthUseCase) GenerateTokenFromUserID(ctxReq context.Context, data model.RequestToken) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.RequestToken) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetClientApp provides a mock function with given fields: clientID, clientSecret
func (_m *AuthUseCase) GetClientApp(clientID string, clientSecret string) <-chan usecase.ResultUseCase {
	ret := _m.Called(clientID, clientSecret)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(clientID, clientSecret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetJTIToken provides a mock function with given fields: ctxReq, token, request
func (_m *AuthUseCase) GetJTIToken(ctxReq context.Context, token string, request string) (string, interface{}, error) {
	ret := _m.Called(ctxReq, token, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctxReq, token, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 interface{}
	if rf, ok := ret.Get(1).(func(context.Context, string, string) interface{}); ok {
		r1 = rf(ctxReq, token, request)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(interface{})
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctxReq, token, request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Logout provides a mock function with given fields: ctxReq, token
func (_m *AuthUseCase) Logout(ctxReq context.Context, token string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailWelcomeMember provides a mock function with given fields: ctxReq, data
func (_m *AuthUseCase) SendEmailWelcomeMember(ctxReq context.Context, data model.AccessTokenResponse) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.AccessTokenResponse) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ValidateBasicAuth provides a mock function with given fields: ctxReq, clientID, clientSecret
func (_m *AuthUseCase) ValidateBasicAuth(ctxReq context.Context, clientID string, clientSecret string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, clientID, clientSecret)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, clientID, clientSecret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// VerifyCaptcha provides a mock function with given fields: ctxReq, data
func (_m *AuthUseCase) VerifyCaptcha(ctxReq context.Context, data model.GoogleCaptcha) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.GoogleCaptcha) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// VerifyTokenMember provides a mock function with given fields: ctxReq, token
func (_m *AuthUseCase) VerifyTokenMember(ctxReq context.Context, token string) usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token)

	var r0 usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token)
	} else {
		r0 = ret.Get(0).(usecase.ResultUseCase)
	}

	return r0
}

// VerifyTokenMemberB2b provides a mock function with given fields: ctxReq, token, transaction_type, member_type
func (_m *AuthUseCase) VerifyTokenMemberB2b(ctxReq context.Context, token string, transaction_type string, member_type string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, transaction_type, member_type)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, transaction_type, member_type)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}
