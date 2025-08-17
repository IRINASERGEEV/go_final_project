# Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM ubuntu:latest

WORKDIR /app

# Копируем бинарник и web
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

EXPOSE 7540

# Переменные окружения (значения по умолчанию)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD=""

CMD ["./server"]