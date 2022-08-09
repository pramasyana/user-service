package delivery

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	errResultInvalidToken = "result is not token"
	errBasicAuth          = "basic auth is invalid"
)

// AuthHandlerV3 model v3
type AuthHandlerV3 struct {
	AuthUseCase           usecase.AuthUseCase
	GoogleAuthRedirectURL string
}

// NewHTTPHandler v3 function for initialise *AuthHandlerV3
func NewHTTPHandler(authUseCase usecase.AuthUseCase, gURL string) *AuthHandlerV3 {
	return &AuthHandlerV3{AuthUseCase: authUseCase, GoogleAuthRedirectURL: gURL}
}

// MountRoute as is
func (h *AuthHandlerV3) MountRoute(group *echo.Group) {
	group.POST("", h.GetAccessToken)
	group.POST("/verify-b2b", h.VerifyTokenB2b)
	group.POST("/check-email", h.CheckEmailV3)
}

// GetAccessToken function for getting access token
func (h *AuthHandlerV3) GetAccessToken(c echo.Context) error {
	ctx := "AuthHandlerV3-GetAccessToken"

	// parse client id and secret

	payload := model.RequestToken{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	clientID, clientSecret, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, "invalid basic auth").JSON(c)
	}
	statusCode, err := h.ValidateData(c, &payload, ctx, clientID, clientSecret)
	if err != nil {
		return shared.NewHTTPResponse(statusCode, err.Error()).JSON(c)
	}
	mode := payload.Mode
	payload.Version = helper.Version3

	tokenResultV3 := <-h.AuthUseCase.GenerateToken(context.WithValue(c.Request().Context(), middleware.ContextKeyClientIP, payload.IP), mode, payload)
	if tokenResultV3.Error != nil {
		if tokenResultV3.HTTPStatus == http.StatusForbidden {
			token, ok := tokenResultV3.Result.(model.MFAResponse)
			if !ok {
				return shared.NewHTTPResponse(http.StatusInternalServerError, errResultInvalidToken).JSON(c)
			}
			return shared.NewHTTPResponse(http.StatusForbidden, tokenResultV3.Error.Error(), token).JSON(c)
		}

		if tokenResultV3.Error.Error() == model.ErrorAccountInActiveBahasa {
			result := model.InactiveResponse{}
			result.Status = memberModel.InactiveString
			return shared.NewHTTPResponse(tokenResultV3.HTTPStatus, tokenResultV3.Error.Error(), result).JSON(c)
		}

		return shared.NewHTTPResponse(tokenResultV3.HTTPStatus, tokenResultV3.Error.Error()).JSON(c)
	}

	if check, ok := tokenResultV3.Result.(model.AuthV3TokenResponse); ok {
		return shared.NewHTTPResponse(http.StatusOK, "Get Social Media Info", check).JSON(c)
	}

	tokenV3, ok := tokenResultV3.Result.(model.RequestToken)
	if !ok {
		err := errors.New(errResultInvalidToken)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	response := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       tokenV3.UserID,
		Email:        strings.ToLower(tokenV3.Email),
		FirstName:    tokenV3.FirstName,
		LastName:     tokenV3.LastName,
		FullName:     tokenV3.FullName,
		NewMember:    tokenV3.NewMember,
		HasPassword:  tokenV3.HasPassword,
		Token:        tokenV3.Token,
		RefreshToken: tokenV3.RefreshToken,
		ExpiredTime:  tokenV3.ExpiredAt.Format(time.RFC3339),
		MemberType:   tokenV3.MemberType,
		Department:   tokenV3.Department,
		JobTitle:     tokenV3.JobTitle,
		Mobile:       tokenV3.Mobile,
		AccountID:    tokenV3.AccountID,
	}

	if response.NewMember {
		_ = <-h.AuthUseCase.SendEmailWelcomeMember(c.Request().Context(), response)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Auth Response V3", response).JSON(c)
}

