package delivery

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	errBasicAuth            = "basic auth is invalid"
	errInvalidAuth          = "authentication is invalid"
	validateBasicAuth       = "validate_basic_auth"
	msgErrorResultToken     = "result is not token"
	labelGenerateTokenError = "GenerateTokenErrorParse"
	backendDeviceLogin      = "WEB"
	backendRequestFrom      = "sturgeon-backend"
	backendDeviceID         = "STGBackend"
	scopeCastTokenError     = "cast_token_result_error"
	scopeTokenResultError   = "token_result_error"
)

// HTTPAuthHandler model
type HTTPAuthHandler struct {
	AuthUseCase           usecase.AuthUseCase
	GoogleAuthRedirectURL string
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(authUseCase usecase.AuthUseCase, gURL string) *HTTPAuthHandler {
	return &HTTPAuthHandler{AuthUseCase: authUseCase, GoogleAuthRedirectURL: gURL}
}

// Mount function for mounting routes
func (h *HTTPAuthHandler) Mount(group *echo.Group) {
	group.POST("", h.GetAccessToken)
	group.POST("/verify", h.VerifyToken)
	group.POST("/verify-member", h.VerifyTokenMember)
	group.POST("/logout", h.Logout)
	group.POST("/check-email", h.CheckEmail)
	group.POST("/verify-captcha", h.VerifyCaptcha)
	group.POST("/client-app", h.CreateClientApp)
	group.GET("/oauth2callback", h.AuthCallback)
}

// MountAdmin admin only
func (h *HTTPAuthHandler) MountAdmin(group *echo.Group) {
	group.GET("/auth", h.GetAccessTokenFromUserID)
}

// AuthCallback from google
func (h *HTTPAuthHandler) AuthCallback(c echo.Context) error {
	data := model.RequestToken{
		GrantType:   model.AuthTypeGoogleBackend,
		Code:        c.QueryParam("code"),
		RedirectURI: h.GoogleAuthRedirectURL,
		DeviceLogin: backendDeviceLogin,
		RequestFrom: backendRequestFrom,
		DeviceID:    backendDeviceID,
	}
	tokenGenerator := <-h.AuthUseCase.GenerateToken(c.Request().Context(), "code", data)
	if tokenGenerator.Error != nil {
		if tokenGenerator.HTTPStatus == http.StatusForbidden {
			token, ok := tokenGenerator.Result.(model.MFAResponse)
			if !ok {
				err := errors.New(msgErrorResultToken)
				return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
			}
			return shared.NewHTTPResponse(http.StatusForbidden, tokenGenerator.Error.Error(), token).JSON(c)
		}

		if tokenGenerator.Error.Error() == model.ErrorAccountInActiveBahasa {
			result := model.InactiveResponse{}
			result.Status = memberModel.InactiveString
			return shared.NewHTTPResponse(tokenGenerator.HTTPStatus, tokenGenerator.Error.Error(), result).JSON(c)
		}else if tokenGenerator.Error.Error() == model.ErrorNewAccountBahasa{
			result := model.InactiveResponse{}
			result.Status = memberModel.NewString
			return shared.NewHTTPResponse(tokenGenerator.HTTPStatus, tokenGenerator.Error.Error(), result).JSON(c)
		}

		return shared.NewHTTPResponse(tokenGenerator.HTTPStatus, tokenGenerator.Error.Error()).JSON(c)
	}
	token, ok := tokenGenerator.Result.(model.RequestToken)
	if !ok {
		err := errors.New(msgErrorResultToken)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	res := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       token.UserID,
		Email:        strings.ToLower(token.Email),
		FirstName:    token.FirstName,
		LastName:     token.LastName,
		FullName:     token.FullName,
		NewMember:    token.NewMember,
		HasPassword:  token.HasPassword,
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiredTime:  token.ExpiredAt.Format(time.RFC3339),
		MemberType:   token.MemberType,
		Department:   token.Department,
		JobTitle:     token.JobTitle,
		Mobile:       token.Mobile,
	}
	return shared.NewHTTPResponse(http.StatusOK, "Token Response", res).JSON(c)
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
		return shared.NewHTTPResponse(http.StatusBadRequest, clientAppResult.Error.Error()).JSON(c)
	}

	clientApp, ok := clientAppResult.Result.(*model.ClientApp)

	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "result is not client app")
	}
	clientApp.Secret = base64.StdEncoding.EncodeToString([]byte(clientApp.ClientID + ":" + clientApp.ClientSecret))

	return shared.NewHTTPResponse(http.StatusOK, "Client App Response", clientApp).JSON(c)
}

