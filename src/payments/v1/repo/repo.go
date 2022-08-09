package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/payments/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// MemberRepository interface abstraction
type PaymentsRepository interface {
	AddUpdatePayment(ctxReq context.Context, payments model.Payments) <-chan ResultRepository
	FindPaymentByEmailChannelMethod(ctxReq context.Context, email, channel, method string) ResultRepository
}
