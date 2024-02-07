package config

import (
	"runtime"
	"time"

	"github.com/gofiber/storage/redis/v3"
)

type Cache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Redis struct {
	*redis.Storage
}

var cfg = NewConfig()

func getRedisConfig(cfg *Config) *redis.Config {
	return &redis.Config{
		Host:      cfg.REDIS_HOST,
		Port:      cfg.REDIS_PORT,
		Username:  cfg.REDIS_USERNAME,
		Password:  cfg.REDIS_PASSWORD,
		Database:  cfg.REDIS_DATABASE,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	}
}

func NewRedis() *Redis {
	return &Redis{
		Storage: redis.New(*getRedisConfig(cfg)),
	}
}

func (r *Redis) Set(key string, value string) error {
	expiration := time.Duration(1) * time.Minute

	return r.Storage.Set(key, []byte("true"), expiration)
}

func (r *Redis) Get(key string) (string, error) {
	value, err := r.Storage.Get(key)

	if err != nil {
		return "", err
	}

	if string(value) == "" {
		return "", nil
	}

	return string(value), nil
}

func (r *Redis) Delete(key string) error {
	return r.Storage.Delete(key)
}
