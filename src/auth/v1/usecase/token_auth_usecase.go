package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/rsa"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	sessionInfoModel "github.com/Bhinneka/user-service/src/session/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// GenerateToken function for generating token based on grant type
func (au *AuthUseCaseImpl) GenerateToken(ctxReq context.Context, mode string, data model.RequestToken) <-chan ResultUseCase {
	ctx := "AuthUseCase-GenerateToken"

	outputGT := make(chan ResultUseCase, 1)
	if mode == "" {
		mode = "code"
	}

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputGT)

		tags["email"] = data.Email
		tags["grantType"] = data.GrantType
		tags["memberType"] = data.MemberType
		// default MFA Enabled
		data.MFAEnabled = false

		claimsGT := token.Claim{}
		claimsGT.Issuer = model.Bhinneka
		claimsGT.DeviceID = data.DeviceID
		claimsGT.DeviceLogin = data.DeviceLogin
		claimsGT.Audience = data.Audience
		claimsGT.Subject = data.Audience
		claimsGT.CustomToken = data.TokenBela
		var (
			redisUserID string
			err         error
			httpStatus  int
		)

		if valErr := au.validateInput(&data); valErr != nil {
			outputGT <- ResultUseCase{Error: valErr, HTTPStatus: http.StatusBadRequest}
			return
		}

		claimsGT.MemberType = data.MemberType

		// process the authentication based on grant type
		httpStatus, redisUserID, err = au.parseGlobalData(ctxReq, mode, &data, &claimsGT)
		if err != nil {
			tags[helper.TextResponse] = err.Error()
			outputGT <- ResultUseCase{Error: err, HTTPStatus: httpStatus}
			return
		}
		if resp := au.showOnlyProfileResponse(ctxReq, &data, httpStatus); resp != nil {
			outputGT <- ResultUseCase{Result: resp}
			return
		}

		// only for active mfa (multi factor authentication) & grant (password, fb, google) always get challange otp
		// not applicable for isAdmin = true
		if mfaRes, err := au.showMFAResponse(ctxReq, &data, claimsGT); mfaRes != nil {
			outputGT <- ResultUseCase{Result: mfaRes, Error: err, HTTPStatus: http.StatusForbidden}
			return
		}

		// generate token based on claims
		tokenResult := <-au.AccessTokenGenerator.GenerateAccessToken(claimsGT)
		data.Token = tokenResult.AccessToken.AccessToken
		data.ExpiredAt = tokenResult.AccessToken.ExpiredAt

		// special condition to pass when the grant type is anonymous
		// then generate refresh token
		if redisUserID == "" {
			redisUserID = model.DefaultSubject
		}

		// redis key format: STG-USR123-ASX1234-WEB
		redisLoginKey := strings.Join([]string{"STG", redisUserID, data.DeviceID, claimsGT.DeviceLogin}, "-")
		httpStatus, err = au.generateRefreshToken(ctxReq, redisLoginKey, &data, &claimsGT)
		if err != nil {
			outputGT <- ResultUseCase{Error: err, HTTPStatus: httpStatus}
			return
		}

		now := time.Now()
		expired := data.ExpiredAt.Sub(now)

		paramRedis := &model.LoginSessionRedis{
			Key:         redisLoginKey,
			Token:       data.Token,
			ExpiredTime: expired,
		}

		tags["rediskey"] = redisLoginKey
		tags["token"] = data.Token

		if err := au.saveTokenToRedis(ctxReq, paramRedis, data, tokenResult.AccessToken.JTI); err != nil {
			outputGT <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		outputGT <- ResultUseCase{Result: data}
	})

	return outputGT
}

