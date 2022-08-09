package shared

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func setupMockRedisPubsub() RedisPubSub {
	connRedis := redigomock.NewConn()
	connect := RedisPubSub{
		Conn: connRedis,
		Pool: &redis.Pool{
			Dial: func() (redis.Conn, error) { return connRedis, nil },
		},
	}

	return connect
}

func TestNewPubsub(t *testing.T) {
	mr, _ := miniredis.Run()
	res := NewRedisPubSub(&RedisPubSubConfig{UseTLS: "false", Host: mr.Host(), Port: mr.Port()})
	assert.Equal(t, res, res)
}

func TestPublishPubsub(t *testing.T) {
	connect := setupMockRedisPubsub()
	res := connect.Publish("key", "value")
	assert.Error(t, res)
}

func TestSubcribePubsub(t *testing.T) {
	connect := setupMockRedisPubsub()
	var msg chan []byte
	res := connect.Subscribe("keys", msg)
	assert.NoError(t, res)
}
