package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr, pass string, db int) *RedisStore {
	c := redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db})
	return &RedisStore{
		client: c,
		ctx:    context.Background(),
	}
}

func (s *RedisStore) Incr(key string) (int64, error) {
	return s.client.Incr(s.ctx, key).Result()
}

func (s *RedisStore) Expire(key string, expiration time.Duration) error {
	return s.client.Expire(s.ctx, key, expiration).Err()
}

func (s *RedisStore) Exists(key string) (bool, error) {
	exists, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (s *RedisStore) Set(key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(s.ctx, key, value, expiration).Err()
}
