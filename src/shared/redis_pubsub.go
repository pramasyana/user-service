package shared

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// PubSub interface
type PubSub interface {
	Publish(string, string) error
	Subscribe(string, chan []byte) error
}

// RedisPubSub struct
type RedisPubSub struct {
	Pool *redis.Pool
	Conn redis.Conn
}

// RedisPubSubConfig input for constructor
type RedisPubSubConfig struct {
	Host     string
	Password string
	Port     string
	UseTLS   string
	UseDB    string
}

// NewRedisPubSub return new service
func NewRedisPubSub(conf *RedisPubSubConfig) *RedisPubSub {
	if conf == nil {
		log.Fatal("config is required")
	}

	var tlsSecured bool

	tlsSecured, err := strconv.ParseBool(conf.UseTLS)
	if err != nil {
		log.Fatal(err)
	}

	redispool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port), redis.DialPassword(conf.Password), redis.DialUseTLS(tlsSecured))
		},
	}

	// Get a connection
	conn := redispool.Get()
	defer conn.Close()
	// Test the connection
	_, err = conn.Do("PING")
	if err != nil {
		log.Fatalf("can't connect :\n%v", err)
	}

	return &RedisPubSub{
		Pool: redispool,
		Conn: conn,
	}
}

// Publish publish key value
func (s *RedisPubSub) Publish(key string, value string) error {
	conn := s.Pool.Get()

	_, err := conn.Do("PUBLISH", key, value)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe subscribe
func (s *RedisPubSub) Subscribe(key string, msg chan []byte) error {
	rc := s.Pool.Get()
	rc.Do("CONFIG", "SET", "notify-keyspace-events", "Ex")
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.PSubscribe(key); err != nil {
		return err
	}

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				msg <- v.Data
			}
		}
	}()
	return nil
}
