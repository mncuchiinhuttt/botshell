package watcher

import (
	"context"
	"testing"
	"time"
)

func TestPaneWatcherLifecycle(t *testing.T) {
	watcher := New(nil, "dummy-session", func(session string) {})

	if !watcher.IsEnabled() {
		t.Errorf("expected watcher to be enabled by default")
	}

	watcher.SetEnabled(false)
	if watcher.IsEnabled() {
		t.Errorf("expected watcher to be disabled")
	}

	watcher.SetWaitingCommand()
	if watcher.changeCount != 2 {
		t.Errorf("expected changeCount to be primed to 2, got %d", watcher.changeCount)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go watcher.Start(ctx)
	time.Sleep(20 * time.Millisecond)
	cancel()
}
