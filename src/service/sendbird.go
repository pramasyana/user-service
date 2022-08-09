package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

const (
	urlString      = "%s%s%s%s"
	urlGroup       = "/v3/users/"
	apiTokenHeader = "Api-Token"
	urlStringThree = "%s%s%s"
	metaData       = "/metadata"
	tokenString    = "/token"
)

// SendbirdService data structure
type SendbirdService struct {
	BaseURL  *url.URL
	APIToken string
}

// NewSendbirdService function for initializing sendbird service
func NewSendbirdService() (*SendbirdService, error) {
	var (
		sendbird SendbirdService
		err      error
		ok       bool
	)

	sendbird.APIToken, ok = os.LookupEnv("SENDBIRD_API_TOKEN")
	if !ok {
		return &sendbird, errors.New("you need to specify SENDBIRD_API_TOKEN in the environment variable")
	}

	baseURL, ok := os.LookupEnv("SENDBIRD_SERVICE_URL")
	if !ok {
		return &sendbird, errors.New("you need to specify SENDBIRD_SERVICE_URL in the environment variable")
	}

	sendbird.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		return &sendbird, errors.New("error parsing sendbird services url")
	}

	return &sendbird, nil
}

// CheckUserSenbird function for getting sendbird user
func (p *SendbirdService) CheckUserSenbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-FindSendbirdUser"

	var responseV4 serviceModel.SendbirdStringResponse
	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}
	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)

	// request
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodGet, uri, nil, &responseV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseV4}
	}

	if responseV4.Error {
		e := errors.New(responseV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "get_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseV4}

	}

	return serviceModel.ServiceResult{Result: responseV4}
}

// CheckUserSenbirdV4 function for getting sendbird user
func (p *SendbirdService) CheckUserSenbirdV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult {
	ctx := "Service-FindSendbirdUserV4"

	var response serviceModel.SendbirdStringResponseV4
	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}
	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)

	// request
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodGet, uri, nil, &response, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: response}
	}

	if response.Error {
		e := errors.New(response.Message)
		helper.SendErrorLog(ctxReq, ctx, "get_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: response}

	}

	return serviceModel.ServiceResult{Result: response}
}

// CreateUserSendbird function for create sendbird user
func (p *SendbirdService) CreateUserSendbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-CreateSendbirdUser"

	var response serviceModel.SendbirdResponse
	var responseUser serviceModel.SendbirdStringResponse
	var userBody serviceModel.User

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	uri := fmt.Sprintf("%s%s", p.BaseURL, urlGroup)

	// bind data to Interface user
	userBody.UserID = params.UserID
	userBody.NickName = params.NickName
	userBody.ProfileURL = params.ProfileURL
	bodyMershal, _ := json.Marshal(userBody)
	payload := strings.NewReader(string(bodyMershal))

	// request for create user
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uri, payload, &responseUser, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUser}
	}

	if responseUser.Error {
		e := errors.New(responseUser.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUser}
	}

	// request for get token from sendbird
	getToken := p.CreateTokenUserSendbird(ctxReq, params)
	token := getToken.Result.(serviceModel.SessionTokenResponse)

	var tokenMetadata serviceModel.SessionTokenRequest
	var metadataBody serviceModel.MetadataRequest
	var metadataResponse serviceModel.MetaDataResponse

	tokenMetadata.ExpiresAt = token.ExpiresAt

	uriMetadata := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, metaData)
	tokenData, _ := json.Marshal(tokenMetadata)
	merchantData, _ := json.Marshal(params.Metadata.Merchant)

	// bind token and merchant to metadata interface
	metadataBody.Metadata.Token = string(tokenData)
	metadataBody.Metadata.Merchant = string(merchantData)
	metadataBody.Metadata.MerchantLogo = params.Metadata.MerchantLogo
	bodyMetadataMershal, _ := json.Marshal(metadataBody)
	payloadMetadata := strings.NewReader(string(bodyMetadataMershal))

	// request for create metadata
	err = helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uriMetadata, payloadMetadata, &metadataResponse, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: err, Result: metadataBody}
	}

	if metadataResponse.Error {
		e := errors.New(metadataResponse.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_metadata_user_sendbird", e, uriMetadata)
		return serviceModel.ServiceResult{Error: e, Result: metadataResponse}
	}

	response.UserID = responseUser.UserID
	response.NickName = responseUser.NickName
	response.ProfileURL = responseUser.ProfileURL
	response.Metadata.Merchant = params.Metadata.Merchant
	response.Metadata.Token = token

	return serviceModel.ServiceResult{Result: response}
}

