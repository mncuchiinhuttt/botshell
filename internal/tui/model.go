package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"botshell/internal/config"
	"botshell/internal/logger"
	"botshell/internal/power"
	"botshell/internal/tmux"
)

type ViewMode int

const (
	ViewLogs ViewMode = iota
	ViewScreens
)

type Model struct {
	cfg        *config.Config
	tmux       *tmux.Client
	powerGuard *power.SleepGuard
	pid        int

	viewMode     ViewMode
	windows      []tmux.Window
	selectedWin  int
	screenBuffer string
	logs         []string

	width  int
	height int
	err    error
}

type tickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewModel(cfg *config.Config, t *tmux.Client, p *power.SleepGuard, pid int) Model {
	m := Model{
		cfg:        cfg,
		tmux:       t,
		powerGuard: p,
		pid:        pid,
		viewMode:   ViewLogs,
		width:      80,
		height:     24,
	}
	m.refreshState()
	return m
}

func (m *Model) refreshState() {
	m.logs = logger.Global().Last(20)
	wins, err := m.tmux.ListWindows(m.cfg.DefaultSession)
	if err == nil {
		m.windows = wins
	}
	if len(m.windows) > 0 {
		if m.selectedWin >= len(m.windows) {
			m.selectedWin = 0
		}
		target := fmt.Sprintf("%s:%d", m.cfg.DefaultSession, m.windows[m.selectedWin].Index)
		scr, err := m.tmux.CapturePane(target, false, -20)
		if err == nil {
			m.screenBuffer = strings.TrimSpace(scr)
		}
	}
}

func (m Model) Init() tea.Cmd {
	return tickEvery()
}
