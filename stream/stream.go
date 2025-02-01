package stream

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/go-i2p/go-forward/config"
	"github.com/go-i2p/go-forward/metrics"
)

// Forward handles bidirectional forwarding between two stream connections
func Forward(ctx context.Context, conn1, conn2 net.Conn, cfg *config.ForwardConfig) error {
	var wg sync.WaitGroup
	errc := make(chan error, 2)

	// Start bidirectional copy
	wg.Add(2)
	go func() {
		defer wg.Done()
		errc <- copyStream(ctx, conn1, conn2, cfg, "1->2")
	}()

	go func() {
		defer wg.Done()
		errc <- copyStream(ctx, conn2, conn1, cfg, "2->1")
	}()

	// Wait for both goroutines and collect errors
	go func() {
		wg.Wait()
		close(errc)
	}()

	// Return first non-nil error or nil if none
	for err := range errc {
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func copyStream(ctx context.Context, dst, src net.Conn, cfg *config.ForwardConfig, label string) error {
	buffer := make([]byte, cfg.BufferSize)
	var m *metrics.StreamMetrics

	if cfg.EnableMetrics {
		m = metrics.NewStreamMetrics(label)
	}

	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cfg.ShutdownSignal:
			return nil
		default:
		}

		// Set read timeout
		if err := src.SetReadDeadline(time.Now().Add(cfg.IdleTimeout)); err != nil {
			return err
		}

		n, err := src.Read(buffer)
		if err != nil {
			return err
		}

		if cfg.EnableMetrics {
			m.AddBytesRead(int64(n))
		}

		// Set write timeout
		if err := dst.SetWriteDeadline(time.Now().Add(cfg.IdleTimeout)); err != nil {
			return err
		}

		_, err = dst.Write(buffer[:n])
		if err != nil {
			return err
		}

		if cfg.EnableMetrics {
			m.AddBytesWritten(int64(n))
		}
	}
}