// CreateUserSendbirdV4 function for create sendbird user
func (p *SendbirdService) CreateUserSendbirdV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult {
	ctx := "Service-CreateSendbirdUser"

	var responseV4 serviceModel.SendbirdStringResponseV4
	var responseUserV4 serviceModel.SendbirdStringResponseV4
	var userBody serviceModel.User

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	uri := fmt.Sprintf("%s%s", p.BaseURL, urlGroup)

	// bind data to Interface user
	userBody.UserID = params.UserID
	userBody.NickName = params.NickName
	userBody.ProfileURL = params.ProfileURL
	bodyMershal, _ := json.Marshal(userBody)
	payload := strings.NewReader(string(bodyMershal))

	// request for create user
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uri, payload, &responseUserV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUserV4}
	}

	if responseUserV4.Error {
		e := errors.New(responseUserV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUserV4}
	}

	// request for get token from sendbird
	getToken := p.CreateTokenUserSendbirdV4(ctxReq, params)
	token := getToken.Result.(serviceModel.SessionTokenResponse)

	var tokenMetadata serviceModel.SessionTokenRequest
	var metadataBodyV4 serviceModel.MetadataRequestV4
	var metadataResponseV4 serviceModel.MetaDataResponseV4

	tokenMetadata.ExpiresAt = token.ExpiresAt

	uriMetadata := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, metaData)
	tokenData, _ := json.Marshal(tokenMetadata)

	// bind token and merchant to metadata interface
	metadataBodyV4.MetadataV4.Token = string(tokenData)
	bodyMetadataMershal, _ := json.Marshal(metadataBodyV4)
	payloadMetadata := strings.NewReader(string(bodyMetadataMershal))

	// request for create metadata
	err = helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uriMetadata, payloadMetadata, &metadataResponseV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: err, Result: metadataBodyV4}
	}

	if metadataResponseV4.Error {
		e := errors.New(metadataResponseV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_metadata_user_sendbird", e, uriMetadata)
		return serviceModel.ServiceResult{Error: e, Result: metadataResponseV4}
	}

	responseV4.UserID = responseUserV4.UserID
	responseV4.NickName = responseUserV4.NickName
	responseV4.ProfileURL = responseUserV4.ProfileURL
	responseV4.MetadataV4.Token = string(tokenData)

	return serviceModel.ServiceResult{Result: responseV4}
}

func (p *SendbirdService) GetUserSendbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-GetSendbirdUser"

	var response serviceModel.SendbirdResponse
	var responseUser serviceModel.SendbirdStringResponse

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)

	// request for create user
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodGet, uri, nil, &responseUser, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUser}
	}

	if responseUser.Error {
		e := errors.New(responseUser.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUser}
	}

	// request for get token from sendbird
	getToken := p.CreateTokenUserSendbird(ctxReq, params)
	token := getToken.Result.(serviceModel.SessionTokenResponse)

	response.UserID = responseUser.UserID
	response.NickName = responseUser.NickName
	response.ProfileURL = responseUser.ProfileURL
	response.Metadata.Merchant = params.Metadata.Merchant
	response.Metadata.Token = token

	return serviceModel.ServiceResult{Result: response}
}

func (p *SendbirdService) GetUserSendbirdV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult {
	ctx := "Service-GetSendbirdUserV4"

	var responseV4 serviceModel.SendbirdResponseV4
	var responseUserV4 serviceModel.SendbirdStringResponse

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)

	// request for create user
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodGet, uri, nil, &responseUserV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird_v4", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUserV4}
	}

	if responseUserV4.Error {
		e := errors.New(responseUserV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "create_user_sendbird_v4", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUserV4}
	}

	// request for get token from sendbird
	getToken := p.CreateTokenUserSendbirdV4(ctxReq, params)
	token := getToken.Result.(serviceModel.SessionTokenResponse)

	responseV4.UserID = responseUserV4.UserID
	responseV4.NickName = responseUserV4.NickName
	responseV4.ProfileURL = responseUserV4.ProfileURL
	responseV4.MetadataV4.Reference = params.MetadataV4.Reference
	responseV4.MetadataV4.Token = token

	return serviceModel.ServiceResult{Result: responseV4}
}

// UpdateUserSendbird function for getting sendbird user
func (p *SendbirdService) UpdateUserSendbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-UpdateSendbirdUser"

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	var bodyMetadata serviceModel.MetadataRequestV1
	var responseMetadata serviceModel.MetaDataResponseV1

	uriMetadata := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, metaData)
	tokenData, _ := json.Marshal(params.Metadata.Token)

	// bind data to Interface Metadata Request
	bodyMetadata.Metadata.Token = string(tokenData)

	bodyMetadataMershal, _ := json.Marshal(bodyMetadata)
	payloadMetadata := strings.NewReader(string(bodyMetadataMershal))

	// request for update metadata
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPut, uriMetadata, payloadMetadata, &responseMetadata, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "update_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: err, Result: bodyMetadata}
	}

	if responseMetadata.Error {
		e := errors.New(responseMetadata.Message)
		helper.SendErrorLog(ctxReq, ctx, "update_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: e, Result: responseMetadata}
	}

	var responseUser serviceModel.SendbirdStringResponse
	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)
	bodyMershal, _ := json.Marshal(params)
	payloadUpdate := strings.NewReader(string(bodyMershal))

	// request for update user
	err = helper.GetHTTPNewRequest(ctxReq, http.MethodPut, uri, payloadUpdate, &responseUser, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "update_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUser}
	}

	if responseUser.Error {
		e := errors.New(responseUser.Message)
		helper.SendErrorLog(ctxReq, ctx, "update_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUser}
	}

	return serviceModel.ServiceResult{Result: responseUser}
}

