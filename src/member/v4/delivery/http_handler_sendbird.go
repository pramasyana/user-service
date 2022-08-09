package delivery

import (
	"context"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	authUc "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo"
)

// MemberHandlerV4 model receiver
type MemberHandlerV4 struct {
	MemberUseCase usecase.MemberUseCase
	AuthUseCase   authUc.AuthUseCase
}

// NewHTTPHandlerV4 function for initialise *MemberHandlerV4
func NewHTTPHandlerV4(memberUseCase usecase.MemberUseCase, authUsecase authUc.AuthUseCase) *MemberHandlerV4 {
	return &MemberHandlerV4{
		MemberUseCase: memberUseCase,
		AuthUseCase:   authUsecase,
	}
}

// MountSendbird function for mounting routes
func (h *MemberHandlerV4) MountSendbird(group *echo.Group) {
	group.GET("/token", h.GetAccessToken)
	group.GET("/check-token", h.CheckAccessToken)
}

// GetAccessToken function for requesting email for forgot password
func (h *MemberHandlerV4) GetAccessToken(c echo.Context) error {
	var paramsV4 serviceModel.SendbirdRequestV4
	UserIDV4, errV4 := middleware.ExtractMemberIDFromToken(c)
	if errV4 != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errV4.Error())
	}

	claimsV4, errV4 := middleware.ExtractClaimsFromToken(c)
	if errV4 != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errV4.Error())
	}

	tokenContextV4 := c.Get("token")
	client := c.QueryParam("client")
	parsingTokenV4 := tokenContextV4.(*jwt.Token)

	paramsV4.UserID = UserIDV4
	paramsV4.ExpiresAt = claimsV4.StandardClaims.ExpiresAt
	paramsV4.Token = parsingTokenV4.Raw
	paramsV4.Client = client

	newCtxV4 := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))

	sendbirdUcV4 := h.MemberUseCase.GetSendbirdTokenV4(newCtxV4, &paramsV4)

	if sendbirdUcV4.Error != nil {
		return shared.NewHTTPResponse(sendbirdUcV4.HTTPStatus, sendbirdUcV4.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Session Token Sendbird Response V4", sendbirdUcV4.Result).JSON(c)
}

// CheckAccessToken function for requesting email for forgot password
func (h *MemberHandlerV4) CheckAccessToken(c echo.Context) error {
	var params serviceModel.SendbirdRequestV4
	UserID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := middleware.ExtractClaimsFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokenContext := c.Get("token")
	client := c.QueryParam("client")
	parsingToken := tokenContext.(*jwt.Token)

	params.UserID = UserID
	params.ExpiresAt = claims.StandardClaims.ExpiresAt
	params.Token = parsingToken.Raw
	params.Client = client

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))

	sendbirdUc := h.MemberUseCase.CheckSendbirdTokenV4(newCtx, &params)

	if sendbirdUc.Error != nil {
		return shared.NewHTTPResponse(sendbirdUc.HTTPStatus, sendbirdUc.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Check Token Sendbird ", sendbirdUc.Result).JSON(c)
}
