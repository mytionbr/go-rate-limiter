package limiter

import "time"

type Store interface {
	Incr(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
	Exists(key string) (bool, error)
	Set(key string, value interface{}, expiration time.Duration) error
}
