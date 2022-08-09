package delivery

import (
	"context"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/log/v1/usecase"
	"github.com/Bhinneka/user-service/src/service"
	"github.com/Bhinneka/user-service/src/shared"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/labstack/echo"
)

const (
	errorBadURL = "invalid query parameters"
)

// HTTPClientHandler DI
type HTTPClientHandler struct {
	ActivityService service.ActivityServices
	UseCase         usecase.LogUsecase
}

// NewHTTPHandler return handler to process audit trail
func NewHTTPHandler(as service.ActivityServices, uc usecase.LogUsecase) *HTTPClientHandler {
	return &HTTPClientHandler{
		ActivityService: as,
		UseCase:         uc,
	}
}

// Mount load all router
func (h *HTTPClientHandler) Mount(group *echo.Group) {
	// URL => /v1/log/*
	group.GET("", h.GetAll)
	group.GET("/:id", h.GetByID)

}

// GetAll return all audit trail from activity service
func (h *HTTPClientHandler) GetAll(ctx echo.Context) error {
	param := new(sharedModel.Parameters)
	if err := ctx.Bind(param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, errorBadURL).JSON(ctx)
	}
	serviceCtx := context.WithValue(ctx.Request().Context(), helper.TextAuthorization, ctx.Request().Header.Get(helper.TextAuthorization))

	ucResult := <-h.UseCase.GetAll(serviceCtx, param)
	if ucResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, ucResult.Error.Error()).JSON(ctx)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Log Response", ucResult.Result, ucResult.Meta).JSON(ctx)
}

// GetByID return single audit trail from activity service
func (h *HTTPClientHandler) GetByID(ctx echo.Context) error {
	logID := ctx.Param("id")
	serviceCtx := context.WithValue(ctx.Request().Context(), helper.TextAuthorization, ctx.Request().Header.Get(helper.TextAuthorization))

	ucResult := <-h.UseCase.GetByID(serviceCtx, logID)
	if ucResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, ucResult.Error.Error()).JSON(ctx)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Single Log Response", ucResult.Result).JSON(ctx)
}
