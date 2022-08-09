package delivery

import (
	"context"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo"
)

// MountSendbird function for mounting routes
func (h *MemberHandlerV3) MountSendbird(group *echo.Group) {
	group.GET("/token", h.GetAccessToken)
	group.GET("/check-token", h.CheckAccessToken)
}

// GetAccessToken function for requesting email for forgot password
func (h *MemberHandlerV3) GetAccessToken(c echo.Context) error {
	var paramsV3 serviceModel.SendbirdRequest
	UserIDV3, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := middleware.ExtractClaimsFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokenContext := c.Get("token")
	parsingToken := tokenContext.(*jwt.Token)

	paramsV3.UserID = UserIDV3
	paramsV3.ExpiresAt = claims.StandardClaims.ExpiresAt
	paramsV3.Token = parsingToken.Raw

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))

	sendbirdUc := h.MemberUseCase.GetSendbirdToken(newCtx, &paramsV3)

	if sendbirdUc.Error != nil {
		return shared.NewHTTPResponse(sendbirdUc.HTTPStatus, sendbirdUc.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Session Token Sendbird Response", sendbirdUc.Result).JSON(c)
}

// CheckAccessToken function for requesting email for forgot password
func (h *MemberHandlerV3) CheckAccessToken(c echo.Context) error {
	var params serviceModel.SendbirdRequest
	UserID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := middleware.ExtractClaimsFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokenContext := c.Get("token")
	parsingToken := tokenContext.(*jwt.Token)

	params.UserID = UserID
	params.ExpiresAt = claims.StandardClaims.ExpiresAt
	params.Token = parsingToken.Raw

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))

	sendbirdUc := h.MemberUseCase.CheckSendbirdToken(newCtx, &params)

	if sendbirdUc.Error != nil {
		return shared.NewHTTPResponse(sendbirdUc.HTTPStatus, sendbirdUc.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Check Token Sendbird ", sendbirdUc.Result).JSON(c)
}
