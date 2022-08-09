package usecase

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

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
	"github.com/stretchr/testify/assert"
)

func TestBulkEmail(t *testing.T) {
	data := []*model.Member{
		{
			Email: "myemail@bhinneka.com",
		},
		{
			Email: "otheremail@bhinneka.com",
		},
		{
			Email: "myemails@domain.com",
		},
		{
			Email: "string",
		},
	}

	muc := MemberUseCaseImpl{}
	errs := muc.BulkValidateEmailAndPhone(context.Background(), data)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "email string is invalid", strings.Join(errs, "-"))
}

func Test_cleanTags(t *testing.T) {
	type args struct {
		in0  context.Context
		data *model.Member
	}
	tests := []struct {
		name           string
		args           args
		wantMaskedData *model.Member
	}{
		{
			name: "Test 1",
			args: args{
				in0:  context.Background(),
				data: &model.Member{},
			},
			wantMaskedData: &model.Member{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMaskedData := cleanTags(tt.args.in0, tt.args.data); !reflect.DeepEqual(gotMaskedData, tt.wantMaskedData) {
				t.Errorf("cleanTags() = %v, want %v", gotMaskedData, tt.wantMaskedData)
			}
		})
	}
}

func TestMemberUseCaseImpl_CheckEmailAndMobileExistence(t *testing.T) {
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
		MerchantEmployeeRead              merchantRepoRead.MerchantEmployeeRepository
	}
	type args struct {
		ctxReq context.Context
		data   *model.Member
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
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
				MerchantEmployeeRead:              tt.fields.MerchantEmployeeRead,
			}
			if got := mu.CheckEmailAndMobileExistence(tt.args.ctxReq, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemberUseCaseImpl.CheckEmailAndMobileExistence() = %v, want %v", got, tt.want)
			}
		})
	}
}
