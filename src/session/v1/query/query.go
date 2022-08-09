package query

import (
	"context"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/session/v1/model"
)

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// SessionInfoQuery interface abstraction
type SessionInfoQuery interface {
	GetListSessionInfo(ctxReq context.Context, params *model.ParamList) <-chan ResultQuery
	GetTotalSessionInfo(params *model.ParamList) <-chan ResultQuery
	GetDetailSessionInfo(ctxReq context.Context, params model.ParametersGetSession) <-chan ResultQuery
	GetHistorySessionInfo(ctxReq context.Context, params *memberModel.ParametersLoginActivity) <-chan ResultQuery
	GetTotalHistorySessionInfo(ctxReq context.Context, params *memberModel.ParametersLoginActivity) <-chan ResultQuery
}
