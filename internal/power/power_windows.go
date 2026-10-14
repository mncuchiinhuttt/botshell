//go:build windows

package power

import (
	"syscall"
)

const (
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
	esContinuous      = 0x80000000
)

var (
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecutionState   = kernel32.NewProc("SetThreadExecutionState")
)

func setThreadExecutionState(flags uintptr) uintptr {
	ret, _, _ := procSetThreadExecutionState.Call(flags)
	return ret
}

func startPlatformGuard(pid int) (*SleepGuard, error) {
	// ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED
	setThreadExecutionState(esContinuous | esSystemRequired | esDisplayRequired)
	return &SleepGuard{active: true}, nil
}

func (s *SleepGuard) stopPlatformGuard() error {
	// Revert by calling with just ES_CONTINUOUS
	setThreadExecutionState(esContinuous)
	s.active = false
	return nil
}
