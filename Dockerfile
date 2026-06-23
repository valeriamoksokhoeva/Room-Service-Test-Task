FROM golang:1.25.11-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

RUN go install github.com/rubenv/sql-migrate/...@latest


FROM alpine:3.24

RUN apk --no-cache add ca-certificates tzdata postgresql-client

RUN adduser -D appuser

WORKDIR /home/appuser

COPY --from=builder --chown=appuser:appuser /app/app .
COPY --from=builder --chown=appuser:appuser /go/bin/sql-migrate /usr/local/bin/sql-migrate
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations
COPY --chown=appuser:appuser entrypoint.sh .

RUN chmod +x entrypoint.sh

USER appuser

ENTRYPOINT ["./entrypoint.sh"]