package usecase

import (
	"context"
	"errors"
	"testing"

	mocksMerchantRepo "github.com/Bhinneka/user-service/mocks/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	"github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/mock"
)

func TestMerchantUseCaseImpl_ValidateMerchantBank(t *testing.T) {
	type fields struct {
		Repository           *sharedRepo.Repository
		MerchantRepo         repo.MerchantRepository
		MerchantAddressRepo  repo.MerchantAddressRepository
		MerchantBankRepo     repo.MerchantBankRepository
		MerchantDocumentRepo repo.MerchantDocumentRepository
		MemberRepoRead       memberRepo.MemberRepository
		UploadService        service.UploadServices
		MerchantService      service.MerchantServices
		TokenGenerator       token.AccessTokenGenerator
		MemberQueryRead      query.MemberQuery
		NotificationService  service.NotificationServices
		QueuePublisher       service.QPublisher
		SendbirdService      service.SendbirdServices
	}
	type args struct {
		ctxReq   context.Context
		data     *model.B2CMerchantCreateInput
		merchant *model.B2CMerchantDataV2
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MerchantBankRepo: func() repo.MerchantBankRepository {
					mocksMerchantBankRepo := new(mocksMerchantRepo.MerchantBankRepository)
					mocksMerchantBankRepo.On("FindActiveMerchantBankByID", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Result: model.B2CMerchantBankData{},
							Error:  nil,
						}),
					)
					return mocksMerchantBankRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				data: &model.B2CMerchantCreateInput{
					BankID: 1,
				},
				merchant: &model.B2CMerchantDataV2{},
			},
		},
		{
			name: "Case 2: Error FindActiveMerchantBankByID",
			fields: fields{
				MerchantBankRepo: func() repo.MerchantBankRepository {
					mocksMerchantBankRepo := new(mocksMerchantRepo.MerchantBankRepository)
					mocksMerchantBankRepo.On("FindActiveMerchantBankByID", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Result: model.B2CMerchantBankData{},
							Error:  errors.New("error"),
						}),
					)
					return mocksMerchantBankRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				data: &model.B2CMerchantCreateInput{
					BankID: 1,
				},
				merchant: &model.B2CMerchantDataV2{},
			},
			wantErr: true,
		},
		{
			name: "Case 2: Error result nil",
			fields: fields{
				MerchantBankRepo: func() repo.MerchantBankRepository {
					mocksMerchantBankRepo := new(mocksMerchantRepo.MerchantBankRepository)
					mocksMerchantBankRepo.On("FindActiveMerchantBankByID", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksMerchantBankRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				data: &model.B2CMerchantCreateInput{
					BankID: 1,
				},
				merchant: &model.B2CMerchantDataV2{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantUseCaseImpl{
				Repository:           tt.fields.Repository,
				MerchantRepo:         tt.fields.MerchantRepo,
				MerchantAddressRepo:  tt.fields.MerchantAddressRepo,
				MerchantBankRepo:     tt.fields.MerchantBankRepo,
				MerchantDocumentRepo: tt.fields.MerchantDocumentRepo,
				MemberRepoRead:       tt.fields.MemberRepoRead,
				UploadService:        tt.fields.UploadService,
				MerchantService:      tt.fields.MerchantService,
				TokenGenerator:       tt.fields.TokenGenerator,
				MemberQueryRead:      tt.fields.MemberQueryRead,
				NotificationService:  tt.fields.NotificationService,
				QueuePublisher:       tt.fields.QueuePublisher,
				SendbirdService:      tt.fields.SendbirdService,
			}
			if err := m.ValidateMerchantBank(tt.args.ctxReq, tt.args.data, tt.args.merchant); (err != nil) != tt.wantErr {
				t.Errorf("MerchantUseCaseImpl.ValidateMerchantBank() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDailyOperationalStaff(t *testing.T) {
	type args struct {
		dailyOperationalStaff string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Case 1: Success",
			args: args{
				dailyOperationalStaff: "",
			},
			wantErr: false,
		},
		{
			name: "Case 2: Error",
			args: args{
				dailyOperationalStaff: "ab",
			},
			wantErr: true,
		},
		{
			name: "Case 3: Error",
			args: args{
				dailyOperationalStaff: "abasdasdasdajhsdjkahsjdkhahsjkhdjkahshdjkahshdjahjkshdjhajkshdahjkshdakshdljkahklsdjaksjdkasjdkljaksjdkaksljdlkaksjdlkajslkjdksjkdjalskjabasdasdasdajhsdjkahsjdkhahsjkhdjkahshdjkahshdjahjkshdjhajkshdahjkshdakshdljkahklsdjaksjdkasjdkljaksjdkaksljdlkaksjdlkajslkjdksjkdjalskjabasdasdasdajhsdjkahsjdkhahsjkhdjkahshdjkahshdjahjkshdjhajkshdahjkshdakshdljkahklsdjaksjdkasjdkljaksjdkaksljdlkaksjdlkajslkjdksjkdjalskjabasdasdasdajhsdjkahsjdkhahsjkhdjkahshdjkahshdjahjkshdjhajkshdahjkshdakshdljkahklsdjaksjdkasjdkljaksjdkaksljdlkaksjdlkajslkjdksjkdjalskj",
			},
			wantErr: true,
		},
		{
			name: "Case 4: Error",
			args: args{
				dailyOperationalStaff: "asdas8#34#$%",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDailyOperationalStaff(tt.args.dailyOperationalStaff); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDailyOperationalStaff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerchantUseCaseImpl_validateMerchantFile(t *testing.T) {
	type fields struct {
		Repository           *sharedRepo.Repository
		MerchantRepo         repo.MerchantRepository
		MerchantAddressRepo  repo.MerchantAddressRepository
		MerchantBankRepo     repo.MerchantBankRepository
		MerchantDocumentRepo repo.MerchantDocumentRepository
		MemberRepoRead       memberRepo.MemberRepository
		UploadService        service.UploadServices
		MerchantService      service.MerchantServices
		TokenGenerator       token.AccessTokenGenerator
		MemberQueryRead      query.MemberQuery
		NotificationService  service.NotificationServices
		QueuePublisher       service.QPublisher
		SendbirdService      service.SendbirdServices
	}
	type args struct {
		data *model.B2CMerchantCreateInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Case 1: Success",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "a",
					NpwpHolderName:    "name",
					Npwp:              "123",
					AccountNumber:     "123",
					AccountHolderName: "nama",
					BankBranch:        "bran",
				},
			},
			wantErr: false,
		},
		{
			name:   "Case 2: Error NpwpFile",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "",
					NpwpHolderName:    "name",
					Npwp:              "123",
					AccountNumber:     "123",
					AccountHolderName: "nama",
					BankBranch:        "bran",
				},
			},
			wantErr: true,
		},
		{
			name:   "Case 3: Error NpwpHolderName",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "1",
					NpwpHolderName:    "",
					Npwp:              "123",
					AccountNumber:     "123",
					AccountHolderName: "nama",
					BankBranch:        "bran",
				},
			},
			wantErr: true,
		},
		{
			name:   "Case 4: Error Npwp",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "1",
					NpwpHolderName:    "asd",
					Npwp:              "",
					AccountNumber:     "123",
					AccountHolderName: "nama",
					BankBranch:        "bran",
				},
			},
			wantErr: true,
		},
		{
			name:   "Case 5: Error AccountNumber",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "1",
					NpwpHolderName:    "asd",
					Npwp:              "123",
					AccountNumber:     "",
					AccountHolderName: "nama",
					BankBranch:        "bran",
				},
			},
			wantErr: true,
		},
		{
			name:   "Case 6: Error AccountHolderName",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "1",
					NpwpHolderName:    "asd",
					Npwp:              "123",
					AccountNumber:     "123asd",
					AccountHolderName: "",
					BankBranch:        "bran",
				},
			},
			wantErr: true,
		},
		{
			name:   "Case 7: Error BankBranch",
			fields: fields{},
			args: args{
				data: &model.B2CMerchantCreateInput{
					NpwpFile:          "1",
					NpwpHolderName:    "asd",
					Npwp:              "123",
					AccountNumber:     "123asd",
					AccountHolderName: "asdasd",
					BankBranch:        "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantUseCaseImpl{
				Repository:           tt.fields.Repository,
				MerchantRepo:         tt.fields.MerchantRepo,
				MerchantAddressRepo:  tt.fields.MerchantAddressRepo,
				MerchantBankRepo:     tt.fields.MerchantBankRepo,
				MerchantDocumentRepo: tt.fields.MerchantDocumentRepo,
				MemberRepoRead:       tt.fields.MemberRepoRead,
				UploadService:        tt.fields.UploadService,
				MerchantService:      tt.fields.MerchantService,
				TokenGenerator:       tt.fields.TokenGenerator,
				MemberQueryRead:      tt.fields.MemberQueryRead,
				NotificationService:  tt.fields.NotificationService,
				QueuePublisher:       tt.fields.QueuePublisher,
				SendbirdService:      tt.fields.SendbirdService,
			}
			if err := m.validateMerchantFile(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("MerchantUseCaseImpl.validateMerchantFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
