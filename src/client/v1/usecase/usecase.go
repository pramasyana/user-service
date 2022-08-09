package usecase

import (
	"context"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sharedDomain "github.com/Bhinneka/user-service/src/shared/model"
)

// ClientUsecase usecase for client
type ClientUsecase interface {
	Logout(ctx context.Context, email string) <-chan sharedDomain.ResultUseCase
	GetMemberByEmail(ctxReq context.Context, email string) (member *memberModel.Member, statusCode int, err error)
}
