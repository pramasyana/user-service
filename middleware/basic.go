package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Config for basic auth strategy
type Config struct {
	username string
	password string
}

// NewConfig function new config
func NewConfig(username, password string) *Config {
	return &Config{username: username, password: password}
}

// BasicAuth function basic auth
func BasicAuth(config *Config) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, _ echo.Context) (bool, error) {

		validate := func(username, password string) bool {
			if config.username == username && config.password == password {
				return true
			}
			return false
		}

		if validate(username, password) {
			return true, nil
		}
		return false, nil
	})
}
