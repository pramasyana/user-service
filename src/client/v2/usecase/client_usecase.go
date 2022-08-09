package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/Bhinneka/user-service/src/client/v1/model"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/service"
	sharedDomain "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/pkg/errors"
)

const (
	msgResultNotMember = "result is not member"
)

// ClientUC implement usecase
type ClientUC struct {
	LoginSessionRepo          authRepo.LoginSessionRepository
	RefreshTokenRepo          authRepo.RefreshTokenRepository
	MemberQueryRead           memberQuery.MemberQuery
	CorporateContactQueryRead corporateQuery.ContactQuery

	//Messaging
	QPublisher service.QPublisher
}

// NewClientUsecase implement client usecase
func NewClientUsecase(loginSession authRepo.LoginSessionRepository, rt authRepo.RefreshTokenRepository, mq memberQuery.MemberQuery, cc corporateQuery.ContactQuery) *ClientUC {
	return &ClientUC{
		LoginSessionRepo:          loginSession,
		RefreshTokenRepo:          rt,
		MemberQueryRead:           mq,
		CorporateContactQueryRead: cc,
	}
}

// Logout remove redis session
func (au *ClientUC) Logout(ctxReq context.Context, email string) <-chan sharedDomain.ResultUseCase {
	ctx := "ClientUseCase-Logout"
	output := make(chan sharedDomain.ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(context.Context, map[string]interface{}) {
		defer close(output)

		member := sharedDomain.B2BContactData{}
		// use came method as
		memberResult := <-au.CorporateContactQueryRead.FindContactMicrositeByEmail(ctxReq, email, authModel.LoginTypeShopcart, authModel.UserTypeMicrositeBela)
		if memberResult.Error != nil {

			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(authModel.ErrorIncorrectMemberTypeMicrosite, authModel.LoginTypeShopcart)
			}
			output <- sharedDomain.ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(sharedDomain.B2BContactData)
		if !ok {
			err := errors.New(msgResultNotMember)
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check redis key
		redisKey := strings.Join([]string{"STG", strconv.Itoa(member.ID), model.DeviceIDBela, authModel.DefaultDeviceLogin}, "-")
		rtKey := strings.Join([]string{"RT", strconv.Itoa(member.ID), model.DeviceIDBela, authModel.DefaultDeviceLogin}, "-")
		// delete rediskey
		delRedisKeys := <-au.LoginSessionRepo.Delete(ctxReq, redisKey)
		if delRedisKeys.Error != nil {
			err := errors.Wrap(delRedisKeys.Error, "invalid token format")
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		//delete refresh token key
		delRTKeys := <-au.RefreshTokenRepo.Delete(ctxReq, rtKey)
		if delRTKeys.Error != nil {
			err := errors.Wrap(delRedisKeys.Error, "error delete refresh token")
			output <- sharedDomain.ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- sharedDomain.ResultUseCase{Error: nil, Result: true}
	})

	return output
}
