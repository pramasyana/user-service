package usecase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// GetListMerchantBank function for getting list of merchant
func (m *MerchantUseCaseImpl) GetListMerchantBank(params *model.ParametersMerchantBank) <-chan ResultUseCase {
	output := make(chan ResultUseCase)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			}
			close(output)
		}()

		params, err := m.restructParamsMerchantBank(params)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantBankResult := <-m.MerchantBankRepo.GetListMerchantBank(params)
		if merchantBankResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if merchantBankResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				merchantBankResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "merchant_bank")
			}

			output <- ResultUseCase{Error: merchantBankResult.Error, HTTPStatus: httpStatus}
			return
		}

		merchantBank := merchantBankResult.Result.(model.ListMerchantBank)

		totalResult := <-m.MerchantBankRepo.GetTotalMerchantBank(params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		merchantBank.TotalData = total
		output <- ResultUseCase{Result: merchantBank}

	}()

	return output
}

// restructParamsMerchantBank function for getting params
func (m *MerchantUseCaseImpl) restructParamsMerchantBank(params *model.ParametersMerchantBank) (*model.ParametersMerchantBank, error) {
	var err error

	// validate all parameters
	paging, err := helper.ValidatePagination(
		helper.PaginationParameters{
			Page:     1, // default
			StrPage:  params.StrPage,
			Limit:    20, // default
			StrLimit: params.StrLimit,
		})

	if err != nil {
		return params, err
	}
	params.Page = paging.Page
	params.Limit = paging.Limit
	params.Offset = paging.Offset

	if len(params.OrderBy) > 0 {
		if !helper.StringInSlice(params.OrderBy, model.AllowedSortFieldsMerchantBank) {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "order by")
			return params, err
		}
	} else {
		params.OrderBy = "id"
	}

	if len(params.Sort) > 0 {
		if !helper.StringInSlice(params.Sort, []string{"asc", "desc"}) {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "sort")
			return params, err
		}
	} else {
		params.Sort = "asc"
	}

	return params, nil
}
