package delivery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bhinneka/golib/jsonschema"
	mocksMerchant "github.com/Bhinneka/user-service/mocks/src/merchant/v2/usecase"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	getMerchantAddressPath       = "/api/v2/merchant/MCH201008161332/warehouses"
	getMerchantAddressIDPath     = "/api/v2/merchant/MCH201008161332/warehouses/ADDRS201201014824935"
	getPublicMerchantAddressPath = "/api/v2/merchant/MCH201008161332/warehouses/public"
)

func TestUpdateMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:           testCaseNegative2,
			token:          tokenUserFailed,
			wantStatusCode: http.StatusBadRequest,
			payload:        "",
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			payload:         bodyUpdate,
		},
		{
			name:           testCaseNegative3,
			token:          tokenUser,
			wantStatusCode: http.StatusBadRequest,
			payload:        "{},",
		},
		{
			name:           testCaseNegative4,
			token:          tokenUser,
			wantStatusCode: http.StatusBadRequest,
			payload:        `{"bankId":1111111111111111111111}`,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault},
			payload:         bodyUpdate,
		},
		{
			name:            testCaseNegative6,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantData{}},
			wantStatusCode:  http.StatusBadRequest,
			payload:         bodyUpdate,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			mockMerchantUsecase.On("UpdateMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMerchantUsecase.On("PublishToKafkaMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, merchantWithPath, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.updateMerchant(c)

			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestDeleteMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:           testCaseNegative2,
			token:          tokenUserFailed,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf("another error")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantData{}},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			mockMerchantUsecase.On("DeleteMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMerchantUsecase.On("PublishToKafkaMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, merchantWithPath, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.deleteMerchant(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestGetMerchants(t *testing.T) {
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
			mockMerchantUsecase.On("GetMerchants", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.getList(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestGetSingleMerchant(t *testing.T) {
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
			url:             merchantWithPath,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             merchantWithPath,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusOK, Result: nil},
			url:             merchantWithPath,
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

			handler.getMerchant(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestRejectMerchantRegistration(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:           testCaseNegative2,
			token:          tokenUserFailed,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf("unknown error")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantData{}},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			mockMerchantUsecase.On("RejectMerchantRegistration", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, merchantWithPath, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.rejectMerchantRegistration(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestRejectMerchantUpgrade(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:           testCaseNegative2,
			token:          tokenUserFailed,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf("unknown error happened")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantData{}},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			mockMerchantUsecase.On("RejectMerchantUpgrade", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, merchantWithPath, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.rejectMerchantUpgrade(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestCreateMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:           testCaseNegative2,
			token:          tokenUserFailed,
			wantStatusCode: http.StatusBadRequest,
			payload:        "",
		},
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusCreated,
			payload:         bodyUpdate,
		},
		{
			name:           testCaseNegative3,
			token:          tokenUser,
			wantStatusCode: http.StatusBadRequest,
			payload:        "{},",
		},
		{
			name:           testCaseNegative4,
			token:          tokenUser,
			wantStatusCode: http.StatusBadRequest,
			payload:        `{"bankId":1111111111111111111111}`,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			payload:         bodyUpdate,
		},
		{
			name:            testCaseNegative6,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantData{}},
			wantStatusCode:  http.StatusBadRequest,
			payload:         bodyUpdate,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			mockMerchantUsecase.On("CreateMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMerchantUsecase.On("PublishToKafkaMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, merchantWithPath, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.createMerchant(c)

			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func generateWarehouseData() []*model.WarehouseData {
	m := []*model.WarehouseData{}

	n := model.WarehouseData{}
	o := model.WarehouseData{}
	m = append(m, &n, &o)
	return m
}
func TestGetWarehouses(t *testing.T) {
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
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{TotalData: 10, WarehouseData: defaultWarehouseData}},
			url:             getMerchantAddressPath,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             getMerchantAddressPath,
		},
		{
			name:            testCasePositive2,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{TotalData: 0, WarehouseData: []*model.WarehouseData{}}},
			url:             getMerchantAddressPath,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Result: model.WarehouseData{}},
			url:             getMerchantAddressPath,
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
			mockWarehouseAddressUsecase.On("GetWarehouseAddresses", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.getMerchantWarehouse(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestGetWarehouseByMerchantID(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.WarehouseData{}},
			url:             getMerchantAddressIDPath,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Error: errDefault, HTTPStatus: http.StatusBadRequest},
			url:             getMerchantAddressIDPath,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{}},
			url:             getMerchantAddressIDPath,
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
			mockWarehouseAddressUsecase.On("GetWarehouseAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			handler.getWarehouseAddress(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
