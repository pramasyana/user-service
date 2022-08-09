package usecase

import (
	"context"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	mockToken "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	mockServices "github.com/Bhinneka/user-service/src/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testDataPublish = []basicInput{
	{
		name:         "Test publish #1",
		merchantData: model.B2CMerchantDataV2{},
		input:        helper.EventProduceCreateMerchant,
	},
	{
		name:         "Test publish #2",
		merchantData: model.B2CMerchantDataV2{},
		input:        helper.EventProduceUpdateMerchant,
	},
	{
		name:         "Test publish #3",
		merchantData: model.B2CMerchantDataV2{},
		input:        helper.EventProduceDeleteMerchant,
		err:          errDefault,
		wantError:    true,
	},
}

func TestPublish(t *testing.T) {

	for _, tc := range testDataPublish {
		merchantServices := mockServices.MerchantServices{}
		svcRepo := localConfig.ServiceRepository{}
		svcShared := localConfig.ServiceShared{
			MerchantService: &merchantServices,
		}
		localQuery := localConfig.ServiceQuery{}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()

		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)

		merchantServices.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.err)

		ucResult := <-m.PublishToKafkaMerchant(ctx, tc.merchantData, tc.input)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}
