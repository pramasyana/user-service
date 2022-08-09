package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
)

const (
	root              = "/api/v1/auth"
	noAuth            = `Basic`
	tokenAdmin        = `Basic c3R1cmdlb246Ymhpbm5la2E=`
	testCasePositive1 = "Testcase #1: Positive"
	testCaseNegative2 = "Testcase #2: Negative, Not authorized"
	testCaseNegative3 = "Testcase #3: Negative"
	testCaseNegative4 = "Testcase #4: Negative"
	testCaseNegative5 = "Testcase #5: Negative"
	failedResponse    = "failed"
	authorization     = "Authorization"
	failedAuth        = "basic auth is invalid"
	getClientApp      = "GetClientApp"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestAuthMount(t *testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.AuthUseCase))
	handler.Mount(e.Group("/"))
	handler.MountClientApp(e.Group("/"))
	assert.Equal(t, handler, handler)
}

func TestCreateClientApp(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
	}{
		{
			name:              testCasePositive1,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.NewClientApp("test")},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              testCaseNegative2,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.AuthTypeAnonymous},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(failedResponse)},
			expectStatusCode:  http.StatusOK,
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("CreateClientApp", mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockAuthUseCase)

			err = handler.CreateClientApp(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}

func TestGetAccessToken(t *testing.T) {
	tests := []struct {
		name                                                      string
		token                                                     string
		expectUseCaseData, expectUseCaseData2, expectUseCaseData3 usecase.ResultUseCase
		expectError, jsonBody                                     bool
		expectStatusCode                                          int
		param1, grantType                                         string
		payload                                                   interface{}
	}{
		{
			name:               testCasePositive1,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusOK,
			grantType:          model.AuthTypeRefreshToken,
		},
		{
			name:               testCaseNegative2,
			token:              noAuth,
			expectUseCaseData:  usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(failedAuth)},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusUnauthorized,
			expectError:        true,
		},
		{
			name:               testCaseNegative3,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(model.ErrorAccountInActiveBahasa)},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			grantType:          model.AuthTypeAzure,
			expectStatusCode:   http.StatusBadRequest,
			expectError:        true,
		},
		{
			name:               testCaseNegative4,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUseCaseData2: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusInternalServerError},
			expectStatusCode:   http.StatusOK,
			grantType:          model.AuthTypeRefreshToken,
		},
		{
			name:               testCaseNegative5,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: nil},
			expectUseCaseData2: usecase.ResultUseCase{Result: true, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusInternalServerError},
			expectStatusCode:   http.StatusOK,
			grantType:          model.AuthTypeRefreshToken,
		},
		{
			name:               "Testcase #6: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.AuthAzureToken{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusInternalServerError,
			grantType:          model.AuthTypeVerifyMFA,
			expectError:        true,
		},
		{
			name:               "Testcase #7: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.MFAResponse{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusForbidden,
			grantType:          model.AuthTypeVerifyMFA,
			expectError:        false,
		},
		{
			name:               "Testcase #8: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusOK,
			grantType:          model.AuthTypeRefreshToken,
			jsonBody:           true,
			payload:            tokenAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))
			mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData2))
			mockAuthUseCase.On("SendEmailWelcomeMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData3))

			e := echo.New()
			var body io.Reader
			valueContent := echo.MIMEApplicationForm
			if tt.jsonBody {
				bodyData, _ := json.Marshal(tt.payload)
				body = strings.NewReader(string(bodyData))
				valueContent = echo.MIMEApplicationJSON
			} else {
				data := url.Values{}
				data.Set("requestFrom", tt.param1)
				data.Set("grantType", tt.grantType)
				body = bytes.NewBufferString(data.Encode())
			}

			req, err := http.NewRequest(echo.POST, root, body)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, valueContent)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase)

			err = handler.GetAccessToken(c)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectStatusCode, rec.Code)
			}
		})
	}
}
