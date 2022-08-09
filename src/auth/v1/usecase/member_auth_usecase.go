package usecase

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// GenerateTokenFromUserID function for generating token based on userID
// this feature used by DOLPHIN in order to implement Manual Order
func (au *AuthUseCaseImpl) GenerateTokenFromUserID(ctxReq context.Context, data model.RequestToken) <-chan ResultUseCase {
	output := make(chan ResultUseCase)
	go func() {
		defer close(output)

		if len(data.UserID) <= 0 {
			output <- ResultUseCase{Error: errors.New("user id cannot be empty"), HTTPStatus: http.StatusBadRequest}
			return
		}

		// get detail member to get email for being set on response
		memberResult := <-au.MemberRepoRead.Load(ctxReq, data.UserID)
		if memberResult.Error != nil {
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusUnauthorized}
			return
		}

		member, ok := memberResult.Result.(memberModel.Member)
		if !ok {
			output <- ResultUseCase{Error: errors.New(msgResultNotMember), HTTPStatus: http.StatusUnauthorized}
			return
		}

		// check member status
		// only active member can login

		if member.Status.String() == memberModel.InactiveString || member.Status.String() == memberModel.NewString {
			output <- ResultUseCase{Error: errors.New(model.ErrorAccountInActiveBahasa), HTTPStatus: http.StatusUnauthorized}
			return
		}

		if member.Status.String() == memberModel.BlockedString {
			output <- ResultUseCase{Error: errors.New(model.ErrorAccountBlockedBahasa), HTTPStatus: http.StatusUnauthorized}
			return
		}

		data.UserID = member.ID
		data.Email = member.Email
		data.FirstName = member.FirstName
		data.LastName = member.LastName
		data.NewMember = false

		// set claims
		claims := token.Claim{}
		claims.Issuer = model.Bhinneka
		claims.DeviceID = data.DeviceID
		claims.DeviceLogin = data.DeviceLogin
		claims.Audience = data.Audience
		claims.Subject = member.ID
		claims.Authorised = true
		claims.IsAdmin = member.IsAdmin
		claims.IsStaff = member.IsStaff
		claims.Email = member.Email
		claims.SignUpFrom = member.SignUpFrom
		claims.CustomToken = data.TokenBela

		// generate token based on claims
		tokenResult := <-au.AccessTokenGenerator.GenerateAccessToken(claims)

		data.Token = tokenResult.AccessToken.AccessToken
		data.ExpiredAt = tokenResult.AccessToken.ExpiredAt
		anonRequest := model.RequestToken{
			GrantType:   model.AuthTypeAnonymous,
			DeviceID:    data.DeviceID,
			DeviceLogin: data.DeviceLogin,
			Audience:    member.ID,
		}
		au.GenerateToken(ctxReq, "", anonRequest)

		output <- ResultUseCase{Result: data}

	}()

	return output
}

// checkMemberStatus function for check status of member
func (au *AuthUseCaseImpl) checkMemberStatus(ctxReq context.Context, member memberModel.Member, Socmed interface{}, grantType string) <-chan ResultUseCase {
	ctx := "AuthUseCase-checkMemberStatus"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		// generate data member from media social data
		mediaData := <-au.generateSocmedData(ctxReq, member, Socmed, grantType, false)

		if mediaData.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeParseLdap, mediaData.Error, member)
			output <- ResultUseCase{HTTPStatus: mediaData.HTTPStatus, Error: mediaData.Error}
			return
		}

		member, ok := mediaData.Result.(memberModel.Member)
		if !ok {
			err := errors.New("result is not member data")
			helper.SendErrorLog(ctxReq, ctx, "parse_member", err, member)
			output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			return
		}

		// check member status
		// only active member can login
		if member.StatusString == memberModel.InactiveString || member.StatusString == memberModel.NewString {
			//remove password and salt
			member.Password = ""
			member.Salt = ""
			member.HasPassword = false

			//activate member and remove token activation
			member.Status = memberModel.StringToStatus(memberModel.ActiveString)
			member.Token = ""
		}

		updateResult := <-au.MemberRepoWrite.Save(ctxReq, member)
		if updateResult.Error != nil {
			err := errors.New("failed to update member data")
			helper.SendErrorLog(ctxReq, ctx, "update_result", err, member)
			output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			return
		}

		if member.StatusString == memberModel.BlockedString {
			err := errors.New(model.ErrorAccountBlockedBahasa)
			output <- ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: err}
			return
		}

		tags["args"] = member

		output <- ResultUseCase{Result: member}
	})

	return output
}

