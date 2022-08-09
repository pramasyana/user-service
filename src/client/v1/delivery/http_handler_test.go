package delivery

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/middleware"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authUsecase "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	authMock "github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
	clientUCMock "github.com/Bhinneka/user-service/src/client/v1/usecase/mocks"
	actSvcMock "github.com/Bhinneka/user-service/src/service/mocks"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	getClientApp        = "GetClientApp"
	generateToken       = "GenerateToken"
	verifyToken         = "VerifyTokenMember"
	logout              = "Logout"
	defaultPath         = "/v1/client/login"
	verifyPath          = "/v1/client/verify?token=sometoken"
	headerClientID      = "X-Client-ID"
	headerClientSecret  = "X-Client-Secret"
	defaultClientID     = "someClientID"
	defaultClientSecret = "someSecret"
)

var (
	errorDefault   = errors.New("someerror")
	defaultPayload = `{"payload":{"userName":"Pian","realName":"Pian","role":"KK","lpseId":"123","isLatihan":false,"time":"01-05-2020 12:00:00","email":"pian@yopmail.com"}}`
	badPayload     = `{"payload":{"userName":"Pian","realName":"Pian","role":"KK","lpseId":"123","isLatihan":false,"time":"01-05-2020 12:00:00","email":"pian@yopmail.com",}}`
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateUsecaseResult(data authUsecase.ResultUseCase) <-chan authUsecase.ResultUseCase {
	output := make(chan authUsecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func createSharedResult(data sharedModel.ResultUseCase) <-chan sharedModel.ResultUseCase {
	output := make(chan sharedModel.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestHTTPClientHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(authMock.AuthUseCase), new(actSvcMock.ActivityServices), new(clientUCMock.ClientUsecase))
	handler.Mount(e.Group("/v1/client"), middleware.BasicAuth(&middleware.Config{}))
}

var testDataLogin = []struct {
	name                       string
	clientID                   string
	clientSecret               string
	expectUseCaseGetClientApp  authUsecase.ResultUseCase
	expectUseCaseGenerateToken authUsecase.ResultUseCase
	expectError                bool
	expectStatusCode           int
	payload                    string
}{
	{
		name:             "testLogin #1",
		clientID:         defaultClientID,
		expectStatusCode: http.StatusBadRequest,
	},
	{
		name:                      "testLogin #2",
		clientID:                  defaultClientID,
		clientSecret:              defaultClientSecret,
		expectUseCaseGetClientApp: authUsecase.ResultUseCase{Error: errorDefault},
		expectStatusCode:          http.StatusUnauthorized,
	},
	{
		name:                       "testLogin #3",
		clientID:                   defaultClientID,
		clientSecret:               defaultClientSecret,
		expectUseCaseGetClientApp:  authUsecase.ResultUseCase{Result: true},
		expectUseCaseGenerateToken: authUsecase.ResultUseCase{Result: authModel.RequestToken{}},
		expectStatusCode:           http.StatusOK,
		payload:                    defaultPayload,
	},
	{
		name:                       "testLogin #4",
		clientID:                   defaultClientID,
		clientSecret:               defaultClientSecret,
		expectUseCaseGetClientApp:  authUsecase.ResultUseCase{Result: true},
		expectUseCaseGenerateToken: authUsecase.ResultUseCase{Result: authModel.RequestToken{}},
		expectStatusCode:           http.StatusBadRequest,
		payload:                    badPayload,
	},
	{
		name:                       "testLogin #5",
		clientID:                   defaultClientID,
		clientSecret:               defaultClientSecret,
		expectUseCaseGetClientApp:  authUsecase.ResultUseCase{Result: true},
		expectUseCaseGenerateToken: authUsecase.ResultUseCase{Error: errorDefault},
		expectStatusCode:           http.StatusUnauthorized,
		payload:                    defaultPayload,
	},
	{
		name:                       "testLogin #6",
		clientID:                   defaultClientID,
		clientSecret:               defaultClientSecret,
		expectUseCaseGetClientApp:  authUsecase.ResultUseCase{Result: true},
		expectUseCaseGenerateToken: authUsecase.ResultUseCase{Result: true},
		expectStatusCode:           http.StatusUnauthorized,
		payload:                    defaultPayload,
	},
}

func TestClientLogin(t *testing.T) {
	for _, tc := range testDataLogin {
		mockAuthUseCase := new(authMock.AuthUseCase)
		mockActivityService := new(actSvcMock.ActivityServices)
		mockClientUsecase := new(clientUCMock.ClientUsecase)
		mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.expectUseCaseGetClientApp))
		mockAuthUseCase.On(generateToken, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.expectUseCaseGenerateToken))
		mockActivityService.On("CreateLog", mock.Anything, mock.Anything).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, defaultPath, bytes.NewBufferString(tc.payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(headerClientID, tc.clientID)
		req.Header.Set(headerClientSecret, tc.clientSecret)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewHTTPHandler(mockAuthUseCase, mockActivityService, mockClientUsecase)
		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}

var testDataLogout = []struct {
	name                 string
	clientID             string
	clientSecret         string
	responseGetClientApp authUsecase.ResultUseCase
	responseLogout       sharedModel.ResultUseCase
	expectError          bool
	expectStatusCode     int
	payload              string
}{
	{
		name:                 "test Logout #1",
		clientID:             defaultClientID,
		clientSecret:         defaultClientSecret,
		responseGetClientApp: authUsecase.ResultUseCase{Error: errorDefault},
		expectStatusCode:     http.StatusUnauthorized,
	},
	{
		name:                 "test logout #2",
		clientID:             defaultClientID,
		clientSecret:         defaultClientSecret,
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		expectStatusCode:     http.StatusBadRequest,
		payload:              badPayload,
	},
	{
		name:                 "test logout #3",
		clientID:             defaultClientID,
		clientSecret:         defaultClientSecret,
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		responseLogout:       sharedModel.ResultUseCase{Error: errorDefault},
		expectStatusCode:     http.StatusBadRequest,
		payload:              defaultPayload,
	},
	{
		name:                 "test logout #4",
		clientID:             defaultClientID,
		clientSecret:         defaultClientSecret,
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		responseLogout:       sharedModel.ResultUseCase{Error: nil},
		expectStatusCode:     http.StatusOK,
		payload:              defaultPayload,
	},
}

func TestClientLogout(t *testing.T) {
	for _, tc := range testDataLogout {
		mockAuthUseCase := new(authMock.AuthUseCase)
		mockActivityService := new(actSvcMock.ActivityServices)
		mockClientUsecase := new(clientUCMock.ClientUsecase)

		mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.responseGetClientApp))
		mockClientUsecase.On(logout, mock.Anything, mock.Anything).Return(createSharedResult(tc.responseLogout))
		handler := NewHTTPHandler(mockAuthUseCase, mockActivityService, mockClientUsecase)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, defaultPath, bytes.NewBufferString(tc.payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(headerClientID, tc.clientID)
		req.Header.Set(headerClientSecret, tc.clientSecret)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Logout(c)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}

var testDataVerifyToken = []struct {
	name                 string
	clientID             string
	clientSecret         string
	responseGetClientApp authUsecase.ResultUseCase
	responseVerifyToken  authUsecase.ResultUseCase
	expectError          bool
	expectStatusCode     int
	token                string
}{
	{
		name:                 "test verify token #1",
		responseGetClientApp: authUsecase.ResultUseCase{Error: errorDefault},
		expectStatusCode:     http.StatusBadRequest,
	},
	{
		name:                 "Test Verify Token #2",
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		expectStatusCode:     http.StatusBadRequest,
		token:                "`someSoktnss`--kk",
	},
	{
		name:                 "test verify token #3",
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		responseVerifyToken:  authUsecase.ResultUseCase{Error: errorDefault},
		expectStatusCode:     http.StatusBadRequest,
	},
	{
		name:                 "test verify token #4",
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		responseVerifyToken:  authUsecase.ResultUseCase{Result: authModel.VerifyResponse{}},
		expectStatusCode:     http.StatusBadRequest,
	},
	{
		name:                 "test verify token #5",
		responseGetClientApp: authUsecase.ResultUseCase{Result: true},
		responseVerifyToken:  authUsecase.ResultUseCase{Result: &authModel.VerifyResponse{}},
		expectStatusCode:     http.StatusOK,
	},
}

func TestVerifyToken(t *testing.T) {
	for _, tc := range testDataVerifyToken {
		mockAuthUseCase := new(authMock.AuthUseCase)
		mockActivityService := new(actSvcMock.ActivityServices)
		mockClientUsecase := new(clientUCMock.ClientUsecase)
		mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.responseGetClientApp))
		mockAuthUseCase.On(verifyToken, mock.Anything, mock.Anything).Return(tc.responseVerifyToken)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, verifyPath, strings.NewReader(tc.token))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, getAuth())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(fmt.Sprintf("/v1/client/verify?token=%s", tc.token))
		handler := NewHTTPHandler(mockAuthUseCase, mockActivityService, mockClientUsecase)
		err := handler.VerifyToken(c)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}

func getAuth() string {
	return `Basic Ymhpbm5la2E6ZGExYzI1ZDgtMzdjOC00MWIxLWFmZTItNDJkZDQ4MjViZmVh`
}
