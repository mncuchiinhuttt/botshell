//go:build linux

package power

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type linuxState struct {
	cmd *exec.Cmd
}

func startPlatformGuard(pid int) (*SleepGuard, error) {
	// 1. Try systemd-inhibit
	inhibitPath, err := exec.LookPath("systemd-inhibit")
	if err == nil {
		cmd := exec.Command(inhibitPath,
			"--what=idle:sleep:handle-lid-switch",
			"--who=BotShell",
			"--why=Remote terminal active",
			"sleep", "infinity",
		)
		if err := cmd.Start(); err == nil {
			return &SleepGuard{
				active:        true,
				platformState: &linuxState{cmd: cmd},
			}, nil
		}
	}

	// Fallback when no inhibitor is present (headless Linux servers usually don't sleep by default)
	return &SleepGuard{
		active:        true,
		platformState: &linuxState{},
	}, nil
}

func (s *SleepGuard) stopPlatformGuard() error {
	state, ok := s.platformState.(*linuxState)
	if !ok || state == nil || state.cmd == nil || state.cmd.Process == nil {
		s.active = false
		return nil
	}

	err := state.cmd.Process.Kill()
	_ = state.cmd.Wait()
	state.cmd = nil
	s.active = false

	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("failed to stop linux inhibitor: %w", err)
	}
	return nil
}