// ValidateBasicAuth function
func (h *HTTPAuthHandler) ValidateBasicAuth(clientID string, clientSecret string) usecase.ResultUseCase {
	// get basic auth first
	basic := <-h.AuthUseCase.GetClientApp(clientID, clientSecret)
	if basic.Error != nil {
		return usecase.ResultUseCase{Error: basic.Error, HTTPStatus: basic.HTTPStatus}
	}

	// rechecking the result of basic auth
	b := basic.Result.(bool)
	if !b {
		return usecase.ResultUseCase{Error: errors.New(errInvalidAuth), HTTPStatus: http.StatusUnauthorized}
	}

	return usecase.ResultUseCase{Result: b}
}

// GetAccessToken function for getting access token
func (h *HTTPAuthHandler) GetAccessToken(c echo.Context) error {
	ctx := "AuthPresenter-GetAccessToken"

	// parse client id and secret
	clientID, clientSecret, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errBasicAuth).JSON(c)
	}

	payload := model.RequestToken{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	data, errMsg, statusCode := h.ValidateData(c, payload, ctx, clientID, clientSecret)
	if errMsg != "" {
		return shared.NewHTTPResponse(statusCode, errMsg).JSON(c)
	}

	mode := payload.Mode

	// requestFrom flag for sending email welcome if request from sturgeon CF
	requestFrom := payload.RequestFrom

	tokenResult := <-h.AuthUseCase.GenerateToken(context.WithValue(c.Request().Context(), middleware.ContextKeyClientIP, c.RealIP()), mode, data)
	if tokenResult.Error != nil {
		if tokenResult.HTTPStatus == http.StatusForbidden {
			token, ok := tokenResult.Result.(model.MFAResponse)
			if !ok {
				err := errors.New(msgErrorResultToken)
				return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
			}
			return shared.NewHTTPResponse(http.StatusForbidden, tokenResult.Error.Error(), token).JSON(c)
		}

		if tokenResult.Error.Error() == model.ErrorAccountInActiveBahasa {
			result := model.InactiveResponse{}
			result.Status = memberModel.InactiveString
			return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error(), result).JSON(c)
		} else if tokenResult.Error.Error() == model.ErrorNewAccountBahasa{
			result := model.InactiveResponse{}
			result.Status = memberModel.NewString
			return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error(), result).JSON(c)
		}

		return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error()).JSON(c)
	}

	reqToken, ok := tokenResult.Result.(model.RequestToken)
	if !ok {
		err := errors.New(msgErrorResultToken)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	res := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       reqToken.UserID,
		Email:        strings.ToLower(reqToken.Email),
		FirstName:    reqToken.FirstName,
		LastName:     reqToken.LastName,
		FullName:     reqToken.FullName,
		NewMember:    reqToken.NewMember,
		HasPassword:  reqToken.HasPassword,
		Token:        reqToken.Token,
		RefreshToken: reqToken.RefreshToken,
		ExpiredTime:  reqToken.ExpiredAt.Format(time.RFC3339),
		MemberType:   reqToken.MemberType,
		Department:   reqToken.Department,
		JobTitle:     reqToken.JobTitle,
		Mobile:       reqToken.Mobile,
	}

	h.sendWelcomeEmail(c.Request().Context(), requestFrom, ctx, res)

	return shared.NewHTTPResponse(http.StatusOK, "Auth Response", res).JSON(c)
}

// SendWelcomeEmail form for generate token login
func (h *HTTPAuthHandler) sendWelcomeEmail(ctxReq context.Context, requestFrom, ctx string, res model.AccessTokenResponse) {
	if requestFrom == memberModel.Sturgeon && res.NewMember {
		<-h.AuthUseCase.SendEmailWelcomeMember(ctxReq, res)
	}
}

// ValidateData form for generate token
func (h *HTTPAuthHandler) ValidateData(c echo.Context, data model.RequestToken, ctx, clientID, clientSecret string) (model.RequestToken, string, int) {

	// validate basic auth
	validateBA := h.ValidateBasicAuth(clientID, clientSecret)

	if validateBA.Error != nil {
		return data, validateBA.Error.Error(), validateBA.HTTPStatus
	}

	email := strings.Trim(data.Username, " ")
	emailParam := data.Email
	if emailParam != "" {
		email = strings.Trim(emailParam, " ")
	}

	data.Audience = clientID
	data.DeviceID = strings.Trim(data.DeviceID, " ")
	data.Email = strings.ToLower(email)
	data.IP = c.RealIP()
	data.UserAgent = c.Request().UserAgent()

	// for refreshing token needs old token and old refresh token
	if data.GrantType == model.AuthTypeRefreshToken {
		data.Token = data.OldToken
	}

	return data, "", 0
}

