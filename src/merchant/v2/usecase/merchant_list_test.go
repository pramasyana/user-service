package usecase

import (
	"context"
	"database/sql"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	mockToken "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	mockMerchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo/mocks"
	mockServices "github.com/Bhinneka/user-service/src/service/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedMock "github.com/Bhinneka/user-service/src/shared/repository/mocks"
	serviceMock "github.com/Bhinneka/user-service/src/shared/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var defQParams = model.QueryParameters{}

var testDataGetListMerchants = []struct {
	name               string
	wantError          bool
	input              *model.QueryParameters
	merchantRepoResult merchantRepo.ResultRepository
	merchantRepoTotal  merchantRepo.ResultRepository
}{
	{
		name:               "Test Get List Merchants Negative #1",
		wantError:          true,
		input:              &defQParams,
		merchantRepoResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
	},
	{
		name:               "Test Get List Merchants Positive #1",
		wantError:          false,
		input:              &defQParams,
		merchantRepoResult: merchantRepo.ResultRepository{Result: []model.B2CMerchantDataV2{}},
		merchantRepoTotal:  merchantRepo.ResultRepository{Result: 900},
	},
	{
		name:               "Test Get List Merchants Negative #2",
		wantError:          true,
		input:              &defQParams,
		merchantRepoResult: merchantRepo.ResultRepository{Result: []model.B2CMerchantDataV2{}},
		merchantRepoTotal:  merchantRepo.ResultRepository{Error: errDefault},
	},
	{
		name:      "Test Get List Merchants Negative #3",
		wantError: true,
		input:     &model.QueryParameters{StrPage: "o"},
	},
}

func TestGetListMerchants(t *testing.T) {
	for _, tc := range testDataGetListMerchants {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		svcRepo := localConfig.ServiceRepository{
			MerchantRepository: &merchantRepoMock,
		}
		svcShared := localConfig.ServiceShared{}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()
		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantRepoMock.On("GetMerchants", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoResult))
		merchantRepoMock.On("GetTotalMerchant", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoTotal))
		ucResult := <-m.GetMerchants(ctx, tc.input)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

// Test Get Merchant By ID

var testDataGetListMerchantbyID = []basicInput{
	{
		name:                      "Test Get Merchant Positive #1",
		wantError:                 false,
		input:                     defUserID,
		input2:                    defUserID,
		isAttachment:              "false",
		merchantRepoResult:        merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
		merchantAddressRepoResult: merchantRepo.ResultRepository{Result: model.Maps{ID: "MAPS001"}},
		merchantDocRepoResult: merchantRepo.ResultRepository{Result: model.ListB2CMerchantDocument{
			MerchantDocument: []model.B2CMerchantDocumentData{{ID: defUserID}},
		}},
		privacy: "private",
	},
	{
		name:               "Test Get List Merchant Negative #1",
		wantError:          true,
		input:              defUserID,
		input2:             defUserID,
		isAttachment:       "false",
		merchantRepoResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
		privacy:            "private",
	},
}

func TestGetMerchantByID(t *testing.T) {
	for _, tc := range testDataGetListMerchantbyID {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		merchantDocRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
		merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
		svcRepo := localConfig.ServiceRepository{
			MerchantRepository:         &merchantRepoMock,
			MerchantDocumentRepository: &merchantDocRepoMock,
			MerchantAddressRepository:  &merchantAddressRepoMock,
		}
		svcShared := localConfig.ServiceShared{}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()

		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.merchantRepoResult)
		merchantAddressRepoMock.On("FindAddressMaps", mock.Anything, mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantAddressRepoResult))
		merchantDocRepoMock.On("GetListMerchantDocument", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantDocRepoResult))
		ucResult := <-m.GetMerchantByID(ctx, tc.input, tc.privacy, tc.isAttachment)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

type basicInput struct {
	name                      string
	wantError                 bool
	input                     string
	input2                    string
	isAttachment              string
	merchantRepoResult        merchantRepo.ResultRepository
	serviceResult             serviceModel.ServiceResult
	merchantDocRepoResult     merchantRepo.ResultRepository
	merchantAddressRepoResult merchantRepo.ResultRepository
	merchantRepoResult2       merchantRepo.ResultRepository
	merchantData              model.B2CMerchantDataV2
	err                       error
	privacy                   string
}

var testDataGetListMerchantbyUserID = []basicInput{
	{
		name:               "Test Get Merchant By User ID Positive #1",
		wantError:          false,
		input:              defUserID,
		input2:             defUserID,
		merchantRepoResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defUserID}},
		serviceResult:      serviceModel.ServiceResult{Error: errDefault},
		merchantDocRepoResult: merchantRepo.ResultRepository{Result: model.ListB2CMerchantDocument{
			MerchantDocument: []model.B2CMerchantDocumentData{{ID: defUserID}},
		}},
	},
	{
		name:               "Test Get Merchant By User ID  Negative #1",
		wantError:          true,
		input:              defUserID,
		input2:             defUserID,
		merchantRepoResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
	},
}

func TestGetMerchantByUserID(t *testing.T) {
	for _, tc := range testDataGetListMerchantbyUserID {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		merchantServices := mockServices.MerchantServices{}
		merchantDocRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
		merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
		svcRepo := localConfig.ServiceRepository{
			MerchantRepository:         &merchantRepoMock,
			MerchantDocumentRepository: &merchantDocRepoMock,
			MerchantAddressRepository:  &merchantAddressRepoMock,
		}
		svcShared := localConfig.ServiceShared{
			MerchantService: &merchantServices,
		}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()

		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantRepoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.merchantRepoResult)
		merchantAddressRepoMock.On("FindAddressMaps", mock.Anything, mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantAddressRepoResult))
		merchantServices.On("FindMerchantServiceByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(serviceMock.ServiceResult(tc.serviceResult))
		merchantDocRepoMock.On("GetListMerchantDocument", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantDocRepoResult))
		ucResult := <-m.GetMerchantByUserID(ctx, tc.input, tc.input2)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

var testDataMerchantByName = []basicInput{
	{
		name:               "Test Get Merchant By Name #1",
		input:              defName,
		wantError:          true,
		merchantRepoResult: merchantRepo.ResultRepository{Result: "some name"},
	},
	{
		name:      "Test Get Merchant By Name #2",
		wantError: true,
		input:     "",
	},
	{
		name:               "Test Get Merchant By Name #1",
		input:              defName,
		wantError:          false,
		merchantRepoResult: merchantRepo.ResultRepository{},
	},
}

func TestGetMerchantByName(t *testing.T) {
	for _, tc := range testDataMerchantByName {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}

		svcRepo := localConfig.ServiceRepository{
			MerchantRepository: &merchantRepoMock,
		}
		svcShared := localConfig.ServiceShared{}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()

		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantRepoMock.On("FindMerchantByName", mock.Anything, mock.Anything).Return(tc.merchantRepoResult)

		ucResult := <-m.CheckMerchantName(ctx, tc.input)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}
