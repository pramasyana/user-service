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
	sharedMock "github.com/Bhinneka/user-service/src/shared/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMerchantUseCaseImpl_GetMerchantsPublic(t *testing.T) {
	var defQParams = model.QueryParametersPublic{}

	var testDataGetListMerchants = []struct {
		name               string
		wantError          bool
		input              *model.QueryParametersPublic
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
			merchantRepoResult: merchantRepo.ResultRepository{Result: []model.B2CMerchantDataPublic{}},
			merchantRepoTotal:  merchantRepo.ResultRepository{Result: 900},
		},
		{
			name:               "Test Get List Merchants Negative #2",
			wantError:          true,
			input:              &defQParams,
			merchantRepoResult: merchantRepo.ResultRepository{Result: []model.B2CMerchantDataPublic{}},
			merchantRepoTotal:  merchantRepo.ResultRepository{Error: errDefault},
		},
		{
			name:      "Test Get List Merchants Negative #3",
			wantError: true,
			input:     &model.QueryParametersPublic{StrPage: "o"},
		},
	}

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
		merchantRepoMock.On("GetMerchantsPublic", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoResult))
		merchantRepoMock.On("GetTotalMerchant", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoTotal))
		ucResult := <-m.GetMerchantsPublic(ctx, tc.input)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}

}