// GenerateToken function for generating token based on grant type
func (au *AuthUseCaseImpl) GenerateTokenB2b(ctxReq context.Context, mode string, data model.RequestToken) <-chan ResultUseCase {
	ctx := "AuthUseCase-GenerateToken"

	output := make(chan ResultUseCase, 1)
	if mode == "" {
		mode = "code"
	}

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["email"] = data.Email
		tags["grantType"] = data.GrantType
		tags["memberType"] = data.MemberType
		// default MFA Enabled
		data.MFAEnabled = false

		claims := token.Claim{}
		claims.Issuer = model.Bhinneka
		claims.DeviceID = data.DeviceID
		claims.DeviceLogin = data.DeviceLogin
		claims.Audience = data.Audience
		claims.Subject = data.Audience
		claims.CustomToken = data.TokenBela
		var (
			redisUserID string
			err         error
			httpStatus  int
		)

		claims.MemberType = data.MemberType

		// process the authentication based on grant type
		httpStatus, redisUserID, err = au.parsePasswordTypeMicrositeBela(ctxReq, &data, &claims)
		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: httpStatus}
			return
		}
		if resp := au.showOnlyProfileResponse(ctxReq, &data, httpStatus); resp != nil {
			output <- ResultUseCase{Result: resp}
			return
		}

		// only for active mfa (multi factor authentication) & grant (password, fb, google) always get challange otp
		// not applicable for isAdmin = true
		if mfaRes, err := au.showMFAResponse(ctxReq, &data, claims); mfaRes != nil {
			output <- ResultUseCase{Result: mfaRes, Error: err, HTTPStatus: http.StatusForbidden}
			return
		}

		// generate token based on claims
		tokenResult := <-au.AccessTokenGenerator.GenerateAccessToken(claims)
		data.Token = tokenResult.AccessToken.AccessToken
		data.ExpiredAt = tokenResult.AccessToken.ExpiredAt

		// special condition to pass when the grant type is anonymous
		// then generate refresh token
		if redisUserID == "" {
			redisUserID = model.DefaultSubject
		}

		// redis key format: STG-USR123-ASX1234-WEB
		redisLoginKey := strings.Join([]string{"STG", redisUserID, data.DeviceID, claims.DeviceLogin}, "-")

		now := time.Now()
		expired := data.ExpiredAt.Sub(now)

		paramRedis := &model.LoginSessionRedis{
			Key:         redisLoginKey,
			Token:       data.Token,
			ExpiredTime: expired,
		}

		tags["rediskey"] = redisLoginKey
		tags["token"] = data.Token

		if err := au.saveTokenToRedis(ctxReq, paramRedis, data, tokenResult.AccessToken.JTI); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

func (au *AuthUseCaseImpl) saveTokenToRedis(ctxReq context.Context, paramRedis *model.LoginSessionRedis, data model.RequestToken, jti string) error {
	saveRedis := <-au.LoginSessionRepo.Save(ctxReq, paramRedis)
	if saveRedis.Error != nil {
		return saveRedis.Error
	}
	if err := au.saveSessionInfo(data, jti); err != nil {
		return err
	}
	return nil
}

func (au *AuthUseCaseImpl) showMFAResponse(ctxReq context.Context, data *model.RequestToken, claims token.Claim) (interface{}, error) {
	var redisMFAKey string

	if data.MFAEnabled {
		// redis key format: `mfa-otp-USR123-ABC-WEB`
		redisMFAKey = model.MFATokenKeyRedis
	} else if data.GrantType == model.AuthTypeLDAP && data.NarwhalMFAEnabled {
		// specific for ldap
		redisMFAKey = model.NarwhalMFATokenKeyRedis
	}
	// means mfa for both account and narwhal are disabled
	if redisMFAKey == "" {
		return nil, nil
	}

	redisLoginKey := strings.Join([]string{redisMFAKey, data.UserID, data.DeviceID, claims.DeviceLogin}, "-")
	expired, _ := time.ParseDuration("5m")
	mix := data.Email + "-" + data.UserID + "-" + data.DeviceID + "-" + data.DeviceLogin
	mfaToken := helper.GenerateTokenByString(mix)
	paramRedis := &model.LoginSessionRedis{
		Key:         redisLoginKey,
		Token:       mfaToken,
		ExpiredTime: expired,
	}
	saveMFA := <-au.LoginSessionRepo.Save(ctxReq, paramRedis)
	if saveMFA.Error != nil {
		helper.SendErrorLog(ctxReq, "saveTokenMFA", "save_token_mfa", saveMFA.Error, paramRedis)
	}

	// encode parameter memberID
	base64EncodeEmail := base64.URLEncoding.EncodeToString([]byte(data.UserID))
	mfaTokenCombine := mfaToken + "-" + base64EncodeEmail

	mfaResponse := model.MFAResponse{}
	mfaResponse.MFARequired = true
	mfaResponse.MFAToken = mfaTokenCombine

	return mfaResponse, errors.New(memberModel.ErrorMFARequired)

}
func (au *AuthUseCaseImpl) showOnlyProfileResponse(ctxReq context.Context, data *model.RequestToken, httpStatus int) interface{} {
	ctx := "AuthUseCase-showOnlyProfileResponse"
	if data.Version == helper.Version3 && httpStatus == http.StatusForbidden {
		resp := model.AuthV3TokenResponse{
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
		}

		if data.GrantType == model.AuthTypeGoogle {
			resp.GoogleID = data.UserID
		} else if data.GrantType == model.AuthTypeFacebook {
			resp.FacebookID = data.UserID
		} else if data.GrantType == model.AuthTypeApple {
			resp.AppleID = data.UserID
		}

		// store key di redis
		payload := model.LoginSessionRedis{
			Key:         fmt.Sprintf("STG:%s:%s", data.Email, data.UserID),
			Token:       base64.StdEncoding.EncodeToString([]byte(data.UserID)),
			ExpiredTime: time.Minute * 5,
		}
		saveTemp := <-au.LoginSessionRepo.Save(ctxReq, &payload)
		if saveTemp.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_redis", saveTemp.Error, payload)
		}
		return resp
	}
	return nil
}

