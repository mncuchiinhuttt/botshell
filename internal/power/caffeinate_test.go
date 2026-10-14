package power

import (
	"os"
	"testing"
	"time"
)

func TestStartSleepGuardLifecycle(t *testing.T) {
	guard, err := StartSleepGuard(os.Getpid())
	if err != nil {
		t.Fatalf("StartSleepGuard failed: %v", err)
	}

	if !guard.IsActive() {
		t.Errorf("expected guard to be active")
	}

	time.Sleep(50 * time.Millisecond)

	if err := guard.Stop(); err != nil {
		t.Fatalf("guard.Stop failed: %v", err)
	}

	if guard.IsActive() {
		t.Errorf("expected guard to be inactive after Stop")
	}

	// Calling Stop again should be safe (idempotent)
	if err := guard.Stop(); err != nil {
		t.Errorf("subsequent Stop should not fail, got: %v", err)
	}
}
