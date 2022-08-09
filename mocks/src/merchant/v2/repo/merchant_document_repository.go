// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/merchant/v2/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
)

// MerchantDocumentRepository is an autogenerated mock type for the MerchantDocumentRepository type
type MerchantDocumentRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0
func (_m *MerchantDocumentRepository) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindMerchantDocumentByParam provides a mock function with given fields: ctxReq, param
func (_m *MerchantDocumentRepository) FindMerchantDocumentByParam(ctxReq context.Context, param *model.B2CMerchantDocumentQueryInput) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantDocumentQueryInput) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// GetListMerchantDocument provides a mock function with given fields: ctxReq, params
func (_m *MerchantDocumentRepository) GetListMerchantDocument(ctxReq context.Context, params *model.B2CMerchantDocumentQueryInput) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantDocumentQueryInput) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// InsertNewMerchantDocument provides a mock function with given fields: ctxReq, param
func (_m *MerchantDocumentRepository) InsertNewMerchantDocument(ctxReq context.Context, param *model.B2CMerchantDocumentData) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.B2CMerchantDocumentData) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// ResetRejectedDocument provides a mock function with given fields: ctxReq, param
func (_m *MerchantDocumentRepository) ResetRejectedDocument(ctxReq context.Context, param model.B2CMerchantDocumentData) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.B2CMerchantDocumentData) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *MerchantDocumentRepository) Save(_a0 model.B2CMerchantDocument) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.B2CMerchantDocument) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMerchantDocumentGWS provides a mock function with given fields: _a0
func (_m *MerchantDocumentRepository) SaveMerchantDocumentGWS(_a0 model.B2CMerchantDocument) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.B2CMerchantDocument) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMerchantDocument provides a mock function with given fields: ctxReq, id, param
func (_m *MerchantDocumentRepository) UpdateMerchantDocument(ctxReq context.Context, id string, param *model.B2CMerchantDocumentData) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, id, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.B2CMerchantDocumentData) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, id, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}
