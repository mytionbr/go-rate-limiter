type Limiter struct {
	store Store
	rateLimitIP int64
	rateLimitToken int64
	blockDuration time.Duration
}

func Newlimiter(store Store, rateLimitIP, rateLimitToken int64, blockDuration time.Duration) *Limiter {
	return &Limiter{
		store: store,
		rateLimitIP: rateLimitIP,
		rateLimitToken: rateLimitToken,
		blockDuration: blockDuration,
	}
}

func (l *Limiter) Middleware(next http.Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "ip:" + r.RemoteAddr
		limit := l.rateLimitIP

		if token := r.Header.Get("API_KEY"); token != "" {
			key = "token:" + token
			limit = l.rateLimitToken
		}
		blockKey := key + ":block"

		if blocked, err := l.store.Exists(blockKey); blocked  {
			http.Error(w,
			"you have reached the maximum number of requests or actions allowed within a certain time frame",
			http.StatusTooManyRequests)
			)
			return
		} 
	})
} 