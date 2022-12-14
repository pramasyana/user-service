// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/document/v2/model"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/Bhinneka/user-service/src/document/v2/repo"
)

// DocumentTypeRepository is an autogenerated mock type for the DocumentTypeRepository type
type DocumentTypeRepository struct {
	mock.Mock
}

// AddDocumentType provides a mock function with given fields: ctxReq, data
func (_m *DocumentTypeRepository) AddDocumentType(ctxReq context.Context, data model.DocumentType) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.DocumentType) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// FindDocumentTypeByParam provides a mock function with given fields: ctxReq, param
func (_m *DocumentTypeRepository) FindDocumentTypeByParam(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.DocumentTypeParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// GetListDocumentType provides a mock function with given fields: ctxReq, param
func (_m *DocumentTypeRepository) GetListDocumentType(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.DocumentTypeParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// GetTotalDocumentType provides a mock function with given fields: ctxReq, params
func (_m *DocumentTypeRepository) GetTotalDocumentType(ctxReq context.Context, params *model.DocumentTypeParameters) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, params)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, *model.DocumentTypeParameters) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}

// UpdateDocumentType provides a mock function with given fields: ctxReq, data
func (_m *DocumentTypeRepository) UpdateDocumentType(ctxReq context.Context, data model.DocumentType) <-chan repo.ResultRepository {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan repo.ResultRepository
	if rf, ok := ret.Get(0).(func(context.Context, model.DocumentType) <-chan repo.ResultRepository); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan repo.ResultRepository)
		}
	}

	return r0
}
