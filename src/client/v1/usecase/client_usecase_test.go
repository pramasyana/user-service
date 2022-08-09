package usecase

import (
	"database/sql"
	"testing"

	"github.com/Bhinneka/user-service/config/redis"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
)

func TestNewClientUsecase(t *testing.T) {
	type args struct {
		loginSession authRepo.LoginSessionRepository
		rt           authRepo.RefreshTokenRepository
		mq           memberQuery.MemberQuery
	}
	tests := []struct {
		name string
		args args
		want *ClientUC
	}{
		{
			name: "Case 1: Success",
			args: args{
				loginSession: authRepo.NewLoginSessionRepositoryRedis(redis.Conn{}),
				rt:           authRepo.NewRefreshTokenRepositoryRedis(redis.Conn{}),
				mq:           memberQuery.NewMemberQueryPostgres(sql.OpenDB(nil)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewClientUsecase(tt.args.loginSession, tt.args.rt, tt.args.mq)
		})
	}
}
