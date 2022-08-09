package repo

import (
	"testing"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared/mocks"
	"github.com/stretchr/testify/assert"
)

func setupRepoTokenActivationRepoRedis() *TokenActivationRepoRedis {
	client := mocks.InitFakeRedis()
	return &TokenActivationRepoRedis{
		client: client,
	}
}

func TestNewTokenActivationRepoRedis(t *testing.T) {
	t.Run("POSITIVE_NEW_TOKEN_ACTIVATION_REPO_REDIS", func(t *testing.T) {
		fredis := mocks.InitFakeRedis()
		repoTokenRedis := NewTokenActivationRepoRedis(fredis)
		assert.NotNil(t, repoTokenRedis)
	})
}

func TestTokenActivationRepoRedisSave(t *testing.T) {
	t.Run("NEGATIVE_TOKEN_ACTIVATION_REPO_REDIS_SAVE", func(t *testing.T) {
		repoTokenRedis := setupRepoTokenActivationRepoRedis()
		result := <-repoTokenRedis.Save(&model.TokenActivation{})
		assert.Error(t, result.Error)
	})
}

func TestTokenActivationRepoRedisLoad(t *testing.T) {
	t.Run("NEGATIVE_TOKEN_ACTIVATION_REPO_REDIS_LOAD", func(t *testing.T) {
		repoTokenRedis := setupRepoTokenActivationRepoRedis()
		result := <-repoTokenRedis.Load("userload1")
		assert.Error(t, result.Error)
	})
}

func TestTokenActivationRepoRedisDelete(t *testing.T) {
	t.Run("NEGATIVE_TOKEN_ACTIVATION_REPO_REDIS_DELETE", func(t *testing.T) {
		repoTokenRedis := setupRepoTokenActivationRepoRedis()
		result := <-repoTokenRedis.Delete("userdelete1")
		assert.Error(t, result.Error)
	})
}
