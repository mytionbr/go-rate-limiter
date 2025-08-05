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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!!!")
	})

	log.Println("Servidor rodando na porta 8080")
}
