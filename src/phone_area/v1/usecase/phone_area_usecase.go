package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
	"github.com/Bhinneka/user-service/src/phone_area/v1/query"
)

// PhoneAreaUseCaseImpl data structure
type PhoneAreaUseCaseImpl struct {
	PhoneAreaQueryRead query.PhoneAreaQuery
}

// NewPhoneAreaUseCase function for initialise phone area use case implementation
func NewPhoneAreaUseCase(phoneAreaQueryRead query.PhoneAreaQuery) PhoneAreaUseCase {

	return &PhoneAreaUseCaseImpl{
		PhoneAreaQueryRead: phoneAreaQueryRead,
	}
}

// GetAllPhoneArea function for getting list of phone area
func (uc *PhoneAreaUseCaseImpl) GetAllPhoneArea(ctxReq context.Context) <-chan ResultUseCase {
	ctx := "PhoneAreaUseCase-GetAllPhoneArea"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		phoneAreaResult := <-uc.PhoneAreaQueryRead.FindAll(ctxReq)
		if phoneAreaResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			output <- ResultUseCase{Error: phoneAreaResult.Error, HTTPStatus: httpStatus}
			return
		}

		phoneArea, ok := phoneAreaResult.Result.([]model.PhoneArea)
		if !ok {
			err := errors.New("result is not list of phone area")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}
		tags["args"] = phoneArea
		output <- ResultUseCase{Result: phoneArea, HTTPStatus: http.StatusOK}
	})

	return output
}

// GetTotalPhoneArea function for getting total phone area
func (uc *PhoneAreaUseCaseImpl) GetTotalPhoneArea(ctxReq context.Context) <-chan ResultUseCase {
	ctx := "PhoneAreaUseCase-GetTotalPhoneArea"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		totalPhoneAreaResult := <-uc.PhoneAreaQueryRead.Count(ctxReq)
		if totalPhoneAreaResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			output <- ResultUseCase{Error: totalPhoneAreaResult.Error, HTTPStatus: httpStatus}
			return
		}

		total, ok := totalPhoneAreaResult.Result.(model.TotalPhoneArea)
		if !ok {
			err := errors.New("result is not total phone area")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		tags["args"] = total

		output <- ResultUseCase{Result: total}
	})

	return output
}
