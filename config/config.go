package config

import "time"

// ForwardConfig holds common configuration for both stream and packet forwarding
type ForwardConfig struct {
	BufferSize     int           // Size of the buffer for reading/writing
	IdleTimeout    time.Duration // Maximum idle time before closing connection
	MaxPacketSize  int           // Maximum packet size (for PacketConn)
	EnableMetrics  bool          // Enable forwarding metrics
	ShutdownSignal chan struct{} // Signal for graceful shutdown
}

// DefaultConfig returns a ForwardConfig with sensible defaults
func DefaultConfig() *ForwardConfig {
	return &ForwardConfig{
		BufferSize:     32 * 1024, // 32KB buffer
		IdleTimeout:    30 * time.Second,
		MaxPacketSize:  65507, // Max UDP packet size
		EnableMetrics:  true,
		ShutdownSignal: make(chan struct{}),
	}
}
