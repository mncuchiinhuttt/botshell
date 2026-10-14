package power

// Toggle switches the Always-Awake state.
func (s *SleepGuard) Toggle(pid int) (bool, error) {
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
	s.cmd = guard.cmd
	s.mu.Unlock()
	return true, nil
}
