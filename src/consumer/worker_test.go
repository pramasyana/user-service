package consumer

import (
	"context"
	"testing"

	mocksMerchantUsecase "github.com/Bhinneka/user-service/mocks/src/merchant/v2/usecase"
	memberUC "github.com/Bhinneka/user-service/src/member/v1/usecase"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantUC "github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	"github.com/stretchr/testify/mock"
)

func generateUsecaseResult(data merchantUC.ResultUseCase) <-chan merchantUC.ResultUseCase {
	output := make(chan merchantUC.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateMemberUsecaseResult(data memberUC.ResultUseCase) <-chan memberUC.ResultUseCase {
	output := make(chan memberUC.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func Test_sendEmailMerchant(t *testing.T) {
	type args struct {
		ctxReq          context.Context
		payload         interface{}
		merchantUsecase merchantUC.MerchantUseCase
		eventType       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success SendEmailMerchantAdd",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailMerchantAdd", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailMerchantAdd",
			},
		},
		{
			name: "Case 2: Success SendEmailMerchantRejectRegistration",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailMerchantRejectRegistration", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailMerchantRejectRegistration",
			},
		},
		{
			name: "Case 3: Success SendEmailMerchantRejectUpgrade",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailMerchantRejectUpgrade", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailMerchantRejectUpgrade",
			},
		},
		{
			name: "Case 4: Success SendEmailActivation",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailActivation", mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailActivation",
			},
		},
		{
			name: "Case 5: Success SendEmailApproval",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailApproval", mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailApproval",
			},
		},
		{
			name: "Case 6: Success SendEmailMerchantUpgrade",
			args: args{
				ctxReq:  context.Background(),
				payload: merchantModel.MerchantPayloadEmail{},
				merchantUsecase: func() merchantUC.MerchantUseCase {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantUseCase)
					mocksMerchantUsecase.On("SendEmailMerchantUpgrade", mock.Anything, mock.Anything, mock.Anything).Return(generateUsecaseResult(merchantUC.ResultUseCase{
						Error: nil,
					}))

					return mocksMerchantUsecase
				}(),
				eventType: "SendEmailMerchantUpgrade",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendEmailMerchant(tt.args.ctxReq, tt.args.payload, tt.args.merchantUsecase, tt.args.eventType)
		})
	}
}
