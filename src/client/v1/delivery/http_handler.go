package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	au "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/client/v1/model"
	cu "github.com/Bhinneka/user-service/src/client/v1/usecase"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	invalidAuth             = "invalid client authorization"
	malformedBody           = "malformed payload"
	deviceIDBela            = "belaDID"
	msgErrorResultToken     = "result is not token"
	labelGenerateTokenError = "GenerateTokenErrorParse"
)

// HTTPClientHandler DI
type HTTPClientHandler struct {
	AuthUseCase     au.AuthUseCase
	ActivityService service.ActivityServices
	ClientUsecase   cu.ClientUsecase
}

// NewHTTPHandler return client handler
func NewHTTPHandler(authUseCase au.AuthUseCase, activityService service.ActivityServices, clientUC cu.ClientUsecase) *HTTPClientHandler {
	return &HTTPClientHandler{
		AuthUseCase:     authUseCase,
		ActivityService: activityService,
		ClientUsecase:   clientUC,
	}
}

// Mount return echo group
func (h *HTTPClientHandler) Mount(group *echo.Group, mf echo.MiddlewareFunc) {
	// URL => /v1/client/*
	vc := middleware.ValidateClient()
	group.POST("/login", h.Login, vc)
	group.POST("/logout", h.Logout, vc)
	group.GET("/verify", h.VerifyToken, mf)
}

// Login client login
func (h *HTTPClientHandler) Login(c echo.Context) (err error) {
	// ctx := "ClientAuthDelivery-GetAccessToken"
	clientID, clientSecret := helper.ExtractClientCred(c)
	valid := <-h.AuthUseCase.GetClientApp(clientID, clientSecret)
	if valid.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, invalidAuth).JSON(c)
	}

	param := new(model.LKPPUser)
	if err = c.Bind(param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, malformedBody).JSON(c)
	}
	ctxReq := c.Request().Context()
	tokenResult := <-h.AuthUseCase.GenerateToken(ctxReq, "", authModel.RequestToken{
		GrantType:   authModel.AuthTypePassword,
		MemberType:  authModel.UserTypeClientMicrosite,
		Email:       param.Payload.Email,
		DeviceID:    deviceIDBela,
		DeviceLogin: authModel.DefaultDeviceLogin,
		FirstName:   helper.DefaultLKPPName,
	})
	if tokenResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, tokenResult.Error.Error()).JSON(c)
	}
	token, ok := tokenResult.Result.(authModel.RequestToken)
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, msgErrorResultToken).JSON(c)
	}

	// return specific values only
	resp := authModel.ClientResponse{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		Email:        token.Email,
		ExpiredAt:    token.ExpiredAt,
	}
	jsReq, _ := json.Marshal(param)
	jsResp, _ := json.Marshal(resp)
	ctxHeader := strings.Join([]string{helper.TextBearer, token.Token}, " ")
	ctxLog := context.WithValue(ctxReq, helper.TextAuthorization, ctxHeader)
	ctxLog = context.WithValue(ctxLog, middleware.ContextKeyClientIP, c.RealIP())

	go h.ActivityService.CreateLog(ctxLog, serviceModel.Payload{
		Module: "ClientLogin",
		Action: helper.TextInsertUpper,
		Logs: []serviceModel.Log{
			{
				Field:    "Request",
				OldValue: string(jsReq),
				NewValue: string(jsResp),
			},
		},
		Target: clientID,
	})

	return shared.NewHTTPResponse(http.StatusOK, "Client Auth Response", resp).JSON(c)
}

// Logout client logout
func (h *HTTPClientHandler) Logout(c echo.Context) error {
	clientID, clientSecret := helper.ExtractClientCred(c)
	valid := <-h.AuthUseCase.GetClientApp(clientID, clientSecret)
	if valid.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, invalidAuth).JSON(c)
	}

	param := new(model.LKPPUser)
	if err := c.Bind(param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, malformedBody).JSON(c)
	}
	ctxReq := c.Request().Context()

	hh := <-h.ClientUsecase.Logout(ctxReq, param.Payload.Email)
	if hh.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, hh.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Client Logout Response", map[string]interface{}{"success": true, "error": nil}).JSON(c)
}

// VerifyToken verify given token
func (h *HTTPClientHandler) VerifyToken(c echo.Context) error {
	param := new(model.QueryParam)
	if err := c.Bind(param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, malformedBody).JSON(c)
	}
	ctxReq := c.Request().Context()

	tokenValidation := h.AuthUseCase.VerifyTokenMember(ctxReq, param.Token)
	if tokenValidation.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid token").JSON(c)
	}
	resp, ok := tokenValidation.Result.(*authModel.VerifyResponse)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid response").JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, "success verify token", h.cleanTokenResponse(resp)).JSON(c)
}

func (h *HTTPClientHandler) cleanTokenResponse(input *authModel.VerifyResponse) model.TokenResponse {
	return model.TokenResponse{
		Token:      input.Token,
		Exp:        input.Exp,
		IssueAt:    input.Iat,
		Issuer:     input.Iss,
		Email:      input.Email,
		SignUpFrom: input.SignUpFrom,
	}
}
