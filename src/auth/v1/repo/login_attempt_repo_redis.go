package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// AttemptRepositoryRedis data structure
type AttemptRepositoryRedis struct {
	client redis.Client
}

// NewAttemptRepositoryRedis function for initializing AttemptRepositoryRedis
func NewAttemptRepositoryRedis(cl redis.Client) *AttemptRepositoryRedis {
	return &AttemptRepositoryRedis{cl}
}

// Save function for saving refresh token into redis DB
func (repo *AttemptRepositoryRedis) Save(ctxReq context.Context, data *model.LoginAttempt) <-chan ResultRepository {
	ctx := "AttemptRepositoryRedis-Save"

	outputSave := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputSave)
		tags[helper.TagKey] = data.Key

		_, err := repo.client.Set(data.Key, data.Attempt, data.LoginAttemptAge)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_to_redis", err, data)
			outputSave <- ResultRepository{Error: err}
			return
		}

		outputSave <- ResultRepository{Result: data}

	})

	return outputSave
}

// Load function for loading refresh token by its key
func (repo *AttemptRepositoryRedis) Load(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "AttemptRepositoryRedis-Load"

	outputLoad := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputLoad)

		val, err := repo.client.Get(key)
		tags[helper.TagKey] = key

		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "load_from_redis", err, key)
			}
			outputLoad <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		outputLoad <- ResultRepository{Result: model.LoginAttempt{Attempt: val}}
		tags[helper.TextResponse] = model.LoginAttempt{Attempt: val}

	})

	return outputLoad
}

// Delete function for delete login attempt data into redis DB
func (repo *AttemptRepositoryRedis) Delete(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "AttemptRepositoryRedis-Delete"

	outputDelete := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputDelete)
		tags[helper.TagKey] = key

		if _, err := repo.client.Del(key); err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "delete_from_redis", err, key)
			}
			outputDelete <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		outputDelete <- ResultRepository{Error: nil}
		tags[helper.TextResponse] = nil
	})

	return outputDelete
}
