# builder stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o main .

# runner stage
FROM alpine:3.23.3

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app
COPY --from=builder /app/main .

ENV PORT=8888

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

EXPOSE 8888

CMD ["./main"]