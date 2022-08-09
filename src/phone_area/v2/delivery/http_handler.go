package delivery

import (
	"errors"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
	"github.com/Bhinneka/user-service/src/phone_area/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// HTTPPhoneAreaHandlerV2 model
type HTTPPhoneAreaHandlerV2 struct {
	PhoneAreaUseCase usecase.PhoneAreaUseCase
}

// NewHTTPHandler v2 function for initialise *HTTPAuthHandler
func NewHTTPHandler(phoneAreaUseCase usecase.PhoneAreaUseCase) *HTTPPhoneAreaHandlerV2 {
	return &HTTPPhoneAreaHandlerV2{PhoneAreaUseCase: phoneAreaUseCase}
}

// MountPhoneArea v2 function for mounting routes
func (h *HTTPPhoneAreaHandlerV2) MountPhoneArea(group *echo.Group) {
	group.GET("", h.GetAllPhoneArea)
}

// GetAllPhoneArea v2 function for getting list of phone area
func (h *HTTPPhoneAreaHandlerV2) GetAllPhoneArea(c echo.Context) error {
	ctx := "PhoneAreaPresenter-GetAllPhoneArea"

	// get list phone area
	phoneAreaResult := <-h.PhoneAreaUseCase.GetAllPhoneArea(c.Request().Context())
	if phoneAreaResult.Error != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, phoneAreaResult.Error, nil)
		response := shared.NewHTTPResponse(phoneAreaResult.HTTPStatus, phoneAreaResult.Error.Error())
		return response.JSON(c)
	}

	phoneArea, ok := phoneAreaResult.Result.([]model.PhoneArea)
	if !ok {
		err := errors.New("result is not list of phone area")
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, err, phoneAreaResult.Result)
		response := shared.NewHTTPResponse(http.StatusBadRequest, err.Error())
		return response.JSON(c)
	}

	// get total phone area
	totalPhoneAreaResult := <-h.PhoneAreaUseCase.GetTotalPhoneArea(c.Request().Context())
	if totalPhoneAreaResult.Error != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, totalPhoneAreaResult.Error, nil)
		response := shared.NewHTTPResponse(totalPhoneAreaResult.HTTPStatus, totalPhoneAreaResult.Error.Error())
		return response.JSON(c)
	}

	total, ok := totalPhoneAreaResult.Result.(model.TotalPhoneArea)
	if !ok {
		err := errors.New("result is not total phone area")
		helper.SendErrorLog(c.Request().Context(), ctx, helper.TextPhoneArea, err, totalPhoneAreaResult.Result)
		response := shared.NewHTTPResponse(http.StatusBadRequest, err.Error())
		return response.JSON(c)
	}

	var meta shared.Meta
	meta.Page = 1
	meta.Limit = total.TotalData
	meta.TotalPages = 1
	meta.TotalRecords = total.TotalData

	response := shared.NewHTTPResponse(http.StatusOK, model.MessageSuccess, phoneArea, meta)
	return response.JSON(c)
}
