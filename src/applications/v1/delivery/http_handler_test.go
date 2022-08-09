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

	"github.com/Bhinneka/user-service/src/applications/v1/model"
	"github.com/Bhinneka/user-service/src/applications/v1/usecase"
	"github.com/Bhinneka/user-service/src/applications/v1/usecase/mocks"
)

const (
	root       = "/api/v1/applications"
	tokenAdmin = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZW1haWwiOiJ5ZGVyYW5hQGdtYWlsLmNvbSIsImV4cCI6MTU2OTkxOTYyMiwiaWF0IjoxNTY5OTEyNDIyLCJpc3MiOiJiaGlubmVrYS5jb20iLCJqdGkiOiI0NDYwZmEzN2MwNzk2MGQyNTBkNzk5OTZiZDg1MTgxM2JkMDQ1NDcwIiwic3RhZmYiOmZhbHNlLCJzdWIiOiJVU1IxODA1NjIyMDYifQ.H8N-gLWGzPCTbbJGfs-w7gVRX0Ahxa96wKQXVg-OhXucGXbRUiyB2clWlIrVfOeifgTq_-qCpAevez4WFx22wZIh5-OJ9atYlMfcVjpjlOTSgaJE7K8gyw5fbleDwtai9DcYjIaWCF15jmlTBgilP09DCngWXUTEcgTUf9Y22nVeokRmyAoZ2TL5RgjCjM3JBIeqWiHg5iTBVmZwSIubxfuDt1X1mjO8LTV787jOrmYAl6JokC46ENY7NZX-gsB6v7_b4EfjrUFdy0pUX5agrUvtHwI3-t7sKev9rytXK_Zio3au5iewPbYGp3_BL0E_YqAsoKiSAeaggba7u4vttQ`
	nonAdmin   = "Bearer "
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

func TestGetApplicationListMount(t *testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.ApplicationsUseCase))
	handler.MountInfo(e.Group("/"))
	assert.Equal(t, handler, handler)
}

func TestGetApplicationList(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
	}{
		{
			name:  "Testcase #1: Positive",
			token: tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.ApplicationList{
				Data: []model.Application{{}},
			}},
			expectStatusCode: http.StatusOK,
		},
		{
			name:              "Testcase #2: Positive, empty application",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.ApplicationList{}},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              "Testcase #3: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: 500, Error: fmt.Errorf("pq: error")},
			expectStatusCode:  http.StatusInternalServerError,
		},
		{
			name:              "Testcase #4: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: "", HTTPStatus: 400, Error: fmt.Errorf("something error")},
			expectStatusCode:  http.StatusBadRequest,
		},
		{
			name:              "Testcase #5: Negative, Not authorized",
			token:             nonAdmin,
			expectUseCaseData: usecase.ResultUseCase{HTTPStatus: 401, Error: fmt.Errorf("Unauthorized")},
			expectStatusCode:  http.StatusUnauthorized,
		},
		{
			name:              "Testcase #6: Negative",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: nil},
			expectStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApplicationUseCase := new(mocks.ApplicationsUseCase)
			mockApplicationUseCase.On("GetApplicationsList", mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockApplicationUseCase)

			err = handler.GetApplicationsList(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}
