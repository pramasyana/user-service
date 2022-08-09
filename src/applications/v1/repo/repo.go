package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/applications/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// ApplicationRepository interface abstraction
type ApplicationRepository interface {
	Save(ctxReq context.Context, data model.Application) <-chan ResultRepository
	Update(ctxReq context.Context, data model.Application) <-chan ResultRepository
	FindApplicationByID(ctxReq context.Context, id string) <-chan ResultRepository
	Delete(ctxReq context.Context, id string) <-chan ResultRepository
	GetListApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultRepository
	GetTotalApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultRepository
}
