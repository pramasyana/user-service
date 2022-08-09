// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/Bhinneka/user-service/src/document/v2/model"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/Bhinneka/user-service/src/document/v2/usecase"
)

// DocumentUseCase is an autogenerated mock type for the DocumentUseCase type
type DocumentUseCase struct {
	mock.Mock
}

// AddUpdateDocument provides a mock function with given fields: ctxReq, data
func (_m *DocumentUseCase) AddUpdateDocument(ctxReq context.Context, data model.DocumentData) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.DocumentData) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// AddUpdateDocumentType provides a mock function with given fields: ctxReq, data
func (_m *DocumentUseCase) AddUpdateDocumentType(ctxReq context.Context, data model.DocumentType) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, data)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, model.DocumentType) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// DeleteDocument provides a mock function with given fields: ctxReq, documentID, memberID
func (_m *DocumentUseCase) DeleteDocument(ctxReq context.Context, documentID string, memberID string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, documentID, memberID)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, documentID, memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetDetailDocument provides a mock function with given fields: ctxReq, documentID, memberID
func (_m *DocumentUseCase) GetDetailDocument(ctxReq context.Context, documentID string, memberID string) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, documentID, memberID)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, documentID, memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetListDocument provides a mock function with given fields: ctxReq, param
func (_m *DocumentUseCase) GetListDocument(ctxReq context.Context, param *model.DocumentParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.DocumentParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetListDocumentType provides a mock function with given fields: ctxReq, param
func (_m *DocumentUseCase) GetListDocumentType(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq, param)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context, *model.DocumentTypeParameters) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}

// GetRequiredDocument provides a mock function with given fields: ctxReq
func (_m *DocumentUseCase) GetRequiredDocument(ctxReq context.Context) <-chan usecase.ResultUseCase {
	ret := _m.Called(ctxReq)

	var r0 <-chan usecase.ResultUseCase
	if rf, ok := ret.Get(0).(func(context.Context) <-chan usecase.ResultUseCase); ok {
		r0 = rf(ctxReq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan usecase.ResultUseCase)
		}
	}

	return r0
}
