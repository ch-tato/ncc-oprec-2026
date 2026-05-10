# builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# runner stage
FROM alpine:3.23.3

RUN adduser -D -g '' appuser
WORKDIR /app
RUN mkdir -p data && chown -R appuser:appuser /app
USER appuser

COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/templates ./templates
COPY --from=builder --chown=appuser:appuser /app/static ./static

ENV PORT=8880

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/ || exit 1

EXPOSE 8880

CMD ["./main"]