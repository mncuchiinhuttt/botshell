package watcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"botshell/internal/logger"
	"botshell/internal/tmux"
)

// NotifyFunc is called when a completed task is detected.
type NotifyFunc func(session string)

// PaneWatcher monitors tmux panes for task completion.
type PaneWatcher struct {
	mu           sync.Mutex
	tmux         *tmux.Client
	session      string
	notifyFn     NotifyFunc
	lastHash     string
	changeCount  int
	idleCount    int
	waitingCmd   bool
	enabled      bool
	lastNotifyAt time.Time
}

// New creates a new PaneWatcher instance.
func New(t *tmux.Client, session string, notifyFn NotifyFunc) *PaneWatcher {
	return &PaneWatcher{
		tmux:     t,
		session:  session,
		notifyFn: notifyFn,
		enabled:  true,
	}
}

// SetWaitingCommand marks that a command was just sent and should be watched until completion.
func (w *PaneWatcher) SetWaitingCommand() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.waitingCmd = true
	w.changeCount = 2 // primed
	w.idleCount = 0
}

// SetEnabled toggles auto notification.
func (w *PaneWatcher) SetEnabled(enabled bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.enabled = enabled
}

// IsEnabled returns current auto-notify status.
func (w *PaneWatcher) IsEnabled() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.enabled
}

// Start begins the monitoring polling loop.
func (w *PaneWatcher) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.check()
		}
	}
}

func (w *PaneWatcher) check() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.enabled {
		return
	}

	content, err := w.tmux.CapturePane(w.session, false, -50)
	if err != nil {
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	if w.lastHash == "" {
		w.lastHash = hash
		return
	}

	if hash != w.lastHash {
		// Screen is actively changing
		w.lastHash = hash
		w.changeCount++
		w.idleCount = 0
		return
	}

	// Screen is identical to previous sample
	if w.changeCount >= 3 || w.waitingCmd {
		w.idleCount++

		// When output has been stable for 3 seconds after active execution
		if w.idleCount >= 3 {
			// Cooldown of 5 seconds to avoid repeated notifications on the same screen
			if time.Since(w.lastNotifyAt) > 5*time.Second {
				w.lastNotifyAt = time.Now()
				logger.Global().Add("Task finished in tmux session '%s'. Sending auto-capture notification...", w.session)
				if w.notifyFn != nil {
					go w.notifyFn(w.session)
				}
			}

			// Reset counters
			w.changeCount = 0
			w.idleCount = 0
			w.waitingCmd = false
		}
	}
}
