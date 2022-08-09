package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	"github.com/golang-jwt/jwt"
)

func (au *AuthUseCaseImpl) generateAppleToken(ctxReq context.Context, claims *token.Claim, data *model.RequestToken) (httpStatus int, redisUserID string, err error) {
	ctx := "AuthUseCase-generateAppleToken"

	validateNameError := au.validateAppleParameter(data)
	if validateNameError != nil {
		return http.StatusBadRequest, "", validateNameError
	}

	// getting apple token and get response data
	appleToken := <-au.AuthQueryOAuth.GetAppleToken(ctxReq, data.Code, data.RedirectURI, data.ClientID)
	if appleToken.Error != nil {
		if strings.Contains(appleToken.Error.Error(), "Bad Request") {
			appleToken.Error = errors.New("invalid apple credentials")
		}
		helper.SendErrorLog(ctxReq, ctx, "getting_apple_token", appleToken.Error, data)
		return http.StatusUnauthorized, "", appleToken.Error
	}

	apple, ok := appleToken.Result.(model.AppleResponse)
	if !ok {
		err := errors.New("result is not apple response")
		helper.SendErrorLog(ctxReq, ctx, "error_parse_apple_data", err, appleToken.Result)
		return http.StatusInternalServerError, "", err
	}

	// parse data from apple token & parameter name from frontend callback
	appleProfile, err := au.parseTokenApple(ctxReq, apple.IDToken, data)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeErrorParseSocmed, err, apple)
		return http.StatusUnauthorized, "", err
	}

	// check member available for create or authorize user status
	dataMember := au.CheckMemberSocmedType(ctxReq, data, appleProfile, appleProfile.Email)
	if dataMember.Error != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeValidateToken, err, appleProfile)
		return dataMember.HTTPStatus, "", dataMember.Error
	}

	currentMember := dataMember.Data

	// append claims data
	claims.Subject = currentMember.ID
	claims.Authorised = true
	claims.IsAdmin = currentMember.IsAdmin
	claims.IsStaff = currentMember.IsStaff
	claims.Email = currentMember.Email
	claims.SignUpFrom = currentMember.SignUpFrom

	// append value of data variable
	data.UserID = currentMember.ID
	data.Email = currentMember.Email
	data.FirstName = currentMember.FirstName
	data.LastName = currentMember.LastName
	data.Mobile = currentMember.Mobile
	data.NewMember = dataMember.NewMember
	data.MFAEnabled = currentMember.MFAEnabled

	if len(currentMember.Password) > 0 {
		data.HasPassword = true
	}
	redisUserID = data.UserID

	return http.StatusOK, redisUserID, nil
}

func (au *AuthUseCaseImpl) validateAppleParameter(data *model.RequestToken) error {
	var (
		validFNChar, validFNLen, validLNChar, validLNLen bool
	)
	if data.FirstName != "" {
		if golib.ValidateAlphabetWithSpace(data.FirstName) {
			validFNChar = true
		}

		if err := golib.ValidateMaxInput(data.FirstName, 25); err == nil {
			validFNLen = true
		}
	}

	if len(data.LastName) != 0 {
		if golib.ValidateAlphabetWithSpace(data.LastName) {
			validLNChar = true
		}

		if err := golib.ValidateMaxInput(data.LastName, 25); err == nil {
			validLNLen = true
		}
	}
	if !validFNChar || !validFNLen {
		data.FirstName = helper.DefaultFirstName
	}
	if !validLNChar || !validLNLen {
		data.LastName = helper.DefaultLastName
	}

	return nil
}
func (au *AuthUseCaseImpl) parseTokenApple(ctxReq context.Context, token string, data *model.RequestToken) (model.AppleProfile, error) {
	appleProfile := model.AppleProfile{}
	claims := jwt.MapClaims{}
	jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != tkn.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", tkn.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if jwtResult == nil && err != nil {
		err := errors.New("invalid token")
		return appleProfile, err
	}
	parsedEmail, ok := claims["email"].(string)
	if !ok {
		return appleProfile, errors.New(msgErrorEmailRegisterAccount)
	}
	appleProfile.Email = parsedEmail

	isPrivateEmail, found := claims["is_private_email"]
	if found {
		appleProfile.IsPrivateEmail = isPrivateEmail.(string)
	}
	appleProfile.Sub = claims["sub"].(string)

	if data.FirstName != "" {
		appleProfile.FirstName = data.FirstName
	}

	if len(data.LastName) != 0 {
		appleProfile.LastName = data.LastName
	}

	if len(appleProfile.Email) == 0 {
		err := errors.New(msgErrorEmailRegisterAccount)
		return appleProfile, err
	}

	return appleProfile, nil
}
