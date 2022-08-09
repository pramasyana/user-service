package query

import (
	"context"

	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// MemberQuery interface abstraction
type MemberQuery interface {
	FindByID(ctxReq context.Context, email string) <-chan ResultQuery
	FindByEmail(ctxReq context.Context, email string) <-chan ResultQuery
	FindByMobile(ctxReq context.Context, mobile string) <-chan ResultQuery
	FindMaxID(ctxReq context.Context) <-chan ResultQuery
	FindByToken(ctxReq context.Context, token string) <-chan ResultQuery
	UpdateBlockedMember(email string) <-chan ResultQuery
	UnblockMember(email string) <-chan ResultQuery
	GetListMembers(ctxReq context.Context, params *model.Parameters) <-chan ResultQuery
	GetTotalMembers(params *model.Parameters) <-chan ResultQuery
	UpdateLastTokenAttempt(ctxReq context.Context, email string) <-chan ResultQuery
	BulkFindByEmail(ctxReq context.Context, emails []string) <-chan ResultQuery
}

// MemberMFAQuery interface abstraction
type MemberMFAQuery interface {
	FindMFASettings(ctxReq context.Context, uid string) <-chan ResultQuery
	FindNarwhalMFASettings(ctxReq context.Context, uid string) <-chan ResultQuery
}
