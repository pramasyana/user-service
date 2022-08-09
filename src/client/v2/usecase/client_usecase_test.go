package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Bhinneka/user-service/config/redis"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/service"
	"github.com/stretchr/testify/mock"

	mocksRepoLoginSession "github.com/Bhinneka/user-service/mocks/src/auth/v1/repo"
	mocksRepoCorporateContact "github.com/Bhinneka/user-service/mocks/src/corporate/v2/query"
	sharedDomain "github.com/Bhinneka/user-service/src/shared/model"
)

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

func generateUsecaseResult(data corporateQuery.ResultQuery) <-chan corporateQuery.ResultQuery {
	output := make(chan corporateQuery.ResultQuery, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateResultRepository(data authRepo.ResultRepository) <-chan authRepo.ResultRepository {
	output := make(chan authRepo.ResultRepository, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestNewClientUsecase(t *testing.T) {
	type args struct {
		loginSession authRepo.LoginSessionRepository
		rt           authRepo.RefreshTokenRepository
		mq           memberQuery.MemberQuery
		cc           corporateQuery.ContactQuery
	}
	tests := []struct {
		name string
		args args
		want *ClientUC
	}{
		{
			name: "Success",
			args: args{
				loginSession: authRepo.NewLoginSessionRepositoryRedis(redis.Conn{}),
				rt:           authRepo.NewRefreshTokenRepositoryRedis(redis.Conn{}),
				mq:           memberQuery.NewMemberQueryPostgres(sql.OpenDB(nil)),
				cc:           corporateQuery.NewContactQueryPostgres(sql.OpenDB(nil)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewClientUsecase(tt.args.loginSession, tt.args.rt, tt.args.mq, tt.args.cc)
		})
	}
}

func TestClientUC_Logout(t *testing.T) {
	type fields struct {
		LoginSessionRepo          authRepo.LoginSessionRepository
		RefreshTokenRepo          authRepo.RefreshTokenRepository
		MemberQueryRead           memberQuery.MemberQuery
		CorporateContactQueryRead corporateQuery.ContactQuery
		QPublisher                service.QPublisher
	}
	type args struct {
		ctxReq context.Context
		email  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				LoginSessionRepo: func() authRepo.LoginSessionRepository {
					mocksRepoLoginSession := new(mocksRepoLoginSession.LoginSessionRepository)
					mocksRepoLoginSession.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoLoginSession
				}(),
				RefreshTokenRepo: func() authRepo.RefreshTokenRepository {
					mocksRepoRefreshToken := new(mocksRepoLoginSession.RefreshTokenRepository)
					mocksRepoRefreshToken.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoRefreshToken
				}(),
				MemberQueryRead: memberQuery.NewMemberQueryPostgres(sql.OpenDB(nil)),
				QPublisher:      &service.KafkaPublisherImpl{},
				CorporateContactQueryRead: func() corporateQuery.ContactQuery {
					mocksRepoCorporateContact := new(mocksRepoCorporateContact.ContactQuery)
					mocksRepoCorporateContact.On("FindContactMicrositeByEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						generateUsecaseResult(corporateQuery.ResultQuery{
							Result: sharedDomain.B2BContactData{},
							Error:  nil,
						}),
					)
					return mocksRepoCorporateContact
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				email:  "a@gmail.com",
			},
		},
		{
			name: "Case 2: Error FindContactMicrositeByEmail",
			fields: fields{
				CorporateContactQueryRead: func() corporateQuery.ContactQuery {
					mocksRepoCorporateContact := new(mocksRepoCorporateContact.ContactQuery)
					mocksRepoCorporateContact.On("FindContactMicrositeByEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						generateUsecaseResult(corporateQuery.ResultQuery{
							Result: sharedDomain.B2BContactData{},
							Error:  errors.New("error"),
						}),
					)
					return mocksRepoCorporateContact
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				email:  "a@gmail.com",
			},
		},
		{
			name: "Case 3: Success FindContactMicrositeByEmail & failed convert",
			fields: fields{
				CorporateContactQueryRead: func() corporateQuery.ContactQuery {
					mocksRepoCorporateContact := new(mocksRepoCorporateContact.ContactQuery)
					mocksRepoCorporateContact.On("FindContactMicrositeByEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						generateUsecaseResult(corporateQuery.ResultQuery{
							Error: nil,
						}),
					)
					return mocksRepoCorporateContact
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				email:  "a@gmail.com",
			},
		},
		{
			name: "Case 4: Error Login session",
			fields: fields{
				LoginSessionRepo: func() authRepo.LoginSessionRepository {
					mocksRepoLoginSession := new(mocksRepoLoginSession.LoginSessionRepository)
					mocksRepoLoginSession.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: errors.New("error"),
						}),
					)
					return mocksRepoLoginSession
				}(),
				RefreshTokenRepo: func() authRepo.RefreshTokenRepository {
					mocksRepoRefreshToken := new(mocksRepoLoginSession.RefreshTokenRepository)
					mocksRepoRefreshToken.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoRefreshToken
				}(),
				MemberQueryRead: memberQuery.NewMemberQueryPostgres(sql.OpenDB(nil)),
				QPublisher:      &service.KafkaPublisherImpl{},
				CorporateContactQueryRead: func() corporateQuery.ContactQuery {
					mocksRepoCorporateContact := new(mocksRepoCorporateContact.ContactQuery)
					mocksRepoCorporateContact.On("FindContactMicrositeByEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						generateUsecaseResult(corporateQuery.ResultQuery{
							Result: sharedDomain.B2BContactData{},
							Error:  nil,
						}),
					)
					return mocksRepoCorporateContact
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				email:  "a@gmail.com",
			},
		},
		{
			name: "Case 5: Error refresh token",
			fields: fields{
				LoginSessionRepo: func() authRepo.LoginSessionRepository {
					mocksRepoLoginSession := new(mocksRepoLoginSession.LoginSessionRepository)
					mocksRepoLoginSession.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: nil,
						}),
					)
					return mocksRepoLoginSession
				}(),
				RefreshTokenRepo: func() authRepo.RefreshTokenRepository {
					mocksRepoRefreshToken := new(mocksRepoLoginSession.RefreshTokenRepository)
					mocksRepoRefreshToken.On("Delete", mock.Anything, mock.Anything).Return(
						generateResultRepository(authRepo.ResultRepository{
							Error: errors.New("error"),
						}),
					)
					return mocksRepoRefreshToken
				}(),
				MemberQueryRead: memberQuery.NewMemberQueryPostgres(sql.OpenDB(nil)),
				QPublisher:      &service.KafkaPublisherImpl{},
				CorporateContactQueryRead: func() corporateQuery.ContactQuery {
					mocksRepoCorporateContact := new(mocksRepoCorporateContact.ContactQuery)
					mocksRepoCorporateContact.On("FindContactMicrositeByEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						generateUsecaseResult(corporateQuery.ResultQuery{
							Result: sharedDomain.B2BContactData{},
							Error:  nil,
						}),
					)
					return mocksRepoCorporateContact
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				email:  "a@gmail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := &ClientUC{
				LoginSessionRepo:          tt.fields.LoginSessionRepo,
				RefreshTokenRepo:          tt.fields.RefreshTokenRepo,
				MemberQueryRead:           tt.fields.MemberQueryRead,
				CorporateContactQueryRead: tt.fields.CorporateContactQueryRead,
				QPublisher:                tt.fields.QPublisher,
			}
			au.Logout(tt.args.ctxReq, tt.args.email)
		})
	}
}
