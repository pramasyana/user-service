package mocks

import (
	"fmt"
	"time"
)

func InitFakeRedis() *FakeRedis {
	return &FakeRedis{}
}

type FakeRedis struct {
	GetFunc  func(string) (string, error)
	SetFunc  func(string, string, time.Duration) (string, error)
	DelFunc  func(string) (int64, error)
	PingFunc func() (string, error)
	ExpFunc  func(string, time.Duration) (bool, error)
}

func (r *FakeRedis) Get(key string) (string, error) {
	if r.GetFunc != nil {
		return r.GetFunc(key)
	}
	return "", fmt.Errorf("Get %s Error", key)
}

func (r *FakeRedis) Del(key string) (int64, error) {
	if r.DelFunc != nil {
		return r.DelFunc(key)
	}
	return int64(0), fmt.Errorf("Del %s Error", key)
}

func (r *FakeRedis) Set(key, value string, exp time.Duration) (string, error) {
	if r.SetFunc != nil {
		return r.SetFunc(key, value, exp)
	}
	return "", fmt.Errorf("Set %s Error", key)
}
func (r *FakeRedis) Ping() (string, error) {
	if r.PingFunc != nil {
		return r.PingFunc()
	}
	return "", fmt.Errorf("Ping %s Error", "unable to ping")
}

func (r *FakeRedis) Expire(key string, exp time.Duration) (bool, error) {
	if r.ExpFunc != nil {
		return r.ExpFunc(key, exp)
	}
	return false, fmt.Errorf("Exp %s Error", key)
}

func (r *FakeRedis) Keys(key string) ([]string, error) {
	if r.ExpFunc != nil {
		return r.Keys(key)
	}
	return []string{""}, fmt.Errorf("Exp %s Error", key)
}
