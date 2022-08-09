package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	successText = "success"
)

var testData = []struct {
	name    string
	code    int
	message string
	param   interface{}
	success bool
}{
	{
		"#1",
		200,
		successText,
		map[string]interface{}{"data": "something"},
		true,
	},
	{
		"#2",
		504,
		"error",
		MultiError{},
		false,
	},
	{
		"#3",
		200,
		successText,
		Meta{},
		true,
	},
	{
		"#4",
		200,
		successText,
		&Meta{},
		true,
	},
}

func TestHTTPResponse(t *testing.T) {
	for _, tt := range testData {
		app := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := app.NewContext(req, rec)
		res := NewHTTPResponse(tt.code, tt.message, tt.param)
		res.SetSuccess(tt.success)
		err := res.JSON(c)
		assert.NoError(t, err)
	}
}
