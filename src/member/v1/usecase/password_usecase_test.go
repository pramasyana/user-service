package usecase

import (
	"context"
	"testing"
	"time"

	mocksRepoMember "github.com/Bhinneka/user-service/mocks/src/member/v1/repo"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	authUsecase "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepoRead "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	service "github.com/Bhinneka/user-service/src/service"
	sessionQuery "github.com/Bhinneka/user-service/src/session/v1/query"
	shippingRepo "github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
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

func TestMemberUseCaseImpl_SyncPassword(t *testing.T) {
	type fields struct {
		MemberRepoRead                    repo.MemberRepository
		MemberRepoWrite                   repo.MemberRepository
		MemberMFARepoWrite                repo.MemberMFARepository
		MemberRepoRedis                   repo.MemberRepositoryRedis
		TokenActivationRepo               repo.TokenActivationRepository
		ShippingAddressRepo               shippingRepo.ShippingAddressRepository
		LoginAttemptRepo                  authRepo.AttemptRepository
		LoginSessionRedis                 authRepo.LoginSessionRepository
		MemberQueryRead                   query.MemberQuery
		MemberQueryWrite                  query.MemberQuery
		MemberMFAQueryRead                query.MemberMFAQuery
		SessionQueryRead                  sessionQuery.SessionInfoQuery
		StaticService                     service.StaticServices
		UploadService                     service.UploadServices
		ActivityService                   service.ActivityServices
		QPublisher                        service.QPublisher
		Hash                              model.PasswordHasher
		TokenActivationExpiration         time.Duration
		ResendActivationAttemptAge        string
		ResendActivationAttemptAgeRequest string
		Topic                             string
		IsProductionStage                 bool
		SturgeonCFUrl                     string
		B2cCFUrl                          string
		AccessTokenGenerator              token.AccessTokenGenerator
		NotificationService               service.NotificationServices
		SendbirdService                   service.SendbirdServices
		AuthUseCase                       authUsecase.AuthUseCase
		CorporateContactQueryRead         corporateQuery.ContactQuery
		CorporateAccContactQueryRead      corporateQuery.AccountContactQuery
		MerchantRepoRead                  merchantRepoRead.MerchantRepository
	}
	type args struct {
		ctxReq      context.Context
		token       string
		uid         string
		oldPassword string
		newPassword string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MemberRepoRead: func() repo.MemberRepository {
					mocksRepoMember := new(mocksRepoMember.MemberRepository)
					mocksRepoMember.On("Load", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoMember
				}(),
			},
			args: args{
				ctxReq:      context.Background(),
				token:       testing.CoverMode(),
				uid:         testing.CoverMode(),
				oldPassword: "Bhinneka123!",
				newPassword: "Bhinneka123!@#",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &MemberUseCaseImpl{
				MemberRepoRead:                    tt.fields.MemberRepoRead,
				MemberRepoWrite:                   tt.fields.MemberRepoWrite,
				MemberMFARepoWrite:                tt.fields.MemberMFARepoWrite,
				MemberRepoRedis:                   tt.fields.MemberRepoRedis,
				TokenActivationRepo:               tt.fields.TokenActivationRepo,
				ShippingAddressRepo:               tt.fields.ShippingAddressRepo,
				LoginAttemptRepo:                  tt.fields.LoginAttemptRepo,
				LoginSessionRedis:                 tt.fields.LoginSessionRedis,
				MemberQueryRead:                   tt.fields.MemberQueryRead,
				MemberQueryWrite:                  tt.fields.MemberQueryWrite,
				MemberMFAQueryRead:                tt.fields.MemberMFAQueryRead,
				SessionQueryRead:                  tt.fields.SessionQueryRead,
				StaticService:                     tt.fields.StaticService,
				UploadService:                     tt.fields.UploadService,
				ActivityService:                   tt.fields.ActivityService,
				QPublisher:                        tt.fields.QPublisher,
				Hash:                              tt.fields.Hash,
				TokenActivationExpiration:         tt.fields.TokenActivationExpiration,
				ResendActivationAttemptAge:        tt.fields.ResendActivationAttemptAge,
				ResendActivationAttemptAgeRequest: tt.fields.ResendActivationAttemptAgeRequest,
				Topic:                             tt.fields.Topic,
				IsProductionStage:                 tt.fields.IsProductionStage,
				SturgeonCFUrl:                     tt.fields.SturgeonCFUrl,
				B2cCFUrl:                          tt.fields.B2cCFUrl,
				AccessTokenGenerator:              tt.fields.AccessTokenGenerator,
				NotificationService:               tt.fields.NotificationService,
				SendbirdService:                   tt.fields.SendbirdService,
				AuthUseCase:                       tt.fields.AuthUseCase,
				CorporateContactQueryRead:         tt.fields.CorporateContactQueryRead,
				CorporateAccContactQueryRead:      tt.fields.CorporateAccContactQueryRead,
				MerchantRepoRead:                  tt.fields.MerchantRepoRead,
			}
			mu.SyncPassword(tt.args.ctxReq, tt.args.token, tt.args.oldPassword, tt.args.newPassword)
		})
	}
}

