package redis

import "time"

// Client interface contract
type Client interface {
	Get(key string) (string, error)
	Del(key string) (int64, error)
	Set(key string, value string, duration time.Duration) (string, error)
	Ping() (string, error)
	Expire(key string, exp time.Duration) (bool, error)
	Keys(key string) ([]string, error)
}
