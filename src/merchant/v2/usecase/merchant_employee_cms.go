package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/golang-jwt/jwt"
)

// CmsGetAllMerchantEmployee return all merchants by given parameters
func (m *MerchantUseCaseImpl) CmsGetAllMerchantEmployee(ctxReq context.Context, tokenCMS string, paramsCMS *model.QueryCmsMerchantEmployeeParameters) <-chan ResultUseCase {
	ctx := "MerchantUseCase-CmsGetAllMerchantEmployee"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claimsCMS := jwt.MapClaims{}
		tokenCMS = strings.Replace(tokenCMS, StringBearer, "", -1)
		jwtResultCMS, err := jwt.ParseWithClaims(tokenCMS, claimsCMS, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResultCMS == nil && err != nil) || len(claimsCMS) == 0 {
			output <- ResultUseCase{Error: errors.New(ErrorInvalidToken)}
			return
		}

		paging, errCMS := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1,
				StrPage:  paramsCMS.StrPage,
				Limit:    10,
				StrLimit: paramsCMS.StrLimit,
			},
		)
		if errCMS != nil {
			tags[helper.TextResponse] = errCMS.Error()
			output <- ResultUseCase{Error: errCMS, HTTPStatus: http.StatusBadRequest}
			return
		}

		paramsCMS.Offset = paging.Offset
		paramsCMS.Page = paging.Page
		paramsCMS.Limit = paging.Limit
		if paramsCMS.Status == "" {
			paramsCMS.Status = "INVITED,ACTIVE,INACTIVE"
		}

		filter := model.QueryMerchantEmployeeParameters(*paramsCMS)
		merchantEmployeeResult := <-m.MerchantEmployeeRepo.GetAllMerchantEmployees(ctxReq, &filter)
		if merchantEmployeeResult.Error != nil {
			tags[helper.TextResponse] = merchantEmployeeResult.Error.Error()
			output <- ResultUseCase{Error: merchantEmployeeResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchantEmployeeCMS := merchantEmployeeResult.Result.([]model.B2CMerchantEmployeeData)

		merchantEmployeeQueryCMS := <-m.MerchantEmployeeRepo.GetTotalMerchantEmployees(ctxReq, &filter)
		if merchantEmployeeQueryCMS.Error != nil {
			output <- ResultUseCase{Error: merchantEmployeeQueryCMS.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		totalDataCMS := merchantEmployeeQueryCMS.Result.(int)

		output <- ResultUseCase{Result: merchantEmployeeCMS, TotalData: totalDataCMS}
	})

	return output
}
