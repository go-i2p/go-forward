package metrics

import (
	"sync/atomic"
	"time"
)

// BaseMetrics contains common metrics for both stream and packet connections
type BaseMetrics struct {
	bytesRead    atomic.Int64
	bytesWritten atomic.Int64
	startTime    time.Time
	label        string
}

// StreamMetrics holds metrics specific to stream connections
type StreamMetrics struct {
	BaseMetrics
	activeStreams atomic.Int32
}

// PacketMetrics holds metrics specific to packet connections
type PacketMetrics struct {
	BaseMetrics
	packetsReceived atomic.Int64
	packetsSent     atomic.Int64
	packetErrors    atomic.Int64
}

// NewStreamMetrics creates a new StreamMetrics instance
func NewStreamMetrics(label string) *StreamMetrics {
	m := &StreamMetrics{}
	m.label = label
	m.startTime = time.Now()
	return m
}

// NewPacketMetrics creates a new PacketMetrics instance
func NewPacketMetrics(label string) *PacketMetrics {
	m := &PacketMetrics{}
	m.label = label
	m.startTime = time.Now()
	return m
}

// BaseMetrics methods

func (m *BaseMetrics) AddBytesRead(n int64) {
	m.bytesRead.Add(n)
}

func (m *BaseMetrics) AddBytesWritten(n int64) {
	m.bytesWritten.Add(n)
}

func (m *BaseMetrics) GetBytesRead() int64 {
	return m.bytesRead.Load()
}

func (m *BaseMetrics) GetBytesWritten() int64 {
	return m.bytesWritten.Load()
}

func (m *BaseMetrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

func (m *BaseMetrics) GetLabel() string {
	return m.label
}

// StreamMetrics methods

func (m *StreamMetrics) IncrementActiveStreams() {
	m.activeStreams.Add(1)
}

func (m *StreamMetrics) DecrementActiveStreams() {
	m.activeStreams.Add(-1)
}

func (m *StreamMetrics) GetActiveStreams() int32 {
	return m.activeStreams.Load()
}

// GetStats returns a snapshot of current stream metrics
func (m *StreamMetrics) GetStats() StreamStats {
	uptime := m.GetUptime()
	bytesRead := m.GetBytesRead()
	bytesWritten := m.GetBytesWritten()

	return StreamStats{
		Label:          m.GetLabel(),
		Uptime:         uptime,
		BytesRead:      bytesRead,
		BytesWritten:   bytesWritten,
		ActiveStreams:  m.GetActiveStreams(),
		BytesPerSecond: float64(bytesRead+bytesWritten) / uptime.Seconds(),
	}
}

// PacketMetrics methods

func (m *PacketMetrics) AddPacketReceived() {
	m.packetsReceived.Add(1)
}

func (m *PacketMetrics) AddPacketSent() {
	m.packetsSent.Add(1)
}

func (m *PacketMetrics) AddPacketError() {
	m.packetErrors.Add(1)
}

func (m *PacketMetrics) GetPacketsReceived() int64 {
	return m.packetsReceived.Load()
}

func (m *PacketMetrics) GetPacketsSent() int64 {
	return m.packetsSent.Load()
}

func (m *PacketMetrics) GetPacketErrors() int64 {
	return m.packetErrors.Load()
}

// GetStats returns a snapshot of current packet metrics
func (m *PacketMetrics) GetStats() PacketStats {
	uptime := m.GetUptime()
	bytesRead := m.GetBytesRead()
	bytesWritten := m.GetBytesWritten()
	packetsReceived := m.GetPacketsReceived()
	packetsSent := m.GetPacketsSent()

	return PacketStats{
		Label:            m.GetLabel(),
		Uptime:           uptime,
		BytesRead:        bytesRead,
		BytesWritten:     bytesWritten,
		PacketsReceived:  packetsReceived,
		PacketsSent:      packetsSent,
		PacketErrors:     m.GetPacketErrors(),
		PacketsPerSecond: float64(packetsReceived+packetsSent) / uptime.Seconds(),
		BytesPerSecond:   float64(bytesRead+bytesWritten) / uptime.Seconds(),
	}
}

// Stats structures for metrics reporting

type StreamStats struct {
	Label          string
	Uptime         time.Duration
	BytesRead      int64
	BytesWritten   int64
	ActiveStreams  int32
	BytesPerSecond float64
}

type PacketStats struct {
	Label            string
	Uptime           time.Duration
	BytesRead        int64
	BytesWritten     int64
	PacketsReceived  int64
	PacketsSent      int64
	PacketErrors     int64
	PacketsPerSecond float64
	BytesPerSecond   float64
}

// Registry for global metrics collection
type MetricsRegistry struct {
	streamMetrics map[string]*StreamMetrics
	packetMetrics map[string]*PacketMetrics
}

// Global registry instance
var globalRegistry = &MetricsRegistry{
	streamMetrics: make(map[string]*StreamMetrics),
	packetMetrics: make(map[string]*PacketMetrics),
}

// RegisterStreamMetrics adds stream metrics to the global registry
func RegisterStreamMetrics(metrics *StreamMetrics) {
	globalRegistry.streamMetrics[metrics.GetLabel()] = metrics
}

// RegisterPacketMetrics adds packet metrics to the global registry
func RegisterPacketMetrics(metrics *PacketMetrics) {
	globalRegistry.packetMetrics[metrics.GetLabel()] = metrics
}

// GetAllStreamStats returns stats for all registered stream metrics
func GetAllStreamStats() []StreamStats {
	stats := make([]StreamStats, 0, len(globalRegistry.streamMetrics))
	for _, m := range globalRegistry.streamMetrics {
		stats = append(stats, m.GetStats())
	}
	return stats
}

// GetAllPacketStats returns stats for all registered packet metrics
func GetAllPacketStats() []PacketStats {
	stats := make([]PacketStats, 0, len(globalRegistry.packetMetrics))
	for _, m := range globalRegistry.packetMetrics {
		stats = append(stats, m.GetStats())
	}
	return stats
}
