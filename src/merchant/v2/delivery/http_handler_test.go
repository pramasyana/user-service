package delivery

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/middleware"
	mocksMerchant "github.com/Bhinneka/user-service/mocks/src/merchant/v2/usecase"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/goleak"
)

const (
	root                     = "/api/v2/merchant/me"
	noAuth                   = `Basic`
	noAuthUser               = `Bearer`
	tokenAdmin               = `Basic Ymhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNTo2MjY4NjktNmU2ZTY1LTZiNjEyMC02ZDY1NmUtNzQ2MTcyLTY5MjA2NC02OTZkNjU2ZS03MzY5MDA=`
	tokenUser                = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	tokenUserNotAdmin        = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzMmIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6dHJ1ZSwiZGhUIjoiYzBiNGQxYjRjNDQ3NCIsImRsaSI6IldFQiIsImlhdCI6MTU0NDU0Mjk2MCwiaXNzIjoiYmhpbm5la2EuY29tIiwic3ViIjoiYmhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNSJ9.knglRfO7rHjalwTkG4ZLc-XJjWnijODFfiO2TByA37Y`
	tokenUserFailed          = `beyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`
	testCasePositive1        = "Testcase #1: Positive"
	testCasePositive2        = "Testcase #2: Positive"
	testCaseNegative2        = "Testcase #2: Negative"
	testCaseNegative3        = "Testcase #3: Negative"
	testCaseNegative4        = "Testcase #4: Negative"
	testCaseNegative5        = "Testcase #5: Negative"
	testCaseNegative6        = "Testcase #6: Negative"
	testCaseNegative7        = "Testcase #7: Negative"
	paramBusinessType        = "perorangan"
	paramMerchantName        = "merchant test"
	paramPhoneNumber         = "0812345678"
	paramMerchantDescription = "desc test"
	paramPic                 = "Test PIC"
	paramPicOccupation       = "test"
	paramMobilePhoneNumber   = "08123478959"
	paramAccountNumber       = "12384657887692"
	msgErrorPq               = "pq: error"
	fileBase64Merchant       = "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"
	fileBase64MerchantBlank  = "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"
	jsonSchemaMerchantDir    = "../../../../schema/"
	nameText                 = "name"
	labelText                = "label"
	testinglabel             = "gudang lama"
	testingName              = "Nauval"
	phoneText                = "phone"
	testingPhone             = "02188888"
	mobileTxt                = "mobile"
	testingMobile            = "085283318899"
	addressText              = "address"
	postalCodeText           = "postalCode"
	testingPostalCode        = "17111"
	subdistrictIDText        = "subDistrictId"
	testingSubdistrictID     = "0104040501"
	subdistrictNameText      = "subDistrictName"
	testingSubdistrictName   = "Aren Jaya"
	districtIDText           = "districtId"
	testingDistrictID        = "01040405"
	districtNameText         = "districtName"
	testingDistrictName      = "Bekasi Timur"
	cityIDText               = "cityId"
	testingCityID            = "010404"
	cityNameText             = "cityName"
	testingCityName          = "Bekasi"
	provinceIDText           = "provinceId"
	testingProvinceID        = "0104"
	provinceNameText         = "provinceName"
	testingProvinceName      = "Jawa Barat"
	pqError                  = "pq: error"
	merchantWithPath         = "/api/v2/merchant/MCH201210161001"
	merchantEmployeeWithPath = "/api/v2/merchant/me/employees/USR220410359750525"
	getMerchantPath          = "/api/v2/merchant/list"
	badMerchantPath          = "/api/v2/merchant/list?page=-1&limit=-1"
	bodyUpdate               = `{"merchantName":"some name", "merchantAddress":"Alamat update","merchantType":"REGULAR","merchantDescription":"some description","merchantGroup":"MICRO","upgradeStatus":"status","merchantLogo":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg","additionalEmail":"email@domain.com","businessType":"perorangan","phoneNumber":"0219099","mobilePhoneNumber":"081279090999","companyName":"some company","dailyOperationalStaff":"some value","pic":"MY STAFF","picKtpFile":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg","picOccupation":"someone","genderPic":"FEMALE","storeClosureDate":"2021-12-12 12:12:12","storeReopenDate":"2021-12-15 12:12:12","merchantCategory":"merchantCategory","npwp":"123456789012345","npwpFile":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg","npwpHolderName":"NPWP holder name","merchantCityId":"city ID","merchantDistrictId":"district ID","merchantVillageId":"village ID","zipCode":"12345","storeProvinceId":"province ID","storeCityId":"city ID","storeDistrictId":"district ID","storeActiveShippingDate":"2020-09-10","storeVillageId":"village ID","storeZipCode":"12121","storeIsClosed":false,"storeAddress":"Jalan jalan kemana aja","isPKP":true,"isActive":true,"createdIp":"127.0.0.1","updatedIp":"127.0.0.1","accountNumber":"982120999012","accountHolderNumber":"09090990","accountHolderName":"My Name","bankId":12,"bankBranch":"Cabang Mana","productType":"PHYSIC","documents":[{"documentType":"AktaPendirian-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"},{"documentType":"AktaPerubahan-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"}],"richContent":"", "merchantEmail":"email@domain.com"}`
	rejectUpgradePath        = "/api/v2/merchant/MCH201210161001/reject-upgrade"
)

var (
	errDefault = fmt.Errorf("default error")
)
var payloadMerchantAddress = model.WarehouseData{
	Label:           testinglabel,
	Address:         "ini adalah alamat baris 1",
	Name:            testingName,
	Phone:           testingPhone,
	Mobile:          testingMobile,
	PostalCode:      testingPostalCode,
	SubDistrictID:   testingSubdistrictID,
	SubDistrictName: testingSubdistrictName,
	DistrictID:      testingDistrictID,
	DistrictName:    testingDistrictName,
	WarehousePrimary: model.WarehousePrimary{
		CityID:       testingCityID,
		CityName:     testingCityName,
		ProvinceID:   testingProvinceID,
		ProvinceName: testingProvinceName,
	},
}
var testsAddUpdateMerchantAddressMe = []struct {
	name            string
	token           string
	wantUsecaseData usecase.ResultUseCase
	wantError       bool
	wantStatusCode  int
	address         string
	payload         interface{}
}{
	{
		name:            testCasePositive1,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.WarehouseData{}},
		wantStatusCode:  http.StatusCreated,
		payload:         payloadMerchantAddress,
	},
	{
		name:            testCaseNegative2,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadMerchantAddress,
	},
	{
		name:            testCaseNegative3,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadMerchantAddress,
	},
	{
		name:            testCaseNegative4,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
		wantStatusCode:  http.StatusBadRequest,
		payload: model.WarehouseData{
			Label:           testinglabel,
			Name:            testingName,
			Phone:           testingPhone,
			Mobile:          testingMobile,
			PostalCode:      testingPostalCode,
			SubDistrictID:   testingSubdistrictID,
			SubDistrictName: testingSubdistrictName,
			DistrictID:      testingDistrictID,
			DistrictName:    testingDistrictName,
			WarehousePrimary: model.WarehousePrimary{
				CityID:       testingCityID,
				CityName:     testingCityName,
				ProvinceID:   testingProvinceID,
				ProvinceName: testingProvinceName,
			},
		},
	},
	{
		name:            testCaseNegative5,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{Result: model.MerchantError{}},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadMerchantAddress,
	},
	{
		name:            testCaseNegative6,
		token:           tokenUserFailed,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
		wantStatusCode:  http.StatusBadRequest,
		payload:         payloadMerchantAddress,
	},
	{
		name:            testCaseNegative7,
		token:           tokenUser,
		wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
		wantStatusCode:  http.StatusBadRequest,
		payload:         model.B2CMerchantBank{},
	},
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func generateRSAMerchant() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateTokenMerchant(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSAMerchant(), nil
	})
}

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestHTTPMerchantHandlerMount(*testing.T) {
	e := echo.New()
	handler := NewHTTPHandler(new(mocksMerchant.MerchantUseCase), new(mocksMerchant.MerchantAddressUseCase))
	handler.MountMe(e.Group("/anon"))
	handler.MountMerchant(e.Group("/basic"))
	handler.MountCMS(e.Group("/basic"))
	handler.MountMerchantPublic(e.Group("/api/v2/merchant"))
}

