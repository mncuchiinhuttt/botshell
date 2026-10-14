//go:build darwin

package power

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type darwinState struct {
	cmd *exec.Cmd
}

func startPlatformGuard(pid int) (*SleepGuard, error) {
	caffeinatePath, err := exec.LookPath("caffeinate")
	if err != nil {
		caffeinatePath = "/usr/bin/caffeinate"
	}

	targetPID := pid
	if targetPID <= 0 {
		targetPID = os.Getpid()
	}

	cmd := exec.Command(caffeinatePath, "-dism", "-w", strconv.Itoa(targetPID))
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start caffeinate: %w", err)
	}

	return &SleepGuard{
		active:        true,
		platformState: &darwinState{cmd: cmd},
	}, nil
}

func (s *SleepGuard) stopPlatformGuard() error {
	state, ok := s.platformState.(*darwinState)
	if !ok || state == nil || state.cmd == nil || state.cmd.Process == nil {
		s.active = false
		return nil
	}

	err := state.cmd.Process.Kill()
	_ = state.cmd.Wait()
	state.cmd = nil
	s.active = false

	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	return nil
}