func (au *AuthUseCaseImpl) saveSessionInfo(data model.RequestToken, JTI string) error {
	if data.GrantType != model.AuthTypeAnonymous && data.GrantType != model.AuthTypeRefreshToken {
		params := sessionInfoModel.SessionInfoRequest{
			GrantType:   data.GrantType,
			UserID:      data.UserID,
			Email:       strings.ToLower(data.Email),
			IP:          data.IP,
			UserAgent:   data.UserAgent,
			DeviceID:    data.DeviceID,
			DeviceLogin: data.DeviceLogin,
			JTI:         JTI,
		}

		saveSession := <-au.SessionInfoRepo.SaveSessionInfo(&params)
		if saveSession.Error != nil {
			return saveSession.Error
		}
	}
	return nil
}
func (au *AuthUseCaseImpl) parsePasswordType(ctxReq context.Context, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	if data.MemberType == model.UserTypeCorporate {
		memberData, httpStatus, err := au.RequestTokenPasswordTypeCorporate(ctxReq, data)
		if err != nil {
			return httpStatus, "", err
		}

		claims.Subject = strconv.Itoa(memberData.ID)
		claims.Authorised = true
		claims.IsAdmin = false
		claims.IsStaff = false
		claims.Email = memberData.Email

		data.UserID = strconv.Itoa(memberData.ID)
		data.Email = memberData.Email
		data.FirstName = memberData.FirstName
		data.LastName = memberData.LastName
		data.Mobile = memberData.PhoneNumber
		data.NewMember = false
		data.HasPassword = true
		data.AccountID = memberData.AccountID
		redisUserID = strconv.Itoa(memberData.ID)
	} else if data.MemberType == model.UserTypeMicrositeBela {
		memberData, httpStatus, err := au.RequestTokenPasswordTypeMicrositeBela(ctxReq, data)
		if err != nil {
			return httpStatus, "", err
		}

		claims.Subject = strconv.Itoa(memberData.ID)
		claims.Authorised = true
		claims.IsAdmin = false
		claims.IsStaff = false
		claims.Email = memberData.Email

		data.UserID = strconv.Itoa(memberData.ID)
		data.Email = memberData.Email
		data.FirstName = memberData.FirstName
		data.LastName = memberData.LastName
		data.Mobile = memberData.PhoneNumber
		data.NewMember = false
		data.HasPassword = true
		redisUserID = strconv.Itoa(memberData.ID)

		// save to redis lpseID
		now := time.Now()
		expired := data.ExpiredAt.Sub(now)

		redisLpseIDKey := strings.Join([]string{"STG-BELA", redisUserID, memberData.Email}, "-")
		paramsLpseID := &model.LoginSessionRedis{
			Key:         redisLpseIDKey,
			Token:       data.LpseID,
			ExpiredTime: expired,
		}
		saveRedisLpseID := <-au.LoginSessionRepo.Save(ctxReq, paramsLpseID)
		if saveRedisLpseID.Error != nil {
			return http.StatusInternalServerError, "", saveRedisLpseID.Error
		}

	} else if au.isMicrositeClient(data.MemberType) {
		// validate request grantType password for microsite
		memberData, httpStatus, err := au.RequestTokenPasswordTypeMicrosite(ctxReq, data)
		if err != nil {
			return httpStatus, "", err
		}

		claims.Subject = strconv.Itoa(memberData.ID)
		claims.Authorised = true
		claims.IsAdmin = false
		claims.IsStaff = false
		claims.Email = memberData.Email

		data.UserID = strconv.Itoa(memberData.ID)
		data.Email = memberData.Email
		data.FirstName = memberData.FirstName
		data.LastName = memberData.LastName
		data.Mobile = memberData.PhoneNumber
		data.NewMember = false
		data.HasPassword = true
		redisUserID = strconv.Itoa(memberData.ID)
	} else if data.MemberType == model.UserTypeClientMicrosite { // specific for endpoint /v1/client/login
		memberData, httpStatus, err := au.RequestTokenClientMicrosite(ctxReq, data)
		if err != nil {
			return httpStatus, "", err
		}
		claims.Subject = memberData.ID
		claims.Authorised = true
		claims.IsAdmin = false
		claims.IsStaff = false
		claims.Email = memberData.Email
		claims.SignUpFrom = memberData.SignUpFrom
		claims.Audience = memberData.SignUpFrom

		data.UserID = memberData.ID
		data.Email = memberData.Email
		data.FirstName = memberData.FirstName
		data.LastName = memberData.LastName
		data.Mobile = memberData.Mobile
		data.NewMember = false
		data.HasPassword = true
		redisUserID = memberData.ID
	} else {
		memberData, httpStatus, err := au.RequestTokenPasswordType(ctxReq, data)
		if err != nil {
			return httpStatus, "", err
		}

		claims.Subject = memberData.ID
		claims.Authorised = true
		claims.IsAdmin = memberData.IsAdmin
		claims.IsStaff = memberData.IsStaff
		claims.Email = memberData.Email
		claims.SignUpFrom = memberData.SignUpFrom

		data.UserID = memberData.ID
		data.Email = memberData.Email
		data.FirstName = memberData.FirstName
		data.LastName = memberData.LastName
		data.Mobile = memberData.Mobile
		data.NewMember = false
		data.HasPassword = true
		redisUserID = memberData.ID
		data.MFAEnabled = memberData.MFAEnabled
	}
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parsePasswordTypeMicrositeBela(ctxReq context.Context, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {

	memberData, httpStatus, err := au.RequestTokenPasswordTypeMicrositeBela(ctxReq, data)
	if err != nil {
		return httpStatus, "", err
	}

	claims.Subject = strconv.Itoa(memberData.ID)
	claims.Authorised = true
	claims.IsAdmin = false
	claims.IsStaff = false
	claims.Email = memberData.Email

	data.UserID = strconv.Itoa(memberData.ID)
	data.Email = memberData.Email
	data.FirstName = memberData.FirstName
	data.LastName = memberData.LastName
	data.Mobile = memberData.PhoneNumber
	data.NewMember = false
	data.HasPassword = true
	data.AccountID = memberData.AccountID
	redisUserID = strconv.Itoa(memberData.ID)

	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) isMicrositeClient(memberType string) bool {
	return memberType == model.UserTypeMicrosite || strings.HasPrefix(memberType, "MICROSITE_")
}

func (au *AuthUseCaseImpl) parseAzureType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseAzureType"
	var token string
	token = data.Code
	if mode == "code" {
		azureToken := <-au.AuthQueryOAuth.GetAzureToken(ctxReq, data.Code, data.RedirectURI)
		if azureToken.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_azure_token", azureToken.Error, data)
			if strings.Contains(azureToken.Error.Error(), model.ErrorAzureToken) {
				err = fmt.Errorf(model.ErrorAzureTokenBahasa)
			}

			if strings.Contains(azureToken.Error.Error(), model.ErrorAzureInvalidRedirectURL) {
				err = fmt.Errorf(model.ErrorInvalidRedirectURL)
			}
			return http.StatusUnauthorized, "", err
		}

		accessToken := azureToken.Result.(model.AuthAzureToken)

		token = accessToken.Token
	}

	azureResult := <-au.AuthQueryOAuth.GetDetailAzureMember(ctxReq, token)
	if azureResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_azure_data", azureResult.Error, data)
		return http.StatusUnauthorized, "", azureResult.Error
	}

	azure, ok := azureResult.Result.(model.AzureResponse)
	if !ok {
		err := errors.New("result is not azure response")
		return http.StatusUnauthorized, "", err
	}

	if len(azure.Email) <= 0 {
		err := errors.New(msgErrorEmailRegisterAccount)
		return http.StatusUnauthorized, "", err
	}

	dataMember := au.CheckMemberSocmedType(ctxReq, data, azure, azure.Email)
	if dataMember.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "check_member_sosmed", dataMember.Error, data)
		return dataMember.HTTPStatus, "", dataMember.Error
	}

	member := dataMember.Data
	claims.Subject = member.ID
	claims.Authorised = true
	claims.IsAdmin = member.IsAdmin
	claims.IsStaff = member.IsStaff
	claims.Email = member.Email
	claims.SignUpFrom = member.SignUpFrom

	data.UserID = member.ID
	data.Email = member.Email
	data.FirstName = member.FirstName
	data.LastName = member.LastName
	data.Mobile = member.Mobile
	data.NewMember = dataMember.NewMember
	data.HasPassword = false
	if len(member.Password) > 0 && !dataMember.NewMember {
		data.HasPassword = true
	}
	redisUserID = member.ID
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parseFacebookType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseFacebookType"
	var token string
	token = data.Code

	if mode == "code" {
		facebookToken := <-au.AuthQueryOAuth.GetFacebookToken(ctxReq, data.Code, data.RedirectURI)
		if facebookToken.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "getting_facebook_token", facebookToken.Error, data)
			return http.StatusUnauthorized, "", facebookToken.Error
		}

		accessToken := facebookToken.Result.(model.AuthFacebookToken)
		token = accessToken.AccessToken
	}

	facebookResult := <-au.AuthQueryOAuth.GetDetailFacebookMember(ctxReq, token)
	if facebookResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_facebook_data", facebookResult.Error, data)
		return http.StatusUnauthorized, "", facebookResult.Error
	}

	facebook, ok := facebookResult.Result.(model.FacebookResponse)
	if !ok {
		return http.StatusInternalServerError, "", errors.New("result is not facebook response")
	}

	if len(facebook.Email) <= 0 {
		return http.StatusUnauthorized, "", errors.New(msgErrorEmailRegisterAccount)
	}

	if data.Email != "" && data.Email != facebook.Email {
		return http.StatusUnauthorized, "", errors.New(msgFailedLoginSocmedEmail)
	}

	dataMember := au.CheckMemberSocmedType(ctxReq, data, facebook, facebook.Email)
	err = dataMember.Error
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeValidateRequestTokenFacebook, err, data)
		return dataMember.HTTPStatus, "", err
	}
	// if request coming from version 3 and member not yet register then just return with available data
	if data.Version == helper.Version3 && dataMember.HTTPStatus == http.StatusForbidden {
		data.UserID = facebook.ID
		data.FirstName = facebook.FirstName
		data.LastName = facebook.LastName
		data.Email = facebook.Email
		return http.StatusForbidden, "", nil
	}

	member := dataMember.Data

	claims.Subject = member.ID
	claims.Authorised = true
	claims.IsAdmin = member.IsAdmin
	claims.IsStaff = false
	claims.Email = member.Email
	claims.SignUpFrom = member.SignUpFrom

	data.UserID = member.ID
	data.Email = member.Email
	data.FirstName = member.FirstName
	data.LastName = member.LastName
	data.Mobile = member.Mobile
	data.NewMember = dataMember.NewMember
	data.MFAEnabled = member.MFAEnabled

	if len(member.Password) > 0 {
		data.HasPassword = true
	}
	redisUserID = member.ID
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parseGoogleType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseGoogleType"
	var token string
	token = data.Code
	if mode == "code" {
		googleToken := <-au.AuthQueryOAuth.GetGoogleToken(ctxReq, data.Code, data.RedirectURI)
		if googleToken.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "getting_google_token", err, data)
			if strings.Contains(googleToken.Error.Error(), "Bad Request") {
				googleToken.Error = errors.New("invalid google credentials")
			}
			return http.StatusUnauthorized, "", googleToken.Error
		}

		accessToken := googleToken.Result.(model.AuthGoogleToken)
		token = accessToken.AccessToken
	}

	googleResult := <-au.AuthQueryOAuth.GetDetailGoogleMember(ctxReq, token)
	if googleResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_google_data", googleResult.Error, data)
		return http.StatusBadRequest, "", googleResult.Error
	}

	google, ok := googleResult.Result.(model.GoogleOAuth2Response)
	if !ok {
		return http.StatusInternalServerError, "", errors.New("result is not google response")
	}

	if len(google.Email) <= 0 {
		return http.StatusUnauthorized, "", errors.New(msgErrorEmailRegisterAccount)
	}

	if data.Email != "" && data.Email != google.Email {
		return http.StatusUnauthorized, "", errors.New(msgFailedLoginSocmedEmail)
	}

	dataMember := au.CheckMemberSocmedType(ctxReq, data, google, google.Email)
	err = dataMember.Error
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "check_member_socmed_type", err, google)
		return dataMember.HTTPStatus, "", err
	}
	// if request coming from version 3 and member not yet register then just return with available data
	if data.Version == helper.Version3 && dataMember.HTTPStatus == http.StatusForbidden {
		data.UserID = google.ID
		data.FirstName = google.GivenName
		data.LastName = google.FamilyName
		data.Email = google.Email
		return http.StatusForbidden, "", nil
	}

	member := dataMember.Data

	claims.Subject = member.ID
	claims.Authorised = true
	claims.IsAdmin = member.IsAdmin
	claims.IsStaff = member.IsStaff
	claims.Email = member.Email
	claims.SignUpFrom = member.SignUpFrom

	data.UserID = member.ID
	data.Email = member.Email
	data.FirstName = member.FirstName
	data.LastName = member.LastName
	data.Mobile = member.Mobile
	data.NewMember = dataMember.NewMember
	data.MFAEnabled = member.MFAEnabled

	if len(member.Password) > 0 {
		data.HasPassword = true
	}
	redisUserID = member.ID
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parseGoogleOAuthType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseGoogleOAuthType"
	var token string
	token = data.Code

	googleResult := <-au.AuthQueryOAuth.GetGoogleTokenInfo(ctxReq, token)
	if googleResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_google_data", googleResult.Error, data)
		return http.StatusBadRequest, "", googleResult.Error
	}

	google, ok := googleResult.Result.(model.GoogleOAuthToken)
	if !ok {
		return http.StatusInternalServerError, "", errors.New("result is not google response")
	}

	if len(google.Email) <= 0 {
		return http.StatusUnauthorized, "", errors.New(msgErrorEmailRegisterAccount)
	}

	dataMember := au.CheckMemberSocmedType(ctxReq, data, google, google.Email)
	err = dataMember.Error
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "check_member_socmed_type", err, google)
		return dataMember.HTTPStatus, "", err
	}

	// if request coming from version 3 and member not yet register then just return with available data
	if data.Version == helper.Version3 && dataMember.HTTPStatus == http.StatusForbidden {
		data.UserID = google.Sub
		data.FirstName = google.GivenName
		data.LastName = google.FamilyName
		data.Email = google.Email
		return http.StatusForbidden, "", nil
	}

	member := dataMember.Data

	claims.Subject = member.ID
	claims.Authorised = true
	claims.IsAdmin = member.IsAdmin
	claims.IsStaff = member.IsStaff
	claims.Email = member.Email
	claims.SignUpFrom = member.SignUpFrom

	data.UserID = member.ID
	data.Email = member.Email
	data.FirstName = member.FirstName
	data.LastName = member.LastName
	data.Mobile = member.Mobile
	data.NewMember = dataMember.NewMember
	data.MFAEnabled = member.MFAEnabled

	if len(member.Password) > 0 {
		data.HasPassword = true
	}
	redisUserID = member.ID
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parseRefreshTokenType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseRefreshTokenType"
	pubKey, err := rsa.InitPublicKey()
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeGetPublicKey, err, data)
		return http.StatusInternalServerError, "", err
	}

	oldClaims, err := token.VerifyTokenIgnoreExpiration(pubKey, data.Token)
	if err != nil {
		return http.StatusUnauthorized, "", err
	}

	RTKey := getRefreshTokenKey(oldClaims.Subject, oldClaims.DeviceID, oldClaims.DeviceLogin)

	refreshResult := <-au.RefreshTokenRepo.Load(ctxReq, RTKey)
	if refreshResult.Error != nil {
		return http.StatusUnauthorized, "", errors.New(model.ErrorGetToken)
	}

	rt, ok := refreshResult.Result.(model.RefreshToken)
	if !ok {
		return http.StatusUnauthorized, "", errors.New("result is not a refresh token")
	}

	if !rt.Match(data.RefreshToken) {
		return http.StatusUnauthorized, "", errors.New(model.ErrorRefreshToken)
	}

	if err := au.checkSubjectType(ctxReq, data, claims, oldClaims); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "check_token_subject_type", err, oldClaims)
		return http.StatusUnauthorized, "", err
	}

	claims.DeviceID = oldClaims.DeviceID
	claims.DeviceLogin = oldClaims.DeviceLogin
	claims.Subject = oldClaims.Subject
	claims.Authorised = true
	claims.IsAdmin = oldClaims.IsAdmin
	claims.Email = oldClaims.Email
	claims.MemberType = oldClaims.MemberType
	data.MemberType = oldClaims.MemberType

	redisUserID = oldClaims.Subject
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) checkSubjectType(ctxReq context.Context, data *model.RequestToken, claims *token.Claim, oldToken *token.BearerClaims) error {
	ctx := "AuthUseCaseImpl-checkSubjectType"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{"token": data.Email}
	defer tr.Finish(tags)

	if oldToken.MemberType == "" || oldToken.MemberType == "personal" {
		if err := au.assignRefreshTokenPersonal(ctxReq, data, claims, oldToken); err != nil {
			return err
		}
	} else {
		if err := au.assignRefreshTokenCorporate(ctxReq, data, claims, oldToken); err != nil {
			return err
		}
	}

	data.DeviceID = oldToken.DeviceID
	data.UserID = oldToken.Subject
	data.NewMember = false

	return nil
}

