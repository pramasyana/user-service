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

func generateResultRepository(data repo.ResultRepository) <-chan repo.ResultRepository {
	output := make(chan repo.ResultRepository, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestMerchantUseCaseImpl_ValidateDocument(t *testing.T) {
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
		document *model.B2CMerchantDocumentInput
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
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("FindMerchantDocumentByParam", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Result: model.B2CMerchantDocumentData{},
							Error:  nil,
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				document: &model.B2CMerchantDocumentInput{
					DocumentType: "apa",
				},
			},
			wantErr: false,
		},
		{
			name: "Case 1: Success",
			fields: fields{
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("FindMerchantDocumentByParam", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Result: model.B2CMerchantDocumentData{
								DocumentValue: "ini sisinya",
								MerchantID:    "MRCNT0001",
							},
							Error: nil,
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				document: &model.B2CMerchantDocumentInput{
					DocumentType: "apa",
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
			if err := m.ValidateDocument(tt.args.ctxReq, tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("MerchantUseCaseImpl.ValidateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerchantUseCaseImpl_InsertUpdateDocument(t *testing.T) {
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
		ctxReq           context.Context
		document         model.B2CMerchantDocumentInput
		merchantDocument model.B2CMerchantDocumentData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success InsertNewMerchantDocument",
			fields: fields{
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("InsertNewMerchantDocument", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq:           context.Background(),
				document:         model.B2CMerchantDocumentInput{},
				merchantDocument: model.B2CMerchantDocumentData{},
			},
		},
		{
			name: "Case 2: Error InsertNewMerchantDocument",
			fields: fields{
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("InsertNewMerchantDocument", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: errors.New("error"),
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq:           context.Background(),
				document:         model.B2CMerchantDocumentInput{},
				merchantDocument: model.B2CMerchantDocumentData{},
			},
		},
		{
			name: "Case 3: Success UpdateMerchantDocument",
			fields: fields{
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("UpdateMerchantDocument", mock.Anything, mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				document: model.B2CMerchantDocumentInput{
					ID: "DCM001",
				},
				merchantDocument: model.B2CMerchantDocumentData{},
			},
		},
		{
			name: "Case 4: Error UpdateMerchantDocument",
			fields: fields{
				MerchantDocumentRepo: func() repo.MerchantDocumentRepository {
					mocksMerchantDocumentRepo := new(mocksMerchantRepo.MerchantDocumentRepository)
					mocksMerchantDocumentRepo.On("UpdateMerchantDocument", mock.Anything, mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: errors.New("Error"),
						}),
					)
					return mocksMerchantDocumentRepo
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				document: model.B2CMerchantDocumentInput{
					ID: "DCM001",
				},
				merchantDocument: model.B2CMerchantDocumentData{},
			},
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
			m.InsertUpdateDocument(tt.args.ctxReq, tt.args.document, tt.args.merchantDocument)
		})
	}
}
