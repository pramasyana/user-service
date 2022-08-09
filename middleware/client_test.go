package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var testDataVC = []struct {
	name         string
	clientID     string
	clientSecret string
	expectError  bool
}{
	{
		name:        "test vc #1",
		clientID:    "ss",
		expectError: true,
	},
	{
		name:         "test vc #2",
		clientID:     "ss",
		clientSecret: "cs",
		expectError:  false,
	},
}

func TestValidateClient(t *testing.T) {
	for _, tc := range testDataVC {
		headers := map[string]string{
			"X-Client-ID":          tc.clientID,
			echo.HeaderContentType: echo.MIMEApplicationJSON,
			"X-Client-Secret":      tc.clientSecret,
		}
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/v1/client/login", nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		he := echo.NewHTTPError(http.StatusUnauthorized)
		next := ctx.JSON(http.StatusOK, ctx.String(http.StatusOK, "bhinneka.com"))

		handler := func(echo.Context) error {
			if tc.expectError {
				return he
			}
			return next
		}
		vc := ValidateClient()(handler)
		err := vc(ctx)
		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

	}

}
