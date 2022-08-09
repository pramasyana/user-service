package delivery

import (
	"errors"
	"math"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/applications/v1/model"
	"github.com/Bhinneka/user-service/src/applications/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// HTTPApplicationsHandler structure
type HTTPApplicationsHandler struct {
	ApplicationsUseCase usecase.ApplicationsUseCase
}

// NewHTTPHandler function for initialise *HTTPApplicationsHandler
func NewHTTPHandler(applicationsUseCase usecase.ApplicationsUseCase) *HTTPApplicationsHandler {
	return &HTTPApplicationsHandler{ApplicationsUseCase: applicationsUseCase}
}

// MountInfo function for mounting routes
func (h *HTTPApplicationsHandler) MountInfo(group *echo.Group) {
	group.POST("", h.AddApplication)
	group.PUT("/:id", h.UpdateApplication)
	group.DELETE("/:id", h.DeleteApplication)
	group.GET("", h.GetApplicationList)
}

// AddApplication function for add new application data
func (h *HTTPApplicationsHandler) AddApplication(c echo.Context) error {
	newApp := model.Application{}
	if err := c.Bind(&newApp); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	// Add app usecase process
	saveResult := <-h.ApplicationsUseCase.AddUpdateApplication(c.Request().Context(), newApp)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	newAppResult, ok := saveResult.Result.(model.Application)
	if !ok {
		err := errors.New("result is not proper response")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Success Application Added", newAppResult).JSON(c)

}

// UpdateApplication function for add new application data
func (h *HTTPApplicationsHandler) UpdateApplication(c echo.Context) error {
	app := model.Application{}
	app.ID = c.Param("id")

	if err := c.Bind(&app); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	// Add app usecase process
	updateResult := <-h.ApplicationsUseCase.AddUpdateApplication(c.Request().Context(), app)
	if updateResult.Error != nil {
		return shared.NewHTTPResponse(updateResult.HTTPStatus, updateResult.Error.Error()).JSON(c)
	}

	updatedApp, ok := updateResult.Result.(model.Application)
	if !ok {
		err := errors.New("result is not proper response")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success Application Updated", updatedApp).JSON(c)

}

// DeleteApplication function for removing application
func (h *HTTPApplicationsHandler) DeleteApplication(c echo.Context) error {
	id := c.Param("id")

	// Delete application usecase process
	result := <-h.ApplicationsUseCase.DeleteApplication(c.Request().Context(), id)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success Application Deleted").JSON(c)
}

// GetApplicationList function for getting list of application
func (h *HTTPApplicationsHandler) GetApplicationList(c echo.Context) error {
	ctx := "ApplicationsPresenter-GetApplicationList"

	params := model.ParametersApplication{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
	}

	applicationsResult := <-h.ApplicationsUseCase.GetListApplication(c.Request().Context(), &params)
	if applicationsResult.Error != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, "get_list_application", applicationsResult.Error, params)
		return shared.NewHTTPResponse(applicationsResult.HTTPStatus, applicationsResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	applications, ok := applicationsResult.Result.(model.ListApplication)
	if !ok {
		err := errors.New("result is not list of application")
		helper.SendErrorLog(c.Request().Context(), ctx, "parse_application", applicationsResult.Error, params)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(applications.TotalData) / float64(params.Limit))

	if len(applications.Application) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, "Success Get All Application", make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: applications.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, "Success Get All Application", applications.Application, meta).JSON(c)
}