func (au *AuthUseCaseImpl) assignRefreshTokenPersonal(ctxReq context.Context, data *model.RequestToken, claims *token.Claim, oldClaims *token.BearerClaims) error {
	ctx := "AuthUseCaseImpl-assignRefreshTokenPersonal"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{"token": data.Email}
	defer tr.Finish(tags)

	memberResult := <-au.MemberRepoRead.Load(ctxReq, oldClaims.Subject)
	if memberResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "load_member", memberResult.Error, oldClaims)
		return memberResult.Error
	}
	member, ok := memberResult.Result.(memberModel.Member)
	if !ok {
		return errors.New(msgResultNotMember)
	}
	// append value of data variable
	data.UserID = oldClaims.Subject
	data.Email = member.Email
	data.FirstName = member.FirstName
	data.LastName = member.LastName
	data.Mobile = member.Mobile
	data.DeviceID = oldClaims.DeviceID
	data.NewMember = false
	data.FullName = member.FirstName + " " + member.LastName
	data.MemberType = model.UserTypePersonal
	data.JobTitle = member.JobTitle
	if len(member.Password) > 0 {
		data.HasPassword = true
	}
	return nil
}

func (au *AuthUseCaseImpl) assignRefreshTokenCorporate(ctxReq context.Context, data *model.RequestToken, claims *token.Claim, oldClaims *token.BearerClaims) error {
	ctx := "AuthUseCaseImpl-assignRefreshTokenCorporate"
	if strings.HasPrefix(oldClaims.Subject, "USR") {
		return errors.New(msgUserIdNotValid)
	}
	contactResult := <-au.CorporateContactQueryRead.FindByID(ctxReq, oldClaims.Subject)
	if contactResult.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, "load_contact", contactResult.Error, oldClaims)
		return contactResult.Error
	}
	contact, ok := contactResult.Result.(sharedModel.B2BContactData)
	if !ok {
		return errors.New(msgResultNotMember)
	}
	data.UserID = oldClaims.Subject
	data.Email = contact.Email
	data.FirstName = contact.FirstName
	data.LastName = contact.LastName
	data.DeviceID = oldClaims.DeviceID
	data.NewMember = false
	data.MemberType = claims.MemberType
	data.JobTitle = contact.JobTitle
	if len(contact.Password) > 0 {
		data.HasPassword = true
	}
	return nil
}

