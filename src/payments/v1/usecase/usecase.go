package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/payments/v1/model"
)

type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.PaymentsError
}

// MemberUseCase interface abstraction
type PaymentsUseCase interface {
	AddUpdatePayments(ctxReq context.Context, data *model.Payments) <-chan ResultUseCase
	CompareHeaderAndBody(ctxReq context.Context, data *model.Payments, basicAuth string) <-chan ResultUseCase
	GetPaymentDetail(ctxReq context.Context, data *model.Payments) <-chan ResultUseCase
}
