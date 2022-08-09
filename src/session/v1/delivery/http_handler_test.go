package delivery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Bhinneka/user-service/src/session/v1/model"
	"github.com/Bhinneka/user-service/src/session/v1/usecase"
	"github.com/Bhinneka/user-service/src/session/v1/usecase/mocks"
)

const (
	root       = "/api/v1/session"
	tokenAdmin = `Basic Ymhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNTo2MjY4NjktNmU2ZTY1LTZiNjEyMC02ZDY1NmUtNzQ2MTcyLTY5MjA2NC02OTZkNjU2ZS03MzY5MDA=`
	nonAdmin   = "Basic "
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
func generateUsecaseResultSession(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestGetSessionInfoListMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.SessionInfoUseCase))
	handler.MountInfo(e.Group("/"))
}

func TestGetSessionInfoList(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
	}{
		{
			name:  "Testcase Session #1: Positive",
			token: tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.SessionInfoList{
				Data: []model.SessionInfoResponse{{}},
			}},
			expectStatusCode: http.StatusOK,
		},
		{
			name:              "Testcase Session #2: Positive, empty session info",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.SessionInfoList{}},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              "Testcase Session #3: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: 500, Error: fmt.Errorf("pq: error")},
			expectStatusCode:  http.StatusInternalServerError,
		},
		{
			name:              "Testcase Session #4: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: ""},
			expectStatusCode:  http.StatusBadRequest,
		},
		{
			name:              "Testcase Session #5: Negative, Not authorized",
			token:             nonAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: 401, Error: fmt.Errorf("Unauthorized")},
			expectStatusCode:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionInfoUseCase := new(mocks.SessionInfoUseCase)
			mockSessionInfoUseCase.On("GetSessionInfoList", mock.Anything).Return(generateUsecaseResultSession(tt.expectUseCaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockSessionInfoUseCase)

			err = handler.GetSessionInfoList(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}
