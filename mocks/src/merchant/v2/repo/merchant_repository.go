// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/merchant/v2/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
)

// MerchantRepository is an autogenerated mock type for the MerchantRepository type
type MerchantRepository struct {
	mock.Mock
}

// AddUpdateMerchant provides a mock function with given fields: ctxReq, data
func (_m *MerchantRepository) AddUpdateMerchant(ctxReq context.Context, data model.B2CMerchantDataV2) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// ClearRejectUpgrade provides a mock function with given fields: ctxReq, data
func (_m *MerchantRepository) ClearRejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: _a0
func (_m *MerchantRepository) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindMerchantByEmail provides a mock function with given fields: ctxReq, uid
func (_m *MerchantRepository) FindMerchantByEmail(ctxReq context.Context, uid string) repo.ResultRepository {
	ret := _m.Called(ctxReq, uid)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, uid)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// FindMerchantByID provides a mock function with given fields: ctxReq, id, uid
func (_m *MerchantRepository) FindMerchantByID(ctxReq context.Context, id string, uid string) repo.ResultRepository {
	ret := _m.Called(ctxReq, id, uid)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, id, uid)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// FindMerchantByName provides a mock function with given fields: ctxReq, uid
func (_m *MerchantRepository) FindMerchantByName(ctxReq context.Context, uid string) repo.ResultRepository {
	ret := _m.Called(ctxReq, uid)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, uid)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// FindMerchantBySlug provides a mock function with given fields: ctxReq, slug
func (_m *MerchantRepository) FindMerchantBySlug(ctxReq context.Context, slug string) repo.ResultRepository {
	ret := _m.Called(ctxReq, slug)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, slug)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// FindMerchantByUser provides a mock function with given fields: ctxReq, uid
func (_m *MerchantRepository) FindMerchantByUser(ctxReq context.Context, uid string) repo.ResultRepository {
	ret := _m.Called(ctxReq, uid)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, uid)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// GetMerchants provides a mock function with given fields: ctxReq, params
func (_m *MerchantRepository) GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.QueryParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// GetMerchantsPublic provides a mock function with given fields: ctxReq, params
func (_m *MerchantRepository) GetMerchantsPublic(ctxReq context.Context, params *model.QueryParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.QueryParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// GetTotalMerchant provides a mock function with given fields: ctxReq, params
func (_m *MerchantRepository) GetTotalMerchant(ctxReq context.Context, params *model.QueryParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.QueryParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// LoadCompanySize provides a mock function with given fields: ctxReq, id
func (_m *MerchantRepository) LoadCompanySize(ctxReq context.Context, id int) repo.ResultRepository {
	ret := _m.Called(ctxReq, id)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, int) repo.ResultRepository); ok {
		r0 = rf(ctxReq, id)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// LoadLegalEntity provides a mock function with given fields: ctxReq, id
func (_m *MerchantRepository) LoadLegalEntity(ctxReq context.Context, id int) repo.ResultRepository {
	ret := _m.Called(ctxReq, id)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, int) repo.ResultRepository); ok {
		r0 = rf(ctxReq, id)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// LoadMerchant provides a mock function with given fields: ctxReq, uid, privacy
func (_m *MerchantRepository) LoadMerchant(ctxReq context.Context, uid string, privacy string) repo.ResultRepository {
	ret := _m.Called(ctxReq, uid, privacy)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, uid, privacy)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// LoadMerchantByVanityURL provides a mock function with given fields: ctxReq, vanityURL
func (_m *MerchantRepository) LoadMerchantByVanityURL(ctxReq context.Context, vanityURL string) repo.ResultRepository {
	ret := _m.Called(ctxReq, vanityURL)

	var r0 repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.ResultRepository); ok {
		r0 = rf(ctxReq, vanityURL)
	} else {
		r0 = ret.Get(0).(repo.ResultRepository)
	}

	return r0
}

// RejectUpgrade provides a mock function with given fields: ctxReq, data, reasonReject
func (_m *MerchantRepository) RejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, reasonReject string) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data, reasonReject)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDataV2, string) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data, reasonReject)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *MerchantRepository) Save(_a0 model.B2CMerchant) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.B2CMerchant) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMerchantGWS provides a mock function with given fields: ctxReq, data
func (_m *MerchantRepository) SaveMerchantGWS(ctxReq context.Context, data model.B2CMerchant) error {
	ret := _m.Called(ctxReq, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchant) error); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SoftDelete provides a mock function with given fields: ctxReq, merchantID
func (_m *MerchantRepository) SoftDelete(ctxReq context.Context, merchantID string) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, merchantID)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, merchantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// UpdateMerchantGWS provides a mock function with given fields: ctxReq, data
func (_m *MerchantRepository) UpdateMerchantGWS(ctxReq context.Context, data model.B2CMerchant) error {
	ret := _m.Called(ctxReq, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchant) error); ok {
		r0 = rf(ctxReq, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
