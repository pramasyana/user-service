package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/payments/v1/model"
	"github.com/Bhinneka/user-service/src/payments/v1/repo"
)

// PaymentsUseCaseImpl data structure
type PaymentsUseCaseImpl struct {
	PaymentsRepoRead  repo.PaymentsRepository
	PaymentsRepoWrite repo.PaymentsRepository
}

// NewPaymentsUseCase function for initialise Payments use case implementation
func NewPaymentsUseCase(
	repository localConfig.ServiceRepository,
	query localConfig.ServiceQuery,
) PaymentsUseCase {

	return &PaymentsUseCaseImpl{
		PaymentsRepoRead:  repository.PaymentsRepository,
		PaymentsRepoWrite: repository.PaymentsRepository,
	}
}

func (pu *PaymentsUseCaseImpl) AddUpdatePayments(ctxReq context.Context, data *model.Payments) <-chan ResultUseCase {
	ctx := "PaymentUseCase-AddPayment"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		payments := pu.PaymentsRepoRead.FindPaymentByEmailChannelMethod(ctxReq, data.Email, data.Channel, data.Method)
		if payments.Result != nil {
			payment := payments.Result.(model.Payments)
			data.ID = payment.ID
		} else {
			t := time.Now()
			data.ID = "TKN" + t.Format(helper.FormatYmdhis)
		}

		saveResult := <-pu.PaymentsRepoWrite.AddUpdatePayment(ctxReq, *data)
		if saveResult.Error != nil {
			err := saveResult.Error
			tracer.SetError(ctxReq, saveResult.Error)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		// check sign up from process

		response := model.SuccessResponse{
			ID:        data.ID,
			Email:     data.Email,
			Message:   helper.SuccessMessage,
			Channel:   data.Channel,
			Method:    data.Method,
			Token:     data.Token,
			ExpiredAt: data.ExpiredAt.Format(time.RFC3339),
		}

		// return the token to be sent to email notification service
		output <- ResultUseCase{Result: response}
	})

	return output
}

func (pu *PaymentsUseCaseImpl) CompareHeaderAndBody(ctxReq context.Context, data *model.Payments, basicAuth string) <-chan ResultUseCase {
	ctx := "PaymentUseCase-CompareHeaderAndBody"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		decode, errs := hex.DecodeString(basicAuth)
		if errs != nil {
			tags[helper.TextResponse] = errs
			output <- ResultUseCase{Error: errs, HTTPStatus: http.StatusBadRequest}
			return
		}

		err, text := golib.Decrypt(decode, os.Getenv("STURGEON_MACKAREL"))
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		datas := strings.Split(text, ":")
		if data.Email != datas[0] {
			err := errors.New(model.ErrorTokenInvalid)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if data.Channel != datas[1] {
			err := errors.New(model.ErrorTokenInvalid)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if data.Method != datas[2] {
			err := errors.New(model.ErrorTokenInvalid)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		output <- ResultUseCase{Error: nil}
	})
	return output
}

func (pu *PaymentsUseCaseImpl) GetPaymentDetail(ctxReq context.Context, data *model.Payments) <-chan ResultUseCase {
	ctx := "PaymentUseCase-GetPaymentDetail"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		paymentsResult := pu.PaymentsRepoRead.FindPaymentByEmailChannelMethod(ctxReq, data.Email, data.Channel, data.Method)
		if paymentsResult.Result == nil {
			err := errors.New("Data Payment Token Not Found")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		payment := paymentsResult.Result.(model.Payments)

		tags[helper.TextMerchantIDCamel] = payment
		tags[helper.TextEmail] = payment.Email
		timeLocation, _ := time.LoadLocation("Asia/Jakarta")
		expiredAt := payment.ExpiredAt.In(timeLocation)

		response := model.SuccessResponse{
			ID:        payment.ID,
			Email:     payment.Email,
			Message:   helper.SuccessMessage,
			Channel:   payment.Channel,
			Method:    payment.Method,
			Token:     payment.Token,
			ExpiredAt: expiredAt.Format(time.RFC3339),
		}
		output <- ResultUseCase{Result: response}

	})

	return output
}
