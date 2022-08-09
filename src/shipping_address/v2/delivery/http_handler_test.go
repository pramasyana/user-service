package delivery

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/usecase"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/usecase/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root                   = "/api/v2/shipping-address"
	tokenUser              = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	tokenUserFailedID      = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20ifQ.stRqFGMoWfuqMQA666SmU9lRKkoEgmUZ5pe84yYWdiU`
	tokenUserFailed        = `beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	jsonSchemaDir          = "../../../../schema/"
	testCasePositive1      = "Testcase #1: Positive"
	testCasePositive2      = "Testcase #2: Positive"
	testCaseNegative2      = "Testcase #2: Negative"
	testCaseNegative3      = "Testcase #3: Negative"
	testCaseNegative4      = "Testcase #4: Negative"
	testCaseNegative5      = "Testcase #5: Negative"
	testCaseNegative6      = "Testcase #6: Negative"
	testingText            = "testing"
	pqError                = "pq: error"
	memberID               = "memberId"
	extText                = "ext"
	testingExt             = "021"
	labelText              = "label"
	testingLabel           = "testing1"
	testingMemberID        = "234234"
	nameText               = "name"
	testingName            = "Nauval"
	phoneText              = "phone"
	testingPhone           = "02188888"
	mobileTxt              = "mobile"
	testingMobile          = "085283318899"
	street1Text            = "street1"
	testingStreet1         = "alamat baris 1 yaa"
	street2Text            = "street2"
	testingStreet2         = "alamat baris 2 yaa"
	postalCodeText         = "postalCode"
	testingPostalCode      = "17111"
	subdistrictIDText      = "subDistrictId"
	testingSubdistrictID   = "0104040501"
	subdistrictNameText    = "subDistrictName"
	testingSubdistrictName = "Aren Jaya"
	districtIDText         = "districtId"
	testingDistrictID      = "01040405"
	districtNameText       = "districtName"
	testingDistrictName    = "Bekasi Timur"
	cityIDText             = "cityId"
	testingCityID          = "010404"
	cityNameText           = "cityName"
	testingCityName        = "Bekasi"
	provinceIDText         = "provinceId"
	testingProvinceID      = "0104"
	provinceNameText       = "provinceName"
	testingProvinceName    = "Jawa Barat"
)

var testsAddUpdateShipping = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
	street1         string
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressData{}},
		wantStatusCode:  http.StatusCreated,
		street1:         testingStreet1,
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusUnauthorized,
		street1:         testingStreet2,
	},
	{
		name:            testCaseNegative3,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "ini adalah alamat baris 3",
	},
	{
		name:            testCaseNegative4,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "",
	},
	{
		name:            testCaseNegative5,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressError{}},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "ini adalah alamat baris 4",
	},
	{
		name:            testCaseNegative6,
		token:           tokenUserFailedID,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
	},
}

var testsAddUpdateShippingMe = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
	street1         string
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressData{}},
		wantStatusCode:  http.StatusCreated,
		street1:         "ini adalah alamat baris 1",
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "ini adalah alamat baris 2",
	},
	{
		name:            testCaseNegative3,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "ini adalah alamat baris 3",
	},
	{
		name:            testCaseNegative4,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "",
	},
	{
		name:            testCaseNegative5,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressError{}},
		wantStatusCode:  http.StatusBadRequest,
		street1:         "ini adalah alamat baris 4",
	},
	{
		name:            testCaseNegative6,
		token:           tokenUserFailedID,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
		wantStatusCode:  http.StatusBadRequest,
	},
}

var testsGetShipping = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressData{}},
		wantStatusCode:  http.StatusOK,
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

var testsGetShippingMe = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.ShippingAddressData{}},
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

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateUsecaseResultShipping(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateRSAShipping() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
   	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateTokenShipping(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSAShipping(), nil
	})
}

func TestHTTPShippingAddressHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocks.ShippingAddressUseCase))
	handler.MountMe(e.Group("/anon"))
	handler.MountShippingAddress(e.Group("/anon"))
}

