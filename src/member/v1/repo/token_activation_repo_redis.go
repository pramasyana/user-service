package repo

import (
	"context"

	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// TokenActivationRepoRedis data structure
type TokenActivationRepoRedis struct {
	client redis.Client
}

// NewTokenActivationRepoRedis function for initializing member repository redis
func NewTokenActivationRepoRedis(cl redis.Client) *TokenActivationRepoRedis {
	return &TokenActivationRepoRedis{cl}
}

// Save function for saving token data into redis DB
func (repo *TokenActivationRepoRedis) Save(ta *model.TokenActivation) <-chan ResultRepository {
	ctx := "TokenActivationRepoRedis-Save"

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		cl := repo.client

		_, err := cl.Set(ta.ID, ta.Value, 0)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, "save_redis_key", err, ta)
			output <- ResultRepository{Error: err}
			return
		}

		_, err = cl.Expire(ta.ID, ta.TTL)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, "redis_expired", err, ta)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	}()

	return output
}

// Load function for loading token data by its ID
func (repo *TokenActivationRepoRedis) Load(id string) <-chan ResultRepository {
	ctx := "TokenActivationRepoRedis-Load"

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		cl := repo.client

		val, err := cl.Get(id)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, "get_keys", err, id)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: model.TokenActivation{ID: id, Value: val}}
	}()

	return output
}

// Delete function for delete token data into redis DB
func (repo *TokenActivationRepoRedis) Delete(id string) <-chan ResultRepository {
	ctx := "TokenActivationRepoRedis-Delete"

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		cl := repo.client

		_, err := cl.Del(id)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, "redis_delete_key", err, id)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	}()

	return output
}
