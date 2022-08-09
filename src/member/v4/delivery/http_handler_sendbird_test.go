package delivery

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	mocksMember "github.com/Bhinneka/user-service/mocks/src/member/v1/usecase"
	authUc "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	authMock "github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
	usecase "github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testCasePositive1 = "Testcase #1: Positive"
)

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

func TestMemberHandlerV4_MountSendbird(t *testing.T) {

	e := echo.New()
	handler := NewHTTPHandlerV4(new(mocksMember.MemberUseCase), new(authMock.AuthUseCase))
	handler.MountSendbird(e.Group("/api/V4"))

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
			args:          e.Group("/api/V4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MemberHandlerV4{
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
	root                  = "/api/v4/sendbird"
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

func TestMemberHandlerV4_GetAccessToken(t *testing.T) {
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

			mockMemberUseCase.On("GetSendbirdTokenV4", mock.Anything, mock.Anything).Return(tt.expectedGetClientApp)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)

			handler := NewHTTPHandlerV4(mockMemberUseCase, mockAuthUseCase)

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

func TestMemberHandlerV4_CheckAccessToken(t *testing.T) {
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

			mockMemberUseCase.On("CheckSendbirdTokenV4", mock.Anything, mock.Anything).Return(tt.expectedGetClientApp)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(authorization, tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)

			handler := NewHTTPHandlerV4(mockMemberUseCase, mockAuthUseCase)

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
