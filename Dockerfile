FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:3.24.3

RUN apk --no-cache add ca-certificates tzdata && adduser -D -u 1001 -s /sbin/nologin appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/main .

USER appuser

EXPOSE 8083

CMD ["./main"]