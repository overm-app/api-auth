FROM golang:1.27.3-alpine AS builder

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod verify

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:3.24.3

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

RUN adduser -D -u 1001 -s /bin/sh appuser

COPY --from=builder /app/main .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8081

CMD ["./main"]