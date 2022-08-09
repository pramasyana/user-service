package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// RefreshTokenRepositoryRedis data structure
type RefreshTokenRepositoryRedis struct {
	client redis.Client
}

// NewRefreshTokenRepositoryRedis function for initializing RefreshTokenRepositoryRedis
func NewRefreshTokenRepositoryRedis(cl redis.Client) *RefreshTokenRepositoryRedis {
	return &RefreshTokenRepositoryRedis{cl}
}

// Save function for saving refresh token into redis DB
func (repo *RefreshTokenRepositoryRedis) Save(ctxReq context.Context, refreshToken *model.RefreshToken) <-chan ResultRepository {
	ctx := "RefreshTokenRepositoryRedis-Save"

	outputSaveRefreshToken := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputSaveRefreshToken)

		tags[helper.TagKey] = refreshToken.ID

		if _, err := repo.client.Set(refreshToken.ID, refreshToken.Token, refreshToken.RefreshTokenAge); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_to_redis", err, refreshToken)
			outputSaveRefreshToken <- ResultRepository{Error: err}
			return
		}

		outputSaveRefreshToken <- ResultRepository{Result: refreshToken}

	})

	return outputSaveRefreshToken
}

// Load function for loading refresh token by its key
func (repo *RefreshTokenRepositoryRedis) Load(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "RefreshTokenRepositoryRedis-Load"

	outputLoadRT := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputLoadRT)

		val, err := repo.client.Get(key)
		tags[helper.TagKey] = key

		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "laod_from_redis", err, key)
			}
			outputLoadRT <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		outputLoadRT <- ResultRepository{Result: model.RefreshToken{Token: val}}
		tags[helper.TextResponse] = model.RefreshToken{Token: val}
	})

	return outputLoadRT
}

// Delete function for deleting refresh token by its key
func (repo *RefreshTokenRepositoryRedis) Delete(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "RefreshTokenRepositoryRedis-Delete"

	outputDeleteRT := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputDeleteRT)
		tags[helper.TagKey] = key

		if _, err := repo.client.Del(key); err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "delete_from_redis", err, key)
			}
			outputDeleteRT <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		outputDeleteRT <- ResultRepository{Result: nil}
		tags[helper.TextResponse] = nil

	})

	return outputDeleteRT
}
