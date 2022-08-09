package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bhinneka/user-service/helper"
	mocksMember "github.com/Bhinneka/user-service/mocks/src/member/v1/usecase"
	authUc "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	authMock "github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
	usecase "github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemberHandlerV3_MountSendbird(t *testing.T) {

	e := echo.New()
	handler := NewHTTPHandlerV3(new(mocksMember.MemberUseCase), new(authMock.AuthUseCase))
	handler.Mount(e.Group("/api/v3"))

	tests := []struct {
		name          string
		AuthUseCase   authUc.AuthUseCase
		MemberUseCase usecase.MemberUseCase
		args          *echo.Group
	}{
		{
			name:          "positif case",
			AuthUseCase:   new(authMock.AuthUseCase),
			MemberUseCase: new(mocksMember.MemberUseCase),
			args:          e.Group("/api/v3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MemberHandlerV3{
				MemberUseCase: tt.MemberUseCase,
				AuthUseCase:   tt.AuthUseCase,
			}
			h.MountSendbird(tt.args)
		})
	}
}

type VerifyMock struct{}
type Body map[string]interface{}

const (
	root                  = "/api/v3/sendbird"
	tokenAdminJwt         = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
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

func TestMemberHandlerV3_GetAccessToken(t *testing.T) {
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
			token:                tokenAdminJwt,
			expectUseCaseData:    usecase.ResultUseCase{Result: VerifyMock{}},
			expectedGetClientApp: usecase.ResultUseCase{Result: true}, // dont return error
			expectStatusCode:     http.StatusOK,
			payload:              Body{helper.TextToken: tokenAdminJwt},
		},

		{
			name:                 testCasePositive1,
			token:                tokenAdminJwt,
			expectUseCaseData:    usecase.ResultUseCase{Result: false},
			expectedGetClientApp: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest}, // dont return errorfalse}, // dont return error
			expectStatusCode:     http.StatusBadRequest,
			payload:              Body{helper.TextToken: tokenAdminJwt},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUseCase := new(mocksMember.MemberUseCase)
			mockAuthUseCase := new(authMock.AuthUseCase)

			mockMemberUseCase.On("GetSendbirdToken", mock.Anything, mock.Anything).Return(tt.expectedGetClientApp)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)

			handler := NewHTTPHandlerV3(mockMemberUseCase, mockAuthUseCase)

			err := handler.GetAccessToken(c)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}

func TestMemberHandlerV3_CheckAccessToken(t *testing.T) {
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
			token:                tokenAdminJwt,
			expectUseCaseData:    usecase.ResultUseCase{Result: VerifyMock{}},
			expectedGetClientApp: usecase.ResultUseCase{Result: true}, // dont return error
			expectStatusCode:     http.StatusOK,
			payload:              Body{helper.TextToken: tokenAdminJwt},
		},

		{
			name:                 testCasePositive1,
			token:                tokenAdminJwt,
			expectUseCaseData:    usecase.ResultUseCase{Result: false},
			expectedGetClientApp: usecase.ResultUseCase{Error: fmt.Errorf(failedResponse), HTTPStatus: http.StatusBadRequest}, // dont return errorfalse}, // dont return error
			expectStatusCode:     http.StatusBadRequest,
			payload:              Body{helper.TextToken: tokenAdminJwt},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUseCase := new(mocksMember.MemberUseCase)
			mockAuthUseCase := new(authMock.AuthUseCase)

			mockMemberUseCase.On("CheckSendbirdToken", mock.Anything, mock.Anything).Return(tt.expectedGetClientApp)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)

			handler := NewHTTPHandlerV3(mockMemberUseCase, mockAuthUseCase)

			err := handler.CheckAccessToken(c)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}
