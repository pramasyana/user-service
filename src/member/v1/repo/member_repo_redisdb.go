package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// MemberRepoRedis data structure
type MemberRepoRedis struct {
	client redis.Client
}

// NewMemberRepoRedis function for initializing member repository redis
func NewMemberRepoRedis(cl redis.Client) *MemberRepoRedis {
	return &MemberRepoRedis{cl}
}

// Save function for saving token data into redis DB
func (repo *MemberRepoRedis) Save(member *model.MemberRedis) <-chan ResultRepository {
	ctx := "MemberRepoRedis-Save"

	output := make(chan ResultRepository)
	ctxReq := context.Background()
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		cl := repo.client
		key := fmt.Sprintf("%s:%s", model.ForgotPasswordKeyRedis, member.ID)
		tags[helper.TagKey] = key

		redisData, err := json.Marshal(member)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "marshal_json", err, key)
			output <- ResultRepository{Error: err}
			return
		}

		if _, err = cl.Set(key, string(redisData), member.TTL); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_to_redis", err, string(redisData))
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}

// Load function for loading token data by its ID
func (repo *MemberRepoRedis) Load(uid string) <-chan ResultRepository {
	ctx := "MemberRepoRedis-Load"

	output := make(chan ResultRepository)
	ctxReq := context.Background()
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		cl := repo.client

		key := fmt.Sprintf("%s:%s", model.ForgotPasswordKeyRedis, uid)
		tags[helper.TagKey] = key
		redisData := model.MemberRedis{}

		val, err := cl.Get(key)
		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "get_redis", err, key)
			}

			output <- ResultRepository{Error: err}
			return
		}

		var jsonBlob = []byte(val)
		e := json.Unmarshal(jsonBlob, &redisData)
		if e != nil {
			helper.SendErrorLog(ctxReq, ctx, "json_unmarshal", err, val)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: redisData}
	})

	return output
}

// Delete function for delete token data into redis DB
func (repo *MemberRepoRedis) Delete(uid string) <-chan ResultRepository {
	ctx := "MemberRepoRedis-Delete"

	output := make(chan ResultRepository)
	ctxReq := context.Background()
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		cl := repo.client
		key := fmt.Sprintf("%s:%s", model.ForgotPasswordKeyRedis, uid)
		tags[helper.TagKey] = key

		if _, err := cl.Del(key); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "del_redis", err, key)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}

// LoadByKey function for loading refresh token by its key
func (repo *MemberRepoRedis) LoadByKey(ctxReq context.Context, key string) <-chan ResultRepository {
	ctx := "MemberRepoRedis-LoadByKey"

	output := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TagKey] = key

		cl := repo.client
		val, err := cl.Get(key)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_redis", err, key)
			output <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultRepository{Result: model.ResendActivationAttempt{Attempt: val}}
		tags[helper.TextResponse] = model.ResendActivationAttempt{Attempt: val}
	})

	return output
}

// SaveResendActivationAttempt function for saving refresh token into redis DB
func (repo *MemberRepoRedis) SaveResendActivationAttempt(ctxReq context.Context, data *model.ResendActivationAttempt) <-chan ResultRepository {
	ctx := "AuthRepo-Save"

	output := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["args"] = data
		cl := repo.client

		if _, err := cl.Set(data.Key, data.Attempt, data.ResendActivationAttemptAge); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_to_redis", err, data)
			output <- ResultRepository{Error: err}
			return
		}

		tags[helper.TextResponse] = data

		output <- ResultRepository{Result: data}

	})

	return output
}

// RevokeAllAccess function for delete token data into redis DB
func (repo *MemberRepoRedis) RevokeAllAccess(ctxReq context.Context, key, activeKey string) <-chan ResultRepository {
	ctx := "MemberRepoRedis-RevokeAllAccess"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		cl := repo.client

		keys := fmt.Sprintf("%s*", key)
		tags["key"] = keys
		val, err := cl.Keys(keys)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "redis_delete_keys", err, keys)
			output <- ResultRepository{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		for _, v := range val {
			if v == activeKey {
				continue
			}
			_, err := cl.Del(v)
			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, "delete_redis", err, v)
				output <- ResultRepository{Error: err}
				tags[helper.TextResponse] = err
			}
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}
