package usecase

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/src/payments/v1/model"
	"github.com/Bhinneka/user-service/src/payments/v1/repo"
	mockPaymentsRepo "github.com/Bhinneka/user-service/src/payments/v1/repo/mocks"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

var (
	errDefault = fmt.Errorf("default error")
	defInput   = model.Payments{
		Channel:   "b2c",
		Method:    "KREDIVO",
		Email:     "tes@getnada.com",
		Token:     "random token",
		ExpiredAt: time.Now(),
	}
	defInput2 = model.Payments{
		Email:     "tes1234@getnada.com",
		Channel:   "b2b",
		Method:    "kredivo",
		Token:     "214123",
		ExpiredAt: time.Now(),
	}
)

func generateRepoResult(data repo.ResultRepository) <-chan repo.ResultRepository {
	output := make(chan repo.ResultRepository)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}
func TestPaymentsUseCaseImpl_AddUpdatePayments(t *testing.T) {
	//defInput2 := defInput

	var testDataAddPayment = []struct {
		name       string
		wantError  bool
		payment    *model.Payments
		repoResult repo.ResultRepository
		saveResult repo.ResultRepository
	}{
		{
			name:    "Test Add Payment #1", // all passed
			payment: &defInput,
			repoResult: repo.ResultRepository{Result: model.Payments{
				ID:        defInput.ID,
				Email:     defInput.Email,
				Channel:   defInput.Channel,
				Method:    defInput.Method,
				Token:     defInput.Token,
				ExpiredAt: defInput.ExpiredAt,
			}},
			saveResult: repo.ResultRepository{Error: nil},
		},
		{
			name:       "Test Add Payment #2", //Empty Find Payment in Database
			payment:    &defInput,
			repoResult: repo.ResultRepository{Error: errDefault},
		},
		{
			name:    "Test Add Payment #3", //Failed Save Payment
			payment: &defInput,
			repoResult: repo.ResultRepository{Result: model.Payments{
				ID:        defInput.ID,
				Email:     defInput.Email,
				Channel:   defInput.Channel,
				Method:    defInput.Method,
				Token:     defInput.Token,
				ExpiredAt: defInput.ExpiredAt,
			}},
			saveResult: repo.ResultRepository{Error: errDefault},
			wantError:  true,
		},
	}
	for _, tc := range testDataAddPayment {
		paymentRepoMock := mockPaymentsRepo.PaymentsRepository{}
		mockDB, sqlMock, _ := sqlMock.New()
		sqlMock.ExpectBegin()

		defer mockDB.Close()

		svcRepo := localConfig.ServiceRepository{
			PaymentsRepository: &paymentRepoMock,
			Repository:         &repository.Repository{WriteDB: mockDB},
		}

		svcQuery := localConfig.ServiceQuery{}
		p := NewPaymentsUseCase(svcRepo, svcQuery)
		ctxReq := context.Background()

		paymentRepoMock.On("FindPaymentByEmailChannelMethod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.repoResult)
		paymentRepoMock.On("AddUpdatePayment", mock.Anything, mock.Anything).Return(generateRepoResult(tc.saveResult))

		ucResult := <-p.AddUpdatePayments(ctxReq, tc.payment)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
			sqlMock.ExpectRollback()
		} else {
			assert.NoError(t, ucResult.Error)
			sqlMock.ExpectCommit()
		}
	}
}

