package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/document/v2/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// DocumentRepository interface abstraction
type DocumentRepository interface {
	AddDocument(ctxReq context.Context, data model.DocumentData) <-chan ResultRepository
	UpdateDocument(ctxReq context.Context, data model.DocumentData) <-chan ResultRepository
	FindDocumentByParam(ctxReq context.Context, param *model.DocumentParameters) <-chan ResultRepository
	DeleteDocumentByID(ctxReq context.Context, documentID string) <-chan ResultRepository
	GetListDocument(ctxReq context.Context, param *model.DocumentParameters) <-chan ResultRepository
	GetTotalDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository
	GetDetailDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository
}

// DocumentTypeRepository interface abstraction
type DocumentTypeRepository interface {
	AddDocumentType(ctxReq context.Context, data model.DocumentType) <-chan ResultRepository
	UpdateDocumentType(ctxReq context.Context, data model.DocumentType) <-chan ResultRepository
	FindDocumentTypeByParam(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan ResultRepository
	GetListDocumentType(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan ResultRepository
	GetTotalDocumentType(ctxReq context.Context, params *model.DocumentTypeParameters) <-chan ResultRepository
}
