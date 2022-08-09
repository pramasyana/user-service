package delivery

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
)

type VerifyMock struct{}

type Body map[string]interface{}

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
	badAuth               = "Basic IiI6IiI="
)

var (
	errDefault = errors.New("some error")
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
func generateRSA() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSA(), nil
	})
}

func TestAuthMount(t *testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.AuthUseCase), googleAuthRedirectURL)
	handler.Mount(e.Group("/"))
	handler.Mount(e.Group("/logout"))
	handler.MountAdmin(e.Group("/admin"))
}

func TestAuthClientMount(t *testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.AuthUseCase), googleAuthRedirectURL)
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
			expectUseCaseData: usecase.ResultUseCase{Result: nil},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(failedResponse)},
			expectStatusCode:  http.StatusBadRequest,
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

			c.Set(authorization, tt.token)
			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

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
			expectUseCaseData:  usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(failedAuth)},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusUnauthorized,
		},
		{
			name:               testCaseNegative3,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(model.ErrorAccountInActiveBahasa)},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.RequestToken{
				NewMember: true,
			}},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectUseCaseData3: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf("pq: error"),
			},
			expectStatusCode: http.StatusOK,
			param1:           "sturgeon",
			grantType:        model.AuthTypeAzure,
		},
		{
			name:               testCaseNegative5,
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.MFAResponse{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusForbidden,
			grantType:          model.AuthTypeVerifyMFA,
		},
		{
			name:               "Testcase #6: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.AuthAzureToken{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusInternalServerError,
			grantType:          model.AuthTypeVerifyMFA,
		},
		{
			name:               "Testcase #7: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Result: model.RequestToken{}},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectStatusCode:   http.StatusBadRequest,
			grantType:          model.AuthTypeRefreshToken,
			jsonBody:           true,
			payload:            tokenAdmin,
		},
		{
			name:               "Testcase #11: Negative",
			token:              tokenAdmin,
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectUseCaseData:  usecase.ResultUseCase{Result: model.RefreshToken{}},
			expectStatusCode:   http.StatusInternalServerError,
			grantType:          model.AuthTypeRefreshToken,
		},
		{
			name:               "Testcase #8: Negative",
			token:              tokenAdmin,
			expectUseCaseData2: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest},
			expectStatusCode:   http.StatusBadRequest,
			grantType:          model.AuthTypeRefreshToken,
			jsonBody:           true,
			payload:            `{"grantType":"google",}`,
		},
		{
			name:               "Testcase #9: Negative",
			token:              tokenAdmin,
			expectUseCaseData2: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest},
			expectStatusCode:   http.StatusBadRequest,
			grantType:          model.AuthTypePassword,
		},
		{
			name:               "Testcase #10: Negative",
			token:              tokenAdmin,
			expectUseCaseData:  usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			expectUseCaseData2: usecase.ResultUseCase{Result: true},
			expectUseCaseData3: usecase.ResultUseCase{},
			expectStatusCode:   http.StatusBadRequest,
			param1:             "sturgeon",
			grantType:          model.AuthTypeAzure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))
			mockAuthUseCase.On("SendEmailWelcomeMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData3))
			mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData2))

			e := echo.New()
			var body io.Reader
			valueContent := echo.MIMEApplicationForm
			if tt.jsonBody {
				bodyData, _ := json.Marshal(tt.payload)
				body = strings.NewReader(string(bodyData))
				valueContent = echo.MIMEApplicationJSON
			} else {
				data := url.Values{}
				data.Set("email", "abc@email.co")
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

			err := handler.GetAccessToken(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}

func TestGetAccessTokenFromUserID(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
	}{
		{
			name:              testCasePositive1,
			token:             tokenAdminJwt,
			expectUseCaseData: usecase.ResultUseCase{Result: model.RequestToken{}},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(failedAuth)},
			expectStatusCode:  http.StatusUnauthorized,
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdminJwt,
			expectUseCaseData: usecase.ResultUseCase{Result: &model.RequestToken{}},
			expectStatusCode:  http.StatusInternalServerError,
		},
		{
			name:              testCaseNegative5,
			token:             tokenAdminJwt,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(failedAuth)},
			expectStatusCode:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GenerateTokenFromUserID", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err := handler.GetAccessTokenFromUserID(c)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		name                 string
		token                string
		expectUseCaseData    usecase.ResultUseCase
		expectedGetClientApp usecase.ResultUseCase
		expectError          bool
		expectStatusCode     int
		payload              interface{}
	}{
		{
			name:                 testCasePositive1,
			token:                tokenAdmin,
			expectUseCaseData:    usecase.ResultUseCase{Result: VerifyMock{}},
			expectedGetClientApp: usecase.ResultUseCase{Result: true}, // dont return error
			expectStatusCode:     http.StatusOK,
			payload:              Body{helper.TextToken: tokenAdmin},
		},
		{
			name:                 testCaseNegative3,
			token:                tokenAdmin,
			expectUseCaseData:    usecase.ResultUseCase{Result: VerifyMock{}},
			expectedGetClientApp: usecase.ResultUseCase{Error: fmt.Errorf("some error")},
			expectStatusCode:     http.StatusUnauthorized,
			payload:              Body{helper.TextToken: tokenAdmin},
		},
		{
			name:                 testCasePositive1,
			token:                tokenAdmin,
			expectUseCaseData:    usecase.ResultUseCase{Error: fmt.Errorf("some error happened")},
			expectedGetClientApp: usecase.ResultUseCase{Result: true},
			expectStatusCode:     http.StatusUnauthorized,
			payload:              Body{helper.TextToken: tokenAdmin},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GetClientApp", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectedGetClientApp))
			mockAuthUseCase.On("VerifyTokenMember", mock.Anything, tokenAdmin).Return(tt.expectUseCaseData)

			e := echo.New()
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.VerifyToken(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}

func TestVerifyTokenMember(t *testing.T) {
	tests := []struct {
		name                 string
		token                string
		expectUseCaseData    usecase.ResultUseCase
		expectedGetClientApp usecase.ResultUseCase
		expectError          bool
		expectStatusCode     int
		payload              interface{}
	}{
		{
			name:              testCasePositive1,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: VerifyMock{}},
			expectStatusCode:  http.StatusOK,
			payload:           Body{helper.TextToken: tokenAdmin},
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(unAuthorized)},
			expectStatusCode:  http.StatusUnauthorized,
			payload:           Body{helper.TextToken: tokenAdmin},
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: VerifyMock{}},
			expectStatusCode:  http.StatusBadRequest,
			payload:           tokenAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GetClientApp", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectedGetClientApp))
			mockAuthUseCase.On("VerifyTokenMember", mock.Anything, tokenAdmin).Return(tt.expectUseCaseData)

			e := echo.New()
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.VerifyTokenMember(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name                      string
		token                     string
		expectUseCaseData         usecase.ResultUseCase
		expectUseCaseGetClientApp usecase.ResultUseCase
		expectError               bool
		expectStatusCode          int
		payload                   interface{}
	}{
		{
			name:                      testCasePositive1,
			token:                     tokenAdmin,
			expectUseCaseData:         usecase.ResultUseCase{Result: model.Logout{}},
			expectUseCaseGetClientApp: usecase.ResultUseCase{Result: true},
			expectStatusCode:          http.StatusOK,
			payload:                   Body{helper.TextToken: tokenAdmin},
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			expectStatusCode:  http.StatusUnauthorized,
			payload:           Body{helper.TextToken: tokenAdmin},
		},
		{
			name:                      testCaseNegative3,
			token:                     tokenAdmin,
			expectUseCaseData:         usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			expectStatusCode:          http.StatusBadRequest,
			payload:                   Body{helper.TextToken: tokenAdmin},
			expectUseCaseGetClientApp: usecase.ResultUseCase{Result: true},
		},
		{
			name:                      testCaseNegative4,
			token:                     tokenAdmin,
			expectUseCaseData:         usecase.ResultUseCase{Result: model.Logout{}},
			expectStatusCode:          http.StatusBadRequest,
			payload:                   tokenAdmin,
			expectUseCaseGetClientApp: usecase.ResultUseCase{Result: true},
		},
		{
			name:                      testCaseNegative5,
			token:                     tokenAdmin,
			expectUseCaseData:         usecase.ResultUseCase{Result: model.Logout{}},
			expectStatusCode:          http.StatusBadRequest,
			payload:                   tokenAdmin,
			expectUseCaseGetClientApp: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
		},
		{
			name:                      testCaseNegative5,
			token:                     tokenAdmin,
			expectUseCaseData:         usecase.ResultUseCase{Result: model.Logout{}},
			expectStatusCode:          http.StatusUnauthorized,
			payload:                   tokenAdmin,
			expectUseCaseGetClientApp: usecase.ResultUseCase{Result: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("Logout", mock.Anything, tokenAdmin).Return(generateUsecaseResult(tt.expectUseCaseData))
			mockAuthUseCase.On(getClientApp, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseGetClientApp))

			e := echo.New()
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.Logout(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}

func TestCheckEmail(t *testing.T) {
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
			payload:           model.CheckEmail{},
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(unAuthorized)},
			expectStatusCode:  http.StatusUnauthorized,
			payload:           model.CheckEmail{},
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
			mockAuthUseCase.On("CheckEmail", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.CheckEmail(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}

func TestVerifyCaptcha(t *testing.T) {
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
			expectUseCaseData: usecase.ResultUseCase{Result: model.GoogleCaptchaResponse{}},
			expectStatusCode:  http.StatusOK,
			payload:           model.GoogleCaptcha{},
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(unAuthorized)},
			expectStatusCode:  http.StatusUnauthorized,
			payload:           model.GoogleCaptcha{},
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.GoogleCaptchaResponse{}, Error: fmt.Errorf(failedResponse)},
			expectStatusCode:  http.StatusBadRequest,
			payload:           model.GoogleCaptcha{},
		},
		{
			name:              testCaseNegative4,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.GoogleCaptchaResponse{}},
			expectStatusCode:  http.StatusBadRequest,
			payload:           tokenAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("VerifyCaptcha", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err = handler.VerifyCaptcha(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}

func TestCallback(t *testing.T) {
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
			expectUseCaseData: usecase.ResultUseCase{Result: model.RequestToken{}},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              testCaseNegative2,
			token:             noAuth,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(failedAuth)},
			expectStatusCode:  http.StatusUnauthorized,
		},
		{
			name:              testCaseNegative3,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(model.ErrorAccountInActiveBahasa)},
			expectStatusCode:  http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.RequestToken{
				NewMember: true,
			}},
			expectStatusCode: http.StatusOK,
		},
		{
			name:              testCaseNegative5,
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.MFAResponse{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectStatusCode:  http.StatusForbidden,
		},
		{
			name:              "Testcase #6: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.AuthAzureToken{}, Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusForbidden},
			expectStatusCode:  http.StatusInternalServerError,
		},
		{
			name:              "Testcase #11: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.RefreshToken{}},
			expectStatusCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.AuthUseCase)
			mockAuthUseCase.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set(authorization, tt.token)
			handler := NewHTTPHandler(mockAuthUseCase, googleAuthRedirectURL)

			err := handler.AuthCallback(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)
		})
	}
}
