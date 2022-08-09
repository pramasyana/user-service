package delivery

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	mocksMember "github.com/Bhinneka/user-service/mocks/src/member/v1/usecase"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authUsecase "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	authMock "github.com/Bhinneka/user-service/src/auth/v1/usecase/mocks"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	defEmailInput          = `{"email":"myemail@bhinneka.com"}`
	msgErrorPq             = "pq: error"
	testCasePositive1      = "Testcase #1: Positive"
	testCasePositive2      = "Testcase #2: Positive"
	testCaseNegative2      = "Testcase #2: Negative"
	testCaseNegative3      = "Testcase #3: Negative"
	testCaseNegative4      = "Testcase #4: Negative"
	testCaseNegative5      = "Testcase #5: Negative"
	testCaseNegative6      = "Testcase #6: Negative"
	testCaseNegative7      = "Testcase #6: Negative"
	usecaseCheckEmail      = "CheckEmailAndMobileExistence"
	labelMember            = "member"
	tokenAdmin             = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	defInputRegister       = `{"firstName":"pian","lastName":"zm","email":"pian.mutakin@bhinneka.com","password":"passwordKeren","rePassword":"passwordKeren","gender":"M","dob":"2012-09-01","mobile":"081387788803","registerType":"personal","signUpFrom":"sturgeon","socialMedia":{"facebookId":"someFacebookId","googleId":"someGoogleId","appleId":"someAppleId"}}`
	defInputRegisterNormal = `{"firstName":"pian","lastName":"zm","email":"pian.mutakin@bhinneka.com","password":"passwordKeren","rePassword":"passwordKeren","gender":"M","dob":"2012-09-01","mobile":"081387788803","registerType":"personal","socialMedia":{"facebookId":"someFacebookId","googleId":"someGoogleId","appleId":"someAppleId"}}`
)

func generateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSA(), nil
	})
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

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func usecaseResultAuth(data authUsecase.ResultUseCase) <-chan authUsecase.ResultUseCase {
	output := make(chan authUsecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestMountHandlerV3(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandlerV3(new(mocksMember.MemberUseCase), new(authMock.AuthUseCase))
	handler.Mount(e.Group("/api/v3"))
}

func TestHandlerForgotPasswordV3(t *testing.T) {
	tests := []struct {
		name                                  string
		token                                 string
		wantUCForgotPassword, wantUCSendEmail usecase.ResultUseCase
		wantError                             bool
		wantStatusCode                        int
		param1                                string
	}{
		{
			name:                 testCasePositive1,
			token:                tokenAdmin,
			wantUCForgotPassword: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:       http.StatusOK,
			param1:               defEmailInput,
		},
		{
			name:           testCaseNegative2,
			token:          tokenAdmin,
			wantStatusCode: http.StatusBadRequest,
			param1:         `{"email":"",}`,
			wantError:      true,
		},
		{
			name:                 testCaseNegative3,
			token:                tokenAdmin,
			param1:               defEmailInput,
			wantUCForgotPassword: usecase.ResultUseCase{Error: fmt.Errorf("email required"), HTTPStatus: http.StatusBadRequest},
			wantStatusCode:       http.StatusOK,
		},
		{
			name:                 testCaseNegative4,
			token:                tokenAdmin,
			param1:               defEmailInput,
			wantUCForgotPassword: usecase.ResultUseCase{Result: "interface{}"},
			wantStatusCode:       http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockAuthUsecase := new(authMock.AuthUseCase)
			mockMemberUsecase.On("ForgotPassword", mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.wantUCForgotPassword))
			mockMemberUsecase.On("SendEmailForgotPassword", mock.Anything, mock.Anything).Return(generateUsecaseResult(tc.wantUCSendEmail))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v3/auth", bytes.NewBufferString(tc.param1))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tc.token)
			c.Set("token", token)
			handler := NewHTTPHandlerV3(mockMemberUsecase, mockAuthUsecase)

			err := handler.ForgotPassword(c)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestRegisterMemberV3(t *testing.T) {
	tests := []struct {
		name                                                      string
		token                                                     string
		usecaseCheckEmailPhone, usecaseRegister, usecaseSendEmail usecase.ResultUseCase
		usecaseGenerateToken                                      authUsecase.ResultUseCase
		wantError                                                 bool
		wantStatusCode                                            int
		payload                                                   string
	}{
		{
			name:                 testCasePositive1,
			token:                tokenAdmin,
			usecaseRegister:      usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:       http.StatusCreated,
			payload:              defInputRegister,
			usecaseGenerateToken: authUsecase.ResultUseCase{Result: authModel.RequestToken{}},
		},
		{
			name:                   testCaseNegative2,
			token:                  tokenAdmin,
			usecaseCheckEmailPhone: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:         http.StatusBadRequest,
			payload:                defInputRegister,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
			payload:        defInputRegister,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
			payload:         defInputRegister,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			usecaseSendEmail: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode:       http.StatusCreated,
			payload:              defInputRegister,
			usecaseGenerateToken: authUsecase.ResultUseCase{Result: authModel.RequestToken{}},
		},
		{
			name:                 testCasePositive2,
			token:                tokenAdmin,
			usecaseRegister:      usecase.ResultUseCase{Result: model.SuccessResponse{}},
			usecaseSendEmail:     usecase.ResultUseCase{HTTPStatus: http.StatusOK},
			wantStatusCode:       http.StatusCreated,
			payload:              defInputRegisterNormal,
			usecaseGenerateToken: authUsecase.ResultUseCase{Result: authModel.RequestToken{}},
		},
		{
			name:           testCaseNegative6,
			token:          tokenAdmin,
			wantStatusCode: http.StatusBadRequest,
			payload:        `{"email":"someEmail",}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockAuthUsecase := new(authMock.AuthUseCase)
			mockMemberUsecase.On(usecaseCheckEmail, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseCheckEmailPhone))
			mockMemberUsecase.On("RegisterMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseRegister))
			mockMemberUsecase.On("SendEmailRegisterMember", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseSendEmail))
			mockMemberUsecase.On("SendEmailWelcomeMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseSendEmail))
			mockAuthUsecase.On("GenerateToken", mock.Anything, mock.Anything, mock.Anything).Return(usecaseResultAuth(tt.usecaseGenerateToken))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v3/register", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandlerV3(mockMemberUsecase, mockAuthUsecase)

			err := handler.RegisterMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestAddMemberV3(t *testing.T) {
	tests := []struct {
		name                                                      string
		token                                                     string
		usecaseCheckEmailPhone, usecaseRegister, usecaseSendEmail usecase.ResultUseCase
		usecaseGenerateToken                                      authUsecase.ResultUseCase
		wantError                                                 bool
		wantStatusCode                                            int
		payload                                                   string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
			payload:         defInputRegister,
		},
		{
			name:           testCaseNegative2,
			token:          tokenAdmin,
			wantStatusCode: http.StatusBadRequest,
			payload:        `{"email":"someEmails",}`,
		},
		{
			name:                   testCaseNegative3,
			token:                  tokenAdmin,
			usecaseCheckEmailPhone: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:         http.StatusBadRequest,
			payload:                defInputRegister,
		},
		{
			name:  testCaseNegative4,
			token: tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
			payload:        defInputRegister,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			usecaseRegister: usecase.ResultUseCase{Result: true},
			wantStatusCode:  http.StatusInternalServerError,
			payload:         defInputRegister,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockAuthUsecase := new(authMock.AuthUseCase)
			mockMemberUsecase.On(usecaseCheckEmail, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseCheckEmailPhone))
			mockMemberUsecase.On("RegisterMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseRegister))
			mockMemberUsecase.On("SendEmailAddMember", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.usecaseSendEmail))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v3/member", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandlerV3(mockMemberUsecase, mockAuthUsecase)

			err := handler.AddMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