// UpdateUserSendbird function for getting sendbird user
func (p *SendbirdService) UpdateUserSendbirdV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult {
	ctx := "Service-UpdateSendbirdUserV4"

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}

	var bodyMetadataV4 serviceModel.MetadataRequestV4
	var responseMetadataV4 serviceModel.MetaDataResponseV4

	uriMetadata := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, metaData)
	tokenData, _ := json.Marshal(params.MetadataV4.Token)

	// bind data to Interface Metadata Request
	bodyMetadataV4.MetadataV4.Token = string(tokenData)

	bodyMetadataMershal, _ := json.Marshal(bodyMetadataV4)
	payloadMetadata := strings.NewReader(string(bodyMetadataMershal))

	// request for update metadata
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPut, uriMetadata, payloadMetadata, &responseMetadataV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "update_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: err, Result: bodyMetadataV4}
	}

	if responseMetadataV4.Error {
		e := errors.New(responseMetadataV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "update_metadata_user_sendbird", err, uriMetadata)
		return serviceModel.ServiceResult{Error: e, Result: responseMetadataV4}
	}

	var responseUserV4 serviceModel.SendbirdStringResponseV4
	uri := fmt.Sprintf(urlStringThree, p.BaseURL, urlGroup, params.UserID)
	bodyMershal, _ := json.Marshal(params)
	payloadUpdate := strings.NewReader(string(bodyMershal))

	// request for update user
	err = helper.GetHTTPNewRequest(ctxReq, http.MethodPut, uri, payloadUpdate, &responseUserV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "update_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseUserV4}
	}

	if responseUserV4.Error {
		e := errors.New(responseUserV4.Message)
		helper.SendErrorLog(ctxReq, ctx, "update_user_sendbird", e, uri)
		return serviceModel.ServiceResult{Error: e, Result: responseUserV4}
	}

	return serviceModel.ServiceResult{Result: responseUserV4}
}

// CreateTokenUserSendbird function for getting session token sendbird user
func (p *SendbirdService) CreateTokenUserSendbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-GetTokenSendbirdUser"

	var response serviceModel.SessionTokenResponse

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}
	uri := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, tokenString)

	type RequestBody struct {
		ExpiresAt int64 `json:"expires_at"`
	}

	var body RequestBody
	exp := params.ExpiresAt * 1000
	body.ExpiresAt = exp
	bodyMershal, _ := json.Marshal(body)
	payload := strings.NewReader(string(bodyMershal))

	// request for update session token
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uri, payload, &response, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: response}
	}

	if response.Error {
		e := errors.New(response.Message)
		if response.Code != 400201 {
			helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", e, uri)
		}
		return serviceModel.ServiceResult{Error: e, Result: response}
	}

	return serviceModel.ServiceResult{Result: response}
}

// CreateTokenUserSendbirdV4 function for getting session token sendbird user
func (p *SendbirdService) CreateTokenUserSendbirdV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult {
	ctx := "Service-GetTokenSendbirdUserV4"

	var responseV4 serviceModel.SessionTokenResponse

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}
	uri := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, tokenString)

	type RequestBody struct {
		ExpiresAt int64 `json:"expires_at"`
	}

	var bodyV4 RequestBody
	exp := params.ExpiresAt * 1000
	bodyV4.ExpiresAt = exp
	bodyMershal, _ := json.Marshal(bodyV4)
	payload := strings.NewReader(string(bodyMershal))

	// request for update session token
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uri, payload, &responseV4, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: responseV4}
	}

	if responseV4.Error {
		e := errors.New(responseV4.Message)
		if responseV4.Code != 400201 {
			helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", e, uri)
		}
		return serviceModel.ServiceResult{Error: e, Result: responseV4}
	}

	return serviceModel.ServiceResult{Result: responseV4}
}

// GetTokenUserSendbird function for getting session token sendbird user
func (p *SendbirdService) GetTokenUserSendbird(ctxReq context.Context, params *serviceModel.SendbirdRequest) serviceModel.ServiceResult {
	ctx := "Service-GetTokenSendbirdUser"

	var response serviceModel.SessionTokenResponse

	// generate headers
	headers := map[string]string{
		apiTokenHeader: p.APIToken,
	}
	uri := fmt.Sprintf(urlString, p.BaseURL, urlGroup, params.UserID, tokenString)
	payload := strings.NewReader(string("{}"))
	// request for update session token
	err := helper.GetHTTPNewRequest(ctxReq, http.MethodPost, uri, payload, &response, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", err, uri)
		return serviceModel.ServiceResult{Error: err, Result: response}
	}

	if response.Error {
		e := errors.New(response.Message)
		if response.Code != 400201 {
			helper.SendErrorLog(ctxReq, ctx, "get_token_user_sendbird", e, uri)
		}
		return serviceModel.ServiceResult{Error: e, Result: response}
	}

	return serviceModel.ServiceResult{Result: response}
}
