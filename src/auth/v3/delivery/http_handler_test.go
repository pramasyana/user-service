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

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	root                  = "/api/v2/auth"
	noAuth                = `Basic`
	tokenAdmin            = `Basic c3R1cmdlb246Ymhpbm5la2E=`
	tokenAdminJwt         = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	testCasePositive1     = "Testcase #1: Positive"
	testCaseNegative2     = "Testcase #2: Negative, Not authorized"
	testCaseNegative3     = "Testcase #3: Negative"
	testCaseNegative4     = "Testcase #4: Negative"
	testCaseNegative5     = "Testcase #5: Negative"
	getClientApp          = "GetClientApp"
	failedResponse        = "failed"
	authorization         = "Authorization"
	unAuthorized          = "Unauthorized"
	failedAuth            = "basic auth is invalid"
	googleAuthRedirectURL = "http://localhost:8081/api/v2/auth/oauth2callback"
	jsonSchemaMerchantDir = "../../../../schema/"
)

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestAuthV3Mount(t *testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.AuthUseCase), googleAuthRedirectURL)
	handler.MountRoute(e.Group("/api/v3/auth"))
	assert.NotNil(t, handler)
}

func TestGetAccessTokenV3(t *testing.T) {
	tests := []struct {
		name                                                                       string
		token                                                                      string
		expectUCGenerateToken, expectUCValidateBasicAuth, expectUCSendEmailWelcome usecase.ResultUseCase
		expectError, jsonBody                                                      bool
		expectStatusCode                                                           int
		param1, grantType                                                          string
		payload                                                                    interface{}
	}{
		{
			name:                      testCasePositive1,
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusOK,
			grantType:                 model.AuthTypeRefreshToken,
		},
		{
			name:                      testCaseNegative2,
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(failedAuth)},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusBadRequest,
		},
		{
			name:                      testCaseNegative3,
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(model.ErrorAccountInActiveBahasa)},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenAdmin,
			expectUCGenerateToken: usecase.ResultUseCase{Result: model.RequestToken{
				NewMember: true,
			}},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectUCSendEmailWelcome: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf("pq: error"),
			},
			expectStatusCode: http.StatusOK,
			param1:           "sturgeon",
			grantType:        model.AuthTypeAzure,
		},
		{
			name:                      testCaseNegative5,
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.MFAResponse{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusForbidden,
			grantType:                 model.AuthTypeVerifyMFA,
		},
		{
			name:                      "Testcase #6: Negative",
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.AuthAzureToken{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusInternalServerError,
			grantType:                 model.AuthTypeVerifyMFA,
		},
		{
			name:                      "Testcase #7: Negative",
			token:                     tokenAdmin,
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusBadRequest,
			grantType:                 model.AuthTypeRefreshToken,
			jsonBody:                  true,
			payload:                   tokenAdmin,
		},
		{
			name:                      "Testcase #8: Negative",
			token:                     tokenAdmin,
			expectUCValidateBasicAuth: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest},
			expectStatusCode:          http.StatusBadRequest,
			grantType:                 model.AuthTypeRefreshToken,
			jsonBody:                  true,
			payload:                   `{"grantType":"google",}`,
		},
		{
			name:                      "Testcase #9: Negative",
			token:                     tokenAdmin,
			expectUCValidateBasicAuth: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest},
			expectStatusCode:          http.StatusBadRequest,
			grantType:                 model.AuthTypeRefreshToken,
		},
		{
			name:                      "Testcase #10: Positive",
			token:                     tokenAdmin,
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.AuthV3TokenResponse{}},
			expectStatusCode:          http.StatusOK,
			grantType:                 model.AuthTypeRefreshToken,
		},
		{
			name:                      "Testcase #11: Negative",
			token:                     tokenAdmin,
			expectUCValidateBasicAuth: usecase.ResultUseCase{Result: true},
			expectUCGenerateToken:     usecase.ResultUseCase{Result: model.RefreshToken{}},
			expectStatusCode:          http.StatusInternalServerError,
			grantType:                 model.AuthTypeRefreshToken,
		},
		{
			name:             "Testcase #12: Negative",
			token:            noAuth,
			expectStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("ValidateBasicAuth", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUCValidateBasicAuth))
			mockAuthUseCase.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUCGenerateToken))
			mockAuthUseCase.On("SendEmailWelcomeMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUCSendEmailWelcome))

			e := echo.New()
			var body io.Reader
			valueContent := echo.MIMEApplicationForm
			if tt.jsonBody {
				bodyData, _ := json.Marshal(tt.payload)
				body = strings.NewReader(string(bodyData))
				valueContent = echo.MIMEApplicationJSON
			} else {
				data := url.Values{}
				data.Set("email", "abc@email.com")
				data.Set("requestFrom", tt.param1)
				data.Set("grantType", tt.grantType)
				body = bytes.NewBufferString(data.Encode())
			}

			req := httptest.NewRequest(echo.POST, root, body)
			req.Header.Set(echo.HeaderContentType, valueContent)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			// parse client id and secret
			err := handler.GetAccessToken(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}

func TestAuthHandlerV3_CheckEmailV3(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)

	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
		payload           interface{}
	}{
		{
			name:              testCasePositive1,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.CheckEmail{}},
			expectStatusCode:  http.StatusOK,
			payload:           model.CheckEmailPayload{},
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(unAuthorized)},
			expectStatusCode:  http.StatusUnauthorized,
			payload:           model.CheckEmailPayload{},
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("Bad Request")},
			expectStatusCode:  http.StatusBadRequest,
			payload:           model.AuthFacebookToken{},
		},
		{
			name:              testCaseNegative4,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.CheckEmail{}},
			expectStatusCode:  http.StatusBadRequest,
			payload:           tokenAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("CheckEmailV3", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.CheckEmailV3(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}
