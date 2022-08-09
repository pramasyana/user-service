package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/golang-jwt/jwt"
)

// GetDataUser func
func (mu *MemberUseCaseImpl) GetDataUser(ctxReq context.Context, params *serviceModel.SendbirdRequest) *serviceModel.SendbirdRequest {
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-GetDataUser")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)

	tags[helper.TextArgs] = params
	var response serviceModel.SendbirdRequest
	response.UserID = params.UserID
	switch params.MemberType {
	case authModel.UserTypeCorporate:
		tags["userID"] = params.UserID
		emailCorporateResult := <-mu.CorporateContactQueryRead.FindByID(ctxReq, params.UserID)
		if emailCorporateResult.Result != nil {
			corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
			response.NickName = corporateContact.FirstName + " " + corporateContact.LastName
			response.ProfileURL = corporateContact.Avatar
		}
	case authModel.UserTypeMicrosite:
		tags["userID"] = params.UserID
		contactID, _ := strconv.Atoi(params.UserID)
		emailMicrositeResult := <-mu.CorporateAccContactQueryRead.FindAccountMicrositeByContactID(contactID)
		if emailMicrositeResult.Result != nil {
			micrositeContact := emailMicrositeResult.Result.(sharedModel.B2BContactData)
			response.NickName = micrositeContact.FirstName + " " + micrositeContact.LastName
			response.ProfileURL = micrositeContact.Avatar

		}
	case authModel.UserTypePersonal:
		tags["userID"] = params.UserID
		emailResult := <-mu.MemberQueryRead.FindByID(ctxReq, params.UserID)
		if emailResult.Result != nil {
			member := emailResult.Result.(model.Member)
			response.NickName = member.FirstName + " " + member.LastName
			response.ProfileURL = member.ProfilePicture
			merchantResult := mu.MerchantRepoRead.FindMerchantByUser(ctxReq, params.UserID)
			if merchantResult.Result != nil {
				merchantData := merchantResult.Result.(merchantModel.B2CMerchantDataV2)
				response.Metadata.Merchant.IsMerchant = "true"
				response.Metadata.Merchant.MerchantID = merchantData.ID
				response.Metadata.Merchant.MerchantName = merchantData.MerchantName
				response.Metadata.MerchantLogo = merchantData.MerchantLogo
			}
		}
	}
	tags[helper.TextResponse] = response
	return &response
}

// GetDataUser func
func (mu *MemberUseCaseImpl) GetDataUserV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) *serviceModel.SendbirdRequestV4 {
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-GetDataUserV4")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)

	tags[helper.TextArgs] = params
	var responseV4 serviceModel.SendbirdRequestV4
	responseV4.UserID = params.UserID
	switch params.MemberType {
	case authModel.UserTypeCorporate:
		tags["userID"] = params.UserID
		emailCorporateResult := <-mu.CorporateContactQueryRead.FindByID(ctxReq, params.UserID)
		if emailCorporateResult.Result != nil {
			corporateContact := emailCorporateResult.Result.(sharedModel.B2BContactData)
			responseV4.NickName = corporateContact.FirstName + " " + corporateContact.LastName
			responseV4.ProfileURL = corporateContact.Avatar
		}
	case authModel.UserTypeMicrosite:
		tags["userID"] = params.UserID
		contactID, _ := strconv.Atoi(params.UserID)
		emailMicrositeResult := <-mu.CorporateAccContactQueryRead.FindAccountMicrositeByContactID(contactID)
		if emailMicrositeResult.Result != nil {
			micrositeContact := emailMicrositeResult.Result.(sharedModel.B2BContactData)
			responseV4.NickName = micrositeContact.FirstName + " " + micrositeContact.LastName
			responseV4.ProfileURL = micrositeContact.Avatar
		}
	case authModel.UserTypePersonal:
		tags["userID"] = params.UserID
		emailResult := <-mu.MemberQueryRead.FindByID(ctxReq, params.UserID)
		if emailResult.Result != nil {
			member := emailResult.Result.(model.Member)
			responseV4.NickName = member.FirstName + " " + member.LastName
			responseV4.ProfileURL = member.ProfilePicture
			merchantResult := mu.MerchantRepoRead.FindMerchantByUser(ctxReq, params.UserID)
			if merchantResult.Result != nil {
				merchantData := merchantResult.Result.(merchantModel.B2CMerchantDataV2)
				if params.Client == "seller" {
					responseV4.UserID = merchantData.ID
					responseV4.NickName = merchantData.MerchantName
					responseV4.ProfileURL = merchantData.MerchantLogo.String
					responseV4.MetadataV4.Reference = merchantData.UserID
				} else {
					responseV4.MetadataV4.Reference = merchantData.ID
				}
			}
		}
	}
	tags[helper.TextResponse] = responseV4
	return &responseV4
}

