package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

// ValidateClient validate x-client-id and x-secret
func ValidateClient() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			headerReq := c.Request().Header

			clientID := headerReq.Get("X-Client-ID")
			clientSecret := headerReq.Get("X-Client-Secret")
			if clientID == "" || clientSecret == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