func (au *AuthUseCaseImpl) parseLDAPType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseLDAPType"
	tr := tracer.StartTrace(ctxReq, ctx)
	defer tr.Finish(nil)
	if errValidation := jsonschema.ValidateTemp("auth_get_token_ldap_params", data); errValidation != nil {
		return http.StatusBadRequest, "", errValidation
	}

	ldapResult, err := au.AuthServices.Auth(ctxReq, data.Email, data.Password)

	if err != nil {
		return http.StatusUnauthorized, "", err
	}

	claims.Subject = ldapResult.ObjectID
	claims.Authorised = true
	claims.IsAdmin = false
	claims.IsStaff = false
	claims.Email = ""

	data.UserID = ""
	data.Email = ""
	data.FirstName = ldapResult.FirstName
	data.LastName = ldapResult.LastName
	data.FullName = ldapResult.DisplayName
	data.NewMember = false
	data.Department = ldapResult.Department
	data.JobTitle = ldapResult.JobTitle
	data.HasPassword = true
	redisUserID = ldapResult.ObjectID

	if len(ldapResult.Email) > 0 {
		dataMember := au.CheckMemberSocmedType(tr.NewChildContext(), data, ldapResult, ldapResult.Email)
		if dataMember.Error != nil {
			return dataMember.HTTPStatus, "", dataMember.Error
		}

		member := dataMember.Data

		claims.Subject = member.ID
		claims.Authorised = true
		claims.IsAdmin = member.IsAdmin
		claims.IsStaff = member.IsStaff
		claims.Email = member.Email
		claims.SignUpFrom = member.SignUpFrom

		data.UserID = member.ID
		data.Email = member.Email
		data.FirstName = member.FirstName
		data.LastName = member.LastName
		data.Mobile = member.Mobile
		data.FullName = member.FirstName + " " + member.LastName
		data.NewMember = false
		data.Department = ldapResult.Department
		data.JobTitle = ldapResult.JobTitle
		data.NarwhalMFAEnabled = member.AdminMFAEnabled

		if len(member.Password) > 0 {
			data.HasPassword = true
		}
		redisUserID = member.ID
	}
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) parseVerifyMFAType(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-GenerateToken-parseVerifyMFAType"
	memberID, newToken, err := au.parseTokenMFA(data.MFAToken)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "parse_mfa_token", err, data)
		return http.StatusBadRequest, "", err
	}
	var redisOtpKey string
	if data.GrantType == model.AuthTypeVerifyMFANarwhal {
		redisOtpKey = model.NarwhalMFATokenKeyRedis
	} else {
		redisOtpKey = model.MFATokenKeyRedis
	}

	redisLoginKey := strings.Join([]string{redisOtpKey, memberID, data.DeviceID, claims.DeviceLogin}, "-") //mfa-otp-USR123-ABC-WEB
	getTokenRedis := <-au.LoginSessionRepo.Load(ctxReq, redisLoginKey)
	if getTokenRedis.Error != nil {
		if getTokenRedis.Error.Error() != helper.ErrorRedis {
			helper.SendErrorLog(ctxReq, ctx, "mfa_load_redis", getTokenRedis.Error, data)
		}

		return http.StatusBadRequest, "", errors.New(memberModel.ErrorTokenMFA)
	}

	tokenRedis := ""
	if existingToken, ok := getTokenRedis.Result.(model.LoginSessionRedis); ok {
		tokenRedis = existingToken.Token
	}

	if tokenRedis != newToken {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "existing token")
		return http.StatusUnauthorized, "", err
	}

	memberResult := <-au.MemberQueryRead.FindByID(ctxReq, memberID)
	if memberResult.Result == nil {
		return http.StatusUnauthorized, "", fmt.Errorf(helper.ErrorDataNotFound, "token")
	}

	memberData, ok := memberResult.Result.(memberModel.Member)
	if !ok {
		return http.StatusInternalServerError, "", errors.New(msgResultNotMember)
	}
	var mfaKeyDB string
	if data.GrantType == model.AuthTypeVerifyMFANarwhal {
		mfaKeyDB = memberData.MFAAdminKey
	} else {
		mfaKeyDB = memberData.MFAKey
	}

	mfaKey, _ := base64.URLEncoding.DecodeString(mfaKeyDB)
	validateOTP := <-au.verifyMFACode(ctxReq, string(mfaKey), data.OTP)
	if validateOTP.Error != nil {
		return validateOTP.HTTPStatus, "", validateOTP.Error
	}

	claims.Subject = memberData.ID
	claims.Authorised = true
	claims.IsAdmin = memberData.IsAdmin
	claims.IsStaff = memberData.IsStaff
	claims.Email = memberData.Email
	claims.SignUpFrom = memberData.SignUpFrom

	data.UserID = memberData.ID
	data.Email = memberData.Email
	data.FirstName = memberData.FirstName
	data.LastName = memberData.LastName
	data.Mobile = memberData.Mobile
	data.NewMember = false
	if memberData.Password != "" {
		data.HasPassword = true
	}
	redisUserID = memberData.ID
	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) generateRefreshToken(ctxReq context.Context, redisLoginKey string, data *model.RequestToken, claims *token.Claim) (httpStatus int, err error) {
	ctx := "AuthUseCase-GenerateToken-generateRefreshToken"
	if data.GrantType != model.AuthTypeAnonymous {

		refreshToken := helper.RandomStringBase64(38)
		var rtAge time.Duration

		emails := strings.Split(au.EmailSpecialTokenAge, ",")
		if golib.StringInSlice(claims.Email, emails) {
			rtAge, _ = time.ParseDuration(au.SpecialRefreshTokenAge)
		} else {
			rtAge, _ = time.ParseDuration(au.RefreshTokenAge)
		}

		newRTokenKey := getRefreshTokenKey(claims.Subject, claims.DeviceID, claims.DeviceLogin)
		rToken := model.NewRefreshToken(newRTokenKey, refreshToken, rtAge)

		refreshTokenResult := <-au.RefreshTokenRepo.Save(ctxReq, rToken)
		if refreshTokenResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_refresh_token", refreshTokenResult.Error, claims)
			return http.StatusInternalServerError, errors.New("failed to generate refresh token")
		}

		data.RefreshToken = rToken.Token

		// update last login when generate token successful
		au.AuthQueryDB.UpdateLastLogin(data.UserID)
	}
	return http.StatusOK, nil
}

