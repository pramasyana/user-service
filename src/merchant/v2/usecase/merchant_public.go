package usecase

import (
	"context"
	"net/http"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"

	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// GetMerchants return all merchants by given parameters
func (m *MerchantUseCaseImpl) GetMerchantsPublic(ctxReq context.Context, params *model.QueryParametersPublic) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchantsPublic"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1,
				StrPage:  params.StrPage,
				Limit:    10,
				StrLimit: params.StrLimit,
			},
		)
		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		params.Page = paging.Page
		params.Limit = paging.Limit
		param := params.RestructParamPublic()
		mr := <-m.MerchantRepo.GetMerchantsPublic(ctxReq, &param)
		if mr.Error != nil {
			tags[helper.TextResponse] = mr.Error.Error()
			output <- ResultUseCase{Error: mr.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchants := mr.Result.([]model.B2CMerchantDataPublic)

		merchantQuery := <-m.MerchantRepo.GetTotalMerchant(ctxReq, &param)
		if merchantQuery.Error != nil {
			output <- ResultUseCase{Error: merchantQuery.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		totalData := merchantQuery.Result.(int)

		output <- ResultUseCase{Result: merchants, TotalData: totalData}
	})

	return output
}
