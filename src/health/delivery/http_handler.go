package delivery

import (
	"net/http"

	"github.com/Bhinneka/user-service/src/health/model"
	"github.com/Bhinneka/user-service/src/health/usecase"
	"github.com/labstack/echo"
)

// HTTPHealthHandler model
type HTTPHealthHandler struct {
	healthUseCase usecase.HealthUseCase
}

// NewHTTPHandler for initializing http health model
func NewHTTPHandler(healthUseCase usecase.HealthUseCase) *HTTPHealthHandler {
	return &HTTPHealthHandler{healthUseCase: healthUseCase}
}

// Mount function for mounting routes
func (h *HTTPHealthHandler) Mount(group *echo.Group) {
	group.GET("", h.Ping)
}

// Ping function for checking service
func (h *HTTPHealthHandler) Ping(c echo.Context) error {

	pingResult := <-h.healthUseCase.Ping()

	if pingResult.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot ping server")
	}

	ping, ok := pingResult.Result.(*model.Health)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "result is not ping")
	}

	return c.JSON(http.StatusOK, ping)

}
