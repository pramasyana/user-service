package delivery

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	mocksMember "github.com/Bhinneka/user-service/mocks/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root                  = "/api/v2/me"
	tokenAdmin            = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	nonAdmin              = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6ZmFsc2UsImRpZCI6ImMwYjRkMWI0YzQ0NzQiLCJkbGkiOiJXRUIiLCJleHAiOjE2Mjk5MTc4MTksImlhdCI6MTUyOTkxNjkxOSwiaXNzIjoiYmhpbm5la2EuY29tIiwic3ViIjoiYmhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNSJ9.ZJjhtpsj5VVDvsb_BMQsurBl3KrdRkf-0Q1Q5nfktAe36DwRKjJCFQmXnwPZgq_kfro9ThFDuZR7YDzRFjZfbHzAyKjD3cY7CDI6HPMDeggw52WFmrNa9RviOqs1MNCPumQw6DJgiioEJ58lPskDoxS-35q2zFyeNWEu_wLr8S7sT7gVtOQQ5JeL5xzPPOuft10PPQUFXbOLgu4Lopi7cWl4cw93ZyMvw2xDw3Ga-nQmubtSzqJOq0XYYASKQA4MjcaSRdOmdt5MM9p3tnwU6-6nmxS6WcxiHHypSZR2gHltAh_CITKbxBaLHMBz1l79wWsC9rGhyhPJ0f9Q4cR_Cg"
	tokenNonAdmin         = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6dHJ1ZSwiZGlkIjoiYzBiNGQxYjRjNDQ3NCIsImRsaSI6IldFQiIsImlhdCI6MTU0NDU0Mjk2MCwiaXNzIjoiYmhpbm5la2EuY29tIiwic3ViIjoiYmhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNSIsImp0aSI6IjEwMmYxOWQwLTQyNjgtNDgwNC1iYmQ4LTBiMzlkMDBlZGZkMSIsImV4cCI6MTU4OTk2MDM1MH0.9TizKfQj48oHAZTTrUA5htYcZdI9ly1eKiko99CUZvI`
	tokenUserFailedID     = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20ifQ.stRqFGMoWfuqMQA666SmU9lRKkoEgmUZ5pe84yYWdiU`
	tokenfailed           = "beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6ZmFsc2UsImRpZCI6ImMwYjRkMWI0YzQ0NzQiLCJkbGkiOiJXRUIiLCJleHAiOjE2Mjk5MTc4MTksImlhdCI6MTUyOTkxNjkxOSwiaXNzIjoiYmhpbm5la2EuY29tIiwic3ViIjoiYmhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNSJ9.ZJjhtpsj5VVDvsb_BMQsurBl3KrdRkf-0Q1Q5nfktAe36DwRKjJCFQmXnwPZgq_kfro9ThFDuZR7YDzRFjZfbHzAyKjD3cY7CDI6HPMDeggw52WFmrNa9RviOqs1MNCPumQw6DJgiioEJ58lPskDoxS-35q2zFyeNWEu_wLr8S7sT7gVtOQQ5JeL5xzPPOuft10PPQUFXbOLgu4Lopi7cWl4cw93ZyMvw2xDw3Ga-nQmubtSzqJOq0XYYASKQA4MjcaSRdOmdt5MM9p3tnwU6-6nmxS6WcxiHHypSZR2gHltAh_CITKbxBaLHMBz1l79wWsC9rGhyhPJ0f9Q4cR_Cg"
	fileBase64            = "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"
	fileBase64Member      = "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"
	testCasePositive1     = "Testcase #1: Positive"
	testCasePositive2     = "Testcase #2: Positive"
	testCaseNegative2     = "Testcase #2: Negative"
	testCaseNegative3     = "Testcase #3: Negative"
	testCaseNegative4     = "Testcase #4: Negative"
	testCaseNegative5     = "Testcase #5: Negative"
	testCaseNegative6     = "Testcase #6: Negative"
	testCaseNegative7     = "Testcase #6: Negative"
	getDetailMemberByID   = "GetDetailMemberByID"
	optionalParamSturgeon = "sturgeon"
	paramRequestForm      = "requestFrom"
	paramPhone            = "0812345678"
	labelMember           = "member"
	msgErrorPq            = "pq: error"
	usecaseCheckEmail     = "CheckEmailAndMobileExistence"
	provID                = "provinceId"
)

