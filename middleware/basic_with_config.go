package middleware

import (
	"net/http"

	clientQuery "github.com/Bhinneka/user-service/src/client/v1/query"
	"github.com/labstack/echo"
)

// BasicAuth function basic auth
func BasicAuthWithConfig(query clientQuery.ClientQuery) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientID, clientSecret, ok := c.Request().BasicAuth()
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid header authorization")
			}
			valid, err := query.Validate(c.Request().Context(), clientID, clientSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			if !valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid basic auth")
			}

			return next(c)
		}
	}
}
