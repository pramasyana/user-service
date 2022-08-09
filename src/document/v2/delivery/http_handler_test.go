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
	"github.com/Bhinneka/user-service/src/document/v2/model"
	"github.com/Bhinneka/user-service/src/document/v2/usecase"
	"github.com/Bhinneka/user-service/src/document/v2/usecase/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root              = "/api/v2/document"
	tokenUser         = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	tokenUserFailedID = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20ifQ.stRqFGMoWfuqMQA666SmU9lRKkoEgmUZ5pe84yYWdiU`
	tokenUserFailed   = `beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	jsonSchemaDir     = "../../../../schema/"
	testCasePositive1 = "Testcase #1: Positive"
	testCasePositive2 = "Testcase #2: Positive"
	testCaseNegative2 = "Testcase #2: Negative"
	testCaseNegative3 = "Testcase #3: Negative"
	testCaseNegative4 = "Testcase #4: Negative"
	testCaseNegative5 = "Testcase #5: Negative"
	testCaseNegative6 = "Testcase #6: Negative"
	testCaseNegative7 = "Testcase #7: Negative"
	title             = "testing"
	documentType      = "documentType"
	errorDefault      = "pq: error"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

var payloadTest = model.DocumentTypePayload{
	DocumentType: "NPWP",
	IsB2c:        "true",
	IsB2b:        "true",
	IsActive:     "true",
}

var payloadDocument = model.DocumentData{
	DocumentType: "NPWP",
	DocumentFile: "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/KTP-file-1591178959.png",
	Title:        "title",
	Number:       "123456",
}

var tests = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
	title           string
	payload         interface{}
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentData{}},
		wantStatusCode:  http.StatusCreated,
		payload:         payloadDocument,
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadDocument,
	},
	{
		name:            testCaseNegative3,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadDocument,
	},
	{
		name:            testCaseNegative4,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentError{}},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadDocument,
	},
	{
		name:            testCaseNegative5,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		payload:         model.DocumentData{},
	},
	{
		name:            testCaseNegative6,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentData{}},
		wantStatusCode:  http.StatusBadRequest,
		payload:         tokenUser,
	},
}

var addUpdateDocumentTypeTest = []struct {
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
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentType{}},
		wantStatusCode:  http.StatusCreated,
		payload:         payloadTest,
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusUnauthorized,
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
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentError{}},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadTest,
	},
	{
		name:            testCaseNegative5,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadTest,
	},
	{
		name:            testCaseNegative6,
		token:           tokenUserFailedID,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadTest,
	},
	{
		name:            testCaseNegative7,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentType{}},
		wantStatusCode:  http.StatusBadRequest,
		payload:         tokenUser,
	},
}

func generateUsecaseResultDocument(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateRSADocument() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
   	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateTokenDocument(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSADocument(), nil
	})
}

func TestHTTPDocumentHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.DocumentUseCase))
	handler.MountMe(e.Group("/anon"))
	handler.MountDocumentType(e.Group("/anon"))
}

func TestHTTPDocumentgHandlerAddDocument(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("AddUpdateDocument", context.Background(), mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.AddDocumentMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPDocumentgHandlerUpdateDocument(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("AddUpdateDocument", context.Background(), mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.UpdateDocumentMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPDocumentHandlerDeleteDocument(t *testing.T) {
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
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("DeleteDocument", context.Background(), mock.Anything, mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.DeleteDocumentMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPDocumentHandlerGetDocument(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListDocument{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListDocument{
				Document: []*model.DocumentData{{}},
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
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(errorDefault),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("GetListDocument", mock.Anything, mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))
			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.GetDocumentMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPDocumentHandlerGetDetailDocument(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.DocumentData{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
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
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(errorDefault),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("GetDetailDocument", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))
			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.GetDocumentDetailMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
func TestHTTPDocumentgHandlerAddDocumentType(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range addUpdateDocumentTypeTest {
		t.Run(tt.name, func(t *testing.T) {
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("AddUpdateDocumentType", context.Background(), mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.AddDocumentType(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPDocumentgHandlerUpdateDocumentType(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range addUpdateDocumentTypeTest {
		if tt.wantStatusCode == http.StatusCreated {
			tt.wantStatusCode = http.StatusOK
		}
		t.Run(tt.name, func(t *testing.T) {

			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("AddUpdateDocumentType", context.Background(), mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.UpdateDocumentType(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPDocumentHandlerGetDocumentType(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListDocumentType{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListDocumentType{
				DocumentType: []*model.DocumentType{{}},
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
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(errorDefault),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("GetListDocumentType", mock.Anything, mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))
			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.GetDocumentType(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPDocumentHandlerGetRequiredDocument(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.RequiredDocuments{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.RequiredDocuments{
				Merchant: model.MerchantType{
					[]model.DocumentRequire{},
					[]model.DocumentRequire{},
					[]model.DocumentRequire{},
					[]model.DocumentRequire{},
					[]model.DocumentRequire{},
				},
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
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(errorDefault),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDocumentUsecase := new(mocks.DocumentUseCase)
			mockDocumentUsecase.On("GetRequiredDocument", mock.Anything, mock.Anything).Return(generateUsecaseResultDocument(tt.wantUsecaseData))
			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenDocument(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockDocumentUsecase)

			err = handler.GetRequiredDocument(c)
			if tt.wantError {
				assert.Error(t, err)
			}
		})
	}
}
