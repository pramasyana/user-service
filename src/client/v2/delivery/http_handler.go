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
	cuV2 "github.com/Bhinneka/user-service/src/client/v2/usecase"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	authInvalid             = "invalid client authorization"
	BodyMalformed           = "malformed payload"
	belaDeviceID            = "belaDID"
	ResultTokenMsgError     = "result is not token"
	GenerateTokenLabelError = "GenerateTokenErrorParse"
)

// HTTPClientHandler DI
type HTTPClientHandler struct {
	AuthUseCase     au.AuthUseCase
	ActivityService service.ActivityServices
	ClientUsecase   cu.ClientUsecase
	ClientV2Usecase cuV2.ClientUsecase
}

// NewHTTPHandler return client handler
func NewHTTPHandler(authUseCase au.AuthUseCase, activityService service.ActivityServices, clientUC cu.ClientUsecase, clientV2UC cuV2.ClientUsecase) *HTTPClientHandler {
	return &HTTPClientHandler{
		AuthUseCase:     authUseCase,
		ActivityService: activityService,
		ClientUsecase:   clientUC,
		ClientV2Usecase: clientV2UC,
	}
}

// Mount return echo group
func (h *HTTPClientHandler) Mount(group *echo.Group, mf echo.MiddlewareFunc) {
	// URL => /v2/client/*
	vc := middleware.ValidateClient()
	group.POST("/login", h.Login, vc)
	group.POST("/logout", h.Logout, vc)
	group.GET("/verify", h.VerifyToken, vc)
}

func (h *HTTPClientHandler) Login(c echo.Context) (err error) {
	// ctx := "ClientAuthDelivery-GetAccessToken"
	clientIDs, clientSecrets := helper.ExtractClientCred(c)
	valids := <-h.AuthUseCase.GetClientApp(clientIDs, clientSecrets)
	if valids.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, authInvalid).JSON(c)
	}

	params := new(model.LKPPUser)
	if err = c.Bind(params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, BodyMalformed).JSON(c)
	}
	ctxReqs := c.Request().Context()
	tokenResults := <-h.AuthUseCase.GenerateToken(ctxReqs, "", authModel.RequestToken{
		GrantType:       authModel.AuthTypePassword,
		MemberType:      authModel.UserTypeMicrositeBela,
		TransactionType: authModel.LoginTypeShopcart,
		Email:           params.Payload.Email,
		DeviceID:        belaDeviceID,
		DeviceLogin:     authModel.DefaultDeviceLogin,
		FirstName:       params.Payload.RealName,
		LpseID:          params.Payload.LpseID,
		TokenBela:       params.TokenBela,
	})
	if tokenResults.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, tokenResults.Error.Error()).JSON(c)
	}
	token, ok := tokenResults.Result.(authModel.RequestToken)
	if !ok {
		return shared.NewHTTPResponse(http.StatusUnauthorized, ResultTokenMsgError).JSON(c)
	}

	// return specific values only
	resp := authModel.ClientResponse{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		Email:        token.Email,
		ExpiredAt:    token.ExpiredAt,
	}
	jsReqs, _ := json.Marshal(params)
	jsResps, _ := json.Marshal(resp)
	ctxHeaders := strings.Join([]string{helper.TextBearer, token.Token}, " ")
	ctxLogs := context.WithValue(ctxReqs, helper.TextAuthorization, ctxHeaders)
	ctxLogs = context.WithValue(ctxLogs, middleware.ContextKeyClientIP, c.RealIP())

	go h.ActivityService.CreateLog(ctxLogs, serviceModel.Payload{
		Module: "ClientLogin",
		Action: helper.TextInsertUpper,
		Logs: []serviceModel.Log{
			{
				Field:    "Request",
				OldValue: string(jsReqs),
				NewValue: string(jsResps),
			},
		},
		Target: clientIDs,
	})

	return shared.NewHTTPResponse(http.StatusOK, "Client Auth Response", resp).JSON(c)
}

// Logout client logout
func (h *HTTPClientHandler) Logout(c echo.Context) error {
	clientID, clientSecret := helper.ExtractClientCred(c)
	valid := <-h.AuthUseCase.GetClientApp(clientID, clientSecret)
	if valid.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, authInvalid).JSON(c)
	}

	param := new(model.LKPPUser)
	if err := c.Bind(param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, BodyMalformed).JSON(c)
	}
	ctxReq := c.Request().Context()

	hh := <-h.ClientV2Usecase.Logout(ctxReq, param.Payload.Email)
	if hh.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, hh.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Client Logout Response", map[string]interface{}{"success": true, "error": nil}).JSON(c)
}

// VerifyToken verify given token
func (h *HTTPClientHandler) VerifyToken(c echo.Context) error {
	params := new(model.QueryParam)
	if err := c.Bind(params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, BodyMalformed).JSON(c)
	}
	ctxReqs := c.Request().Context()

	tokenValidations := h.AuthUseCase.VerifyTokenMember(ctxReqs, params.Token)
	if tokenValidations.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid token").JSON(c)
	}
	resps, ok := tokenValidations.Result.(*authModel.VerifyResponse)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid response").JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, "success verify token", h.cleanTokenResponse(resps)).JSON(c)
}

func (h *HTTPClientHandler) cleanTokenResponse(inputs *authModel.VerifyResponse) model.TokenResponse {
	return model.TokenResponse{
		Token:      inputs.Token,
		Exp:        inputs.Exp,
		IssueAt:    inputs.Iat,
		Issuer:     inputs.Iss,
		Email:      inputs.Email,
		SignUpFrom: inputs.SignUpFrom,
	}
}