// GetAccessTokenFromUserID function for getting access token
func (h *HTTPAuthHandler) GetAccessTokenFromUserID(c echo.Context) error {
	userID := c.QueryParam("userId")
	if err := middleware.ExtractClaimsIsAdmin(c); err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errInvalidAuth).JSON(c)
	}

	// set data for request token parameter
	data := model.RequestToken{}
	data.Audience = "dolphin"
	data.DeviceID = "DOLPHIN-MANUAL-ORDER"
	data.DeviceLogin = "WEB"
	data.UserID = userID

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	tokenResult := <-h.AuthUseCase.GenerateTokenFromUserID(newCtx, data)

	if tokenResult.Error != nil {
		return shared.NewHTTPResponse(tokenResult.HTTPStatus, tokenResult.Error.Error()).JSON(c)
	}

	token, ok := tokenResult.Result.(model.RequestToken)
	if !ok {
		err := errors.New(msgErrorResultToken)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	res := model.AccessTokenResponse{
		ID:           golib.RandomString(8),
		UserID:       token.UserID,
		Email:        strings.ToLower(token.Email),
		FirstName:    token.FirstName,
		LastName:     token.LastName,
		NewMember:    token.NewMember,
		HasPassword:  token.HasPassword,
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiredTime:  token.ExpiredAt.Format(time.RFC3339),
	}

	return shared.NewHTTPResponse(http.StatusOK, "Auth Response", res).JSON(c)
}

// VerifyToken for verify token
func (h *HTTPAuthHandler) VerifyToken(c echo.Context) error {
	var payload struct {
		Token string `json:"token" form:"token"`
	}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	if c.Request().Header.Get(helper.TextAuthorization) != "" {
		clientName, clientSecret, ok := c.Request().BasicAuth()
		if ok {
			valid := <-h.AuthUseCase.GetClientApp(clientName, clientSecret)
			if valid.Error != nil {
				return shared.NewHTTPResponse(http.StatusUnauthorized, valid.Error.Error()).JSON(c)
			}
		}
	}

	res := h.AuthUseCase.VerifyTokenMember(c.Request().Context(), payload.Token)
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, res.Error.Error()).JSON(c)
	}
	if res.Result == nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, "error").JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success verify token", res.Result).JSON(c)
}

// VerifyTokenMember for verify token for member user
func (h *HTTPAuthHandler) VerifyTokenMember(c echo.Context) error {
	return h.VerifyToken(c)
}

// Logout function user logout
func (h *HTTPAuthHandler) Logout(c echo.Context) error {
	// parse client id and secret
	clientID, clientSecret, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errBasicAuth).JSON(c)
	}

	// validate basic auth
	validateBA := h.ValidateBasicAuth(clientID, clientSecret)

	if validateBA.Error != nil {
		return shared.NewHTTPResponse(validateBA.HTTPStatus, validateBA.Error.Error()).JSON(c)
	}

	var payload struct {
		Token string `json:"token" form:"token"`
	}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	resLogout := <-h.AuthUseCase.Logout(c.Request().Context(), payload.Token)
	if resLogout.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, resLogout.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success to logout", resLogout.Result).JSON(c)
}

// CheckEmail for verify email
func (h *HTTPAuthHandler) CheckEmail(c echo.Context) error {
	// parse client id and secret
	_, _, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errBasicAuth).JSON(c)
	}

	var payload struct {
		Email string `json:"email" form:"email"`
	}

	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	res := <-h.AuthUseCase.CheckEmail(c.Request().Context(), payload.Email)
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success check email", res.Result).JSON(c)
}

// VerifyCaptcha for verify email
func (h *HTTPAuthHandler) VerifyCaptcha(c echo.Context) error {
	// parse client id and secret
	_, _, ok := c.Request().BasicAuth()
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, errBasicAuth).JSON(c)
	}

	payload := model.GoogleCaptcha{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload.RemoteIP = c.RealIP()

	res := <-h.AuthUseCase.VerifyCaptcha(c.Request().Context(), payload)
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error(), res.Result).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success verify captcha", res.Result).JSON(c)
}
