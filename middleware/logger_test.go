package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderXRealIP, "127.0.0.1")
	req.Header.Set(echo.HeaderXRequestID, "My-IP")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := echo.HandlerFunc(func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.String(http.StatusOK, "user-service bhinneka.com"))
	})

	errorHandler := echo.HandlerFunc(func(c echo.Context) error {
		err := errors.New("error")
		c.Error(err)
		return err
	})

	mw := Logger(handler)
	err := mw(c)
	assert.NoError(t, err)

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderXRequestID, "My-IP")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	mw = Logger(handler)
	err = mw(c)
	assert.NoError(t, err)

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	mw = Logger(handler)
	err = mw(c)
	assert.NoError(t, err)

	mw = Logger(errorHandler)
	err = mw(c)
	assert.NoError(t, err)

}
