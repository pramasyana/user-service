package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

const (
	errBasicAuth      = "basic auth is invalid"
	errInvalidAuth    = "authentication is invalid"
	validateBasicAuth = "validate_basic_auth"
)

func (au *AuthUseCaseImpl) validateInput(data *model.RequestToken) error {
	// validate if memberType parameter exist
	allowedMemberType := []string{
		model.UserTypeCorporate,
		model.UserTypePersonal,
		model.UserTypeMicrosite,
		model.UserTypeClientMicrosite,
	}
	if data.MemberType != "" && !golib.StringInSlice(data.MemberType, allowedMemberType, false) && !au.isMicrositeClient(data.MemberType) {
		err := errors.New("cannot find member type")
		return err
	}

	// validates `deviceId` when the `grantType` is not refreshed token
	if data.GrantType != model.AuthTypeRefreshToken && len(data.DeviceID) <= 0 {
		err := errors.New("device id is required")
		return err
	}

	// validate `code` when `grantType` as social media login
	if (data.GrantType == model.AuthTypeFacebook || data.GrantType == model.AuthTypeGoogle || data.GrantType == model.AuthTypeAzure || data.GrantType == model.AuthTypeApple) && len(data.Code) == 0 {
		err := fmt.Errorf("%s code is required", strings.ToLower(data.GrantType))
		return err
	}

	// validate redirect uri for azure
	if data.GrantType == model.AuthTypeAzure && len(data.RedirectURI) == 0 {
		err := fmt.Errorf("%s is required", "redirect uri")
		return err
	}

	//  validates `deviceLogin` when the `grantType` is not refreshed token
	if data.GrantType != model.AuthTypeRefreshToken && !helper.StringInSlice(data.DeviceLogin, []string{"WEB", "MOBILE", "APPS"}) {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "device login")
		return err
	}
	// set personal as default `memberType` when `memberType` parameters is blank
	if data.MemberType == "" {
		data.MemberType = "personal"
	}
	return nil
}

// ValidateBasicAuth move from handler to usecase
func (au *AuthUseCaseImpl) ValidateBasicAuth(ctxReq context.Context, clientID, clientSecret string) <-chan ResultUseCase {
	// validate basic auth
	ctx := "AuthUseCaseImpl-ValidateBasicAuth"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		if len(clientID) == 0 || len(clientSecret) == 0 {
			output <- ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errors.New(errBasicAuth)}
			return
		}
		tags["clientId"] = clientID

		// get basic auth first
		basic := <-au.GetClientApp(clientID, clientSecret)
		if basic.Error != nil {
			output <- ResultUseCase{HTTPStatus: basic.HTTPStatus, Error: errors.New(errBasicAuth)}
			return
		}

		// rechecking the result of basic auth
		b := basic.Result.(bool)
		if !b {
			output <- ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errors.New(errInvalidAuth)}
			return
		}
	})

	return output
}
