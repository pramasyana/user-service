package usecase

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/Bhinneka/user-service/src/client/v1/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/service"
	sharedDomain "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/pkg/errors"
)

// ClientUC implement usecase
type ClientUC struct {
	LoginSessionRepo authRepo.LoginSessionRepository
	RefreshTokenRepo authRepo.RefreshTokenRepository
	MemberQueryRead  memberQuery.MemberQuery

	//Messaging
	QPublisher service.QPublisher
}

// NewClientUsecase implement client usecase
func NewClientUsecase(loginSession authRepo.LoginSessionRepository, rt authRepo.RefreshTokenRepository, mq memberQuery.MemberQuery) *ClientUC {
	return &ClientUC{
		LoginSessionRepo: loginSession,
		RefreshTokenRepo: rt,
		MemberQueryRead:  mq,
	}
}

// Logout remove redis session
func (au *ClientUC) Logout(ctxReq context.Context, email string) <-chan sharedDomain.ResultUseCase {
	ctx := "ClientUseCase-Logout"
	output := make(chan sharedDomain.ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(context.Context, map[string]interface{}) {
		defer close(output)

		member, _, err := au.GetMemberByEmail(ctxReq, email)
		if err != nil {
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if member.SignUpFrom != model.BelaLKPP {
			err := errors.New("User is not signUpFrom Bela LKPP")
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}
		// check redis key
		redisKey := strings.Join([]string{"STG", member.ID, model.DeviceIDBela, authModel.DefaultDeviceLogin}, "-")
		rtKey := strings.Join([]string{"RT", member.ID, model.DeviceIDBela, authModel.DefaultDeviceLogin}, "-")
		// delete rediskey
		delRedisKey := <-au.LoginSessionRepo.Delete(ctxReq, redisKey)
		if delRedisKey.Error != nil {
			err := errors.Wrap(delRedisKey.Error, "invalid token format")
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		//delete refresh token key
		delRTKey := <-au.RefreshTokenRepo.Delete(ctxReq, rtKey)
		if delRTKey.Error != nil {
			err := errors.Wrap(delRedisKey.Error, "error delete refresh token")
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- sharedDomain.ResultUseCase{Error: nil, Result: true}
	})

	return output
}

// GetMemberByEmail reusable get member by email
func (au *ClientUC) GetMemberByEmail(ctxReq context.Context, email string) (member *memberModel.Member, statusCode int, err error) {
	ctx := "ClientUseCase-GetMemberByEmail"

	tc := tracer.StartTrace(ctxReq, ctx)
	defer tc.Finish(map[string]interface{}{"email": email})

	memberData := <-au.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberData.Error != nil {
		if memberData.Error == sql.ErrNoRows {
			memberData.Error = errors.New(authModel.ErrorUserLKPPBelaNotFoundBahasa)
		} else {
			tracer.SetError(ctxReq, memberData.Error)
		}
		return nil, http.StatusBadRequest, memberData.Error
	}

	memberModel, ok := memberData.Result.(memberModel.Member)
	if !ok {
		err := errors.New("user is not member")
		return nil, http.StatusBadRequest, err
	}
	return &memberModel, http.StatusOK, nil
}
