// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/member/v1/model"
	mock "github.com/stretchr/testify/mock"

	servicemodel "github.com/Bhinneka/user-service/src/service/model"

	usecase "github.com/Bhinneka/user-service/src/member/v1/usecase"
)

// MemberUseCase is an autogenerated mock type for the MemberUseCase type
type MemberUseCase struct {
	mock.Mock
}

// ActivateMFASettingV3 provides a mock function with given fields: ctxReq, activateData
func (_m *MemberUseCase) ActivateMFASettingV3(ctxReq context.Context, activateData model.MFAActivateSettings) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, activateData)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.MFAActivateSettings) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, activateData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ActivateMFASettings provides a mock function with given fields: ctxReq, activateData
func (_m *MemberUseCase) ActivateMFASettings(ctxReq context.Context, activateData model.MFAActivateSettings) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, activateData)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.MFAActivateSettings) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, activateData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ActivateMember provides a mock function with given fields: ctxReq, token, requestFrom
func (_m *MemberUseCase) ActivateMember(ctxReq context.Context, token string, requestFrom string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, requestFrom)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, requestFrom)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ActivateMerchantEmployee provides a mock function with given fields: ctxReq, token
func (_m *MemberUseCase) ActivateMerchantEmployee(ctxReq context.Context, token string) <-chan usecase.ResultUseCase {
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

// ActivateNewPassword provides a mock function with given fields: ctxReq, token, password, rePassword
func (_m *MemberUseCase) ActivateNewPassword(ctxReq context.Context, token string, password string, rePassword string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, password, rePassword)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, password, rePassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// AddNewPassword provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) AddNewPassword(ctxReq context.Context, data model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// BulkImportMember provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) BulkImportMember(ctxReq context.Context, data []*model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// BulkValidateEmailAndPhone provides a mock function with given fields: ctxReq, members
func (_m *MemberUseCase) BulkValidateEmailAndPhone(ctxReq context.Context, members []*model.Member) []string {
	ret := _m.Called(ctxReq, members)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Member) []string); ok {
		r0 = rf(ctxReq, members)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// ChangeForgotPassword provides a mock function with given fields: ctxReq, token, newPassword, rePassword, requestFrom
func (_m *MemberUseCase) ChangeForgotPassword(ctxReq context.Context, token string, newPassword string, rePassword string, requestFrom string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, newPassword, rePassword, requestFrom)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, newPassword, rePassword, requestFrom)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ChangePassword provides a mock function with given fields: ctxReq, token, oldPassword, newPassword
func (_m *MemberUseCase) ChangePassword(ctxReq context.Context, token string, oldPassword string, newPassword string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, oldPassword, newPassword)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, oldPassword, newPassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CheckEmailAndMobileExistence provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) CheckEmailAndMobileExistence(ctxReq context.Context, data *model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CheckSendbirdToken provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) CheckSendbirdToken(ctxReq context.Context, params *servicemodel.SendbirdRequest) usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *servicemodel.SendbirdRequest) usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		r0 = ret.Get(0).(usecase.ResultUseCase)
	}

	return r0
}

// CheckSendbirdTokenV4 provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) CheckSendbirdTokenV4(ctxReq context.Context, params *servicemodel.SendbirdRequestV4) usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *servicemodel.SendbirdRequestV4) usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		r0 = ret.Get(0).(usecase.ResultUseCase)
	}

	return r0
}

// Clients provides a mock function with given fields: ctxReq, token
func (_m *MemberUseCase) Clients(ctxReq context.Context, token string) <-chan usecase.ResultUseCase {
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

// DisabledMFASetting provides a mock function with given fields: ctxReq, userID, requestFrom
func (_m *MemberUseCase) DisabledMFASetting(ctxReq context.Context, userID string, requestFrom string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, userID, requestFrom)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, userID, requestFrom)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ForgotPassword provides a mock function with given fields: ctxReq, email
func (_m *MemberUseCase) ForgotPassword(ctxReq context.Context, email string) <-chan usecase.ResultUseCase {
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

// GenerateMFASettings provides a mock function with given fields: ctxReq, userID, requestFrom
func (_m *MemberUseCase) GenerateMFASettings(ctxReq context.Context, userID string, requestFrom string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, userID, requestFrom)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, userID, requestFrom)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetDetailMemberByEmail provides a mock function with given fields: email
func (_m *MemberUseCase) GetDetailMemberByEmail(email string) <-chan usecase.ResultUseCase {
	ret := _m.Called(email)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetDetailMemberByID provides a mock function with given fields: ctxReq, uid
func (_m *MemberUseCase) GetDetailMemberByID(ctxReq context.Context, uid string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetListMembers provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) GetListMembers(ctxReq context.Context, params *model.Parameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Parameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetLoginActivity provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) GetLoginActivity(ctxReq context.Context, params *model.ParametersLoginActivity) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.ParametersLoginActivity) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMFASettings provides a mock function with given fields: ctxReq, uid
func (_m *MemberUseCase) GetMFASettings(ctxReq context.Context, uid string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetNarwhalMFASettings provides a mock function with given fields: ctxReq, uid
func (_m *MemberUseCase) GetNarwhalMFASettings(ctxReq context.Context, uid string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetProfileComplete provides a mock function with given fields: ctxReq, uid
func (_m *MemberUseCase) GetProfileComplete(ctxReq context.Context, uid string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetSendbirdToken provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) GetSendbirdToken(ctxReq context.Context, params *servicemodel.SendbirdRequest) usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *servicemodel.SendbirdRequest) usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		r0 = ret.Get(0).(usecase.ResultUseCase)
	}

	return r0
}

