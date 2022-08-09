package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/applications/v1/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
}

// ApplicationsUseCase interface abstraction
type ApplicationsUseCase interface {
	GetApplicationsList() <-chan ResultUseCase
	AddUpdateApplication(ctxReq context.Context, data model.Application) <-chan ResultUseCase
	DeleteApplication(ctxReq context.Context, id string) <-chan ResultUseCase
	GetListApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultUseCase
}
