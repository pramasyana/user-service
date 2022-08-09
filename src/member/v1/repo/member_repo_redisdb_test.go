package repo

import (
	"context"
	"testing"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared/mocks"
	"github.com/stretchr/testify/assert"
)

func setupRepoMemberRedis() *MemberRepoRedis {
	client := mocks.InitFakeRedis()
	return &MemberRepoRedis{
		client: client,
	}
}

func TestNewMemberRepoRedis(t *testing.T) {
	t.Run("POSITIVE_NEW_MEMBER_REPO_REDIS", func(t *testing.T) {
		fredis := mocks.InitFakeRedis()
		r := NewMemberRepoRedis(fredis)
		assert.NotNil(t, r)
	})
}

func TestMemberRepoRedisSave(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_SAVE", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.Save(&model.MemberRedis{})
		assert.Error(t, result.Error)
	})
}

func TestMemberRepoRedisLoad(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_LOAD", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.Load("usr123")
		assert.Error(t, result.Error)
	})
}

func TestMemberRepoRedisDelete(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_DELETE", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.Delete("usr1234")
		assert.Error(t, result.Error)
	})
}

func TestMemberRepoRedisLoadByKey(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_LOAD_BY_KEY", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.LoadByKey(context.Background(), "load_by_key")
		assert.Error(t, result.Error)
	})
}

func TestMemberRepoRedisSaveResendActivationAttempt(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_ACTIVATION_ATTEMPT", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.SaveResendActivationAttempt(context.Background(), &model.ResendActivationAttempt{})
		assert.Error(t, result.Error)
	})
}

func TestMemberRepoRedisRevokeAllAccess(t *testing.T) {
	t.Run("NEGATIVE_MEMBER_REPO_REDIS_REVOKE_ALL_ACCESS", func(t *testing.T) {
		r := setupRepoMemberRedis()
		result := <-r.RevokeAllAccess(context.Background(), "revoke_key", "revoke_active_key")
		assert.Error(t, result.Error)
	})
}
