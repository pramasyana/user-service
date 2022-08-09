package usecase

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// GetNarwhalMFASettings function for getting detail member based on member id
func (mu *MemberUseCaseImpl) GetNarwhalMFASettings(ctxReq context.Context, uid string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetNarwhalMFASettings"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = uid
		if !strings.Contains(uid, usrFormat) {
			output <- ResultUseCase{Error: fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID), HTTPStatus: http.StatusBadRequest}
			return
		}

		memberNarwhal := <-mu.MemberMFAQueryRead.FindNarwhalMFASettings(ctxReq, uid)
		if memberNarwhal.Error != nil {
			tracer.SetError(ctxReq, memberNarwhal.Error)
			output <- ResultUseCase{Error: memberNarwhal.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		result := memberNarwhal.Result.(model.MFAAdminSettings)
		output <- ResultUseCase{Result: result}
	})

	return output
}
