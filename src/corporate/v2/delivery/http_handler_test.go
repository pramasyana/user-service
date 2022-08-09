package delivery

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/corporate/v2/usecase"
	"github.com/Bhinneka/user-service/src/corporate/v2/usecase/mocks"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root              = "/api/v2/corporate"
	tokenUser         = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	tokenUserFailed   = `beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	jsonSchemaDir     = "../../../../schema/"
	testCasePositive1 = "Testcase #1: Positive"
	testCasePositive2 = "Testcase #2: Positive"
	testCaseNegative2 = "Testcase #2: Negative"
	testCaseNegative3 = "Testcase #3: Negative"
	testCaseNegative4 = "Testcase #4: Negative"
	testCaseNegative5 = "Testcase #5: Negative"
	pqError           = "pq: error"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateUsecaseResultCorporate(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateRSACorporate() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
   	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateTokenCorporate(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSACorporate(), nil
	})
}

func TestHTTPCorporateHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.CorporateUseCase))
	handler.MountCorporate(e.Group("/anon"))
}

var testsGetContact = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: sharedModel.ListContact{}},
		wantStatusCode:  http.StatusOK,
	},
	{
		name:  testCasePositive2,
		token: tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: sharedModel.ListContact{
			Contact: []*sharedModel.B2BContactData{{}},
		}},
		wantStatusCode: http.StatusOK,
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusUnauthorized,
	},
	{
		name:  testCaseNegative3,
		token: tokenUser,
		wantUsecaseData: usecase.ResultUseCase{
			HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(pqError),
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

var testsGetContactDetail = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: sharedModel.B2BContactData{}},
		wantStatusCode:  http.StatusOK,
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
			HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(pqError),
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

func TestHTTPCorporateHandlerGetContact(t *testing.T) {
	for _, tt := range testsGetContact {
		t.Run(tt.name, func(t *testing.T) {
			mockCorporateUsecase := new(mocks.CorporateUseCase)
			mockCorporateUsecase.On("GetAllListContact", mock.Anything, mock.Anything).Return(generateUsecaseResultCorporate(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenCorporate(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockCorporateUsecase)

			err = handler.GetContactList(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPCorporateHandlerGetDetailContact(t *testing.T) {
	for _, tt := range testsGetContactDetail {
		t.Run(tt.name, func(t *testing.T) {
			mockCorporateUsecase := new(mocks.CorporateUseCase)
			mockCorporateUsecase.On("GetDetailContact", mock.Anything, mock.Anything).Return(generateUsecaseResultCorporate(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenCorporate(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockCorporateUsecase)

			err = handler.GetContactDetail(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

var testsImportContact = []struct {
	name            string
	payload         string
	wantUsecaseData error
	wantError       bool
	wantStatusCode  int
}{
	{
		name:            "Test #1",
		wantUsecaseData: nil,
		wantStatusCode:  http.StatusOK,
		payload:         `{"file":"ZW1haWwsZmlyc3ROYW1lLGxhc3ROYW1lLHBob25lTnVtYmVyLGFjY291bnRJZCx0cmFuc2FjdGlvblR5cGUNCmFndXMub2tlQGdldG5hZGEuY29tLEFndXMsLDA4ODg4ODg4OCxCMkJBQ0MwMDAwMDEsc2hvcGNhcnQ="}`,
	},
	{
		name:            "Test #2",
		wantUsecaseData: nil,
		wantStatusCode:  http.StatusBadRequest,
		payload:         `{"file":,"ZW1haWwsZmlyc3ROYW1lLGxhc3ROYW1lLHBob25lTnVtYmVyLGFjY291bnRJZCx0cmFuc2FjdGlvblR5cGUNCmFndXMub2tlQGdldG5hZGEuY29tLEFndXMsLDA4ODg4ODg4OCxCMkJBQ0MwMDAwMDEsc2hvcGNhcnQ="}`,
	},
	{
		name:           "Test #3",
		wantStatusCode: http.StatusBadRequest,
		payload:        `{"file":"aaZW1haWwsZmlyc3ROYW1lLGxhc3ROYW1lLHBob25lTnVtYmVyLGFjY291bnRJZCx0cmFuc2FjdGlvblR5cGUNCmFndXMub2tlQGdldG5hZGEuY29tLEFndXMsLDA4ODg4ODg4OCxCMkJBQ0MwMDAwMDEsc2hvcGNhcnQ="}`,
	},
	{
		name:            "Test #4",
		wantStatusCode:  http.StatusBadRequest,
		payload:         `{"file":"ZW1haWwsZmlyc3ROYW1lLGxhc3ROYW1lLHBob25lTnVtYmVyLGFjY291bnRJZCx0cmFuc2FjdGlvblR5cGUNCmFndXMub2tlQGdldG5hZGEuY29tLEFndXMsLDA4ODg4ODg4OCxCMkJBQ0MwMDAwMDEsc2hvcGNhcnQ="}`,
		wantUsecaseData: errors.New("some error"),
	},
}

func TestImportContact(t *testing.T) {
	os.Setenv("NO_TOKEN", "1")
	for _, tc := range testsImportContact {
		t.Run(tc.name, func(t *testing.T) {
			mockCorporateUsecase := new(mocks.CorporateUseCase)
			mockCorporateUsecase.On("ImportContact", mock.Anything, mock.Anything).Return(nil, tc.wantUsecaseData)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/v2/corporate/contact/import", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockCorporateUsecase)

			err := handler.ImportContact(c)
			if tc.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tc.wantStatusCode, rec.Code)
		})
	}
}
