# Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o gvalkey ./cmd/server

# Runner
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gvalkey .
RUN chown -R nobody:nobody /app
USER nobody
EXPOSE 6379

CMD ["./gvalkey"]
