FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o go-rate-limiter main.go
EXPOSE 8080
CMD ["./go-rate-limiter"]