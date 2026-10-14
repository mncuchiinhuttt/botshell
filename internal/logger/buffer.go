package logger

import (
	"fmt"
	"sync"
	"time"
)

// RingBuffer stores the most recent log entries for TUI display.
type RingBuffer struct {
	mu       sync.RWMutex
	capacity int
	entries  []string
}

var globalBuffer = NewRingBuffer(500)

// NewRingBuffer creates a ring buffer with given maximum capacity.
func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 200
	}
	return &RingBuffer{
		capacity: capacity,
		entries:  make([]string, 0, capacity),
	}
}

// Global returns the shared global log buffer.
func Global() *RingBuffer {
	return globalBuffer
}

// Add appends a new message to the buffer with a timestamp prefix.
func (r *RingBuffer) Add(format string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("[%s] %s", timestamp, msg)

	if len(r.entries) >= r.capacity {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, entry)
}

// GetAll returns a copy of all current log entries.
func (r *RingBuffer) GetAll() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]string, len(r.entries))
	copy(copied, r.entries)
	return copied
}

// Last returns the most recent N log entries.
func (r *RingBuffer) Last(n int) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if n <= 0 || len(r.entries) == 0 {
		return nil
	}
	if n > len(r.entries) {
		n = len(r.entries)
	}

	start := len(r.entries) - n
	copied := make([]string, n)
	copy(copied, r.entries[start:])
	return copied
}
