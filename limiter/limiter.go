package limiter

import (
	"net"
	"net/http"
	"time"
)

type Limiter struct {
	store          Store
	rateLimitIP    int64
	rateLimitToken int64
	blockDuration  time.Duration
}

func NewLimiter(store Store, rateLimitIP, rateLimitToken, blockDuration int) *Limiter {
	return &Limiter{
		store:          store,
		rateLimitIP:    int64(rateLimitIP),
		rateLimitToken: int64(rateLimitToken),
		blockDuration:  time.Duration(blockDuration) * time.Second,
	}
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		limit := l.rateLimitIP

		key := "ip:" + host

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
