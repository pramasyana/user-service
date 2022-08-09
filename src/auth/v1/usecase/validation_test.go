package usecase

import (
	"net/url"
	"testing"

	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/query"
	"github.com/Bhinneka/user-service/src/auth/v1/repo"
	authServices "github.com/Bhinneka/user-service/src/auth/v1/service"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepoRead "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	sessionInfoRepo "github.com/Bhinneka/user-service/src/session/v1/repo"
)

func TestAuthUseCaseImpl_validateInput(t *testing.T) {
	type fields struct {
		ClientAppRepoRead            repo.ClientAppRepository
		ClientAppRepoWrite           repo.ClientAppRepository
		AuthQueryOAuth               query.AuthQueryOA
		AuthQueryDB                  query.AuthQuery
		MemberRepoRead               memberRepo.MemberRepository
		MemberRepoWrite              memberRepo.MemberRepository
		MemberQueryRead              memberQuery.MemberQuery
		MemberQueryWrite             memberQuery.MemberQuery
		CorporateContactQueryRead    corporateQuery.ContactQuery
		CorporateAccContactQueryRead corporateQuery.AccountContactQuery
		RefreshTokenRepo             repo.RefreshTokenRepository
		LoginAttemptRepo             repo.AttemptRepository
		LoginSessionRepo             repo.LoginSessionRepository
		AccessTokenGenerator         token.AccessTokenGenerator
		Hash                         memberModel.PasswordHasher
		RefreshTokenAge              string
		SpecialRefreshTokenAge       string
		EmailSpecialTokenAge         string
		LoginAttemptAge              string
		QPublisher                   service.QPublisher
		Topic                        string
		IsProductionStage            bool
		AuthServices                 authServices.LDAPService
		SessionInfoRepo              sessionInfoRepo.SessionInfoRepository
		MerchantRepoRead             merchantRepoRead.MerchantRepository
		GoogleVerifyCaptchaURL       *url.URL
		StaticService                service.StaticServices
		ActivityService              service.ActivityServices
		B2cCFUrl                     string
		NotificationService          service.NotificationServices
	}
	type args struct {
		data *model.RequestToken
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Case 1: Error empty memberType",
			args: args{
				data: &model.RequestToken{
					MemberType: "percobaan",
				},
			},
			wantErr: true,
		},
		{
			name: "Case 2: Error empty device",
			args: args{
				data: &model.RequestToken{
					MemberType: model.UserTypeCorporate,
					GrantType:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "Case 3: Error empty code",
			args: args{
				data: &model.RequestToken{
					MemberType: model.UserTypeCorporate,
					GrantType:  model.AuthTypeFacebook,
					Code:       "",
					DeviceID:   "BelaID",
				},
			},
			wantErr: true,
		},
		{
			name: "Case 4: Error empty redirect uri",
			args: args{
				data: &model.RequestToken{
					MemberType: model.UserTypeCorporate,
					GrantType:  model.AuthTypeAzure,
					Code:       "PERCOBAAN",
					DeviceID:   "BelaID",
				},
			},
			wantErr: true,
		},
		{
			name: "Case 5: Error empty deviceLogin",
			args: args{
				data: &model.RequestToken{
					MemberType:  model.UserTypeCorporate,
					GrantType:   model.AuthTypePassword,
					Code:        "PERCOBAAN",
					DeviceID:    "BelaID",
					DeviceLogin: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Case 6: Error empty memberType",
			args: args{
				data: &model.RequestToken{
					MemberType:  "",
					GrantType:   model.AuthTypePassword,
					Code:        "PERCOBAAN",
					DeviceID:    "BelaID",
					DeviceLogin: "WEB",
				},
			},
			wantErr: false,
		},
		{
			name: "Case 7: Success",
			args: args{
				data: &model.RequestToken{
					MemberType:  model.UserTypeCorporate,
					GrantType:   model.AuthTypeFacebook,
					DeviceID:    "WEB",
					Code:        "test",
					DeviceLogin: "WEB",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &AuthUseCaseImpl{
				ClientAppRepoRead:            tt.fields.ClientAppRepoRead,
				ClientAppRepoWrite:           tt.fields.ClientAppRepoWrite,
				AuthQueryOAuth:               tt.fields.AuthQueryOAuth,
				AuthQueryDB:                  tt.fields.AuthQueryDB,
				MemberRepoRead:               tt.fields.MemberRepoRead,
				MemberRepoWrite:              tt.fields.MemberRepoWrite,
				MemberQueryRead:              tt.fields.MemberQueryRead,
				MemberQueryWrite:             tt.fields.MemberQueryWrite,
				CorporateContactQueryRead:    tt.fields.CorporateContactQueryRead,
				CorporateAccContactQueryRead: tt.fields.CorporateAccContactQueryRead,
				RefreshTokenRepo:             tt.fields.RefreshTokenRepo,
				LoginAttemptRepo:             tt.fields.LoginAttemptRepo,
				LoginSessionRepo:             tt.fields.LoginSessionRepo,
				AccessTokenGenerator:         tt.fields.AccessTokenGenerator,
				Hash:                         tt.fields.Hash,
				RefreshTokenAge:              tt.fields.RefreshTokenAge,
				SpecialRefreshTokenAge:       tt.fields.SpecialRefreshTokenAge,
				EmailSpecialTokenAge:         tt.fields.EmailSpecialTokenAge,
				LoginAttemptAge:              tt.fields.LoginAttemptAge,
				QPublisher:                   tt.fields.QPublisher,
				Topic:                        tt.fields.Topic,
				IsProductionStage:            tt.fields.IsProductionStage,
				AuthServices:                 tt.fields.AuthServices,
				SessionInfoRepo:              tt.fields.SessionInfoRepo,
				MerchantRepoRead:             tt.fields.MerchantRepoRead,
				GoogleVerifyCaptchaURL:       tt.fields.GoogleVerifyCaptchaURL,
				StaticService:                tt.fields.StaticService,
				ActivityService:              tt.fields.ActivityService,
				B2cCFUrl:                     tt.fields.B2cCFUrl,
				NotificationService:          tt.fields.NotificationService,
			}
			if err := au.validateInput(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AuthUseCaseImpl.validateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