// GetSendbirdTokenV4 provides a mock function with given fields: ctxReq, params
func (_m *MemberUseCase) GetSendbirdTokenV4(ctxReq context.Context, params *servicemodel.SendbirdRequestV4) usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *servicemodel.SendbirdRequestV4) usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		r0 = ret.Get(0).(usecase.ResultUseCase)
	}

	return r0
}

// ImportMember provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) ImportMember(ctxReq context.Context, data *model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// InsertLogMember provides a mock function with given fields: ctxReq, oldData, newData, action
func (_m *MemberUseCase) InsertLogMember(ctxReq context.Context, oldData *model.Member, newData *model.Member, action string) error {
	ret := _m.Called(ctxReq, oldData, newData, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member, *model.Member, string) error); ok {
		r0 = rf(ctxReq, oldData, newData, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MigrateMember provides a mock function with given fields: ctxReq, members
func (_m *MemberUseCase) MigrateMember(ctxReq context.Context, members *model.Members) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, members)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Members) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, members)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ParseMemberData provides a mock function with given fields: ctxReq, input
func (_m *MemberUseCase) ParseMemberData(ctxReq context.Context, input []byte) ([]*model.Member, error) {
	ret := _m.Called(ctxReq, input)

	var r0 []*model.Member
	if rf, ok := ret.Get(0).(func(context.Context, []byte) []*model.Member); ok {
		r0 = rf(ctxReq, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctxReq, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishToKafkaUser provides a mock function with given fields: ctxReq, data, eventType
func (_m *MemberUseCase) PublishToKafkaUser(ctxReq context.Context, data *model.Member, eventType string) error {
	ret := _m.Called(ctxReq, data, eventType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member, string) error); ok {
		r0 = rf(ctxReq, data, eventType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegenerateToken provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) RegenerateToken(ctxReq context.Context, data model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// RegisterMember provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) RegisterMember(ctxReq context.Context, data *model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ResendActivation provides a mock function with given fields: ctxReq, email
func (_m *MemberUseCase) ResendActivation(ctxReq context.Context, email string) <-chan usecase.ResultUseCase {
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

// RevokeAccess provides a mock function with given fields: ctxReq, uid, jti
func (_m *MemberUseCase) RevokeAccess(ctxReq context.Context, uid string, jti string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid, jti)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid, jti)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// RevokeAllAccess provides a mock function with given fields: ctxReq, uid, token
func (_m *MemberUseCase) RevokeAllAccess(ctxReq context.Context, uid string, token string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, uid, token)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, uid, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailAddMember provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) SendEmailAddMember(ctxReq context.Context, data model.SuccessResponse) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.SuccessResponse) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailForgotPassword provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) SendEmailForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.SuccessResponse) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailRegisterMember provides a mock function with given fields: ctxReq, data, registerType
func (_m *MemberUseCase) SendEmailRegisterMember(ctxReq context.Context, data model.SuccessResponse, registerType string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, registerType)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.SuccessResponse, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, registerType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailSuccessForgotPassword provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) SendEmailSuccessForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.SuccessResponse) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailWelcomeMember provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) SendEmailWelcomeMember(ctxReq context.Context, data model.SuccessResponse) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.SuccessResponse) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SyncPassword provides a mock function with given fields: ctxReq, token, oldPassword, newPassword
func (_m *MemberUseCase) SyncPassword(ctxReq context.Context, token string, oldPassword string, newPassword string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, oldPassword, newPassword)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, oldPassword, newPassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdateDetailMemberByID provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) UpdateDetailMemberByID(ctxReq context.Context, data model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdatePassword provides a mock function with given fields: ctxReq, token, uid, oldPassword, newPassword
func (_m *MemberUseCase) UpdatePassword(ctxReq context.Context, token string, uid string, oldPassword string, newPassword string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, uid, oldPassword, newPassword)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, uid, oldPassword, newPassword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdateProfileName provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) UpdateProfileName(ctxReq context.Context, data model.ProfileName) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.ProfileName) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdateProfilePicture provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) UpdateProfilePicture(ctxReq context.Context, data model.ProfilePicture) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.ProfilePicture) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ValidateEmailAndPhone provides a mock function with given fields: ctxReq, data
func (_m *MemberUseCase) ValidateEmailAndPhone(ctxReq context.Context, data *model.Member) error {
	ret := _m.Called(ctxReq, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Member) error); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateEmailDomain provides a mock function with given fields: ctxReq, email
func (_m *MemberUseCase) ValidateEmailDomain(ctxReq context.Context, email string) <-chan usecase.ResultUseCase {
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

// ValidateToken provides a mock function with given fields: ctxReq, token
func (_m *MemberUseCase) ValidateToken(ctxReq context.Context, token string) <-chan usecase.ResultUseCase {
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
