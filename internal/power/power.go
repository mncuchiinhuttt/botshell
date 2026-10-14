package power

import "sync"

// SleepGuard controls system power assertion across macOS, Linux, and Windows.
type SleepGuard struct {
	mu     sync.Mutex
	active bool
	// platform-specific state managed in OS-specific files
	platformState any
}

// IsActive returns whether sleep prevention is currently active.
func (s *SleepGuard) IsActive() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}
