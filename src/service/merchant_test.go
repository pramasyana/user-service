package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Bhinneka/bhinneka-go-sdk"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/service/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"
)

const (
	defaultSalmonURL  = "http://merchant.bhinnekatesting.com"
	defaultUserID     = "USR09"
	defaultMerchantID = "MCH079"
	defaultToken      = "Bearer token"
	constMerchantURL  = "MERCHANT_SERVICE_URL"
	constGql          = "MERCHANT_SERVICE_GRAPHQL_URL"
)

var (
	defaultGQLURL = fmt.Sprintf("%s/graphql", defaultSalmonURL)
)

func TestInitMerchantService(t *testing.T) {
	mockPub := new(mocks.QPublisher)
	mockActivity := new(mocks.ActivityServices)
	_, err := NewMerchantService(mockPub, mockActivity)
	assert.Error(t, err)

	os.Setenv(constMerchantURL, badURL)
	_, err = NewMerchantService(mockPub, mockActivity)
	assert.Error(t, err)

	os.Setenv(constGql, badURL)
	_, err = NewMerchantService(mockPub, mockActivity)
	assert.Error(t, err)

	os.Setenv(constMerchantURL, defaultSalmonURL)
	_, err = NewMerchantService(mockPub, mockActivity)
	assert.Error(t, err)

	os.Setenv(constGql, defaultGQLURL)
	_, err = NewMerchantService(mockPub, mockActivity)
	assert.NoError(t, err)
}

type basicType struct {
	name            string
	wantError       bool
	statusCode      int
	serviceResponse interface{}
}

var testDataMerchantLegacy = []basicType{
	{
		name:            "FindMerchantServiceByID Legacy Test #1",
		statusCode:      http.StatusOK,
		serviceResponse: serviceModel.ResponseMerchantService{},
		wantError:       false,
	},
	{
		name:            "FindMerchantServiceByID Legacy Test #2",
		statusCode:      http.StatusOK,
		serviceResponse: []byte(`1`),
		wantError:       true,
	},
}

var testDataMerchantGraphQL = []basicType{
	{
		name:            "Merchant Test GraphQL #1",
		statusCode:      http.StatusOK,
		serviceResponse: serviceModel.ResponseGWSMerchant{Code: http.StatusOK, Data: serviceModel.MerchantDetailGWS{GetMerchantDetail: serviceModel.Merchantdetail{Success: true, Code: http.StatusOK}}},
		wantError:       false,
	},
	{
		name:            "Merchant Test GraphQL #2",
		statusCode:      http.StatusOK,
		serviceResponse: []byte(`1`),
		wantError:       true,
	},
	{
		name:            "Merchant Test GraphQL #3",
		statusCode:      http.StatusBadRequest,
		serviceResponse: serviceModel.ResponseGWSMerchant{Code: http.StatusOK, Data: serviceModel.MerchantDetailGWS{GetMerchantDetail: serviceModel.Merchantdetail{Success: false}}},
		wantError:       true,
	},
}

func TestFailedURL(t *testing.T) {
	mockPub := new(mocks.QPublisher)
	mockActivity := new(mocks.ActivityServices)
	os.Setenv(constMerchantURL, "https:////some-bad-url.@")
	os.Setenv(constGql, "https:///some-bad-url.com")
	ms, _ := NewMerchantService(mockPub, mockActivity)
	ctx := context.Background()
	sr := <-ms.FindMerchantServiceByID(ctx, defaultUserID, defaultToken, defaultMerchantID)
	assert.Error(t, sr.Error)
}

