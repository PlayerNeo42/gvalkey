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

- Go 1.24.3 or higher

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

## 🔧 Usage

### Starting the Server

```bash
# Start GValkey server on default port (6379)
./gvalkey
```

The server will start listening on `localhost:6379` by default.

## ⚙️ Configuration

GValkey can be configured using environment variables. All configuration options have sensible defaults.

### Environment Variables

| Variable | Description | Default | Valid Values |
|----------|-------------|---------|--------------|
| `GVK_HOST` | Server bind address | `0.0.0.0` | Valid hostname or IP address |
| `GVK_PORT` | Server listen port | `6379` | 1-65535 |
| `GVK_LOG_LEVEL` | Logging level | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |

### Configuration Examples

```bash
# Start server on custom port
GVK_PORT=8080 ./gvalkey

# Start server with debug logging
GVK_LOG_LEVEL=DEBUG ./gvalkey

# Start server on specific interface
export GVK_HOST=127.0.0.1
export GVK_PORT=6380
export GVK_LOG_LEVEL=WARN
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

### Running Tests

```bash
go test ./...
```

### Code Formatting && Linting

```bash
golangci-lint fmt 
golangci-lint run --fix
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Redis](https://redis.io/)