func TestPaymentsUseCaseImpl_GetPaymentDetail(t *testing.T) {
	var testGetDataPayment = []struct {
		name       string
		wantError  bool
		payment    *model.Payments
		repoResult repo.ResultRepository
	}{
		{
			name:    "Test Add Payment #1", // all passed
			payment: &defInput,
			repoResult: repo.ResultRepository{Result: model.Payments{
				ID:        defInput.ID,
				Email:     defInput.Email,
				Channel:   defInput.Channel,
				Method:    defInput.Method,
				Token:     defInput.Token,
				ExpiredAt: defInput.ExpiredAt,
			}},
		},
		{
			name:       "Test Add Payment #2", //Empty Find Payment in Database
			payment:    &defInput,
			repoResult: repo.ResultRepository{Error: errDefault},
			wantError:  true,
		},
	}
	for _, tc := range testGetDataPayment {
		paymentRepoMock := mockPaymentsRepo.PaymentsRepository{}
		mockDB, sqlMock, _ := sqlMock.New()
		sqlMock.ExpectBegin()

		defer mockDB.Close()

		svcRepo := localConfig.ServiceRepository{
			PaymentsRepository: &paymentRepoMock,
			Repository:         &repository.Repository{WriteDB: mockDB},
		}

		svcQuery := localConfig.ServiceQuery{}
		p := NewPaymentsUseCase(svcRepo, svcQuery)
		ctxReq := context.Background()

		paymentRepoMock.On("FindPaymentByEmailChannelMethod", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.repoResult)

		ucResult := <-p.GetPaymentDetail(ctxReq, tc.payment)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
			sqlMock.ExpectRollback()
		} else {
			assert.NoError(t, ucResult.Error)
			sqlMock.ExpectCommit()
		}
	}
}

func TestPaymentsUseCaseImpl_CompareHeaderAndBody(t *testing.T) {
	defInput3 := defInput2
	defInput3.Email = ""
	defInput4 := defInput2
	defInput4.Channel = ""
	defInput5 := defInput2
	defInput5.Method = ""
	var testDataAddPayment = []struct {
		name      string
		basicAuth string
		wantError bool
		decode    []byte
		payment   *model.Payments
	}{
		{
			name:      "Test Add Payment #1", // Failed Decrypt
			payment:   &defInput,
			basicAuth: "9e7295636d60c43ca836f356714661c7ea0e1f5c3a4db5d2d3f2a23804eb742809cc246301c711c5321c2bbbe45afa2f919bf2d993fd8f0e654b90",
			wantError: true,
		},
		{
			name:      "Test Add Payment #2", // Failed Decodestring
			payment:   &defInput2,
			basicAuth: "9e7295636d60c43ca836f356714661c7ea0e1f5c3a4db5d2d3f2a23804eb742809cc246301c711c5321c2bbbe45afa2f919bf2d993fd8f0e654b23re",
			wantError: true,
		},
		{
			name:      "Test Add Payment #3", // Email Mismatch
			payment:   &defInput3,
			basicAuth: "0225c5c10c02aa626a9b0897cf85417a18238eb0b4acb996e519ad225f048acb3434af57b8e894456990757d0a786adbaad1aa36890277056692fe",
			wantError: true,
		},
		{
			name:      "Test Add Payment #4", // Channel Mismatch
			payment:   &defInput4,
			basicAuth: "0225c5c10c02aa626a9b0897cf85417a18238eb0b4acb996e519ad225f048acb3434af57b8e894456990757d0a786adbaad1aa36890277056692fe",
			wantError: true,
		},
		{
			name:      "Test Add Payment #5", // Method Mismatch
			payment:   &defInput5,
			basicAuth: "0225c5c10c02aa626a9b0897cf85417a18238eb0b4acb996e519ad225f048acb3434af57b8e894456990757d0a786adbaad1aa36890277056692fe",
			wantError: true,
		},
		{
			name:      "Test Add Payment #6", // Success
			payment:   &defInput2,
			basicAuth: "0225c5c10c02aa626a9b0897cf85417a18238eb0b4acb996e519ad225f048acb3434af57b8e894456990757d0a786adbaad1aa36890277056692fe",
			wantError: false,
		},
	}
	for _, tc := range testDataAddPayment {
		paymentRepoMock := mockPaymentsRepo.PaymentsRepository{}

		svcRepo := localConfig.ServiceRepository{
			PaymentsRepository: &paymentRepoMock,
		}

		svcQuery := localConfig.ServiceQuery{}
		p := NewPaymentsUseCase(svcRepo, svcQuery)
		ctxReq := context.Background()
		os.Setenv("STURGEON_MACKAREL", "579c6e95aa0b76291b0ee44bb6863e6fb024367d651e464d41fdc425eca462d4")

		ucResult := <-p.CompareHeaderAndBody(ctxReq, tc.payment, tc.basicAuth)
		if tc.wantError {
			assert.Error(t, ucResult.Error)

		} else {
			assert.NoError(t, ucResult.Error)

		}
	}
}
