package delivery

import (
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// GetNarwhalMFASettings function for getting status mfa settings
func (h *HTTPMemberHandler) GetNarwhalMFASettings(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.GetNarwhalMFASettings(c.Request().Context(), memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Get Narwhal MFA Setting", memberResult.Result).JSON(c)
}

// GenerateNarwhalMFASettings function for getting status mfa settings
func (h *HTTPMemberHandler) GenerateNarwhalMFASettings(c echo.Context) error {
	narwhalUserID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.GenerateMFASettings(c.Request().Context(), narwhalUserID, helper.TextNarwhal)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success generate shared key", memberResult.Result).JSON(c)
}

// ActivateNarwhalMFASettings function for activate mfa settings
func (h *HTTPMemberHandler) ActivateNarwhalMFASettings(c echo.Context) error {
	adminUserID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	reqActivation := model.MFAActivateSettings{}
	reqActivation.Otp = c.FormValue("otp")
	reqActivation.SharedKeyText = c.FormValue("sharedKeyText")
	reqActivation.MemberID = adminUserID
	reqActivation.RequestFrom = helper.TextNarwhal

	activationResult := <-h.MemberUseCase.ActivateMFASettings(c.Request().Context(), reqActivation)
	if activationResult.Error != nil {
		return shared.NewHTTPResponse(activationResult.HTTPStatus, activationResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success activate MFA Narwhal").JSON(c)
}

// DisabledNarwhalMFASetting function for getting status mfa settings
func (h *HTTPMemberHandler) DisabledNarwhalMFASetting(c echo.Context) error {
	narwhalUserID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	adminResult := <-h.MemberUseCase.DisabledMFASetting(c.Request().Context(), narwhalUserID, helper.TextNarwhal)
	if adminResult.Error != nil {
		return shared.NewHTTPResponse(adminResult.HTTPStatus, adminResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success disable MFA Narwhal").JSON(c)
}
