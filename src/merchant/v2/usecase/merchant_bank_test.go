package usecase

import (
	"database/sql"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	mockToken "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	mockMerchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo/mocks"
	sharedMock "github.com/Bhinneka/user-service/src/shared/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var defInputBank = model.ParametersMerchantBank{StrPage: "1", StrLimit: "10", OrderBy: "id", Sort: "desc"}
var testDataGetListBank = []struct {
	name               string
	wantError          bool
	input              *model.ParametersMerchantBank
	merchantRepoResult merchantRepo.ResultRepository
	merchantRepoTotal  merchantRepo.ResultRepository
}{
	{
		name:               "Test Get List Bank Negative #1",
		wantError:          true,
		input:              &defInputBank,
		merchantRepoResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
	},
	{
		name:               "Test Get List Bank Positive #1",
		wantError:          false,
		input:              &defInputBank,
		merchantRepoResult: merchantRepo.ResultRepository{Result: model.ListMerchantBank{TotalData: 999}},
		merchantRepoTotal:  merchantRepo.ResultRepository{Result: 900},
	},
	{
		name:      "Test Get List Bank Negative #2",
		wantError: true,
		input:     &model.ParametersMerchantBank{StrPage: "a"},
	},
	{
		name:      "Test Get List Bank Negative #3",
		wantError: true,
		input:     &model.ParametersMerchantBank{Sort: "other"},
	},
	{
		name:      "Test Get List Bank Negative #4",
		wantError: true,
		input:     &model.ParametersMerchantBank{OrderBy: "other sort"},
	},
	{
		name:      "Test Get List Bank Negative #5",
		wantError: true,
		input:     &model.ParametersMerchantBank{},
	},
	{
		name:               "Test Get List Bank Negative #6",
		wantError:          true,
		input:              &defInputBank,
		merchantRepoResult: merchantRepo.ResultRepository{Result: model.ListMerchantBank{TotalData: 999}},
		merchantRepoTotal:  merchantRepo.ResultRepository{Error: errDefault},
	},
}

func TestGetListBank(t *testing.T) {

	for _, tc := range testDataGetListBank {
		merchantBankRepoMock := mockMerchantRepo.MerchantBankRepository{}
		svcRepo := localConfig.ServiceRepository{
			MerchantBankRepository: &merchantBankRepoMock,
		}
		svcShared := localConfig.ServiceShared{}
		tokenGen := mockToken.AccessTokenGenerator{}

		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantBankRepoMock.On("GetListMerchantBank", mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoResult))
		merchantBankRepoMock.On("GetTotalMerchantBank", mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoTotal))
		ucResult := <-m.GetListMerchantBank(tc.input)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}
