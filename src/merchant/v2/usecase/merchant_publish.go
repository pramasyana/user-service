package usecase

import (
	"context"
	"net/http"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// PublishToKafka function for get merchant by merchant ID
func (m *MerchantUseCaseImpl) PublishToKafkaMerchant(ctxReq context.Context, data model.B2CMerchantDataV2, eventType string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-PublishToKafkaMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextMerchantIDCamel] = data.ID

		eventTypeProduce := helper.EventProduceCreateMerchant
		switch eventType {
		case helper.EventProduceCreateMerchant:
			eventTypeProduce = helper.EventProduceCreateMerchant
		case helper.EventProduceUpdateMerchant:
			eventTypeProduce = helper.EventProduceUpdateMerchant
		case helper.EventProduceDeleteMerchant:
			eventTypeProduce = helper.EventProduceDeleteMerchant
		}
		// this is new payload produced by sturgeon
		err := m.MerchantService.PublishToKafkaUserMerchant(ctxReq, &data, eventTypeProduce, producerSturgeon)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tags[helper.TextResponse] = data
		output <- ResultUseCase{Result: data}
	})
	return output
}
