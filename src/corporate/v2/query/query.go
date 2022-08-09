package query

import (
	"context"

	"github.com/Bhinneka/user-service/src/corporate/v2/model"
)

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// ContactQuery interface abstraction
type ContactQuery interface {
	FindByEmail(ctxReq context.Context, email string) <-chan ResultQuery
	FindContactCorporateByEmail(ctxReq context.Context, email string) <-chan ResultQuery
	FindContactMicrositeByEmail(ctxReq context.Context, email, transactionType, memberType string) <-chan ResultQuery
	FindContactByEmail(ctxReq context.Context, email string) <-chan ResultQuery
	FindAccountByMemberType(ctxReq context.Context, memberType string) <-chan ResultQuery
	FindByID(ctxReq context.Context, uid string) <-chan ResultQuery
	GetListContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultQuery
	GetTotalContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultQuery
	GetTransactionType(ctxReq context.Context, email string) <-chan ResultQuery
}

// AccountContactQuery interface abstraction
type AccountContactQuery interface {
	FindByAccountContactID(id int) <-chan ResultQuery
	FindAccountMicrositeByContactID(id int) <-chan ResultQuery
	FindByAccountMicrositeContactID(id int) <-chan ResultQuery
}
