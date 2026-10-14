package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"botshell/internal/logger"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case updateCheckMsg:
		m.checkingUpd = false
		if msg.err != nil {
			m.updateStatus = "offline"
		} else {
			m.updateInfo = msg.info
			if msg.info.HasUpdate {
				m.updateStatus = fmt.Sprintf("UPDATE AVAILABLE: %s", msg.info.LatestVersion)
				logger.Global().Add("New version available: %s (%s)", msg.info.LatestVersion, msg.info.ReleaseURL)
			} else {
				m.updateStatus = "latest"
			}
		}
		return m, nil

	case tickMsg:
		m.refreshState()
		return m, tickEvery()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			logger.Global().Add("Quit command received. Shutting down...")
			return m, tea.Quit

		case "u": // Manual check update
			m.checkingUpd = true
			m.updateStatus = "checking..."
			return m, checkUpdateCmd

		case "a": // Toggle Always-Awake
			if m.powerGuard != nil {
				active, err := m.powerGuard.Toggle(m.pid)
				if err != nil {
					logger.Global().Add("Error toggling awake: %v", err)
				} else {
					logger.Global().Add("Always-Awake changed: %v", active)
				}
			}

		case "tab": // Switch views
			if m.viewMode == ViewLogs {
				m.viewMode = ViewScreens
			} else {
				m.viewMode = ViewLogs
			}
			m.refreshState()

		case "left", "h": // Previous window
			if len(m.windows) > 1 {
				m.selectedWin = (m.selectedWin - 1 + len(m.windows)) % len(m.windows)
				m.refreshState()
			}

		case "right", "l": // Next window
			if len(m.windows) > 1 {
				m.selectedWin = (m.selectedWin + 1) % len(m.windows)
				m.refreshState()
			}

		case "s": // Select window in tmux
			if len(m.windows) > 0 && m.selectedWin < len(m.windows) {
				idx := m.windows[m.selectedWin].Index
				_ = m.tmux.SelectWindow(m.cfg.DefaultSession, idx)
				logger.Global().Add("Selected tmux window [%d] %s", idx, m.windows[m.selectedWin].Name)
				m.refreshState()
			}

		case "b": // Run in background (detach daemon)
			logger.Global().Add("Detaching to background...")
			detachDaemon(m.cfgPath())
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) cfgPath() string {
	return m.cfg.DefaultSession
}

func detachDaemon(session string) {
	bin, err := os.Executable()
	if err != nil {
		return
	}
	cmd := exec.Command(bin, "--daemon")
	cmd.Env = os.Environ()
	_ = cmd.Start()
	fmt.Printf("\n[BotShell running in background with PID %s]\n", strconv.Itoa(cmd.Process.Pid))
}
