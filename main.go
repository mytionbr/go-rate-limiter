package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mytionbr/go-rate-limiter/config"
	"github.com/mytionbr/go-rate-limiter/limiter"
)

func main() {
	cfg := config.Load()
	store := limiter.NewRedisStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	rl := limiter.NewLimiter(store, cfg.RateLimitIP, cfg.RateLimitToken, cfg.BlockDurationSec)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!!!")
	})

	handler := rl.Middleware(mux)
	log.Println("teste Servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
