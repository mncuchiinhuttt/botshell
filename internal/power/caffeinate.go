package power

// StartSleepGuard initiates the sleep prevention guard for the host OS.
func StartSleepGuard(pid int) (*SleepGuard, error) {
	return startPlatformGuard(pid)
}

// Stop terminates the sleep prevention assertion.
func (s *SleepGuard) Stop() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopPlatformGuard()
}

// Toggle flips the sleep guard state.
func (s *SleepGuard) Toggle(pid int) (bool, error) {
	if s == nil {
		return false, nil
	}
	if s.IsActive() {
		if err := s.Stop(); err != nil {
			return false, err
		}
		return false, nil
	}

	guard, err := StartSleepGuard(pid)
	if err != nil {
		return false, err
	}
	s.mu.Lock()
	s.platformState = guard.platformState
	s.active = true
	s.mu.Unlock()
	return true, nil
}
