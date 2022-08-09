package delivery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Bhinneka/user-service/src/health/model"
	"github.com/Bhinneka/user-service/src/health/usecase"
	"github.com/Bhinneka/user-service/src/health/usecase/mocks"
)

const (
	root       = "/health"
	tokenAdmin = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZW1haWwiOiJ5ZGVyYW5hQGdtYWlsLmNvbSIsImV4cCI6MTU2OTkxOTYyMiwiaWF0IjoxNTY5OTEyNDIyLCJpc3MiOiJiaGlubmVrYS5jb20iLCJqdGkiOiI0NDYwZmEzN2MwNzk2MGQyNTBkNzk5OTZiZDg1MTgxM2JkMDQ1NDcwIiwic3RhZmYiOmZhbHNlLCJzdWIiOiJVU1IxODA1NjIyMDYifQ.H8N-gLWGzPCTbbJGfs-w7gVRX0Ahxa96wKQXVg-OhXucGXbRUiyB2clWlIrVfOeifgTq_-qCpAevez4WFx22wZIh5-OJ9atYlMfcVjpjlOTSgaJE7K8gyw5fbleDwtai9DcYjIaWCF15jmlTBgilP09DCngWXUTEcgTUf9Y22nVeokRmyAoZ2TL5RgjCjM3JBIeqWiHg5iTBVmZwSIubxfuDt1X1mjO8LTV787jOrmYAl6JokC46ENY7NZX-gsB6v7_b4EfjrUFdy0pUX5agrUvtHwI3-t7sKev9rytXK_Zio3au5iewPbYGp3_BL0E_YqAsoKiSAeaggba7u4vttQ`
	nonAdmin   = "Bearer "
)

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestPing(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		expectUseCaseData usecase.ResultUseCase
		expectError       bool
		expectStatusCode  int
	}{
		{
			name:              "Testcase #1",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: model.Health{}},
			expectStatusCode:  http.StatusOK,
		},
		{
			name:              "Testcase #2",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Error: errors.New("some error")},
			expectStatusCode:  http.StatusOK,
			expectError:       true,
		},
		{
			name:              "Testcase #3",
			token:             tokenAdmin,
			expectUseCaseData: usecase.ResultUseCase{Result: &model.Health{}},
			expectStatusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHealth := new(mocks.HealthUseCase)
			mockHealth.On("Ping", mock.Anything).Return(generateUsecaseResult(tt.expectUseCaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set(echo.HeaderAuthorization, tt.token)
			handler := NewHTTPHandler(mockHealth)
			handler.Mount(e.Group("health"))

			err := handler.Ping(c)
			if tt.expectError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.expectStatusCode, rec.Code)

		})
	}
}
