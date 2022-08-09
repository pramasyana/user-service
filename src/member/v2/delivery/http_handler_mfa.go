package delivery

import (
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// GetMFASettings function for getting status mfa settings
func (h *HTTPMemberHandler) GetMFASettings(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.GetMFASettings(c.Request().Context(), memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Get MFA Setting", memberResult.Result).JSON(c)
}

// GenerateMFASettings function for getting status mfa settings
func (h *HTTPMemberHandler) GenerateMFASettings(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.GenerateMFASettings(c.Request().Context(), memberID, helper.TextAccount)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success generate shared key", memberResult.Result).JSON(c)
}

// ActivateMFASettings function for activate mfa settings
func (h *HTTPMemberHandler) ActivateMFASettings(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	activateData := model.MFAActivateSettings{}
	activateData.Otp = c.FormValue("otp")
	activateData.SharedKeyText = c.FormValue("sharedKeyText")
	activateData.MemberID = memberID
	activateData.RequestFrom = helper.TextAccount

	activateResult := <-h.MemberUseCase.ActivateMFASettings(c.Request().Context(), activateData)
	if activateResult.Error != nil {
		return shared.NewHTTPResponse(activateResult.HTTPStatus, activateResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, model.SuccessMFAActivation).JSON(c)
}

// DisabledMFASetting function for getting status mfa settings
func (h *HTTPMemberHandler) DisabledMFASetting(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.DisabledMFASetting(c.Request().Context(), memberID, helper.TextAccount)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success disable MFA").JSON(c)
}
