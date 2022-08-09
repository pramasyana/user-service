package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// ValidateToken function for validating token of forgot password
func (mu *MemberUseCaseImpl) ValidateToken(ctxReq context.Context, token string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ValidateToken"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		email, newToken, err := mu.parseToken(token)
		tags[helper.TextEmail] = email
		if err != nil {
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, email)
		if memberResult.Error != nil {
			// when data not found
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, msgErrorMember+email)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tokenResult := <-mu.MemberRepoRedis.Load(member.ID)
		if tokenResult.Error != nil {
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		tokenRedis, ok := tokenResult.Result.(model.MemberRedis)
		if !ok {
			err := errors.New("result is not member from redis")
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		// validate the token from redis and form
		if tokenRedis.Token != newToken {
			err := fmt.Errorf(helper.ErrorParameterInvalid, "token")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusUnauthorized}
			return
		}

		response := model.SuccessResponse{
			ID:          golib.RandomString(8),
			Message:     "token is valid",
			HasPassword: true,
		}

		output <- ResultUseCase{Result: response}

	})

	return output
}

// RegenerateToken function for regenerating token for member activation
func (mu *MemberUseCaseImpl) RegenerateToken(ctxReq context.Context, data model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-RegenerateToken"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = data
		if data.Status.String() == model.ActiveString {
			err := errors.New("member is already active")
			output <- ResultUseCase{Error: err}
			return
		} else if data.Status.String() == model.BlockedString {
			err := errors.New("member is blocked")
			output <- ResultUseCase{Error: err}
			return
		}

		// generate random string
		mix := data.Email + "-" + data.FirstName
		data.Token = helper.GenerateTokenByString(mix)

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, data)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		var hasPassword bool
		if len(data.Password) > 0 {
			hasPassword = true
		}

		response := model.SuccessResponse{
			ID:          data.ID,
			Message:     helper.SuccessMessage,
			Token:       data.Token,
			HasPassword: hasPassword,
			FirstName:   data.FirstName,
			LastName:    data.LastName,
			Email:       strings.ToLower(data.Email),
		}

		// return the token to be sent to email notification service
		output <- ResultUseCase{Result: response}
	})

	return output
}