func (au *AuthUseCaseImpl) parseGlobalData(ctxReq context.Context, mode string, data *model.RequestToken, claims *token.Claim) (httpStatus int, redisUserID string, err error) {
	switch data.GrantType {
	case model.AuthTypeAnonymous:
		claims.Authorised = false
		redisUserID = data.Audience

	case model.AuthTypePassword:
		return au.parsePasswordType(ctxReq, data, claims)

	case model.AuthTypeAzure:
		return au.parseAzureType(ctxReq, mode, data, claims)

	case model.AuthTypeFacebook:
		return au.parseFacebookType(ctxReq, mode, data, claims)

	case model.AuthTypeGoogle, model.AuthTypeGoogleBackend:
		return au.parseGoogleType(ctxReq, mode, data, claims)

	case model.AuthTypeGoogleOAauth:
		return au.parseGoogleOAuthType(ctxReq, mode, data, claims)

	case model.AuthTypeApple:
		return au.generateAppleToken(ctxReq, claims, data)

	case model.AuthTypeRefreshToken:
		return au.parseRefreshTokenType(ctxReq, mode, data, claims)

	case model.AuthTypeLDAP:
		return au.parseLDAPType(ctxReq, mode, data, claims)

	case model.AuthTypeVerifyMFA, model.AuthTypeVerifyMFANarwhal:
		return au.parseVerifyMFAType(ctxReq, mode, data, claims)

	default:
		return httpStatus, "", errors.New("invalid grant type")
	}
	return http.StatusOK, redisUserID, nil
}
