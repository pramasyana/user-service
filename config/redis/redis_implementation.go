package redis

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// Conn base struct
type Conn struct {
	Client *redis.Client
}

// ConnectRedis init redis
func ConnectRedis(redisHost, redisTLS, redisPassword, redisPort, redisDB string) (Client, error) {
	cl, err := GetRedis(redisHost, redisTLS, redisPassword, redisPort, redisDB)
	if err != nil {
		return nil, err
	}
	return Conn{
		Client: cl,
	}, nil
}

// Del key
func (r Conn) Del(key string) (int64, error) {
	return r.Client.Del(key).Result()
}

// Get key
func (r Conn) Get(key string) (string, error) {
	return r.Client.Get(key).Result()
}

// Set redis
func (r Conn) Set(key string, value string, exp time.Duration) (string, error) {
	return r.Client.Set(key, value, exp).Result()
}

// Ping result
func (r Conn) Ping() (string, error) {
	return r.Client.Ping().Result()
}

// Expire result
func (r Conn) Expire(key string, exp time.Duration) (bool, error) {
	return r.Client.Expire(key, exp).Result()
}

// Keys get multi key
func (r Conn) Keys(key string) ([]string, error) {
	return r.Client.Keys(key).Result()
}

// GetRedis function
func GetRedis(redisHost, redisTLS, redisPassword, redisPort, redisDB string) (*redis.Client, error) {
	//Transport Layer Security config,
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// https://godoc.org/crypto/tls#Config

	tlsSecured, err := strconv.ParseBool(redisTLS)
	if err != nil {
		return nil, err
	}

	var conf *tls.Config

	// force checking for unsecured aws redis
	if tlsSecured {
		conf = &tls.Config{
			InsecureSkipVerify: tlsSecured,
		}
	} else {
		conf = nil
	}

	useDB, _ := strconv.Atoi(redisDB)
	cl := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password:  redisPassword,
		DB:        useDB, // use default DB
		TLSConfig: conf,
	})

	return cl, nil
}
