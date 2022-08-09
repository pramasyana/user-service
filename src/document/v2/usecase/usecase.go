package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/document/v2/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.DocumentError
}

// DocumentUseCase interface abstraction
type DocumentUseCase interface {
	AddUpdateDocument(ctxReq context.Context, data model.DocumentData) <-chan ResultUseCase
	DeleteDocument(ctxReq context.Context, documentID string, memberID string) <-chan ResultUseCase
	GetListDocument(ctxReq context.Context, param *model.DocumentParameters) <-chan ResultUseCase
	GetDetailDocument(ctxReq context.Context, documentID string, memberID string) <-chan ResultUseCase
	AddUpdateDocumentType(ctxReq context.Context, data model.DocumentType) <-chan ResultUseCase
	GetListDocumentType(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan ResultUseCase
	GetRequiredDocument(ctxReq context.Context) <-chan ResultUseCase
}
