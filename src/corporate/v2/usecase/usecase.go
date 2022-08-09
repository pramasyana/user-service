package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/corporate/v2/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.CorporateError
}

// CorporateUseCase interface abstraction
type CorporateUseCase interface {
	GetAllListContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultUseCase
	GetDetailContact(ctxReq context.Context, id string) <-chan ResultUseCase
	ImportContact(ctxReq context.Context, content []byte) ([]*model.ContactPayload, error)
}
