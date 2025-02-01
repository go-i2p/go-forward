# go-forward

A Go library providing a clean, consistent API for network traffic forwarding, supporting both stream (TCP-like) and packet (UDP-like) connections with minimal boilerplate.

## Core Features

- 🎯 Simple, consistent API for both connection types
- 🔄 Bidirectional forwarding with a single function call
- 📊 Built-in metrics with zero configuration
- 🛠️ Sensible defaults with optional configuration

## Quick Start

```go
import (
    "context"
    "github.com/go-i2p/go-forward/stream"
    "github.com/go-i2p/go-forward/packet"
)

// Stream forwarding (TCP-like)
stream.Forward(ctx, conn1, conn2, config.DefaultConfig())

// Packet forwarding (UDP-like)
packet.Forward(ctx, conn1, conn2, config.DefaultConfig())
```

## API Overview

The library provides two main functions with identical signatures:

```go
// Stream forwarding
func stream.Forward(
    ctx context.Context,
    conn1, conn2 net.Conn,
    cfg *config.ForwardConfig,
) error

// Packet forwarding
func packet.Forward(
    ctx context.Context,
    conn1, conn2 net.PacketConn,
    cfg *config.ForwardConfig,
) error
```

### Common Patterns

Both connection types share:
- Context-based cancellation
- Identical configuration structure
- Consistent error handling
- Built-in metrics collection
- Automatic resource cleanup

## Usage Examples

### Basic Stream Forwarding

```go
package main

import (
    "context"
    "net"
    "github.com/go-i2p/go-forward/stream"
    "github.com/go-i2p/go-forward/config"
)

func main() {
    // Create connections
    listener, _ := net.Listen("tcp", ":8080")
    conn1, _ := listener.Accept()
    conn2, _ := net.Dial("tcp", "target:8081")

    // Forward with default configuration
    ctx := context.Background()
    err := stream.Forward(ctx, conn1, conn2, config.DefaultConfig())
}
```

### Basic Packet Forwarding

```go
package main

import (
    "context"
    "net"
    "github.com/go-i2p/go-forward/packet"
    "github.com/go-i2p/go-forward/config"
)

func main() {
    // Create connections
    conn1, _ := net.ListenPacket("udp", ":8080")
    conn2, _ := net.ListenPacket("udp", ":8081")

    // Forward with default configuration
    ctx := context.Background()
    err := packet.Forward(ctx, conn1, conn2, config.DefaultConfig())
}
```

### With Custom Configuration

```go
cfg := &config.ForwardConfig{
    BufferSize:  32 * 1024,     // 32KB buffer
    IdleTimeout: time.Minute,   // 1 minute timeout
}
err := stream.Forward(ctx, conn1, conn2, cfg)
```

### With Metrics

```go
// Metrics are enabled by default
cfg := config.DefaultConfig()

// Forward with automatic metric collection
err := stream.Forward(ctx, conn1, conn2, cfg)

// Access metrics
stats := metrics.GetAllStreamStats()
for _, stat := range stats {
    fmt.Printf("Bytes transferred: %d\n", stat.BytesRead + stat.BytesWritten)
}
```

## Configuration

Simple configuration with sensible defaults:

```go
// Use defaults
cfg := config.DefaultConfig()

// Or customize
cfg := &config.ForwardConfig{
    BufferSize:     32 * 1024,        // Buffer size
    IdleTimeout:    time.Minute,      // Idle timeout
    MaxPacketSize:  65507,            // For UDP
    EnableMetrics:  true,             // Metrics collection
    ShutdownSignal: make(chan struct{}), // Graceful shutdown
}
```

## Installation

```bash
go get github.com/go-i2p/go-forward
```

Requires Go 1.23.5 or later.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Run `make fmt` to format code
4. Submit a Pull Request

## License

MIT License - see LICENSE file for details.

## Documentation

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/go-i2p/go-forward).