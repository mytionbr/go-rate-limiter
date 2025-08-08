# go-rate-limiter

Um rate limiter simples em Go, com suporte a Redis, para limitar requisições por IP ou por token (API_KEY).

## Requisitos

- Go 1.22+
- Docker (opcional, para rodar Redis e a aplicação em containers)

## Instalação

Clone o repositório:

```sh
git clone https://github.com/mytionbr/go-rate-limiter.git
cd go-rate-limiter
```

Copie o arquivo `.env.example` para `.env` e ajuste as variáveis conforme necessário:

```sh
cp .env.example .env
```

## Rodando com Docker Compose

Suba o Redis e a aplicação:

```sh
docker-compose up --build
```

A aplicação estará disponível em `http://localhost:8080`.

## Rodando localmente (sem Docker)

Certifique-se de ter um Redis rodando (pode ser local ou em container):

```sh
docker run -p 6379:6379 redis:7-alpine
```

Instale as dependências e rode a aplicação:

```sh
go mod download
go run main.go
```

## Configuração

As configurações são feitas via variáveis de ambiente (veja `.env.example`):

- `RATE_LIMIT_IP`: Limite de requisições por IP por segundo.
- `RATE_LIMIT_TOKEN`: Limite de requisições por token (API_KEY) por segundo.
- `BLOCK_DURATION_SECONDS`: Tempo de bloqueio (em segundos) após exceder o limite.
- `REDIS_ADDR`: Endereço do Redis.
- `REDIS_PASSWORD`: Senha do Redis (opcional).
- `REDIS_DB`: Banco do Redis.
- `PORT`: Porta da aplicação.

## Testes

Para rodar os testes unitários:

```sh
go test ./limiter
```


## Teste de Carga (Load Test)

Você pode usar o [oha](https://github.com/hatoo/oha) para testar a performance do rate limiter:

### Com API_KEY (token)

```sh
oha -n 1000 -c 50 -q 150 -H "API_KEY: abc123" --no-tui --stats-success-breakdown http://localhost:8080/
```

### Por IP (sem API_KEY)

```sh
oha -n 1000 -c 1 -q 20 --no-tui --stats-success-breakdown http://localhost:8080/
```

Arquivos de teste de carga também estão disponíveis (`load-test.yml`, `load-test-token.yml`) e podem ser usados com ferramentas como [Artillery](https://artillery.io/).

## Uso

Faça requisições para `http://localhost:8080/`.  
Para o limite por token, envie o header `API_KEY`:

```sh
curl -H "API_KEY: abc123" http://localhost:8080/
```

---