func TestHTTPShippingHandlerAddShippingAddress(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range testsAddUpdateShipping {

		data := url.Values{}
		data.Set(memberID, "123")
		data.Set(nameText, testingName)
		data.Set(phoneText, testingPhone)
		data.Set(mobileTxt, testingMobile)
		data.Set(street1Text, tt.street1)
		data.Set(street2Text, testingStreet2)
		data.Set(postalCodeText, testingPostalCode)
		data.Set(subdistrictIDText, testingSubdistrictID)
		data.Set(subdistrictNameText, testingSubdistrictName)
		data.Set(districtIDText, testingDistrictID)
		data.Set(districtNameText, testingDistrictName)
		data.Set(cityIDText, testingCityID)
		data.Set(cityNameText, testingCityName)
		data.Set(provinceIDText, testingProvinceID)
		data.Set(provinceNameText, testingProvinceName)
		data.Set(extText, testingExt)
		data.Set(labelText, "testing1")

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("AddShippingAddress", mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.AddShippingAddress(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerAddShippingAddressMe(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range testsAddUpdateShippingMe {

		data := url.Values{}
		data.Set(memberID, "123")
		data.Set(nameText, testingName)
		data.Set(phoneText, testingPhone)
		data.Set(mobileTxt, testingMobile)
		data.Set(street1Text, tt.street1)
		data.Set(street2Text, testingStreet2)
		data.Set(postalCodeText, testingPostalCode)
		data.Set(subdistrictIDText, testingSubdistrictID)
		data.Set(subdistrictNameText, testingSubdistrictName)
		data.Set(districtIDText, testingDistrictID)
		data.Set(districtNameText, testingDistrictName)
		data.Set(cityIDText, testingCityID)
		data.Set(cityNameText, testingCityName)
		data.Set(provinceIDText, testingProvinceID)
		data.Set(provinceNameText, testingProvinceName)
		data.Set(extText, testingExt)
		data.Set(labelText, testingLabel)

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("AddShippingAddress", mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.AddShippingAddressMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})

	}
}

func TestHTTPShippingHandlerDeleteShippingAddress(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusUnauthorized,
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
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("DeleteShippingAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.DeleteShippingAddress(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerDeleteShippingAddressMe(t *testing.T) {
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
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("DeleteShippingAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.DeleteShippingAddressMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerUpdateShippingAddress(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range testsAddUpdateShipping {
		if tt.wantStatusCode == http.StatusCreated {
			tt.wantStatusCode = http.StatusOK
		}
		data := url.Values{}
		data.Set("id", testingMemberID)
		data.Set(nameText, testingName)
		data.Set(phoneText, testingPhone)
		data.Set(mobileTxt, testingMobile)
		data.Set(street1Text, tt.street1)
		data.Set(street2Text, testingStreet2)
		data.Set(postalCodeText, testingPostalCode)
		data.Set(subdistrictIDText, testingSubdistrictID)
		data.Set(subdistrictNameText, testingSubdistrictName)
		data.Set(districtIDText, testingDistrictID)
		data.Set(districtNameText, testingDistrictName)
		data.Set(cityIDText, testingCityID)
		data.Set(cityNameText, testingCityName)
		data.Set(provinceIDText, testingProvinceID)
		data.Set(provinceNameText, testingProvinceName)
		data.Set(extText, testingExt)
		data.Set(labelText, testingLabel)

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("UpdateShippingAddress", mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.UpdateShippingAddress(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerUpdateShippingAddressMe(t *testing.T) {
	jsonschema.Load(jsonSchemaDir)

	for _, tt := range testsAddUpdateShippingMe {
		if tt.wantStatusCode == http.StatusCreated {
			tt.wantStatusCode = http.StatusOK
		}
		data := url.Values{}
		data.Set("id", testingMemberID)
		data.Set(nameText, testingName)
		data.Set(phoneText, testingPhone)
		data.Set(mobileTxt, testingMobile)
		data.Set(street1Text, tt.street1)
		data.Set(street2Text, testingStreet2)
		data.Set(postalCodeText, testingPostalCode)
		data.Set(subdistrictIDText, testingSubdistrictID)
		data.Set(subdistrictNameText, testingSubdistrictName)
		data.Set(districtIDText, testingDistrictID)
		data.Set(districtNameText, testingDistrictName)
		data.Set(cityIDText, testingCityID)
		data.Set(cityNameText, testingCityName)
		data.Set(provinceIDText, testingProvinceID)
		data.Set(provinceNameText, testingProvinceName)
		data.Set(extText, testingExt)
		data.Set(labelText, testingLabel)

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("UpdateShippingAddress", mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.UpdateShippingAddressMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerGetShippingAddress(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListShippingAddress{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListShippingAddress{
				ShippingAddress: []*model.ShippingAddressData{{}},
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
		{
			name:            testCaseNegative5,
			token:           tokenUserFailedID,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("GetListShippingAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))
			mockShippingAddressUsecase2 := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase2.On("GetAllListShippingAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.GetShippingAddressMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

			handler = NewHTTPHandler(mockShippingAddressUsecase2)
			err = handler.GetShippingAddress(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerGetDetailShippingAddressMe(t *testing.T) {
	for _, tt := range testsGetShippingMe {
		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("GetDetailShippingAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.GetShippingAddressDetailMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})

	}
}

func TestHTTPShippingHandlerGetDetailShippingAddress(t *testing.T) {
	for _, tt := range testsGetShipping {
		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("GetDetailShippingAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.GetShippingAddressDetail(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}

func TestHTTPShippingHandlerUpdatePrimaryShippingAddress(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusUnauthorized,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUserFailedID,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf("error")},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		data := url.Values{}
		data.Set(memberID, testingMemberID)

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("UpdatePrimaryShippingAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.UpdateIsPrimary(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerUpdatePrimaryShippingAddressMe(t *testing.T) {
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

		data := url.Values{}
		data.Set(memberID, testingMemberID)

		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("UpdatePrimaryShippingAddressByID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.UpdateIsPrimaryMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPShippingHandlerGetPrimaryShippingAddressMe(t *testing.T) {

	for _, tt := range testsGetShippingMe {
		t.Run(tt.name, func(t *testing.T) {
			mockShippingAddressUsecase := new(mocks.ShippingAddressUseCase)
			mockShippingAddressUsecase.On("GetPrimaryShippingAddress", mock.Anything, mock.Anything).Return(generateUsecaseResultShipping(tt.wantUsecaseData))

			e := echo.New()
			req, err := http.NewRequest(echo.GET, root, nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenShipping(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockShippingAddressUsecase)

			err = handler.GetShippingAddressPrimaryMe(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})
	}
}
