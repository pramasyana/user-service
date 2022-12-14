// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/merchant/v2/model"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/Bhinneka/user-service/src/merchant/v2/usecase"

	v1model "github.com/Bhinneka/user-service/src/member/v1/model"
)

// MerchantUseCase is an autogenerated mock type for the MerchantUseCase type
type MerchantUseCase struct {
	mock.Mock
}

// AddEmployee provides a mock function with given fields: ctxReq, token, email, firstName
func (_m *MerchantUseCase) AddEmployee(ctxReq context.Context, token string, email string, firstName string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, email, firstName)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, email, firstName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// AddMerchant provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) AddMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// AddMerchantPIC provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) AddMerchantPIC(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ChangeMerchantName provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) ChangeMerchantName(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CheckMerchantName provides a mock function with given fields: ctxReq, merchantName
func (_m *MerchantUseCase) CheckMerchantName(ctxReq context.Context, merchantName string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, merchantName)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, merchantName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// ClearRejectUpgrade provides a mock function with given fields: ctxReq, memberID, userAttribute
func (_m *MerchantUseCase) ClearRejectUpgrade(ctxReq context.Context, memberID string, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, memberID, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, memberID, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CmsGetAllMerchantEmployee provides a mock function with given fields: ctxReq, token, params
func (_m *MerchantUseCase) CmsGetAllMerchantEmployee(ctxReq context.Context, token string, params *model.QueryCmsMerchantEmployeeParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.QueryCmsMerchantEmployeeParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// CreateMerchant provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) CreateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// DeleteMerchant provides a mock function with given fields: ctxReq, merchantID, userAttribute
func (_m *MerchantUseCase) DeleteMerchant(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, merchantID, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, merchantID, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetAllMerchantEmployee provides a mock function with given fields: ctxReq, token, params
func (_m *MerchantUseCase) GetAllMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetListMerchantBank provides a mock function with given fields: params
func (_m *MerchantUseCase) GetListMerchantBank(params *model.ParametersMerchantBank) <-chan usecase.ResultUseCase {
	ret := _m.Called(params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(*model.ParametersMerchantBank) <-chan usecase.ResultUseCase); ok {
		r0 = rf(params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchantByID provides a mock function with given fields: ctxReq, id, privacy, isAttachment
func (_m *MerchantUseCase) GetMerchantByID(ctxReq context.Context, id string, privacy string, isAttachment string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, id, privacy, isAttachment)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, id, privacy, isAttachment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchantByUserID provides a mock function with given fields: ctxReq, userID, token
func (_m *MerchantUseCase) GetMerchantByUserID(ctxReq context.Context, userID string, token string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, userID, token)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, userID, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchantByVanityURL provides a mock function with given fields: ctxReq, vanityURL
func (_m *MerchantUseCase) GetMerchantByVanityURL(ctxReq context.Context, vanityURL string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, vanityURL)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, vanityURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchantEmployee provides a mock function with given fields: ctxReq, token, params
func (_m *MerchantUseCase) GetMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchants provides a mock function with given fields: ctxReq, params
func (_m *MerchantUseCase) GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.QueryParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetMerchantsPublic provides a mock function with given fields: ctxReq, params
func (_m *MerchantUseCase) GetMerchantsPublic(ctxReq context.Context, params *model.QueryParametersPublic) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.QueryParametersPublic) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// InsertLogMerchant provides a mock function with given fields: ctxReq, old, new, action
func (_m *MerchantUseCase) InsertLogMerchant(ctxReq context.Context, old model.B2CMerchantDataV2, new model.B2CMerchantDataV2, action string) error {
	ret := _m.Called(ctxReq, old, new, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, model.B2CMerchantDataV2, string) error); ok {
		r0 = rf(ctxReq, old, new, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PublishToKafkaMerchant provides a mock function with given fields: ctxReq, data, eventType
func (_m *MerchantUseCase) PublishToKafkaMerchant(ctxReq context.Context, data model.B2CMerchantDataV2, eventType string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, eventType)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, eventType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// RejectMerchantRegistration provides a mock function with given fields: ctxReq, merchantID, userAttribute
func (_m *MerchantUseCase) RejectMerchantRegistration(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, merchantID, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, merchantID, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// RejectMerchantUpgrade provides a mock function with given fields: ctxReq, merchantID, userAttribute, reasonReject
func (_m *MerchantUseCase) RejectMerchantUpgrade(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute, reasonReject string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, merchantID, userAttribute, reasonReject)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.MerchantUserAttribute, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, merchantID, userAttribute, reasonReject)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SelfUpdateMerchant provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) SelfUpdateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SelfUpdateMerchantPartial provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) SelfUpdateMerchantPartial(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailActivation provides a mock function with given fields: ctxReq, merchant
func (_m *MerchantUseCase) SendEmailActivation(ctxReq context.Context, merchant model.B2CMerchantDataV2) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, merchant)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, merchant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailAdmin provides a mock function with given fields: ctxReq, data, memberName, reasonReject, adminCMS
func (_m *MerchantUseCase) SendEmailAdmin(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string, reasonReject string, adminCMS string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, memberName, reasonReject, adminCMS)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, memberName, reasonReject, adminCMS)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailApproval provides a mock function with given fields: ctxReq, old
func (_m *MerchantUseCase) SendEmailApproval(ctxReq context.Context, old model.B2CMerchantDataV2) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, old)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, old)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantAdd provides a mock function with given fields: ctxReq, data, memberName
func (_m *MerchantUseCase) SendEmailMerchantAdd(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, memberName)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, memberName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantEmployeeLogin provides a mock function with given fields: ctxReq, dataMerchant, dataMember
func (_m *MerchantUseCase) SendEmailMerchantEmployeeLogin(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember v1model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, dataMerchant, dataMember)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, v1model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, dataMerchant, dataMember)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantEmployeeRegister provides a mock function with given fields: ctxReq, dataMerchant, dataMember
func (_m *MerchantUseCase) SendEmailMerchantEmployeeRegister(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember v1model.Member) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, dataMerchant, dataMember)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, v1model.Member) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, dataMerchant, dataMember)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantRejectRegistration provides a mock function with given fields: ctxReq, data, memberName
func (_m *MerchantUseCase) SendEmailMerchantRejectRegistration(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, memberName)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, memberName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantRejectUpgrade provides a mock function with given fields: ctxReq, data, memberName, reasonReject
func (_m *MerchantUseCase) SendEmailMerchantRejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string, reasonReject string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, memberName, reasonReject)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, memberName, reasonReject)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// SendEmailMerchantUpgrade provides a mock function with given fields: ctxReq, data, memberName
func (_m *MerchantUseCase) SendEmailMerchantUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, memberName)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, memberName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdateMerchant provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) UpdateMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpdateMerchantEmployee provides a mock function with given fields: ctxReq, token, params
func (_m *MerchantUseCase) UpdateMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, token, params)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.QueryMerchantEmployeeParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, token, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// UpgradeMerchant provides a mock function with given fields: ctxReq, data, userAttribute
func (_m *MerchantUseCase) UpgradeMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data, userAttribute)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantCreateInput, *model.MerchantUserAttribute) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data, userAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}
