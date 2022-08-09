package delivery

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo"
)

// HTTPAuthHandler model
type HTTPAuthHandler struct {
	AuthUseCase usecase.AuthUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(authUseCase usecase.AuthUseCase) *HTTPAuthHandler {
	return &HTTPAuthHandler{AuthUseCase: authUseCase}
}

// Mount function for mounting routes
func (h *HTTPAuthHandler) Mount(group *echo.Group) {
	group.POST("", h.GetAccessToken)
}

// MountClientApp function
func (h *HTTPAuthHandler) MountClientApp(group *echo.Group) {
	group.POST("", h.CreateClientApp)
}

// CreateClientApp function
func (h *HTTPAuthHandler) CreateClientApp(c echo.Context) error {
	appName := c.FormValue("appName")

	clientAppResult := <-h.AuthUseCase.CreateClientApp(appName)

	if clientAppResult.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, clientAppResult.Error.Error())
	}

	clientApp, ok := clientAppResult.Result.(*model.ClientApp)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "result is not client app")
	}

	return json.NewEncoder(c.Response()).Encode(clientApp)
}

// GetAccessToken function for getting access token
func (h *HTTPAuthHandler) GetAccessToken(c echo.Context) error {
	ctx := "AuthPresenter-GetAccessToken"

	// parse client id and secret
	clientID, clientSecret, _ := c.Request().BasicAuth()
	requestData := model.RequestToken{}
	if err := c.Bind(&requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.ErrorPayload)
	}

	data, errMsg, statusCode := h.ValidateData(c, requestData, ctx, clientID, clientSecret)
	if errMsg != "" {
		return echo.NewHTTPError(statusCode, errMsg)
	}

	mode := data.Mode

	tokenResult := <-h.AuthUseCase.GenerateToken(context.WithValue(c.Request().Context(), middleware.ContextKeyClientIP, c.RealIP()), mode, data)

	if tokenResult.Error != nil {
		if tokenResult.HTTPStatus == http.StatusForbidden {
			token, ok := tokenResult.Result.(model.MFAResponse)
			if !ok {
				err := errors.New("result is not token")
				helper.SendErrorLog(c.Request().Context(), ctx, "parse_mfa", err, nil)
				return echo.NewHTTPError(tokenResult.HTTPStatus, tokenResult.Error.Error())
			}

			payload, err := helper.MarshalConvertOnePayload(&token)
			if err != nil {
				helper.SendErrorLog(c.Request().Context(), ctx, "jsonapi_generate_token", err, nil)
				return echo.NewHTTPError(http.StatusInternalServerError, "error occurred")
			}

			c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
			c.Response().WriteHeader(http.StatusForbidden)
			return json.NewEncoder(c.Response()).Encode(payload)
		}

		return echo.NewHTTPError(tokenResult.HTTPStatus, tokenResult.Error.Error())
	}

	token, ok := tokenResult.Result.(model.RequestToken)
	if !ok {
		err := errors.New("result is not token")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       token.UserID,
		Email:        token.Email,
		FirstName:    token.FirstName,
		LastName:     token.LastName,
		NewMember:    token.NewMember,
		HasPassword:  token.HasPassword,
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiredTime:  token.ExpiredAt.Format(time.RFC3339),
		MemberType:   token.MemberType,
		Mobile:       token.Mobile,
	}

	payload, err := helper.MarshalConvertOnePayload(&res)
	if err != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, "jsonapi_generate_token", err, res)
		return echo.NewHTTPError(http.StatusInternalServerError, "error occurred")
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// ValidateData form for generate token
func (h *HTTPAuthHandler) ValidateData(c echo.Context, data model.RequestToken, ctx, clientID, clientSecret string) (model.RequestToken, string, int) {
	// validate basic auth
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return data, "basic auth is invalid", http.StatusUnauthorized
	}

	// get basic auth first
	basic := <-h.AuthUseCase.GetClientApp(clientID, clientSecret)
	if basic.Error != nil {
		return data, basic.Error.Error(), basic.HTTPStatus
	}

	// rechecking the result of basic auth
	b := basic.Result.(bool)
	if !b {
		return data, "authentication is invalid", http.StatusUnauthorized
	}

	data.Audience = clientID
	data.DeviceID = strings.Trim(data.DeviceID, " ")
	data.Email = strings.Trim(data.Username, " ")
	data.IP = c.RealIP()
	data.UserAgent = c.Request().UserAgent()

	// for refreshing token needs old token and old refresh token
	if data.GrantType == model.AuthTypeRefreshToken {
		data.Token = data.OldToken
	}

	return data, "", 0
}