// AdjustMemberData function for adjust member information
func (au *AuthUseCaseImpl) AdjustMemberData(ctxReq context.Context, responseVerify *model.VerifyResponse) *model.VerifyResponse {

	switch responseVerify.MemberType {
	case model.UserTypeCorporate:
		emailCorporateResult := <-au.CorporateContactQueryRead.FindByID(ctxReq, responseVerify.UserID)
		if emailCorporateResult.Result != nil {
			corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
			responseVerify.FirstName = corporateContact.FirstName
			responseVerify.LastName = corporateContact.LastName
			responseVerify.Mobile = corporateContact.PhoneNumber
			responseVerify.AccountID = corporateContact.AccountID
		}
	case model.UserTypeMicrosite:
		contactID, _ := strconv.Atoi(responseVerify.UserID)
		emailMicrositeResult := <-au.CorporateAccContactQueryRead.FindAccountMicrositeByContactID(contactID)
		if emailMicrositeResult.Result != nil {
			micrositeContact := emailMicrositeResult.Result.(sharedModel.B2BContactData)
			responseVerify.FirstName = micrositeContact.FirstName
			responseVerify.LastName = micrositeContact.LastName
			responseVerify.Mobile = micrositeContact.PhoneNumber
			responseVerify.AccountID = micrositeContact.AccountID
		}
	case model.UserTypePersonal:
		responseVerify.HasPassword = false
		emailResult := <-au.MemberQueryRead.FindByID(ctxReq, responseVerify.UserID)
		if emailResult.Result != nil {
			member := emailResult.Result.(memberModel.Member)
			responseVerify.FirstName = member.FirstName
			responseVerify.LastName = member.LastName
			responseVerify.Mobile = member.Mobile
			if member.Password != "" {
				responseVerify.HasPassword = true
			}
			merchantResult := au.MerchantRepoRead.FindMerchantByUser(ctxReq, responseVerify.UserID)
			if merchantResult.Result != nil {
				merchantData := merchantResult.Result.(merchantModel.B2CMerchantDataV2)
				responseVerify.IsMerchant = true
				responseVerify.MerchantID = merchantData.ID
			}
		}
	}

	responseVerify.MerchantIDs = au.merchantEmployee(ctxReq, responseVerify.Email)

	return responseVerify
}

//AdjustDataPersonal
func (au *AuthUseCaseImpl) merchantEmployee(ctxReq context.Context, email string) []string {
	var merchantIds []string

	emailResult := <-au.MemberQueryRead.FindByEmail(ctxReq, email)
	if emailResult.Result != nil {
		member := emailResult.Result.(memberModel.Member)

		// merchant employee
		filterMerchantEmployee := &merchantModel.QueryMerchantEmployeeParameters{
			MemberID: member.ID,
		}
		merchantEmployeeResult := <-au.MerchantEmployeeRepoRead.GetAllMerchantEmployees(ctxReq, filterMerchantEmployee)
		if merchantEmployeeResult.Result != nil {
			merchantEmployee := merchantEmployeeResult.Result.([]merchantModel.B2CMerchantEmployeeData)
			for _, val := range merchantEmployee {
				merchantIds = append(merchantIds, val.MerchantID)
			}
		}
	}

	return merchantIds
}
