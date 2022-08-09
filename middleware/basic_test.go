package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func getAuth() string {
	return `Basic Ymhpbm5la2E6ZGExYzI1ZDgtMzdjOC00MWIxLWFmZTItNDJkZDQ4MjViZmVh`
}

var testBasicAuth = []struct {
	name      string
	auth      string
	wantError bool
}{
	{
		name: "Test With Valid Auth #1",
		auth: getAuth(),
	},
	{
		name:      "Test With Valid Auth #2",
		auth:      "invalidAuth",
		wantError: true,
	},
	{
		name:      "Test With Invalid Basic Auth",
		auth:      "Basic Ymhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNTo2MjY4NjktNmU2ZTY1LTZiNjEyMC02ZDY1NmUtNzQ2MTcyLTY5MjA2NC02OTZkNjU2ZS03MzY5MDA=",
		wantError: true,
	},
}

func TestBasicAuth(t *testing.T) {
	config := NewConfig("bhinneka", "da1c25d8-37c8-41b1-afe2-42dd4825bfea")

	for _, tc := range testBasicAuth {
		t.Run(tc.auth, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.Header.Set(echo.HeaderAuthorization, tc.auth)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := echo.HandlerFunc(func(c echo.Context) error {
				return c.JSON(http.StatusOK, c.String(http.StatusOK, "bhinneka.com"))
			})

			mw := BasicAuth(config)(handler)
			err := mw(c)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
