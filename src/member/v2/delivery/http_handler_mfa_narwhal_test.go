package delivery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocksMember "github.com/Bhinneka/user-service/mocks/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNarwhalMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFASettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GetNarwhalMFASettings", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetNarwhalMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestGenerateNarwhalMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFAGenerateSettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GenerateMFASettings", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GenerateNarwhalMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestActivateNarwhalMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFAActivateSettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ActivateMFASettings", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ActivateNarwhalMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestDisabledNarwhalMFASetting(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFASettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("DisabledMFASetting", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.DisabledNarwhalMFASetting(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