type Body map[string]interface{}

const (
	jsonSchemaDir = "../../../../schema/"
)

func generateRSA() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSA(), nil
	})
}
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

func TestHTTPMemberHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocksMember.MemberUseCase))
	handler.MountMe(e.Group("/me"))
	handler.MountAdmin(e.Group("/anon"))
	handler.Mount(e.Group(""))
	handler.MountMember(e.Group("/member"))
}

func TestHTTPMemberHandlerFindDetailProfile(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{HasPassword: true}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:  http.StatusUnauthorized,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUserFailedID,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.FindDetailProfile(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerEditProfile(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          paramPhone,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:   testCaseNegative2,
			token:  tokenAdmin,
			param1: paramPhone,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUserFailedID,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("UpdateDetailMemberByID", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("mobile", tt.param1)
			data.Set(provID, "153")
			data.Set("provinceName", "DKI Jakarta")
			data.Set(provID, "153")
			data.Set("cityId", "2111")
			data.Set("cityName", "Jakarta Selatan")
			data.Set("districtId", "4765")
			data.Set("districtName", "Setia Budi")
			data.Set("subDistrictId", "4421")
			data.Set("subDistrictName", "Kuningan Timur")
			data.Set("postalCode", "12950")
			data.Set("street1", "address")
			data.Set("requestFrom", "sturgeon")

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.EditProfile(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPMemberHandlerUpdatePassword(t *testing.T) {
	tests := []struct {
		name                              string
		token                             string
		wantUsecaseData, wantUsecaseData2 usecase.ResultUseCase
		wantError                         bool
		wantStatusCode                    int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative3,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUserFailedID,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("UpdatePassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, "Bearer ajsdhaklsdjkasjd")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.UpdatePassword(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerAddPassword(t *testing.T) {
	tests := []struct {
		name                               string
		token                              string
		wantUsecaseData1, wantUsecaseData2 usecase.ResultUseCase
		wantError                          bool
		wantStatusCode                     int
	}{
		{
			name:             testCasePositive1,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             testCaseNegative2,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{HTTPStatus: 401, Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:   http.StatusUnauthorized,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative4,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:   http.StatusInternalServerError,
		},
		{
			name:             "Testcase #5: Negative, check whether password exists or not",
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{Password: "ppppp", Salt: "ssss"}},
			wantStatusCode:   http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative6,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{HTTPStatus: 500, Error: fmt.Errorf("error")},
			wantStatusCode:   http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative7,
			token:            tokenUserFailedID,
			wantUsecaseData1: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData1))
			mockMemberUsecase.On("AddNewPassword", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.AddPassword(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerRegisterMember(t *testing.T) {
	tests := []struct {
		name                                                 string
		token                                                string
		wantUsecaseData1, wantUsecaseData2, wantUsecaseData3 usecase.ResultUseCase
		wantError                                            bool
		wantStatusCode                                       int
		param1, param2                                       string
	}{
		{
			name:             testCasePositive1,
			token:            tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:   http.StatusCreated,
		},
		{
			name:             testCasePositive2,
			token:            tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:   http.StatusCreated,
			param1:           optionalParamSturgeon,
			param2:           "merchant",
		},
		{
			name:             testCaseNegative3,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:   http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative5,
			token:            tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:   http.StatusInternalServerError,
		},

		{
			name:             testCaseNegative6,
			token:            tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantUsecaseData3: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusCreated,
			param1:         optionalParamSturgeon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("signUpFrom", tt.param1)
			data.Set("registerType", tt.param2)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(usecaseCheckEmail, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData1))
			mockMemberUsecase.On("RegisterMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))
			mockMemberUsecase.On("SendEmailRegisterMember", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData3))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.RegisterMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerActivationMember(t *testing.T) {
	tests := []struct {
		name                              string
		token                             string
		wantUsecaseData, wantUsecaseData2 usecase.ResultUseCase
		wantError                         bool
		wantStatusCode                    int
		param1                            string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCasePositive2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
			param1:          optionalParamSturgeon,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusOK,
			param1:         optionalParamSturgeon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set(paramRequestForm, tt.param1)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ActivateMember", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ActivationMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerForgotPassword(t *testing.T) {
	tests := []struct {
		name                              string
		token                             string
		wantUsecaseData, wantUsecaseData2 usecase.ResultUseCase
		wantError                         bool
		wantStatusCode                    int
		param1                            string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCasePositive2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
			param1:          optionalParamSturgeon,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusOK,
			param1:         optionalParamSturgeon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set(paramRequestForm, tt.param1)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ForgotPassword", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMemberUsecase.On("SendEmailForgotPassword", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ForgotPassword(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerValidateToken(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative3,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ValidateToken", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ValidateToken(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerChangeForgotPassword(t *testing.T) {
	tests := []struct {
		name                              string
		token                             string
		wantUsecaseData, wantUsecaseData2 usecase.ResultUseCase
		wantError                         bool
		wantStatusCode                    int
		param1                            string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCasePositive2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
			param1:          optionalParamSturgeon,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusOK,
			param1:         optionalParamSturgeon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set(paramRequestForm, tt.param1)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ChangeForgotPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ChangeForgotPassword(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerActivateNewPassword(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative3,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ActivateNewPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ActivateNewPassword(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerImportMember(t *testing.T) {
	tests := []struct {
		name                                                 string
		token                                                string
		args                                                 string
		resultInvalidRows                                    []string
		wantUsecaseData1, wantUsecaseData2, wantUsecaseData3 usecase.ResultUseCase
		wantError                                            bool
		wantStatusCode                                       int
		parseDataResult                                      []*model.Member
	}{
		{
			name:             testCasePositive1,
			token:            tokenAdmin,
			wantUsecaseData2: usecase.ResultUseCase{Result: []model.SuccessResponse{}},
			parseDataResult:  []*model.Member{},
			wantStatusCode:   http.StatusCreated,
			args:             fileBase64,
		},
		{
			name:              testCaseNegative2,
			token:             tokenAdmin,
			resultInvalidRows: []string{"email1", "email2"},
			wantStatusCode:    http.StatusBadRequest,
			args:              "",
		},
		{
			name:             testCaseNegative3,
			token:            tokenAdmin,
			wantUsecaseData3: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:   http.StatusBadRequest,
			args:             fileBase64,
		},
		{
			name:              testCaseNegative4,
			token:             tokenAdmin,
			resultInvalidRows: []string{"email1@gmail.com"},
			wantStatusCode:    http.StatusBadRequest,
			args:              fileBase64,
		},
		{
			name:             testCaseNegative5,
			token:            tokenAdmin,
			wantStatusCode:   http.StatusBadRequest,
			args:             fileBase64,
			wantUsecaseData2: usecase.ResultUseCase{Error: fmt.Errorf("some error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("file", tt.args)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ParseMemberData", mock.Anything, mock.Anything).Return(tt.parseDataResult, nil)
			mockMemberUsecase.On("BulkValidateEmailAndPhone", mock.Anything, mock.Anything).Return(tt.resultInvalidRows)
			mockMemberUsecase.On("BulkImportMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))
			mockMemberUsecase.On(usecaseCheckEmail, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ImportMember(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerRegenerateToken(t *testing.T) {
	tests := []struct {
		name                               string
		token                              string
		wantUsecaseData1, wantUsecaseData2 usecase.ResultUseCase
		wantError                          bool
		wantStatusCode                     int
	}{
		{
			name:             testCasePositive1,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             testCaseNegative2,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:   http.StatusUnauthorized,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative4,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:   http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative5,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{Error: fmt.Errorf(msgErrorPq)},
			wantStatusCode:   http.StatusUnauthorized,
		},
		{
			name:             testCaseNegative6,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:   http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData1))
			mockMemberUsecase.On("RegenerateToken", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.RegenerateToken(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerAddNewMember(t *testing.T) {
	tests := []struct {
		name                               string
		token                              string
		wantUsecaseData1, wantUsecaseData2 usecase.ResultUseCase
		wantError                          bool
		wantStatusCode                     int
	}{
		{
			name:             testCasePositive1,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:   http.StatusCreated,
		},
		{
			name:             testCaseNegative2,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:   http.StatusBadRequest,
		},
		{
			name:             testCaseNegative3,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:             testCaseNegative4,
			token:            tokenAdmin,
			wantUsecaseData1: usecase.ResultUseCase{Result: model.Member{}},
			wantUsecaseData2: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:   http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(usecaseCheckEmail, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData1))
			mockMemberUsecase.On("RegisterMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData2))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.AddNewMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerUpdateMember(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          paramPhone,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:   testCaseNegative2,
			token:  tokenAdmin,
			param1: paramPhone,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			param1:          paramPhone,
			wantUsecaseData: usecase.ResultUseCase{Result: model.BlockedString},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("UpdateDetailMemberByID", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("mobile", tt.param1)
			data.Set(provID, "153")
			data.Set("provinceName", "DKI Jakarta")
			data.Set(provID, "153")
			data.Set("cityId", "2111")
			data.Set("cityName", "Jakarta Selatan")
			data.Set("districtId", "4765")
			data.Set("districtName", "Setia Budi")
			data.Set("subDistrictId", "4421")
			data.Set("subDistrictName", "Kuningan Timur")
			data.Set("postalCode", "12950")
			data.Set("street1", "address")

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.UpdateMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			// assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerGetDetailMember(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf(helper.ErrorDataNotFound, labelMember)},
			wantStatusCode:  http.StatusUnauthorized,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetDetailMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerGetMembers(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:  testCasePositive1,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListMembers{
				Members: []*model.Member{{}},
			}},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            "Testcase #2: Positive, empty member",
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListMembers{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GetListMembers", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetMembers(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerMigrateData(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		payload         interface{}
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			payload:         model.Members{},
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:           testCaseNegative2,
			token:          tokenAdmin,
			payload:        "inv",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:    testCaseNegative3,
			token:   tokenAdmin,
			payload: model.Members{},
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:    testCaseNegative4,
			token:   tokenAdmin,
			payload: model.Members{},
			wantUsecaseData: usecase.ResultUseCase{Error: fmt.Errorf(msgErrorPq), ErrorData: []model.MemberError{{
				ID:      "",
				Message: "",
			}}},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("MigrateMember", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err = handler.MigrateData(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerBulkMemberSend(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		args            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusCreated,
			args:            fileBase64Member,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
			args:           "",
		},
		{
			name:            testCaseNegative3,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.BlockedString},
			wantStatusCode:  http.StatusCreated,
			args:            fileBase64Member,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("file", tt.args)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMemberUsecase.On("PublishToKafkaUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.BulkMemberSend(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerMemberSend(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCasePositive2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:            testCaseNegative4,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On(getDetailMemberByID, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMemberUsecase.On("PublishToKafkaUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.MemberSend(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerValidateEmailDomain(t *testing.T) {
	tests := []struct {
		name            string
		email           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			email:           "test@gmail.com",
			wantUsecaseData: usecase.ResultUseCase{Result: model.SuccessResponse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			email: "test@getnada.com",
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: 500, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("email", tt.email)
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ValidateEmailDomain", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ValidateEmailDomain(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerChangeProfilePicture(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/profile_picture/photoqu-1574132451.jpeg",
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:   testCaseNegative2,
			token:  tokenAdmin,
			param1: "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/profile_picture/photoqu-1574132451.jpeg",
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("UpdateProfilePicture", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("file", tt.param1)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ChangeProfilePicture(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPMemberHandlerResendActivation(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          "testemail@email.co",
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ResendActivation", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("file", tt.param1)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ResendActivation(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPMemberHandlerGetMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFASettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GetMFASettings", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerGenerateMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFAGenerateSettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GenerateMFASettings", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GenerateMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPMemberHandlerActivateMFASettings(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFAActivateSettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ActivateMFASettings", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ActivateMFASettings(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerDisabledMFASetting(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.MFASettings{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("DisabledMFASetting", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.DisabledMFASetting(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

var testsNarwhal = []struct {
	name            string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	token           string
	wantStatusCode  int
}{
	{
		name:            testCasePositive1,
		token:           tokenAdmin,
		wantUsecaseData: usecase.ResultUseCase{Result: model.MFASettings{}},
		wantStatusCode:  http.StatusOK,
	},
	{
		name:  testCaseNegative2,
		token: tokenAdmin,
		wantUsecaseData: usecase.ResultUseCase{
			HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:  testCaseNegative3,
		token: tokenfailed,
		wantUsecaseData: usecase.ResultUseCase{
			HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
		},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		name:  testCaseNegative4,
		token: tokenNonAdmin,
		wantUsecaseData: usecase.ResultUseCase{
			HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
		},
		wantStatusCode: http.StatusUnauthorized,
	},
}

func TestHTTPMemberHandlerDisabledMFAMember(t *testing.T) {
	for _, tt := range testsNarwhal {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("DisabledMFASetting", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("memberID", "USR123")

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.DisabledMFAMember(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestDisabledMFANarwhalAdmin(t *testing.T) {
	for _, tt := range testsNarwhal {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("DisabledMFASetting", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("memberID", "USR123")

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.DisabledMFAMemberNarwhal(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandleRevokeAllAccess(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenNonAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("RevokeAllAccess", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.RevokeAllAccess(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandleGetLoginActivity(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.SessionInfoList{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenNonAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:            testCaseNegative5,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.Member{}},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GetLoginActivity", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetLoginActivity(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandleGetProfileComplete(t *testing.T) {
	tests := []struct {
		name            string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		token           string
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfileComplete{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenNonAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("GetProfileComplete", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.GetProfileComplete(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMemberHandlerChangeProfileName(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          "Test Nama",
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:   testCaseNegative2,
			token:  tokenAdmin,
			param1: "Test Nama Error",
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("UpdateProfileName", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			data := url.Values{}
			data.Set("name", tt.param1)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			err := handler.ChangeProfileName(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPMemberHandler_RevokeAccess(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		param1          string
	}{
		{
			name:            testCasePositive1,
			token:           tokenAdmin,
			param1:          "Test Nama",
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:   testCaseNegative2,
			token:  tokenAdmin,
			param1: "Test Nama Error",
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:  testCaseNegative3,
			token: tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  testCaseNegative4,
			token: tokenfailed,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			param1:         "",
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("RevokeAccess", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root+"?sid=1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			handler.RevokeAccess(c)
		})
	}
}

func TestHTTPMemberHandler_SyncPassword(t *testing.T) {
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
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
			payload:         Body{helper.TextToken: tokenAdmin},
		},
		{
			name:            testCaseNegative2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}, Error: errors.New("error")},
			wantStatusCode:  http.StatusOK,
			payload:         Body{helper.TextToken: tokenAdmin},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("SyncPassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			handler.SyncPassword(c)
		})
	}
}

func TestHTTPMemberHandler_ChangePassword(t *testing.T) {
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
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}},
			wantStatusCode:  http.StatusOK,
			payload:         Body{helper.TextToken: tokenAdmin},
		},
		{
			name:            testCaseNegative2,
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ProfilePicture{}, Error: errors.New("error")},
			wantStatusCode:  http.StatusOK,
			payload:         Body{helper.TextToken: tokenAdmin},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMemberUsecase := new(mocksMember.MemberUseCase)
			mockMemberUsecase.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateToken(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMemberUsecase)

			handler.ChangePassword(c)
		})
	}
}
