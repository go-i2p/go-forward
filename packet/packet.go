package packet

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/go-i2p/go-forward/config"
	"github.com/go-i2p/go-forward/metrics"
)

// Forward handles bidirectional forwarding between two packet connections
func Forward(ctx context.Context, conn1, conn2 net.PacketConn, cfg *config.ForwardConfig) error {
	var wg sync.WaitGroup
	errc := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		errc <- forwardPackets(ctx, conn1, conn2, cfg, "1->2")
	}()

	go func() {
		defer wg.Done()
		errc <- forwardPackets(ctx, conn2, conn1, cfg, "2->1")
	}()

	// Wait for both goroutines and collect errors
	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func forwardPackets(ctx context.Context, dst, src net.PacketConn, cfg *config.ForwardConfig, label string) error {
	buffer := make([]byte, cfg.MaxPacketSize)
	var m *metrics.PacketMetrics

	if cfg.EnableMetrics {
		m = metrics.NewPacketMetrics(label)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cfg.ShutdownSignal:
			return nil
		default:
		}

		if err := src.SetReadDeadline(time.Now().Add(cfg.IdleTimeout)); err != nil {
			return err
		}

		n, addr, err := src.ReadFrom(buffer)
		if err != nil {
			return err
		}

		if cfg.EnableMetrics {
			m.AddPacketReceived()
			m.AddBytesRead(int64(n))
		}

		if err := dst.SetWriteDeadline(time.Now().Add(cfg.IdleTimeout)); err != nil {
			return err
		}

		_, err = dst.WriteTo(buffer[:n], addr)
		if err != nil {
			return err
		}

		if cfg.EnableMetrics {
			m.AddPacketSent()
			m.AddBytesWritten(int64(n))
		}
	}
}
