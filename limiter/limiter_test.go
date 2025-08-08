package limiter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type InMemoryStore struct {
	mu  sync.Mutex
	cnt map[string]int64
	blk map[string]time.Time
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{cnt: make(map[string]int64), blk: make(map[string]time.Time)}
}

func (s *InMemoryStore) Incr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cnt[key]++
	return s.cnt[key], nil
}

func (s *InMemoryStore) Expire(key string, expiration time.Duration) error {
	go func() {
		time.Sleep(expiration)
		s.mu.Lock()
		delete(s.cnt, key)
		s.mu.Unlock()
	}()
	return nil
}

func (s *InMemoryStore) Exists(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.blk[key]
	return ok, nil
}

func (s *InMemoryStore) Set(key string, _ interface{}, expiration time.Duration) error {
	s.mu.Lock()
	s.blk[key] = time.Now().Add(expiration)
	s.mu.Unlock()
	go func() {
		time.Sleep(expiration)
		s.mu.Lock()
		delete(s.blk, key)
		s.mu.Unlock()
	}()
	return nil
}

func TestRateLimiter(t *testing.T) {
	store := NewInMemoryStore()
	rl := NewLimiter(store, 2, 5, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:5678"

	for i := 1; i <= 3; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if i <= 2 && rr.Code != http.StatusOK {
			t.Errorf("esperava 200, teve %d na req %d", rr.Code, i)
		}
		if i == 3 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("esperava 429, teve %d na req %d", rr.Code, i)
		}
	}
}

func TestTokenLimiter(t *testing.T) {
	store := NewInMemoryStore()
	rl := NewLimiter(store, 2, 3, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1111"
	req.Header.Set("API_KEY", "tokentest")

	for i := 1; i <= 4; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if i <= 3 && rr.Code != http.StatusOK {
			t.Errorf("esperava 200 no %dº, teve %d", i, rr.Code)
		}
		if i == 4 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("esperava 429 no 4º, teve %d", rr.Code)
		}
	}
}

func TestTokenOverridesIP(t *testing.T) {
	store := NewInMemoryStore()
	rl := NewLimiter(store, 1, 5, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:2222"
	req.Header.Set("API_KEY", "override")

	for i := 1; i <= 5; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("esperava 200 no %dº, teve %d", i, rr.Code)
		}
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("esperava 429 no 6º, teve %d", rr.Code)
	}
}

func TestBlockExpiry(t *testing.T) {
	store := NewInMemoryStore()
	rl := NewLimiter(store, 1, 5, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.6:4444"

	handler.ServeHTTP(httptest.NewRecorder(), req)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("esperava 429 no bloqueio, teve %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("esperava 429 durante bloqueio, teve %d", rr.Code)
	}

	time.Sleep(1100 * time.Millisecond)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("esperava 200 após expirar bloqueio, teve %d", rr.Code)
	}
}

func TestCounterExpiration(t *testing.T) {
	store := NewInMemoryStore()
	rl := NewLimiter(store, 2, 5, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.5:3333"

	handler.ServeHTTP(httptest.NewRecorder(), req)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("esperava 429 imediato, teve %d", rr.Code)
	}

	time.Sleep(1100 * time.Millisecond)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("esperava 200 após expirar, teve %d", rr.Code)
	}
}
