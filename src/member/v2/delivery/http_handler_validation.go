package delivery

import (
	"errors"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// ValidateEmailDomain function for validating email domain
func (h *HTTPMemberHandler) ValidateEmailDomain(c echo.Context) error {
	email := c.FormValue("email")

	validateResult := <-h.MemberUseCase.ValidateEmailDomain(c.Request().Context(), email)
	if validateResult.Error != nil {
		return shared.NewHTTPResponse(validateResult.HTTPStatus, validateResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Validate Email Response", make(helper.EmptySlice, 0)).JSON(c)
}

// ValidateToken function for validating token of forgot password
func (h *HTTPMemberHandler) ValidateToken(c echo.Context) error {
	token := c.FormValue("token")

	validateResult := <-h.MemberUseCase.ValidateToken(c.Request().Context(), token)
	if validateResult.Error != nil {
		return shared.NewHTTPResponse(validateResult.HTTPStatus, validateResult.Error.Error()).JSON(c)
	}

	res, ok := validateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Validate Token Response", res).JSON(c)
}
