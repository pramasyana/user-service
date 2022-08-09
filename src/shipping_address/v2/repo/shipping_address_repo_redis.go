package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
)

// ShippingAddressRepoRedis data structure
type ShippingAddressRepoRedis struct {
	client redis.Client
}

// NewShippingAddressRepoRedis function for initializing shipping address repository redis
func NewShippingAddressRepoRedis(cl redis.Client) *ShippingAddressRepoRedis {
	return &ShippingAddressRepoRedis{cl}
}

// SaveRedisMeta function for saving shipping address data into redis DB
func (s *ShippingAddressRepoRedis) SaveRedisMeta(memberID string, page string, limit string, shippingList model.ListShippingAddress) <-chan error {
	ctx := "ShippingAddressRepo-SaveRedisMeta"

	output := make(chan error)

	go func() {
		defer close(output)

		cl := s.client
		ctxReq := context.Background()
		expire := 24 * time.Hour

		key := fmt.Sprintf("%s_%s_%s_%s_%s", keyShippingDetail, memberID, page, limit, os.Getenv("ENV"))
		shippingDetailsRedis, err := json.Marshal(shippingList)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "json_marshal", err, shippingList)
			output <- err
			return
		}

		_, err = cl.Set(key, string(shippingDetailsRedis), expire)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "save_redis", err, shippingList)
			output <- err
			return
		}

		output <- err

	}()

	return output
}

// LoadRedisMeta function for load redis DB
func (s *ShippingAddressRepoRedis) LoadRedisMeta(memberID string, page string, limit string) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-LoadRedisMeta"

	output := make(chan ResultRepository)

	go func() {
		defer close(output)

		cl := s.client
		ctxReq := context.Background()

		key := fmt.Sprintf("%s_%s_%s_%s_%s", keyShippingDetail, memberID, page, limit, os.Getenv("ENV"))

		response := model.ListShippingAddress{}

		val, err := cl.Get(key)
		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "load_from_redis", err, val)
			}

			output <- ResultRepository{Error: err}
			return
		}

		var jsonBlob = []byte(val)
		e := json.Unmarshal(jsonBlob, &response)
		if e != nil {
			helper.SendErrorLog(ctxReq, ctx, "json_unmarshal", e, val)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: response}
	}()

	return output
}

// DeleteMultipleRedis function for delete redis DB
func (s *ShippingAddressRepoRedis) DeleteMultipleRedis(memberID string) <-chan error {
	ctx := "ShippingAddressRepo-DeleteMultipleRedis"

	output := make(chan error)

	go func() {
		defer close(output)

		cl := s.client
		ctxReq := context.Background()
		key := fmt.Sprintf("%s_%s*", keyShippingDetail, memberID)

		val, err := cl.Keys(key)
		if err != nil {
			if err.Error() != helper.ErrorRedis {
				helper.SendErrorLog(ctxReq, ctx, "load_from_redis", err, key)
			}

			output <- err
			return
		}

		for _, v := range val {
			_, err := cl.Del(v)
			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, "delete_redis", err, v)
				output <- err
				return
			}
		}

		output <- err

	}()

	return output
}
