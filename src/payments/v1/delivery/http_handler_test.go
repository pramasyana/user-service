package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mocksPayment "github.com/Bhinneka/user-service/mocks/src/payments/v1/usecase"
	"github.com/Bhinneka/user-service/src/payments/v1/model"
	"github.com/Bhinneka/user-service/src/payments/v1/usecase"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	headerAuth    = "Basic 0225c5c10c02aa626a9b0897cf85417a18238eb0b4acb996e519ad225f048acb3434af57b8e894456990757d0a786adbaad1aa36890277056692fe"
	TestEmail     = "tes1234@getnada.com"
	TestChannel   = "b2b"
	TestMethod    = "kredivo"
	TestToken     = "randomTokenHere"
	TestExpiredAt = "2022-01-01T23:59:59+07:00"
)

var payloadPayment = model.Payments{
	ID:        "TKN001",
	Email:     TestEmail,
	Channel:   TestChannel,
	Method:    TestMethod,
	Token:     TestToken,
	ExpiredAt: time.Now(),
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
func TestHTTPPaymentsHandler(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocksPayment.PaymentsUseCase))
	handler.MountInfo(e.Group("/basic"))
}
func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestHTTPPaymentsHandler_AddUpdatePayment(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantCompareData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         interface{}
	}{
		{
			name:            "Test Success 1 ",
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}, Error: nil},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload:         payloadPayment,
			wantStatusCode:  http.StatusCreated,
		},
		{
			name:            "Test Failed 1 ", // Compare Header Failed
			token:           headerAuth,
			wantCompareData: usecase.ResultUseCase{Error: assert.AnError},
			payload:         payloadPayment,
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            "Test Failed 2 ", //Add Update Failed
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: nil, Error: assert.AnError},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload:         payloadPayment,
		},
		{
			name:            "Test Failed 2 ", //Add Update Failed
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: nil, Error: assert.AnError},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload:         payloadPayment,
		},
		{
			name:            "Test Failed 3 ", //Payload Failed
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: nil, Error: assert.AnError},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload: model.Parameters{
				Query: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentsUsecase := new(mocksPayment.PaymentsUseCase)
			mockPaymentsUsecase.On("CompareHeaderAndBody", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantCompareData))
			mockPaymentsUsecase.On("AddUpdatePayments", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/payments", strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockPaymentsUsecase)

			err = handler.AddUpdatePayment(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPPaymentsHandler_GetDetailPayment(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantCompareData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         interface{}
	}{
		{
			name:            "Test Success 1 ",
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}, Error: nil},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload:         payloadPayment,
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            "Test Failed 1 ", // Compare Header Failed
			token:           headerAuth,
			wantCompareData: usecase.ResultUseCase{Error: assert.AnError},
			payload:         payloadPayment,
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            "Test Failed 2 ", //Get Detail Failed
			token:           headerAuth,
			wantUsecaseData: usecase.ResultUseCase{Result: nil, Error: assert.AnError},
			wantCompareData: usecase.ResultUseCase{Error: nil},
			payload:         payloadPayment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentsUsecase := new(mocksPayment.PaymentsUseCase)
			mockPaymentsUsecase.On("CompareHeaderAndBody", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantCompareData))
			mockPaymentsUsecase.On("GetPaymentDetail", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			bodyData, err := json.Marshal(tt.payload)
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/payments", strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockPaymentsUsecase)

			err = handler.GetDetailPayment(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
