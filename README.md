<h1 align="center" style="border-bottom: none">
    GValkey
</h1>

<p align="center">
    A Redis-compatible in-memory cache implemented in Go
</p>

<p align="center">
    <a href="https://golang.org/"><img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/PlayerNeo42/gvalkey"></a>
    <a href="https://goreportcard.com/report/github.com/PlayerNeo42/gvalkey"><img alt="Go report" src="https://goreportcard.com/badge/github.com/PlayerNeo42/gvalkey"></a>
    <a href="LICENSE"><img alt="GitHub License" src="https://img.shields.io/github/license/PlayerNeo42/gvalkey"></a>
</p>

## 📋 Table of Contents

- [About](#-about)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Configuration](#-configuration)
- [Supported Commands](#-supported-commands)
- [Development](#-development)
- [License](#-license)

## 🚀 About

GValkey is a lightweight, Redis-compatible in-memory cache written in Go. It implements the Redis Serialization Protocol (RESP) and provides a subset of Redis functionality with a focus on simplicity.

## ✨ Features

- **Redis Protocol Compatible**: Implements RESP (Redis Serialization Protocol)
- **In-Memory**: Fast key-value storage with automatic TTL support
- **Concurrent Safe**: Thread-safe operations using Go's sync.Map
- **Automatic Cleanup**: Background garbage collection for expired keys
- **Graceful Shutdown**: Proper server shutdown handling
- **Structured Logging**: Comprehensive logging with slog
- **Command Validation**: Proper argument validation for Redis commands

## 📦 Installation

### Prerequisites

- Go 1.24.3 or higher (for building from source)
- Docker (for running with Docker)

### Using Docker (Recommended)

```bash
# Pull the latest image
docker pull ghcr.io/playerneo42/gvalkey:latest

# Run GValkey container
docker run -d \
  --name gvalkey \
  -p 6379:6379 \
  ghcr.io/playerneo42/gvalkey:latest
```

### Using Docker Compose

```yaml
services:
  gvalkey:
    image: ghcr.io/playerneo42/gvalkey:latest
    ports:
      - "6379:6379"
    restart: unless-stopped
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/PlayerNeo42/gvalkey
cd gvalkey

# Build the server
go build -o gvalkey ./cmd/server

# Run the server
./gvalkey
```

### Connecting with Redis CLI

Since GValkey implements the Redis protocol, you can use any Redis-compatible client:

```bash
# Using redis-cli
redis-cli -h localhost -p 6379

# Test basic commands
127.0.0.1:6379> SET mykey "Hello, GValkey!"
OK
127.0.0.1:6379> GET mykey
"Hello, GValkey!"
127.0.0.1:6379> DEL mykey
(integer) 1
```

## ⚙️ Configuration

GValkey can be configured using environment variables. All configuration options have sensible defaults.

### Environment Variables

| Variable | Description | Default | Valid Values |
|----------|-------------|---------|--------------|
| `GVK_HOST` | Server bind address | `0.0.0.0` | Valid hostname or IP address |
| `GVK_PORT` | Server listen port | `6379` | 1-65535 |
| `GVK_LOG_LEVEL` | Logging level | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |


## 📝 Supported Commands

| Command | Description | Status |
|---------|-------------|--------|
| `SET key value [EX seconds] [PX milliseconds] [NX\|XX] [GET]` | Set a key-value pair with optional expiration and conditions | ✅ |
| `GET key` | Retrieve value by key | ✅ |
| `DEL key [key ...]` | Delete one or more keys | ✅ |

### SET Command Options

- `EX seconds`: Set expiration in seconds
- `PX milliseconds`: Set expiration in milliseconds
- `NX`: Only set if key doesn't exist
- `XX`: Only set if key already exists
- `GET`: Return the old value when setting

## 🛠️ Development

### Code Formatting && Linting

```bash
golangci-lint fmt 
golangci-lint run --fix
```

### Running Tests

```bash
go test ./...
```

> Or you can do `make check` for both steps.

### Benchmark

```bash
make bench
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Redis](https://redis.io/)
