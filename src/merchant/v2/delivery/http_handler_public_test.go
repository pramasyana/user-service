package delivery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocksMerchant "github.com/Bhinneka/user-service/mocks/src/merchant/v2/usecase"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	publicMerchantPath       = "/api/v2/merchant/MCH220325163007/public"
	publicMerchantPathVanity = "/api/v2/merchant/public/my-merchant-name1"
)

func TestGetPublicWarehouses(t *testing.T) {
	var defaultWarehouseData = generateWarehouseData()
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		url             string
	}{
		{
			name:            testCasePositive1,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{TotalData: 10, WarehouseData: defaultWarehouseData}},
			url:             getPublicMerchantAddressPath,
		},
		{
			name:            testCaseNegative2,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             getPublicMerchantAddressPath,
		},
		{
			name:            testCasePositive2,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{TotalData: 0, WarehouseData: []*model.WarehouseData{}}},
			url:             getPublicMerchantAddressPath,
		},
		{
			name:            testCaseNegative3,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Result: model.WarehouseData{}},
			url:             getPublicMerchantAddressPath,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, tt.url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockWarehouseAddressUsecase.On("GetWarehouseAddresses", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.getPublicMerchantWarehouse(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_GetPublicDetailMerchant(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		url             string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			url:             publicMerchantPath,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             publicMerchantPath,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusOK, Result: nil},
			url:             publicMerchantPath,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, tt.url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			mockMerchantUsecase.On("GetMerchantByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.GetPublicDetailMerchant(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_GetMerchantByVanity(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		url             string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			url:             publicMerchantPathVanity,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             publicMerchantPathVanity,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusOK, Result: nil},
			url:             publicMerchantPathVanity,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, tt.url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			mockMerchantUsecase.On("GetMerchantByVanityURL", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.GetMerchantByVanity(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_GetListMerchantPublic(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		url             string
	}{
		{
			name:           testCaseNegative5,
			token:          tokenUser,
			wantStatusCode: http.StatusBadRequest,
			url:            getMerchantPath + `?page=100000000000000000000`,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			url:             badMerchantPath,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf("something bad")},
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			url:             badMerchantPath + `&status=something`,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf("something bad happenned")},
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			url:             getMerchantPath,
			wantUsecaseData: usecase.ResultUseCase{TotalData: nil},
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: "anything", TotalData: int(10)},
			url:             getMerchantPath,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, tt.url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			mockMerchantUsecase.On("GetMerchantsPublic", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.GetListMerchantPublic(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