func TestHTTPMerchantHandlerGetMerchantBank(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListMerchantBank{
				MerchantBank: []*model.B2CMerchantBankData{{}},
			}},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            "Testcase #2: Positive, empty merchant bank",
			token:           tokenAdmin,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListMerchantBank{}},
			wantStatusCode:  http.StatusOK,
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
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantUsecase.On("GetListMerchantBank", mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("Authorization", tt.token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err := handler.GetMerchantBank(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerAddMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusCreated,
			payload: model.B2CMerchantCreateInput{
				BusinessType:        paramBusinessType,
				MerchantName:        "merchant test",
				PhoneNumber:         "0812345678",
				MerchantDescription: "desc test",
				Pic:                 "Test PIC",
				PicOccupation:       "test",
				MobilePhoneNumber:   "08123478959",
				AccountNumber:       "12384657887692",
				MerchantAddress:     "some string",
			},
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchantCreateInput{},
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchant{},
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload: model.B2CMerchantCreateInput{
				BusinessType:        paramBusinessType,
				MerchantName:        paramMerchantName,
				PhoneNumber:         paramPhoneNumber,
				MerchantDescription: paramMerchantDescription,
				Pic:                 paramPic,
				PicOccupation:       paramPicOccupation,
				MobilePhoneNumber:   "2",
				AccountNumber:       paramAccountNumber,
			},
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload: model.B2CMerchantCreateInput{
				BusinessType:        paramBusinessType,
				MerchantName:        paramMerchantName,
				PhoneNumber:         paramPhoneNumber,
				MerchantDescription: paramMerchantDescription,
				Pic:                 paramPic,
				PicOccupation:       paramPicOccupation,
				MobilePhoneNumber:   paramMobilePhoneNumber,
				AccountNumber:       paramAccountNumber,
			},
		},
		{
			name:            testCaseNegative6,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantBank{}},
			wantStatusCode:  http.StatusBadRequest,
			payload: model.B2CMerchantCreateInput{
				BusinessType:        paramBusinessType,
				MerchantName:        paramMerchantName,
				PhoneNumber:         paramPhoneNumber,
				MerchantDescription: paramMerchantDescription,
				Pic:                 paramPic,
				PicOccupation:       paramPicOccupation,
				MobilePhoneNumber:   paramMobilePhoneNumber,
				AccountNumber:       paramAccountNumber,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("AddMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err = handler.AddMerchant(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerUpgradeMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			payload:         model.B2CMerchantCreateInput{},
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchantCreateInput{},
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchant{},
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchantCreateInput{},
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantBank{}},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.B2CMerchantCreateInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("UpgradeMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err = handler.UpgradeMerchant(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerCheckMerchantName(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.CheckMerchantName{}},
			wantStatusCode:  http.StatusOK,
			payload:         model.CheckMerchantName{},
		},
		{
			name:            testCaseNegative2,
			token:           noAuthUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         model.CheckMerchantName{},
		},
		{
			name:            testCaseNegative3,
			token:           noAuthUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         "asd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("CheckMerchantName", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", tt.token)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err = handler.CheckMerchantName(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerGetMerchantByUserID(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:  testCasePositive1,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ResponseAvailable{
				MerchantServiceAvailable: true,
				MerchantData:             model.B2CMerchantDataV2{},
			}},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{Result: "inv"},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:  testCaseNegative3,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusNotFound, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantUsecase.On("GetMerchantByUserID", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err := handler.GetMerchantByUserID(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerMerchantSend(t *testing.T) {
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
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCaseNegative2,
			token: tokenUserNotAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusUnauthorized,
			args:           "",
		},
		{
			name:  testCaseNegative3,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantBank{}},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("merchantID", "123")
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantUsecase.On("GetMerchantByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMerchantUsecase.On("PublishToKafkaMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err := handler.MerchantSend(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
func TestHTTPMerchantHandlerBulkMerchantSend(t *testing.T) {
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
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			args:            fileBase64Merchant,
		},
		{
			name:  testCaseNegative2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusInternalServerError, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusBadRequest,
			args:           "",
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantBank{}},
			wantStatusCode:  http.StatusOK,
			args:            fileBase64Merchant,
		},
		{
			name:  testCaseNegative4,
			token: tokenUserNotAdmin,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusUnauthorized, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusUnauthorized,
			args:           "",
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantBank{}},
			wantStatusCode:  http.StatusBadRequest,
			args:            "adasda",
		},
		{
			name:  testCaseNegative6,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{
				HTTPStatus: http.StatusBadRequest, Error: fmt.Errorf(msgErrorPq),
			},
			wantStatusCode: http.StatusOK,
			args:           fileBase64Merchant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := url.Values{}
			data.Set("file", tt.args)
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantUsecase.On("GetMerchantByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			mockMerchantUsecase.On("PublishToKafkaMerchant", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err := handler.BulkMerchantSend(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerAddMerchantAddressMe(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)

	for _, tt := range testsAddUpdateMerchantAddressMe {

		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantAddressUsecase.On("AddUpdateWarehouseAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err = handler.AddWarehouse(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})

	}
}

func TestHTTPMerchantHandlerUpdateMerchantAddressMe(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)

	for _, tt := range testsAddUpdateMerchantAddressMe {
		if tt.wantStatusCode == http.StatusCreated {
			tt.wantStatusCode = http.StatusOK
		}
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockMerchantAddressUsecase.On("AddUpdateWarehouseAddress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			bodyData, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(string(bodyData)))
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockMerchantAddressUsecase)

			err = handler.UpdateWarehouse(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})

	}
}

func TestHTTPMerchantHandlerUpdatePrimaryWarehouseAddress(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
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
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		data := url.Values{}

		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockWarehouseAddressUsecase.On("UpdatePrimaryWarehouseAddress", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, bytes.NewBufferString(data.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.UpdateWarehousePrimary(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerGetWarehouse(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:  testCasePositive2,
			token: tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.ListWarehouse{
				WarehouseData: []*model.WarehouseData{{}},
			}},
			wantStatusCode: http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
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
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockWarehouseAddressUsecase.On("GetWarehouseAddresses", mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.GetWarehouse(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandlerGetDetailWarehouse(t *testing.T) {
	testGetWarehouse := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.WarehouseData{}},
			wantStatusCode:  http.StatusOK,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
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
	for _, tt := range testGetWarehouse {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockWarehouseAddressUsecase.On("GetDetailWarehouseAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.GetWarehouseDetail(c)
			if tt.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)

		})

	}
}

func TestHTTPMerchantHandlerDeleteWarehouse(t *testing.T) {
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
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockMerchantUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockWarehouseAddressUsecase.On("DeleteWarehouseAddress", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockMerchantUsecase, mockWarehouseAddressUsecase)

			err := handler.DeleteWarehouse(c)
			if tt.wantError {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

const (
	defUpdateJSON         = `{"companyName":"PT Merah","pic":"sinta adiyasa Update","picOccupation":"lalalala","picKtpFile":"https://bmd-upload.s3.ap-southeast-1.amazonaws.com/merchant/ktp.jpg","dailyOperationalStaff":"daily ops","storeClosureDate":null,"storeReopenDate":null,"storeActiveShippingDate":"2019-01-02","merchantAddress":"merchant address update lagi","merchantVillage":"","merchantVillageId":"village ID","merchantDistrict":"","merchantDistrictId":"","merchantCity":"","merchantCityId":"city ID","merchantProvince":"","merchantProvinceId":"","zipCode":"12345","storeAddress":"Jalan jalan kemana aja","storeVillage":"","storeVillageId":"village ID","storeDistrict":"","storeDistrictId":"","storeCity":"","storeCityId":"","storeProvince":"","storeProvinceId":"","storeZipCode":"12121","phoneNumber":"34124124124124","mobilePhoneNumber":"08574676762","additionalEmail":"email@domain.com","merchantDescription":"merchant description put","merchantLogo":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg","accountHolderName":"bapak budi update put","bankId":12,"bankCode":"","bankName":"test307","bankBranch":"Jakarta","accountNumber":"123232343242","isPKP":true,"npwp":"123456789056789","npwpHolderName":"NPWP sinta update","npwpFile":"https://bmd-upload.s3.ap-southeast-1.amazonaws.com/merchant/npwp.jpg","richContent":"","notificationPreferences":0,"merchantRank":"","acquisitor":"","accountManager":"","launchDev":"","skuLive":null,"mouDate":null,"note":"","businessType":"perorangan","agreementDate":null,"defaultBanner":"","isActive":false,"isClosed":false,"merchantType":"ASSOCIATE","genderPic":"FEMALE","merchantGroup":"MICRO","upgradeStatus":"PENDING_MANAGE","documents":[{"id":"DOC210203152535506","documentType":"AktaPendirian-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"},{"id":"DOC210203152535655","documentType":"AktaPerubahan-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"}],"productType":"NON_PHYSIC","legalEntity":10,"numberOfEmployee":3}`
	defChangeNameJSON     = `{"merchantName":"Testing"}`
	badUpdateJSON         = `{"companyName":"PT Merah","pic":"sinta adiyasa Update","picOccupation":"lalalala","picKtpFile":"https://bmd-upload.s3.ap-southeast-1.amazonaws.com/merchant/ktp.jpg","dailyOperationalStaff":"daily ops","storeClosureDate":null,"storeReopenDate":null,"storeActiveShippingDate":"2019-01-02","merchantAddress":"merchant address update lagi","merchantVillage":"","merchantVillageId":"village ID","merchantDistrict":"","merchantDistrictId":"","merchantCity":"","merchantCityId":"city ID","merchantProvince":"","merchantProvinceId":"","zipCode":"12345","storeAddress":"Jalan jalan kemana aja","storeVillage":"","storeVillageId":"village ID","storeDistrict":"","storeDistrictId":"","storeCity":"","storeCityId":"","storeProvince":"","storeProvinceId":"","storeZipCode":"12121","phoneNumber":"34124124124124","mobilePhoneNumber":"0857467","additionalEmail":"email@domain.com","merchantDescription":"merchant description put","merchantLogo":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg","accountHolderName":"bapak budi update put","bankId":12,"bankCode":"","bankName":"test307","bankBranch":"Jakarta","accountNumber":"123232343242","isPKP":true,"npwp":"123456789056789","npwpHolderName":"NPWP sinta update","npwpFile":"https://bmd-upload.s3.ap-southeast-1.amazonaws.com/merchant/npwp.jpg","richContent":"","notificationPreferences":0,"merchantRank":"","acquisitor":"","accountManager":"","launchDev":"","skuLive":null,"mouDate":null,"note":"","businessType":"perorangan","agreementDate":null,"defaultBanner":"","isActive":false,"isClosed":false,"merchantType":"ASSOCIATE","genderPic":"FEMALE","merchantGroup":"MICRO","upgradeStatus":"PENDING_MANAGE","documents":[{"id":"DOC210203152535506","documentType":"AktaPendirian-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"},{"id":"DOC210203152535655","documentType":"AktaPerubahan-file","documentValue":"https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/npwp-file-1574132451.jpeg"}],"productType":"NON_PHYSIC","legalEntity":10,"numberOfEmployee":3}`
	defAddEmployeeJSON    = `{"email":"email@getnnada.com","firstName":"Test"}`
	defBadAddEmployeeJSON = `{"email":"","firstName":""}`
	defBadParam           = `{"page":"csdcs"}`
)

func TestHandlerSelfUpdateMerchant(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			payload:         defUpdateJSON,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         `{,}`,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         badUpdateJSON,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defUpdateJSON,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defUpdateJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("SelfUpdateMerchant", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.PUT, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.UpdateMerchant(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_UpdateMerchantPartial(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			payload:         defUpdateJSON,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         `{,}`,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defUpdateJSON,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defUpdateJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("SelfUpdateMerchantPartial", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.PATCH, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.UpdateMerchantPartial(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_ChangeMerchantName(t *testing.T) {
	jsonschema.Load(jsonSchemaMerchantDir)
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantDataV2{}},
			wantStatusCode:  http.StatusOK,
			payload:         defChangeNameJSON,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUserFailed,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusUnauthorized, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         `{.}`,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         badUpdateJSON,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defChangeNameJSON,
		},
		{
			name:            testCaseNegative5,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defChangeNameJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("ChangeMerchantName", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.PUT, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.ChangeMerchantName(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_AddEmployee(t *testing.T) {

	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "success"},
			wantStatusCode:  http.StatusOK,
			payload:         defAddEmployeeJSON,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         `{:}`,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defBadAddEmployeeJSON,
		},
		{
			name:            testCaseNegative4,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defAddEmployeeJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("AddEmployee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.AddEmployee(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_ListEmployee(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: "anything", TotalData: int(10)},
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: nil},
			wantStatusCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantEmployeeUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)

			mockMerchantEmployeeUsecase.On("GetAllMerchantEmployee", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantEmployeeUsecase, mockWarehouseAddressUsecase)

			handler.ListEmployee(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_UpdateEmployee(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
		paramMemberId   string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "success"},
			wantStatusCode:  http.StatusOK,
			payload:         defAddEmployeeJSON,
			paramMemberId:   "USR13124324",
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         "{?}",
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defAddEmployeeJSON,
			paramMemberId:   "",
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defAddEmployeeJSON,
			paramMemberId:   "USR444342423",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("UpdateMerchantEmployee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.PUT, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			c.SetParamNames("memberId")
			c.SetParamValues(tt.paramMemberId)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.UpdateEmployee(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_ResendEmailEmployee(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantError       bool
		wantStatusCode  int
		payload         string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{Result: "success"},
			wantStatusCode:  http.StatusOK,
			payload:         defAddEmployeeJSON,
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         "{>}",
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			wantStatusCode:  http.StatusBadRequest,
			payload:         defAddEmployeeJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthUseCase := new(mocksMerchant.MerchantUseCase)
			mockMerchantAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)
			mockAuthUseCase.On("AddEmployee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))

			e := echo.New()
			req := httptest.NewRequest(echo.POST, root, strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			handler := NewHTTPHandler(mockAuthUseCase, mockMerchantAddressUsecase)

			err := handler.ResendEmailEmployee(c)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestHTTPMerchantHandler_GetEmployee(t *testing.T) {
	testData := []struct {
		name            string
		token           string
		wantUsecaseData usecase.ResultUseCase
		wantStatusCode  int
		paramMemberId   string
	}{
		{
			name:            testCasePositive1,
			token:           tokenUser,
			wantStatusCode:  http.StatusOK,
			wantUsecaseData: usecase.ResultUseCase{Result: model.B2CMerchantEmployee{}},
			paramMemberId:   "USR123123424",
		},
		{
			name:            testCaseNegative2,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			paramMemberId:   "",
		},
		{
			name:            testCaseNegative3,
			token:           tokenUser,
			wantStatusCode:  http.StatusBadRequest,
			wantUsecaseData: usecase.ResultUseCase{HTTPStatus: http.StatusBadRequest, Error: errDefault},
			paramMemberId:   "USR543536",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			mockMerchantEmployeeUsecase := new(mocksMerchant.MerchantUseCase)
			mockWarehouseAddressUsecase := new(mocksMerchant.MerchantAddressUseCase)

			e := echo.New()
			req := httptest.NewRequest(echo.GET, root, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			token, _ := generateTokenMerchant(tt.token)
			c.Set("token", token)
			c.SetParamNames("memberId")
			c.SetParamValues(tt.paramMemberId)
			mockMerchantEmployeeUsecase.On("GetMerchantEmployee", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(tt.wantUsecaseData))
			handler := NewHTTPHandler(mockMerchantEmployeeUsecase, mockWarehouseAddressUsecase)

			handler.GetEmployee(c)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
