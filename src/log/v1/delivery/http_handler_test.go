package delivery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	logMock "github.com/Bhinneka/user-service/src/log/v1/usecase/mocks"
	actSvcMock "github.com/Bhinneka/user-service/src/service/mocks"
	"github.com/Bhinneka/user-service/src/shared"
	sharedMock "github.com/Bhinneka/user-service/src/shared/mocks"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	defaultPath = "/v1/log"
	badURL      = "/v1/log?page=10000000000000000000000000"
)

var (
	errorDefault = errors.New("error log")
)

var testDataLog = []struct {
	name             string
	responseUsecase  sharedModel.ResultUseCase
	url              string
	expectStatusCode int
}{
	{
		name:             "test get all log #1",
		expectStatusCode: http.StatusBadRequest,
		url:              badURL,
	},
	{
		name:             "test get all log #2",
		responseUsecase:  sharedModel.ResultUseCase{Error: errorDefault},
		expectStatusCode: http.StatusBadRequest,
		url:              defaultPath,
	},
	{
		name:             "test get all log #3",
		responseUsecase:  sharedModel.ResultUseCase{Result: "something", Meta: shared.Meta{}},
		expectStatusCode: http.StatusOK,
		url:              defaultPath,
	},
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestGetAllLog(t *testing.T) {
	for _, tc := range testDataLog {
		mockUC := new(logMock.LogUsecase)
		mockActivityService := new(actSvcMock.ActivityServices)

		mockUC.On("GetAll", mock.Anything, mock.Anything).Return(sharedMock.CreateUsecaseResult(tc.responseUsecase))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, tc.url, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewHTTPHandler(mockActivityService, mockUC)
		handler.Mount(e.Group("/v1/log"))
		err := handler.GetAll(c)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}

var testDataGetSingleLog = []struct {
	name             string
	responseUsecase  sharedModel.ResultUseCase
	expectStatusCode int
}{
	{
		name:             "test get log #1",
		expectStatusCode: http.StatusBadRequest,
		responseUsecase:  sharedModel.ResultUseCase{Error: errorDefault},
	},
	{
		name:             "test get log #2",
		expectStatusCode: http.StatusOK,
		responseUsecase:  sharedModel.ResultUseCase{Result: "someResult"},
	},
}

func TestGetLogByID(t *testing.T) {
	for _, tc := range testDataGetSingleLog {
		mockUC := new(logMock.LogUsecase)
		mockActivityService := new(actSvcMock.ActivityServices)

		mockUC.On("GetByID", mock.Anything, mock.Anything).Return(sharedMock.CreateUsecaseResult(tc.responseUsecase))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/v1/log/somerandom", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewHTTPHandler(mockActivityService, mockUC)
		err := handler.GetByID(c)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}
