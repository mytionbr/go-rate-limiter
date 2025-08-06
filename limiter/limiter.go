package limiter

import (
	"net/http"
	"time"
)

type Limiter struct {
	store          Store
	rateLimitIP    int64
	rateLimitToken int64
	blockDuration  time.Duration
}

func NewLimiter(store Store, rateLimitIP, rateLimitToken int64, blockDuration int64) *Limiter {
	return &Limiter{
		store:          store,
		rateLimitIP:    rateLimitIP,
		rateLimitToken: rateLimitToken,
		blockDuration:  time.Duration(blockDuration) * time.Second,
	}
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "ip:" + r.RemoteAddr
		limit := l.rateLimitIP

		if token := r.Header.Get("API_KEY"); token != "" {
			key = "token:" + token
			limit = l.rateLimitToken
		}
		blockKey := key + ":block"

		if blocked, _ := l.store.Exists(blockKey); blocked {
			http.Error(w,
				"you have reached the maximum number of requests or actions allowed within a certain time frame",
				http.StatusTooManyRequests)
			return
		}

		count, _ := l.store.Incr(key)
		if count == 1 {
			l.store.Expire(key, time.Second)
		}
		if count > limit {
			l.store.Set(blockKey, 1, l.blockDuration)
			http.Error(w,
				"you have reached the maximum number of requests or actions allowed within a certain time frame",
				http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)

	})
}
