package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	stringLib "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/query"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

const (
	PersonalAccount    = "Akun Personal"
	CorporateAccount   = "Akun Korporasi"
	Corporate          = "corporate"
	Personal           = "personal"
	ErrorEmailNotFound = "Email belum terdaftar"
)

// CheckEmail usecase function for verify email member
func (au *AuthUseCaseImpl) CheckEmail(ctxReq context.Context, email string) <-chan ResultUseCase {
	ctx := "Auth-CheckEmail"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		if email == "" {
			err := fmt.Errorf("email required")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if err := stringLib.ValidateEmail(email); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		response := model.CheckEmailResponse{}
		response.Email = email

		users := make([]model.Users, 0)

		emailResult := <-au.MemberQueryRead.FindByEmail(ctxReq, email)
		if emailResult.Result != nil {
			member := emailResult.Result.(memberModel.Member)
			userType := model.Users{}
			userType.UserType = model.UserTypePersonal
			userType.FirstName = member.FirstName
			userType.LastName = member.LastName
			hasPassword := true
			if member.Password == "" {
				hasPassword = false
			}
			userType.HasPassword = hasPassword
			userType.AccountType = PersonalAccount
			users = append(users, userType)
		}

		emailCorporateResult := <-au.CorporateContactQueryRead.FindByEmail(ctxReq, email)
		if emailCorporateResult.Result != nil {
			corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
			userType := model.Users{}
			userType.UserType = model.UserTypeCorporate
			userType.FirstName = corporateContact.FirstName
			userType.LastName = corporateContact.LastName
			corporateStatus := <-au.CorporateContactQueryRead.GetTransactionType(ctxReq, email)
			transactionType := checkResultGetTransactionType(query.ResultQuery(corporateStatus))
			userType.AccountType = transactionType
			userType.HasPassword = true
			users = append(users, userType)
		}

		if emailResult.Result == nil && emailCorporateResult.Result == nil {
			err := fmt.Errorf("Email belum terdaftar")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		tags[helper.TextResponse] = users

		response.Users = users

		output <- ResultUseCase{Result: response}

	})

	return output
}

// SendEmailWelcomeMember usecase function for send email success registration
func (au *AuthUseCaseImpl) SendEmailWelcomeMember(ctxReq context.Context, data model.AccessTokenResponse) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailWelcomeMember"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, _ map[string]interface{}) {
		defer close(output)

		response := memberModel.SuccessResponse{
			Email:     strings.ToLower(data.Email),
			FirstName: data.FirstName,
			LastName:  data.LastName,
		}

		plEmail := memberModel.MemberEmailQueue{
			Data: response,
		}

		au.QPublisher.QueueJob(ctxReq, plEmail, data.ID, "SendEmailWelcomeMember")

		output <- ResultUseCase{Result: data}
	})

	return output
}

func checkResultGetTransactionType(corporateStatus query.ResultQuery) (AccountType string) {
	if corporateStatus.Result != nil {
		if corporateStatus.Result == "eproc" {
			corporateStatus.Result = "e-Procurement"
		}
		AccountType = fmt.Sprintf("%s - %s", CorporateAccount, corporateStatus.Result)
	}
	return AccountType
}

func (au *AuthUseCaseImpl) CheckEmailV3(ctxReq context.Context, data model.CheckEmailPayload) <-chan ResultUseCase {
	ctx := "Auth-CheckEmailV3"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		response := model.CheckEmailResponse{}
		response.Email = data.Email
		if data.Email == "" {
			err := fmt.Errorf("email required")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if err := stringLib.ValidateEmail(data.Email); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if data.UserType == "" {
			users := make([]model.Users, 0)
			queryPersonal := query.ResultQuery{}
			queryCorporate := query.ResultQuery{}
			users, queryPersonal = au.CheckEmailPersonal(ctxReq, data.Email, users)
			users, queryCorporate = au.CheckEmailCorporate(ctxReq, data.Email, users)

			if queryPersonal.Result == nil && queryCorporate.Result == nil {
				err := fmt.Errorf(ErrorEmailNotFound)
				output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
				return
			}
			tags[helper.TextResponse] = users

			response.Users = users

		} else {
			user, err := au.CheckCorporatePersonal(ctxReq, data)
			if err != nil {
				output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
				return
			}
			tags[helper.TextResponse] = user
			response.Users = user
		}
		output <- ResultUseCase{Result: response}

	})

	return output
}

func (au *AuthUseCaseImpl) CheckCorporatePersonal(ctxReq context.Context, data model.CheckEmailPayload) ([]model.Users, error) {
	users := make([]model.Users, 0)
	queryPersonal := query.ResultQuery{}
	queryCorporate := query.ResultQuery{}
	if data.UserType == Personal {
		users, queryPersonal = au.CheckEmailPersonal(ctxReq, data.Email, users)
		if queryPersonal.Result == nil {
			err := fmt.Errorf(ErrorEmailNotFound)
			return nil, err
		}
	} else if data.UserType == Corporate {
		users, queryCorporate = au.CheckEmailCorporate(ctxReq, data.Email, users)
		if queryCorporate.Result == nil {
			err := errors.New(ErrorEmailNotFound)
			return nil, err
		}
	}
	return users, nil
}

func (au *AuthUseCaseImpl) CheckEmailPersonal(ctxReq context.Context, email string, users []model.Users) ([]model.Users, query.ResultQuery) {
	emailResult := <-au.MemberQueryRead.FindByEmail(ctxReq, email)
	if emailResult.Result != nil {
		member := emailResult.Result.(memberModel.Member)
		userType := model.Users{}
		userType.UserType = model.UserTypePersonal
		userType.FirstName = member.FirstName
		userType.LastName = member.LastName
		hasPassword := true
		if member.Password == "" {
			hasPassword = false
		}
		userType.HasPassword = hasPassword
		userType.AccountType = PersonalAccount
		users = append(users, userType)
	}
	return users, query.ResultQuery(emailResult)
}

func (au *AuthUseCaseImpl) CheckEmailCorporate(ctxReq context.Context, email string, users []model.Users) ([]model.Users, query.ResultQuery) {
	emailCorporateResult := <-au.CorporateContactQueryRead.FindByEmail(ctxReq, email)
	if emailCorporateResult.Result != nil {
		corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
		user := model.Users{}
		user.UserType = model.UserTypeCorporate
		user.FirstName = corporateContact.FirstName
		user.LastName = corporateContact.LastName
		corporateStatus := <-au.CorporateContactQueryRead.GetTransactionType(ctxReq, email)
		transactionType := checkResultGetTransactionType(query.ResultQuery(corporateStatus))
		user.AccountType = transactionType
		user.HasPassword = true
		users = append(users, user)
	}
	return users, query.ResultQuery(emailCorporateResult)
}