func TestFindMerchantServiceByID(t *testing.T) {
	mockPub := new(mocks.QPublisher)
	mockActivity := new(mocks.ActivityServices)
	os.Setenv(constMerchantURL, defaultSalmonURL)
	os.Setenv(constGql, defaultGQLURL)
	ms, _ := NewMerchantService(mockPub, mockActivity)
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testDataMerchantLegacy {
		bhinneka.MockHTTP(http.MethodGet, fmt.Sprintf("%s/api/v1/merchant/me", defaultSalmonURL), tc.statusCode, tc.serviceResponse)
		sr := <-ms.FindMerchantServiceByID(ctx, defaultUserID, defaultToken, defaultMerchantID)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
	os.Setenv("GWS_MERCHANT_ACTIVE", "true")
	for _, tc := range testDataMerchantGraphQL {
		bhinneka.MockHTTP(http.MethodPost, fmt.Sprintf("%s/graphql", defaultGQLURL), tc.statusCode, tc.serviceResponse)
		sr := <-ms.FindMerchantServiceByID(ctx, defaultUserID, defaultToken, defaultMerchantID)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

func TestPublishToKafka(t *testing.T) {
	os.Setenv(constMerchantURL, defaultSalmonURL)
	os.Setenv(constGql, defaultGQLURL)
	os.Setenv("KAFKA_USER_SERVICE_MERCHANT_TOPIC", "anything")
	ctx := context.Background()

	var testData = []struct {
		name            string
		usecaseResponse interface{}
		wantError       bool
	}{
		{
			name:            "Test Publish #2",
			wantError:       true,
			usecaseResponse: fmt.Errorf("error from kafka"),
		},
		{
			name:            "Test Publish #1",
			wantError:       false,
			usecaseResponse: nil,
		},
	}

	for _, tc := range testData {
		mockPub := new(mocks.QPublisher)
		mockActivity := new(mocks.ActivityServices)
		ms, _ := NewMerchantService(mockPub, mockActivity)

		mockPub.On("PublishKafka", ctx, mock.Anything, mock.Anything, mock.Anything).Return(tc.usecaseResponse)
		err := ms.PublishToKafkaUserMerchant(ctx, &merchantModel.B2CMerchantDataV2{}, "UpdateMerchant", "gws")
		if tc.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestInsertLogMerchant(t *testing.T) {
	os.Setenv(constMerchantURL, defaultSalmonURL)
	os.Setenv(constGql, defaultGQLURL)
	os.Setenv("KAFKA_USER_SERVICE_MERCHANT_TOPIC", "anything")
	ctx := context.Background()

	mockPub := new(mocks.QPublisher)
	mockActivity := new(mocks.ActivityServices)
	ms, _ := NewMerchantService(mockPub, mockActivity)

	var testDataInsert = []struct {
		name   string
		input1 merchantModel.B2CMerchantDataV2
		input2 merchantModel.B2CMerchantDataV2
	}{
		{
			name:   "Test Insert Log Merchant #1",
			input1: merchantModel.B2CMerchantDataV2{ID: "1"},
			input2: merchantModel.B2CMerchantDataV2{},
		},
		{
			name:   "Test Insert Log Merchant #2",
			input1: merchantModel.B2CMerchantDataV2{ID: "3"},
			input2: merchantModel.B2CMerchantDataV2{ID: "2"},
		},
	}

	for _, tc := range testDataInsert {
		mockActivity.On("InsertLog", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		ms.InsertLogMerchant(ctx, tc.input1, tc.input2, "update", "merchant")
	}
}

func TestRestruct(t *testing.T) {
	var defDateTime = "2020-12-30T13:28:19+07:00"
	t.Run("RestructMerchantGWS", func(t *testing.T) {
		input := serviceModel.GWSMerchantData{
			UpdatedAt: "2020-12-30T13:28:19.376845+07:00",
		}

		output := RestructMerchantGWS(input)

		assert.Equal(t, defDateTime, *output.LastModified)
	})

	t.Run("RestructMerchant", func(t *testing.T) {
		input := serviceModel.MerchantPayloadCDC{
			Payload: serviceModel.MerchantPayloadData{
				After: merchantModel.B2CMerchant{
					ID: "MCH1701",
				},
			},
		}
		output := RestructMerchant(input)
		assert.Equal(t, "MCH1701", output.ID)
	})

	t.Run("RestructMerchantDocument", func(t *testing.T) {
		input := serviceModel.MerchantDocumentPayloadCDC{
			Payload: serviceModel.MerchantDocumentPayloadData{
				After: serviceModel.B2CMerchantDocument{ID: "DOC101"},
			},
		}
		output := RestructMerchantDocument(input)
		assert.Equal(t, "DOC101", output.ID)
	})

	t.Run("RestructMerchantBank", func(t *testing.T) {
		input := serviceModel.MerchantBankPayloadCDC{
			Payload: serviceModel.MerchantBankPayloadData{
				After: merchantModel.B2CMerchantBank{ID: 1701},
			},
		}

		output := RestructMerchantBank(input)
		assert.Equal(t, 1701, output.ID)
	})

	t.Run("RestructMerchantDocumentGWS", func(t *testing.T) {
		docData := serviceModel.GWSMerchantDocumentData{Type: "ii", Code: "DOC0011"}
		data := serviceModel.GWSMerchantData{UserEmail: "pian.mutakin@bhinneka.com"}

		output := RestructMerchantDocumentGWS(docData, data)
		assert.Equal(t, "DOC0011", output.ID)
	})

	t.Run("RestructMerchantDataGWS", func(t *testing.T) {
		input := serviceModel.GWSMerchantData{
			StoreClosureDate:        defDateTime,
			StoreReopenDate:         defDateTime,
			StoreActiveShippingDate: defDateTime,
			DeletedAt:               defDateTime,
			MouDate:                 defDateTime,
			AgreementDate:           defDateTime,
			CreatedAt:               defDateTime,
			UpdatedAt:               defDateTime,
		}
		RestructMerchantDataGWS(input)
	})

	t.Run("RestructMerchantDocumentDataGWS", func(t *testing.T) {
		docData := serviceModel.GWSMerchantDocumentData{
			Expired: defDateTime,
		}
		data := serviceModel.GWSMerchantData{
			CreatedAt: defDateTime,
			UpdatedAt: defDateTime,
		}

		output := RestructMerchantDocumentDataGWS(docData, data)
		expiredDate, _ := time.Parse(time.RFC3339, defDateTime)
		assert.Equal(t, null.TimeFrom(expiredDate), output.DocumentExpirationDate)
	})

	t.Run("RestructMasterMerchantBankGWS", func(t *testing.T) {
		input := serviceModel.GWSMasterMerchantBank{BankID: 2902}
		output := RestructMasterMerchantBankGWS(input)

		assert.Equal(t, 2902, output.ID)
	})
}

var inputParse = []struct {
	name                             string
	input                            string
	wantError                        bool
	expectedLegalEntity, expectedNOE interface{}
}{
	{
		name:                "positive 1",
		input:               `{"eventType":"merchantCreated","data":{"id":"91Vgv3Y","userId":"USR210138917","userEmail":"merchantoke@yopmail.com","code":"MCH210215152430","name":"TE495qa84","vanityUrl":"te495qa84","richContent":"","additionalEmail":"","launchDev":"","mouDate":"","note":"","agreementDate":"","businessType":"CORPORATE","description":"Desc merchant","companyName":"Nama Perusahaan","storeAddress":"","shippingMethod":[],"storeProvince":{},"storeCity":{},"storeDistrict":{},"storeSubDistrict":{},"storeZipCode":"","storeIsClosed":false,"storeClosureDate":"","storeReopenDate":"","storeActiveShippingDate":"","logo":"","address":"Alamat merchant","province":{},"city":{},"district":{},"subDistrict":{},"zipCode":"","isPkp":true,"picName":"Ujang","picOccupation":"Supervisor","accountManager":"","acquisitor":"","dailyOperationalStaff":"Asep46","phoneNumber":"021123123","mobileNumber":"081234567712","npwpNo":"123451234512345","npwpName":"Ujang","bank":{"id":"Ky1RXEB","bankId":1,"bankCode":"h2632","name":"test398","branch":"Cicendo","accountNo":"1234512345","accountName":"asep"},"documents":[{"id":"71LrM2E","code":"DOC210215152430375","type":"NPWP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""},{"id":"QqD68PY","code":"DOC210215152430370","type":"KTP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""}],"isActive":false,"createdAt":"2021-02-15T15:24:30.362305+07:00","updatedAt":"2021-02-15T15:24:30.362338+07:00","deletedAt":"","merchantType":"REGULAR","upgradeStatus":"","genderPic":"MALE","merchantGroup":"MEDIUM","areaCoverage":null,"createdBy":"USR190102","updatedBy":"USR190102","legalEntity":1, "numberOfEmployee":12}}`,
		wantError:           false,
		expectedLegalEntity: 1,
		expectedNOE:         12,
	},
	{
		name:                "positive 2",
		input:               `{"eventType":"merchantCreated","data":{"id":"91Vgv3Y","userId":"USR210138917","userEmail":"merchantoke@yopmail.com","code":"MCH210215152430","name":"TE495qa84","vanityUrl":"te495qa84","richContent":"","additionalEmail":"","launchDev":"","mouDate":"","note":"","agreementDate":"","businessType":"CORPORATE","description":"Desc merchant","companyName":"Nama Perusahaan","storeAddress":"","shippingMethod":[],"storeProvince":{},"storeCity":{},"storeDistrict":{},"storeSubDistrict":{},"storeZipCode":"","storeIsClosed":false,"storeClosureDate":"","storeReopenDate":"","storeActiveShippingDate":"","logo":"","address":"Alamat merchant","province":{},"city":{},"district":{},"subDistrict":{},"zipCode":"","isPkp":true,"picName":"Ujang","picOccupation":"Supervisor","accountManager":"","acquisitor":"","dailyOperationalStaff":"Asep46","phoneNumber":"021123123","mobileNumber":"081234567712","npwpNo":"123451234512345","npwpName":"Ujang","bank":{"id":"Ky1RXEB","bankId":1,"bankCode":"h2632","name":"test398","branch":"Cicendo","accountNo":"1234512345","accountName":"asep"},"documents":[{"id":"71LrM2E","code":"DOC210215152430375","type":"NPWP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""},{"id":"QqD68PY","code":"DOC210215152430370","type":"KTP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""}],"isActive":false,"createdAt":"2021-02-15T15:24:30.362305+07:00","updatedAt":"2021-02-15T15:24:30.362338+07:00","deletedAt":"","merchantType":"REGULAR","upgradeStatus":"","genderPic":"MALE","merchantGroup":"MEDIUM","areaCoverage":null,"createdBy":"USR190102","updatedBy":"USR190102"}}`,
		wantError:           false,
		expectedLegalEntity: 0,
		expectedNOE:         0,
	},
	{
		name:                "positive 3",
		input:               `{"eventType":"merchantCreated","data":{"id":"91Vgv3Y","userId":"USR210138917","userEmail":"merchantoke@yopmail.com","code":"MCH210215152430","name":"TE495qa84","vanityUrl":"te495qa84","richContent":"","additionalEmail":"","launchDev":"","mouDate":"","note":"","agreementDate":"","businessType":"CORPORATE","description":"Desc merchant","companyName":"Nama Perusahaan","storeAddress":"","shippingMethod":[],"storeProvince":{},"storeCity":{},"storeDistrict":{},"storeSubDistrict":{},"storeZipCode":"","storeIsClosed":false,"storeClosureDate":"","storeReopenDate":"","storeActiveShippingDate":"","logo":"","address":"Alamat merchant","province":{},"city":{},"district":{},"subDistrict":{},"zipCode":"","isPkp":true,"picName":"Ujang","picOccupation":"Supervisor","accountManager":"","acquisitor":"","dailyOperationalStaff":"Asep46","phoneNumber":"021123123","mobileNumber":"081234567712","npwpNo":"123451234512345","npwpName":"Ujang","bank":{"id":"Ky1RXEB","bankId":1,"bankCode":"h2632","name":"test398","branch":"Cicendo","accountNo":"1234512345","accountName":"asep"},"documents":[{"id":"71LrM2E","code":"DOC210215152430375","type":"NPWP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""},{"id":"QqD68PY","code":"DOC210215152430370","type":"KTP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""}],"isActive":false,"createdAt":"2021-02-15T15:24:30.362305+07:00","updatedAt":"2021-02-15T15:24:30.362338+07:00","deletedAt":"","merchantType":"REGULAR","upgradeStatus":"","genderPic":"MALE","merchantGroup":"MEDIUM","areaCoverage":null,"createdBy":"USR190102","updatedBy":"USR190102", "legalEntity":null, "numberOfEmployee":null}}`,
		wantError:           false,
		expectedLegalEntity: 0,
		expectedNOE:         0,
	},
	{
		name:                "positive 4",
		input:               `{"eventType":"merchantCreated","data":{"id":"91Vgv3Y","userId":"USR210138917","userEmail":"merchantoke@yopmail.com","code":"MCH210215152430","name":"TE495qa84","vanityUrl":"te495qa84","richContent":"","additionalEmail":"","launchDev":"","mouDate":"","note":"","agreementDate":"","businessType":"CORPORATE","description":"Desc merchant","companyName":"Nama Perusahaan","storeAddress":"","shippingMethod":[],"storeProvince":{},"storeCity":{},"storeDistrict":{},"storeSubDistrict":{},"storeZipCode":"","storeIsClosed":false,"storeClosureDate":"","storeReopenDate":"","storeActiveShippingDate":"","logo":"","address":"Alamat merchant","province":{},"city":{},"district":{},"subDistrict":{},"zipCode":"","isPkp":true,"picName":"Ujang","picOccupation":"Supervisor","accountManager":"","acquisitor":"","dailyOperationalStaff":"Asep46","phoneNumber":"021123123","mobileNumber":"081234567712","npwpNo":"123451234512345","npwpName":"Ujang","bank":{"id":"Ky1RXEB","bankId":1,"bankCode":"h2632","name":"test398","branch":"Cicendo","accountNo":"1234512345","accountName":"asep"},"documents":[{"id":"71LrM2E","code":"DOC210215152430375","type":"NPWP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""},{"id":"QqD68PY","code":"DOC210215152430370","type":"KTP","value":"https://static.bmdstatic.com/gk/dev/2cacf9f87eacbc5b1ca2a4a5a93ad4ef.png","expired":""}],"isActive":false,"createdAt":"2021-02-15T15:24:30.362305+07:00","updatedAt":"2021-02-15T15:24:30.362338+07:00","deletedAt":"","merchantType":"REGULAR","upgradeStatus":"","genderPic":"MALE","merchantGroup":"MEDIUM","areaCoverage":null,"createdBy":"USR190102","updatedBy":"USR190102", "legalEntity":"", "numberOfEmployee":""}}`,
		wantError:           false,
		expectedLegalEntity: 0,
		expectedNOE:         0,
	},
}

func TestConsumerKafka(t *testing.T) {
	for _, tc := range inputParse {
		t.Run(tc.name, func(t *testing.T) {
			input := []byte(tc.input)
			binder := serviceModel.GWSMerchantPayloadMessage{}
			_ = json.Unmarshal(input, &binder)

			data := RestructMerchantGWS(binder.Data)
			assert.Equal(t, tc.expectedNOE, *data.NumberOfEmployee)
			assert.Equal(t, tc.expectedLegalEntity, *data.LegalEntity)
		})
	}
}
