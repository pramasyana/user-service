package delivery

import (
	"errors"
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
	group.GET("", h.GetApplicationsList)
}

// GetApplicationsList function for get list application data
func (h *HTTPApplicationsHandler) GetApplicationsList(c echo.Context) error {
	ctx := "ApplicationsPresenter-GetListApplications"

	resultList := <-h.ApplicationsUseCase.GetApplicationsList()

	if resultList.Error != nil {
		helper.SendErrorLog(c.Request().Context(), ctx, "get_application_list", resultList.Error, nil)
		res := shared.NewHTTPResponse(resultList.HTTPStatus, resultList.Error.Error())
		return res.JSON(c)
	}

	response, ok := resultList.Result.(model.ApplicationList)

	if !ok {
		err := errors.New("result is not list of application")
		helper.SendErrorLog(c.Request().Context(), ctx, "parse_application_list", resultList.Error, resultList.Result)
		res := shared.NewHTTPResponse(http.StatusBadRequest, err.Error())
		return res.JSON(c)
	}

	res := shared.NewHTTPResponse(http.StatusOK, "Success Get All Application", response.Data)
	return res.JSON(c)
}
