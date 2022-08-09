package delivery

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/applications/v1/model"
	"github.com/Bhinneka/user-service/src/applications/v1/usecase"
	"github.com/Bhinneka/user-service/src/applications/v1/usecase/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root              = "/api/v2/shipping-address"
	tokenUser         = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	tokenUserFailed   = `beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	jsonSchemaDir     = "../../../../schema/"
	testCasePositive1 = "Testcase #1: Positive"
	testCasePositive2 = "Testcase #2: Positive"
	testCaseNegative2 = "Testcase #2: Negative"
	testCaseNegative3 = "Testcase #3: Negative"
	testCaseNegative4 = "Testcase #4: Negative"
	testCaseNegative5 = "Testcase #5: Negative"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateUsecaseResultApplication(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateRSAApplication() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
   	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateTokenApplication(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSAApplication(), nil
	})
}

func TestHTTPApplicationHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.ApplicationsUseCase))
	handler.MountInfo(e.Group("/anon"))
}

func TestHTTPApplicationHandlerAddApplication(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	payloadTest := model.Application{
		Name: "app name",
		URL:  "http://development.shark.bhinneka.com",
		Logo: "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/user_images/pack-shark.svg",
	}
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         interface{}
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Application{}},
			wantStatusCode:  http.StatusCreated,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ApplicationList{}, HTTPStatus: 400, Error: fmt.Errorf("something happened")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Application{}},
			wantStatusCode:  http.StatusBadRequest,
			payload:         tokenUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApplicationUsecase := new(mocks.ApplicationsUseCase)
			mockApplicationUsecase.On("AddUpdateApplication", context.Background(), mock.Anything).Return(generateUsecaseResultApplication(tt.wantUsecaseData))

			bodyData, _ := json.Marshal(tt.payload)
			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenApplication(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockApplicationUsecase)

			err = handler.AddApplication(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPApplicationHandlerUpdateApplication(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	payloadTest := model.Application{
		ID:   "1",
		Name: "app name",
		URL:  "http://development.shark.bhinneka.com",
		Logo: "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/user_images/pack-shark.svg",
	}
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         interface{}
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Application{}},
			wantStatusCode:  http.StatusOK,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ApplicationList{}, HTTPStatus: 400, Error: fmt.Errorf("something bad happenned")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         payloadTest,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Application{}, HTTPStatus: 400, Error: fmt.Errorf("oops")},
			wantStatusCode:  http.StatusBadRequest,
			payload:         tokenUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApplicationUsecase := new(mocks.ApplicationsUseCase)
			mockApplicationUsecase.On("AddUpdateApplication", context.Background(), mock.Anything).Return(generateUsecaseResultApplication(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenApplication(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockApplicationUsecase)

			err = handler.UpdateApplication(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPApplicationHandlerDeleteApplication(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApplicationUsecase := new(mocks.ApplicationsUseCase)
			mockApplicationUsecase.On("DeleteApplication", context.Background(), mock.Anything).Return(generateUsecaseResultApplication(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenApplication(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockApplicationUsecase)

			err = handler.DeleteApplication(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPApplicationHandlerGetApplicationList(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListApplication{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListApplication{
				Application: []*model.Application{{}},
			}},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf("pq: error"),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv", HTTPStatus: 400, Error: fmt.Errorf("someting really bad happened")},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApplicationUsecase := new(mocks.ApplicationsUseCase)
			mockApplicationUsecase.On("GetListApplication", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultApplication(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenApplication(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockApplicationUsecase)

			err = handler.GetApplicationList(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
