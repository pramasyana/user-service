package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/query"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// RequestTokenClientMicrosite from LKPP, bela
func (au *AuthUseCaseImpl) RequestTokenClientMicrosite(ctxReq context.Context, data *model.RequestToken) (*memberModel.Member, int, error) {
	ctx := "AuthUseCase-RequestTokenClientMicrosite"

	tc := tracer.StartTrace(ctxReq, ctx)
	defer tc.Finish(map[string]interface{}{"param": data})

	// create member instead of returning error
	socData := model.MicrositeClient{
		Firstname: data.FirstName,
		Email:     data.Email,
	}

	dataMember := au.CheckMemberSocmedType(ctxReq, data, socData, data.Email)
	if dataMember.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "check_member_sosmed", dataMember.Error, data)
		return nil, http.StatusBadRequest, dataMember.Error
	}
	member := dataMember.Data

	if member.SignUpFrom != signUpFromLKPP || member.Status.String() != memberModel.ActiveString {
		return nil, http.StatusUnauthorized, errors.New("unauthorized")
	}
	if dataMember.NewMember {
		paramEmail := model.AccessTokenResponse{
			Email:     data.Email,
			FirstName: data.FirstName,
		}
		au.SendEmailWelcomeMember(ctxReq, paramEmail)
	}
	return member, http.StatusOK, nil
}

// CheckingContactTypeMicrositeBela function for checking contact if doest exist send to kafka and checking again
func (au *AuthUseCaseImpl) CheckingContactTypeMicrositeBela(ctxReq context.Context, data *model.RequestToken) query.ResultQuery {
	result := query.ResultQuery{}
	b2bAccountResult := <-au.CorporateContactQueryRead.FindAccountByMemberType(ctxReq, data.MemberType)
	if b2bAccountResult.Error != nil {
		result.Error = b2bAccountResult.Error
		return result
	}

	b2bAcount, ok := b2bAccountResult.Result.(sharedModel.B2BAccountCDC)
	if !ok {
		result.Error = errors.New(msgResultNotAccount)
		return result
	}

	payload := corporateModel.ContactPayload{
		Email:           data.Email,
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		AccountID:       b2bAcount.ID,
		TransactionType: data.TransactionType,
		CreatedAt:       time.Now(),
		LpseID:          data.LpseID,
	}
	go au.PublishToKafkaContact(ctxReq, payload, "import")
	time.Sleep(3 * time.Second)

	memberResult := <-au.CorporateContactQueryRead.FindContactMicrositeByEmail(ctxReq, data.Email, data.TransactionType, data.MemberType)
	result = query.ResultQuery(memberResult)
	return result
}

// RequestTokenPasswordTypeMicrosite function for validate request token with grantType is password for microsite
func (au *AuthUseCaseImpl) RequestTokenPasswordTypeMicrositeBela(ctxReq context.Context, data *model.RequestToken) (sharedModel.B2BContactData, int, error) {
	ctx := "AuthUseCase-RequestTokenPasswordTypeMicrosite"

	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	member := sharedModel.B2BContactData{}

	// check contact for get name to send kafka
	checkContactResult := <-au.CorporateContactQueryRead.FindContactByEmail(ctxReq, data.Email)
	if checkContactResult.Error == nil {
		member, ok := checkContactResult.Result.(sharedModel.B2BContactData)
		if !ok {
			err := errors.New(msgResultNotMember)
			return member, http.StatusUnauthorized, err
		}

		if data.FirstName == "" {
			data.FirstName = member.FirstName
			data.LastName = member.LastName
		}
	}

	// send to kafka contact
	memberResult := au.CheckingContactTypeMicrositeBela(ctxReq, data)
	if memberResult.Error != nil {
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(model.ErrorIncorrectMemberTypeMicrosite, data.MemberType)
		}
		return member, http.StatusUnauthorized, memberResult.Error
	}

	member, ok := memberResult.Result.(sharedModel.B2BContactData)
	if !ok {
		err := errors.New(msgResultNotMember)
		return member, http.StatusUnauthorized, err
	}

	key := fmt.Sprintf(keyAttempt, member.Email)

	<-au.LoginAttemptRepo.Delete(ctxReq, key)
	return member, 200, nil
}
