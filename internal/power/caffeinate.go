package power

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

// SleepGuard controls system power assertion via macOS caffeinate.
type SleepGuard struct {
	mu  sync.Mutex
	cmd *exec.Cmd
}

// StartSleepGuard initiates a caffeinate child process bound to the given PID.
// If pid <= 0, the current process PID is used.
func StartSleepGuard(pid int) (*SleepGuard, error) {
	if runtime.GOOS != "darwin" {
		return &SleepGuard{}, nil
	}

	caffeinatePath, err := exec.LookPath("caffeinate")
	if err != nil {
		caffeinatePath = "/usr/bin/caffeinate"
	}

	targetPID := pid
	if targetPID <= 0 {
		targetPID = os.Getpid()
	}

	// -d: display idle sleep prevention
	// -i: system idle sleep prevention
	// -s: system sleep prevention on AC power
	// -m: disk idle sleep prevention
	// -w: waits for process with given PID to exit
	cmd := exec.Command(caffeinatePath, "-dism", "-w", strconv.Itoa(targetPID))

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start caffeinate: %w", err)
	}

	return &SleepGuard{cmd: cmd}, nil
}

// Stop terminates the caffeinate child process if still running.
func (s *SleepGuard) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil || s.cmd.Process == nil {
		return nil
	}

	err := s.cmd.Process.Kill()
	_ = s.cmd.Wait()
	s.cmd = nil

	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	return nil
}

// IsActive returns whether caffeinate process is currently active.
func (s *SleepGuard) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil || s.cmd.Process == nil {
		return false
	}
	// Check if process is still running
	return s.cmd.ProcessState == nil || !s.cmd.ProcessState.Exited()
}
