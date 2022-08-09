package usecase

import (
	"context"

	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// LogUsecase return log usecase
type LogUsecase interface {
	GetAll(ctx context.Context, param *sharedModel.Parameters) <-chan sharedModel.ResultUseCase
	GetByID(ctx context.Context, logID string) <-chan sharedModel.ResultUseCase
}
