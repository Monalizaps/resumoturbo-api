# Etapa 1: compilação
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o resumoturbo .

# Etapa 2: imagem final leve
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/resumoturbo .

EXPOSE 8080

CMD ["./resumoturbo"]
