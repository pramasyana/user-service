package repo

import (
	"context"
	"fmt"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// LoginSessionRepositoryRedis data structure
type LoginSessionRepositoryRedis struct {
	client redis.Client
}

// NewLoginSessionRepositoryRedis function for initializing AttemptRepositoryRedis
func NewLoginSessionRepositoryRedis(cl redis.Client) *LoginSessionRepositoryRedis {
	return &LoginSessionRepositoryRedis{cl}
}

// Save function for saving refresh token into redis DB
func (repo *LoginSessionRepositoryRedis) Save(ctxReq context.Context, data *model.LoginSessionRedis) <-chan ResultRepository {
	ctx := "LoginSessionRepositoryRedis-Save"

	outputLoginSession := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputLoginSession)
		tags[helper.TagKey] = data.Key
		tags[helper.TextToken] = data.Token

		if _, err := repo.client.Set(data.Key, data.Token, data.ExpiredTime); err != nil {
			tags[helper.TagError] = err.Error()
			helper.SendErrorLog(ctxReq, ctx, "save_token_to_redis", err, data)
			outputLoginSession <- ResultRepository{Error: err}
			return
		}

		outputLoginSession <- ResultRepository{Result: data}

	})

	return outputLoginSession
}

// Load function for loading refresh token by its key
func (repo *LoginSessionRepositoryRedis) Load(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "LoginSessionRepositoryRedis-Load"

	outputLoadLoginSess := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputLoadLoginSess)

		val, err := repo.client.Get(key)
		tags[helper.TagKey] = key

		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "load_token_from_redis", err, key)
			}

			outputLoadLoginSess <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		outputLoadLoginSess <- ResultRepository{Result: model.LoginSessionRedis{Token: val}}
		tags[helper.TextResponse] = model.LoginSessionRedis{Token: val, Key: key}

	})

	return outputLoadLoginSess
}

// Delete function for delete login attempt data into redis DB
func (repo *LoginSessionRepositoryRedis) Delete(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "LoginSessionRepositoryRedis-Delete"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TagKey] = key
		if _, err := repo.client.Del(key); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "delete_token_from_redis", err, key)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}

// GetLoginActive function for delete token data into redis DB
func (repo *LoginSessionRepositoryRedis) GetLoginActive(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "LoginSessionRepositoryRedis-GetLoginActive"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		keys := fmt.Sprintf("%s*", key)
		tags[helper.TagKey] = keys
		val, err := repo.client.Keys(keys)
		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "get_redis_key", err, key)
			}
			output <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}
		resultData := []model.LoginSessionRedis{}
		for _, v := range val {
			data := model.LoginSessionRedis{}
			value, _ := repo.client.Get(v)
			data.Key = v
			data.Token = value
			resultData = append(resultData, data)
		}

		output <- ResultRepository{Result: resultData}
	})

	return output
}
