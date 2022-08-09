package usecase

import (
	"context"

	sharedDomain "github.com/Bhinneka/user-service/src/shared/model"
)

// ClientUsecase usecase for client
type ClientUsecase interface {
	Logout(ctx context.Context, email string) <-chan sharedDomain.ResultUseCase
}
