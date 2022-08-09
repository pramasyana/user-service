package usecase

import (
	"context"
	"net/http"

	"github.com/Bhinneka/golib/tracer"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/src/service"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// LogUseCaseImpl DI
type LogUseCaseImpl struct {
	ActivityService service.ActivityServices
}

// NewLogUsecase return implementation
func NewLogUsecase(services localConfig.ServiceShared) *LogUseCaseImpl {
	return &LogUseCaseImpl{
		ActivityService: services.ActivityService,
	}
}

// GetAll retrieve all log from as
func (u *LogUseCaseImpl) GetAll(ctxReq context.Context, param *sharedModel.Parameters) <-chan sharedModel.ResultUseCase {
	ctx := "LogUseCase-GetAll"
	output := make(chan sharedModel.ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(context.Context, map[string]interface{}) {
		defer close(output)
		serviceResult := <-u.ActivityService.GetAll(ctxReq, param)
		if serviceResult.Error != nil {
			output <- sharedModel.ResultUseCase{Error: serviceResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		output <- sharedModel.ResultUseCase{Result: serviceResult.Result, Meta: serviceResult.Meta}
	})
	return output
}

// GetByID retrieve single log by log ID
func (u *LogUseCaseImpl) GetByID(ctxReq context.Context, logID string) <-chan sharedModel.ResultUseCase {
	ctx := "LogUseCase-GetByID"
	output := make(chan sharedModel.ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(context.Context, map[string]interface{}) {
		defer close(output)
		serviceResult := <-u.ActivityService.GetLogByID(ctxReq, logID)
		if serviceResult.Error != nil {
			output <- sharedModel.ResultUseCase{Error: serviceResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		output <- sharedModel.ResultUseCase{Result: serviceResult.Result}
	})
	return output
}
