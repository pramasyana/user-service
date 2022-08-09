package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	HTTPStatus int
	Error      error
}

// AuthUseCase interface abstraction
type AuthUseCase interface {
	CreateClientApp(name string) <-chan ResultUseCase
	GetClientApp(clientID, clientSecret string) <-chan ResultUseCase
	GenerateToken(ctxReq context.Context, mode string, data model.RequestToken) <-chan ResultUseCase
	GenerateTokenB2b(ctxReq context.Context, mode string, data model.RequestToken) <-chan ResultUseCase
	GenerateTokenFromUserID(ctxReq context.Context, data model.RequestToken) <-chan ResultUseCase
	VerifyTokenMember(ctxReq context.Context, token string) (result ResultUseCase)
	VerifyTokenMemberB2b(ctxReq context.Context, token, transaction_type, member_type string) <-chan ResultUseCase
	Logout(ctxReq context.Context, token string) <-chan ResultUseCase
	CheckEmail(ctxReq context.Context, email string) <-chan ResultUseCase
	CheckEmailV3(ctxReq context.Context, data model.CheckEmailPayload) <-chan ResultUseCase
	VerifyCaptcha(ctxReq context.Context, data model.GoogleCaptcha) <-chan ResultUseCase
	SendEmailWelcomeMember(ctxReq context.Context, data model.AccessTokenResponse) <-chan ResultUseCase
	ValidateBasicAuth(ctxReq context.Context, clientID, clientSecret string) <-chan ResultUseCase
	GetJTIToken(ctxReq context.Context, token, request string) (string, interface{}, error)
}
