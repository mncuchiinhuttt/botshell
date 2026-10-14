package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"botshell/internal/version"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cdd6f4")).
			Background(lipgloss.Color("#313244")).
			Padding(0, 1)

	badgeActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#11111b")).
			Background(lipgloss.Color("#a6e3a1")).
			Padding(0, 1)

	badgeAlert = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#11111b")).
			Background(lipgloss.Color("#f38ba8")).
			Padding(0, 1)

	badgeInactive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cdd6f4")).
			Background(lipgloss.Color("#45475a")).
			Padding(0, 1)

	authorBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89b4fa")).
			Background(lipgloss.Color("#1e1e2e")).
			Padding(0, 1)

	tabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cdd6f4")).
			Background(lipgloss.Color("#585b70")).
			Padding(0, 1)

	tabInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6adc8")).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#45475a")).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Padding(0, 1)
)

func (m Model) View() string {
	awakeBadge := badgeInactive.Render("SLEEP PERMITTED")
	if m.powerGuard != nil && m.powerGuard.IsActive() {
		awakeBadge = badgeActive.Render("ALWAYS-AWAKE ON")
	}

	updBadge := badgeInactive.Render("v: " + version.CurrentVersion())
	if m.updateInfo != nil && m.updateInfo.HasUpdate {
		updBadge = badgeAlert.Render("NEW UPDATE: " + m.updateInfo.LatestVersion)
	}

	header := headerStyle.Render("BotShell Control Panel "+version.CurrentVersion()) + " " +
		authorBadge.Render("by "+version.Author) + " " +
		awakeBadge + " " +
		updBadge + " " +
		badgeInactive.Render("Session: "+m.cfg.DefaultSession)
	// Tab header
	tabsView := ""
	if m.viewMode == ViewLogs {
		tabsView = tabActive.Render("[ Logs View ]") + " " + tabInactive.Render("[ Screens View ]")
	} else {
		tabsView = tabInactive.Render("[ Logs View ]") + " " + tabActive.Render("[ Screens View ]")
	}

	content := ""
	if m.viewMode == ViewLogs {
		content = m.renderLogsView()
	} else {
		content = m.renderScreensView()
	}

	footer := footerStyle.Render(
		"[Tab] View | [a] Awake | [s] Star Repo | [u] Check Upd | [Left/Right] Win | [Enter] Select | [b] BG | [q] Quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabsView,
		content,
		footer,
	)
}

func (m Model) renderLogsView() string {
	var b strings.Builder
	b.WriteString("--- Live Activity Logs ---\n")
	if len(m.logs) == 0 {
		b.WriteString("(no logs recorded yet)\n")
	} else {
		for _, logLine := range m.logs {
			b.WriteString(logLine + "\n")
		}
	}
	return boxStyle.Width(m.width - 4).Height(14).Render(b.String())
}

func (m Model) renderScreensView() string {
	var b strings.Builder

	winHeader := "Windows: "
	for i, w := range m.windows {
		label := fmt.Sprintf("[%d:%s]", w.Index, w.Name)
		if i == m.selectedWin {
			label = ">> " + label + " <<"
		}
		if w.Active {
			label = "*" + label
		}
		winHeader += label + " "
	}
	b.WriteString(winHeader + "\n\n")

	if m.screenBuffer == "" {
		b.WriteString("(empty screen buffer)")
	} else {
		b.WriteString(m.screenBuffer)
	}

	return boxStyle.Width(m.width - 4).Height(14).Render(b.String())
}