func TestMemberUseCaseImpl_ChangePassword(t *testing.T) {
	type fields struct {
		MemberRepoRead                    repo.MemberRepository
		MemberRepoWrite                   repo.MemberRepository
		MemberMFARepoWrite                repo.MemberMFARepository
		MemberRepoRedis                   repo.MemberRepositoryRedis
		TokenActivationRepo               repo.TokenActivationRepository
		ShippingAddressRepo               shippingRepo.ShippingAddressRepository
		LoginAttemptRepo                  authRepo.AttemptRepository
		LoginSessionRedis                 authRepo.LoginSessionRepository
		MemberQueryRead                   query.MemberQuery
		MemberQueryWrite                  query.MemberQuery
		MemberMFAQueryRead                query.MemberMFAQuery
		SessionQueryRead                  sessionQuery.SessionInfoQuery
		StaticService                     service.StaticServices
		UploadService                     service.UploadServices
		ActivityService                   service.ActivityServices
		QPublisher                        service.QPublisher
		Hash                              model.PasswordHasher
		TokenActivationExpiration         time.Duration
		ResendActivationAttemptAge        string
		ResendActivationAttemptAgeRequest string
		Topic                             string
		IsProductionStage                 bool
		SturgeonCFUrl                     string
		B2cCFUrl                          string
		AccessTokenGenerator              token.AccessTokenGenerator
		NotificationService               service.NotificationServices
		SendbirdService                   service.SendbirdServices
		AuthUseCase                       authUsecase.AuthUseCase
		CorporateContactQueryRead         corporateQuery.ContactQuery
		CorporateAccContactQueryRead      corporateQuery.AccountContactQuery
		MerchantRepoRead                  merchantRepoRead.MerchantRepository
	}
	type args struct {
		ctxReq      context.Context
		token       string
		uid         string
		oldPassword string
		newPassword string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MemberRepoRead: func() repo.MemberRepository {
					mocksRepoMember := new(mocksRepoMember.MemberRepository)
					mocksRepoMember.On("Load", mock.Anything, mock.Anything).Return(
						generateResultRepository(repo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoMember
				}(),
			},
			args: args{
				ctxReq:      context.Background(),
				token:       testing.CoverMode(),
				uid:         testing.CoverMode(),
				oldPassword: "Bhinneka123!",
				newPassword: "Bhinneka123!@#",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &MemberUseCaseImpl{
				MemberRepoRead:                    tt.fields.MemberRepoRead,
				MemberRepoWrite:                   tt.fields.MemberRepoWrite,
				MemberMFARepoWrite:                tt.fields.MemberMFARepoWrite,
				MemberRepoRedis:                   tt.fields.MemberRepoRedis,
				TokenActivationRepo:               tt.fields.TokenActivationRepo,
				ShippingAddressRepo:               tt.fields.ShippingAddressRepo,
				LoginAttemptRepo:                  tt.fields.LoginAttemptRepo,
				LoginSessionRedis:                 tt.fields.LoginSessionRedis,
				MemberQueryRead:                   tt.fields.MemberQueryRead,
				MemberQueryWrite:                  tt.fields.MemberQueryWrite,
				MemberMFAQueryRead:                tt.fields.MemberMFAQueryRead,
				SessionQueryRead:                  tt.fields.SessionQueryRead,
				StaticService:                     tt.fields.StaticService,
				UploadService:                     tt.fields.UploadService,
				ActivityService:                   tt.fields.ActivityService,
				QPublisher:                        tt.fields.QPublisher,
				Hash:                              tt.fields.Hash,
				TokenActivationExpiration:         tt.fields.TokenActivationExpiration,
				ResendActivationAttemptAge:        tt.fields.ResendActivationAttemptAge,
				ResendActivationAttemptAgeRequest: tt.fields.ResendActivationAttemptAgeRequest,
				Topic:                             tt.fields.Topic,
				IsProductionStage:                 tt.fields.IsProductionStage,
				SturgeonCFUrl:                     tt.fields.SturgeonCFUrl,
				B2cCFUrl:                          tt.fields.B2cCFUrl,
				AccessTokenGenerator:              tt.fields.AccessTokenGenerator,
				NotificationService:               tt.fields.NotificationService,
				SendbirdService:                   tt.fields.SendbirdService,
				AuthUseCase:                       tt.fields.AuthUseCase,
				CorporateContactQueryRead:         tt.fields.CorporateContactQueryRead,
				CorporateAccContactQueryRead:      tt.fields.CorporateAccContactQueryRead,
				MerchantRepoRead:                  tt.fields.MerchantRepoRead,
			}
			mu.ChangePassword(tt.args.ctxReq, tt.args.token, tt.args.oldPassword, tt.args.newPassword)
		})
	}
}
