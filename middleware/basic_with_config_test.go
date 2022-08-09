package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	clientQueryMock "github.com/Bhinneka/user-service/src/client/v1/query/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testBAuth = []struct {
	name      string
	auth      string
	valid     bool
	errValid  error
	wantError bool
}{
	{
		name:  "Test #1",
		auth:  getAuth(),
		valid: true,
	},
	{
		name:      "Test #2",
		auth:      "",
		wantError: true,
	},
	{
		name:      "Test #3",
		auth:      getAuth(),
		errValid:  errors.New("some error"),
		wantError: true,
	},
	{
		name:      "Test #4",
		auth:      getAuth(),
		wantError: true,
	},
}

func TestBasicAuthWithConfig(t *testing.T) {
	for _, tc := range testBAuth {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.Header.Set(echo.HeaderAuthorization, tc.auth)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := echo.HandlerFunc(func(c echo.Context) error {
				return c.JSON(http.StatusOK, c.String(http.StatusOK, "bhinneka.com"))
			})

			mockQuery := new(clientQueryMock.ClientQuery)
			mockQuery.On("Validate", mock.Anything, mock.Anything, mock.Anything).Return(tc.valid, tc.errValid)

			mw := BasicAuthWithConfig(mockQuery)(handler)
			err := mw(c)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