// GetSendbirdToken function for get session token from sendbird
func (mu *MemberUseCaseImpl) GetSendbirdToken(ctxReq context.Context, params *serviceModel.SendbirdRequest) ResultUseCase {
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-GetSendbirdToken")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)
	tags[helper.TextArgs] = params
	var output ResultUseCase
	// extract jti from token
	_, dataToken, err := mu.AuthUseCase.GetJTIToken(ctxReq, params.Token, "")
	if err != nil {
		return ResultUseCase{
			Error:      err,
			HTTPStatus: http.StatusBadRequest}
	}
	claims := dataToken.(jwt.MapClaims)
	memberType := claims["memberType"].(string)
	params.MemberType = memberType
	// get data user
	bodyRequest := mu.GetDataUser(ctxReq, params)
	bodyRequest.ExpiresAt = params.ExpiresAt
	// check userID sendbird
	checkUserSendbird := mu.SendbirdService.CheckUserSenbird(ctxReq, params)
	if checkUserSendbird.Error != nil {
		checkStatusCode := checkUserSendbird.Result.(serviceModel.SendbirdStringResponse)
		// check if user sendbird is not found
		if checkStatusCode.Code == 400201 {
			// create user sendbird
			createUserSendbird := mu.SendbirdService.CreateUserSendbird(ctxReq, bodyRequest)
			if createUserSendbird.Error != nil {
				output = ResultUseCase{Error: createUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: createUserSendbird.Result}
				tags[helper.TextResponse] = output
				return output
			}

			getUserSendbird := mu.SendbirdService.GetUserSendbird(ctxReq, bodyRequest)
			if getUserSendbird.Error != nil {
				output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
				tags[helper.TextResponse] = output
				return output
			}
			dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponse)
			// get user sendbird
			output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
			tags[helper.TextResponse] = output
			return output
		}
		output = ResultUseCase{Error: checkUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: checkUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	// get and update token sendbird
	getTokenUserSendbird := mu.SendbirdService.CreateTokenUserSendbird(ctxReq, params)
	if getTokenUserSendbird.Error != nil {
		output = ResultUseCase{Error: getTokenUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getTokenUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	// bind token & expires_at to request body
	tokenData := getTokenUserSendbird.Result.(serviceModel.SessionTokenResponse)
	bodyRequest.Metadata.Token.ExpiresAt = tokenData.ExpiresAt
	// update data user and metadata (token & merchant)
	outputUpdate := mu.UpdateUserSendbird(ctxReq, bodyRequest)
	if outputUpdate.Error != nil {
		output = ResultUseCase{Error: outputUpdate.Error, HTTPStatus: http.StatusBadRequest, Result: outputUpdate.Result}
		tags[helper.TextResponse] = output
		return output
	}
	tags[helper.TextResponse] = outputUpdate
	getUserSendbird := mu.SendbirdService.GetUserSendbird(ctxReq, bodyRequest)
	if getUserSendbird.Error != nil {
		output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponse)
	output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
	return output
}

// GetSendbirdTokenV4 function for get session token from sendbird
func (mu *MemberUseCaseImpl) GetSendbirdTokenV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) ResultUseCase {
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-GetSendbirdTokenV4")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)
	tags[helper.TextArgs] = params
	var output ResultUseCase
	// extract jti from token
	_, dataTokenV4, err := mu.AuthUseCase.GetJTIToken(ctxReq, params.Token, "")
	if err != nil {
		return ResultUseCase{
			Error:      err,
			HTTPStatus: http.StatusBadRequest}
	}
	claims := dataTokenV4.(jwt.MapClaims)
	memberType := claims["memberType"].(string)
	params.MemberType = memberType
	// get data user
	bodyRequestV4 := mu.GetDataUserV4(ctxReq, params)
	bodyRequestV4.ExpiresAt = params.ExpiresAt
	// check userID sendbird
	checkUserSendbirdV4 := mu.SendbirdService.CheckUserSenbirdV4(ctxReq, bodyRequestV4)
	if checkUserSendbirdV4.Error != nil {
		checkStatusCode := checkUserSendbirdV4.Result.(serviceModel.SendbirdStringResponseV4)
		// check if user sendbird is not found
		if checkStatusCode.Code == 400201 {
			// create user sendbird
			createUserSendbird := mu.SendbirdService.CreateUserSendbirdV4(ctxReq, bodyRequestV4)
			if createUserSendbird.Error != nil {
				output = ResultUseCase{Error: createUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: createUserSendbird.Result}
				tags[helper.TextResponse] = output
				return output
			}
			getUserSendbird := mu.SendbirdService.GetUserSendbirdV4(ctxReq, bodyRequestV4)
			if getUserSendbird.Error != nil {
				output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
				tags[helper.TextResponse] = output
				return output
			}
			dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponseV4)
			// get user sendbird
			output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
			tags[helper.TextResponse] = output
			return output
		}
		output = ResultUseCase{Error: checkUserSendbirdV4.Error, HTTPStatus: http.StatusBadRequest, Result: checkUserSendbirdV4.Result}
		tags[helper.TextResponse] = output
		return output
	}

	// get and update token sendbird
	getTokenUserSendbird := mu.SendbirdService.CreateTokenUserSendbirdV4(ctxReq, bodyRequestV4)
	if getTokenUserSendbird.Error != nil {
		output = ResultUseCase{Error: getTokenUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getTokenUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	// bind token & expires_at to request body
	tokenData := getTokenUserSendbird.Result.(serviceModel.SessionTokenResponse)
	bodyRequestV4.MetadataV4.Token.ExpiresAt = tokenData.ExpiresAt
	// update data user and metadata (token & merchant)
	outputUpdate := mu.UpdateUserSendbirdV4(ctxReq, bodyRequestV4)
	if outputUpdate.Error != nil {
		output = ResultUseCase{Error: outputUpdate.Error, HTTPStatus: http.StatusBadRequest, Result: outputUpdate.Result}
		tags[helper.TextResponse] = output
		return output
	}
	tags[helper.TextResponse] = outputUpdate
	getUserSendbird := mu.SendbirdService.GetUserSendbirdV4(ctxReq, bodyRequestV4)
	if getUserSendbird.Error != nil {
		output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponseV4)
	output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
	return output
}

// CheckSendbirdToken function for get session token from sendbird
func (mu *MemberUseCaseImpl) CheckSendbirdToken(ctxReq context.Context, params *serviceModel.SendbirdRequest) ResultUseCase {
	var output ResultUseCase
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-CheckSendbirdToken")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)

	tags[helper.TextArgs] = params
	// extract jti from token
	_, dataToken, err := mu.AuthUseCase.GetJTIToken(ctxReq, params.Token, "")
	if err != nil {
		return ResultUseCase{
			Error:      err,
			HTTPStatus: http.StatusBadRequest}
	}
	claims := dataToken.(jwt.MapClaims)
	memberType := claims["memberType"].(string)
	params.MemberType = memberType
	// get data user
	bodyRequest := mu.GetDataUser(ctxReq, params)
	bodyRequest.ExpiresAt = params.ExpiresAt
	// get and update token sendbird
	checkTokenUserSendbird := mu.SendbirdService.CreateTokenUserSendbird(ctxReq, params)
	if checkTokenUserSendbird.Error != nil {
		output = ResultUseCase{Error: checkTokenUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: checkTokenUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	// bind token & expires_at to request body
	tokenData := checkTokenUserSendbird.Result.(serviceModel.SessionTokenResponse)
	bodyRequest.Metadata.Token.ExpiresAt = tokenData.ExpiresAt
	// update data user and metadata (token & merchant)
	outputUpdate := mu.UpdateUserSendbird(ctxReq, bodyRequest)
	tags[helper.TextResponse] = outputUpdate
	getUserSendbird := mu.SendbirdService.GetUserSendbird(ctxReq, bodyRequest)
	if getUserSendbird.Error != nil {
		output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponse)
	output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
	return output
}

// CheckSendbirdTokenV4 function for get session token from sendbird
func (mu *MemberUseCaseImpl) CheckSendbirdTokenV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) ResultUseCase {
	var output ResultUseCase
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-CheckSendbirdTokenV4")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)
	tags[helper.TextArgs] = params
	// extract jti from token
	_, dataTokenV4, err := mu.AuthUseCase.GetJTIToken(ctxReq, params.Token, "")
	if err != nil {
		return ResultUseCase{
			Error:      err,
			HTTPStatus: http.StatusBadRequest}
	}
	claims := dataTokenV4.(jwt.MapClaims)
	memberType := claims["memberType"].(string)
	params.MemberType = memberType
	// get data user
	bodyRequestV4 := mu.GetDataUserV4(ctxReq, params)
	bodyRequestV4.ExpiresAt = params.ExpiresAt
	// get and update token sendbird
	checkTokenUserSendbirdV4 := mu.SendbirdService.CreateTokenUserSendbirdV4(ctxReq, bodyRequestV4)
	if checkTokenUserSendbirdV4.Error != nil {
		output = ResultUseCase{Error: checkTokenUserSendbirdV4.Error, HTTPStatus: http.StatusBadRequest, Result: checkTokenUserSendbirdV4.Result}
		tags[helper.TextResponse] = output
		return output
	}
	// bind token & expires_at to request body
	tokenData := checkTokenUserSendbirdV4.Result.(serviceModel.SessionTokenResponse)
	bodyRequestV4.MetadataV4.Token.ExpiresAt = tokenData.ExpiresAt
	// update data user and metadata (token & merchant)
	outputUpdate := mu.UpdateUserSendbirdV4(ctxReq, bodyRequestV4)
	tags[helper.TextResponse] = outputUpdate
	getUserSendbird := mu.SendbirdService.GetUserSendbirdV4(ctxReq, bodyRequestV4)
	if getUserSendbird.Error != nil {
		output = ResultUseCase{Error: getUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: getUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}
	dataResponseSenbird := getUserSendbird.Result.(serviceModel.SendbirdResponseV4)
	output = ResultUseCase{Result: dataResponseSenbird, HTTPStatus: http.StatusOK}
	return output

}

// UpdateUserSendbird func
func (mu *MemberUseCaseImpl) UpdateUserSendbird(ctxReq context.Context, bodyRequest *serviceModel.SendbirdRequest) ResultUseCase {

	var output ResultUseCase
	var err error

	tc := tracer.StartTrace(ctxReq, "MemberUseCase-UpdateUserSendbird")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)

	tags[helper.TextArgs] = bodyRequest

	updateUserSendbird := mu.SendbirdService.UpdateUserSendbird(ctxReq, bodyRequest)
	if updateUserSendbird.Error != nil {
		output = ResultUseCase{Error: updateUserSendbird.Error, HTTPStatus: http.StatusBadRequest, Result: updateUserSendbird.Result}
		tags[helper.TextResponse] = output
		return output
	}

	dataResponse := updateUserSendbird.Result.(serviceModel.SendbirdStringResponse)

	var response serviceModel.SendbirdResponse
	var metadataMerchant serviceModel.Merchant
	var metadataToken serviceModel.SessionTokenResponse

	// casting from token-string to token-interface
	err = json.Unmarshal([]byte(dataResponse.Metadata.Token), &metadataToken)
	if err != nil {
		output = ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest, Result: dataResponse}
		tags[helper.TextResponse] = output
		return output
	}
	response.UserID = dataResponse.UserID
	response.NickName = dataResponse.NickName
	response.ProfileURL = dataResponse.ProfileURL
	response.Metadata.Merchant = metadataMerchant
	response.Metadata.Token = metadataToken
	response.Metadata.MerchantLogo = dataResponse.Metadata.MerchantLogo

	output = ResultUseCase{Result: response, HTTPStatus: http.StatusOK}
	tags[helper.TextResponse] = output

	return output
}

// UpdateUserSendbird func
func (mu *MemberUseCaseImpl) UpdateUserSendbirdV4(ctxReq context.Context, bodyRequest *serviceModel.SendbirdRequestV4) ResultUseCase {

	var output ResultUseCase
	var err error

	tc := tracer.StartTrace(ctxReq, "MemberUseCase-UpdateUserSendbirdV4")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)

	tags[helper.TextArgs] = bodyRequest

	updateUserSendbirdV4 := mu.SendbirdService.UpdateUserSendbirdV4(ctxReq, bodyRequest)
	if updateUserSendbirdV4.Error != nil {
		output = ResultUseCase{Error: updateUserSendbirdV4.Error, HTTPStatus: http.StatusBadRequest, Result: updateUserSendbirdV4.Result}
		tags[helper.TextResponse] = output
		return output
	}

	dataResponseV4 := updateUserSendbirdV4.Result.(serviceModel.SendbirdStringResponseV4)

	var responseV4 serviceModel.SendbirdResponseV4
	var metadataTokenV4 serviceModel.SessionTokenResponse

	// casting from token-string to token-interface
	err = json.Unmarshal([]byte(dataResponseV4.MetadataV4.Token), &metadataTokenV4)
	if err != nil {
		output = ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest, Result: dataResponseV4}
		tags[helper.TextResponse] = output
		return output
	}

	responseV4.UserID = dataResponseV4.UserID
	responseV4.NickName = dataResponseV4.NickName
	responseV4.ProfileURL = dataResponseV4.ProfileURL
	responseV4.MetadataV4.Token = metadataTokenV4

	output = ResultUseCase{Result: responseV4, HTTPStatus: http.StatusOK}
	tags[helper.TextResponse] = output

	return output
}
