package delivery

import (
	"errors"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
	"github.com/Bhinneka/user-service/src/phone_area/v1/usecase"
	"github.com/labstack/echo"
)

// HTTPPhoneAreaHandler model
type HTTPPhoneAreaHandler struct {
	PhoneAreaUseCase usecase.PhoneAreaUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(phoneAreaUseCase usecase.PhoneAreaUseCase) *HTTPPhoneAreaHandler {
	return &HTTPPhoneAreaHandler{PhoneAreaUseCase: phoneAreaUseCase}
}

// MountPhoneArea function for mounting routes
func (h *HTTPPhoneAreaHandler) MountPhoneArea(group *echo.Group) {
	group.GET("", h.GetAllPhoneArea)
}

// GetAllPhoneArea function for getting list of phone area
func (h *HTTPPhoneAreaHandler) GetAllPhoneArea(c echo.Context) error {
	ctx := "PhoneAreaPresenter-GetAllPhoneArea"

	phoneAreaResult := <-h.PhoneAreaUseCase.GetAllPhoneArea(c.Request().Context())
	if phoneAreaResult.Error != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, phoneAreaResult.Error, nil)
		return echo.NewHTTPError(phoneAreaResult.HTTPStatus, phoneAreaResult.Error.Error())
	}

	phoneArea, ok := phoneAreaResult.Result.([]model.PhoneArea)
	if !ok {
		err := errors.New("result is not list of phone area")
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, err, nil)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var response model.PhoneAreaResponse

	response.Status = "success"
	response.Data = phoneArea
	response.Code = http.StatusOK
	response.Message = model.MessageSuccess

	return c.JSON(http.StatusOK, response)
}