// VerifyToken for verify token b2b
func (h *AuthHandlerV3) VerifyTokenB2b(c echo.Context) error {
	var payloads struct {
		Token           string `json:"token" form:"token"`
		TransactionType string `json:"transaction_type" form:"transaction_type"`
		MemberType      string `json:"member_type" form:"member_type"`
	}
	if err := c.Bind(&payloads); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	if c.Request().Header.Get(helper.TextAuthorization) != "" {
		clientNames, clientSecrets, oks := c.Request().BasicAuth()
		if oks {
			valid := <-h.AuthUseCase.GetClientApp(clientNames, clientSecrets)
			if valid.Error != nil {
				return shared.NewHTTPResponse(http.StatusUnauthorized, valid.Error.Error()).JSON(c)
			}
		}
	}

	tokenResult := <-h.AuthUseCase.VerifyTokenMemberB2b(c.Request().Context(), payloads.Token, payloads.TransactionType, payloads.MemberType)
	if tokenResult.Error != nil {
		if tokenResult.HTTPStatus == http.StatusForbidden {
			token, ok := tokenResult.Result.(model.MFAResponse)
			if !ok {
				err := errors.New(errResultInvalidToken)
				return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
			}
			return shared.NewHTTPResponse(http.StatusForbidden, tokenResult.Error.Error(), token).JSON(c)
		}

		if tokenResult.Error.Error() == model.ErrorAccountInActiveBahasa {
			result := model.InactiveResponse{}
			result.Status = memberModel.InactiveString
			return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error(), result).JSON(c)
		}

		return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error()).JSON(c)
	}

	reqTokens, ok := tokenResult.Result.(model.RequestToken)
	if !ok {
		err := errors.New(errResultInvalidToken)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	resp := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       reqTokens.UserID,
		Email:        strings.ToLower(reqTokens.Email),
		FirstName:    reqTokens.FirstName,
		LastName:     reqTokens.LastName,
		FullName:     reqTokens.FullName,
		NewMember:    reqTokens.NewMember,
		HasPassword:  reqTokens.HasPassword,
		Token:        reqTokens.Token,
		RefreshToken: reqTokens.RefreshToken,
		ExpiredTime:  reqTokens.ExpiredAt.Format(time.RFC3339),
		MemberType:   reqTokens.MemberType,
		Department:   reqTokens.Department,
		JobTitle:     reqTokens.JobTitle,
		Mobile:       reqTokens.Mobile,
	}

	return shared.NewHTTPResponse(http.StatusOK, "Auth Response", resp).JSON(c)
}

// ValidateData form for generate token
// return httpStatus and error
func (h *AuthHandlerV3) ValidateData(c echo.Context, payload *model.RequestToken, ctx, clientID, clientSecret string) (httpStatus int, err error) {
	// validate basic auth
	validateBA := <-h.AuthUseCase.ValidateBasicAuth(c.Request().Context(), clientID, clientSecret)

	if validateBA.Error != nil {
		return validateBA.HTTPStatus, validateBA.Error
	}

	email := strings.Trim(payload.Username, " ")
	emailParam := payload.Email
	if emailParam != "" {
		email = strings.Trim(emailParam, " ")
	}

	payload.Audience = clientID
	payload.DeviceID = strings.Trim(payload.DeviceID, " ")
	payload.Email = strings.ToLower(email)
	payload.IP = c.RealIP()
	payload.UserAgent = c.Request().UserAgent()

	// for refreshing token needs old token and old refresh token
	if payload.GrantType == model.AuthTypeRefreshToken {
		payload.Token = payload.OldToken
	}

	return http.StatusOK, nil
}

func (h *AuthHandlerV3) CheckEmailV3(c echo.Context) error {
	// parse client id and secret
	_, _, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errBasicAuth).JSON(c)
	}

	checkEmail := model.CheckEmailPayload{}

	if err := c.Bind(&checkEmail); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	mErr := jsonschema.ValidateTemp("check_email_param", checkEmail)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	res := <-h.AuthUseCase.CheckEmailV3(c.Request().Context(), checkEmail)
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success check email", res.Result).JSON(c)
}
