package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botshell/internal/render"
	"botshell/internal/version"
)
// handleStart responds to /start command.
func (b *Bot) handleStart(msg *tgbotapi.Message) {
	text := fmt.Sprintf(
		"BotShell %s (by %s) Active.\n\n"+
			"Session: `%s`\n"+
			"Always-Awake: %v\n\n"+
			"BotShell Controls:\n"+
			"/screen - View terminal capture as high-res PNG\n"+
			"/text - View terminal capture as raw text\n"+
			"/tabs - List and switch tmux windows\n"+
			"/newtab <name> - Create a new tmux window\n"+
			"/keys <key> - Send named key (e.g. /keys Escape)\n"+
			"/update - Check for newer versions on GitHub\n\n"+
			"Terminal Pass-through:\n"+
			"- Any prompt or shell command will be typed into tmux automatically.\n"+
			"- Any slash command not listed above (e.g. `/exit`, `/help`, `/model`, `/compact`) is automatically forwarded to omp / terminal!\n"+
			"- Or prefix with `//` (e.g. `//exit`) to explicitly force send to terminal.",
		version.CurrentVersion(),
		version.Author,
		b.cfg.DefaultSession,
		b.powerGuard != nil && b.powerGuard.IsActive(),
	)
	m := tgbotapi.NewMessage(msg.Chat.ID, text)
	m.ParseMode = tgbotapi.ModeMarkdown
	_, _ = b.api.Send(m)

	b.sendScreen(msg.Chat.ID)
}

// handleScreen sends current tmux pane as PNG with inline control keyboard.
func (b *Bot) sendScreen(chatID int64) {
	raw, err := b.tmux.CapturePane(b.cfg.DefaultSession, true, -b.cfg.CaptureHistoryLines)
	if err != nil {
		m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error capturing pane: %v", err))
		_, _ = b.api.Send(m)
		return
	}

	pngBytes, err := render.RenderANSIToPNG(raw)
	if err != nil {
		m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error rendering screen: %v", err))
		_, _ = b.api.Send(m)
		return
	}

	photoFile := tgbotapi.FileBytes{
		Name:  "screen.png",
		Bytes: pngBytes,
	}

	photoMsg := tgbotapi.NewPhoto(chatID, photoFile)
	photoMsg.Caption = fmt.Sprintf("Session: %s | %s", b.cfg.DefaultSession, time.Now().Format("15:04:05"))
	photoMsg.ReplyMarkup = BuildScreenKeyboard(b.cfg.DefaultSession, 0)

	_, _ = b.api.Send(photoMsg)
}

// sendTextMode sends the current tmux pane as pure monospace text.
func (b *Bot) sendTextMode(chatID int64) {
	raw, err := b.tmux.CapturePane(b.cfg.DefaultSession, false, -b.cfg.CaptureHistoryLines)
	if err != nil {
		m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error capturing pane: %v", err))
		_, _ = b.api.Send(m)
		return
	}

	clean := strings.TrimSpace(raw)
	if clean == "" {
		clean = "(empty terminal screen)"
	}

	msgText := fmt.Sprintf("```\n%s\n```", clean)
	m := tgbotapi.NewMessage(chatID, msgText)
	m.ParseMode = tgbotapi.ModeMarkdown
	m.ReplyMarkup = BuildScreenKeyboard(b.cfg.DefaultSession, 0)
	_, _ = b.api.Send(m)
}

// sendTabsList sends the list of tmux windows with clickable switch buttons.
func (b *Bot) sendTabsList(chatID int64) {
	windows, err := b.tmux.ListWindows(b.cfg.DefaultSession)
	if err != nil {
		m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error listing windows: %v", err))
		_, _ = b.api.Send(m)
		return
	}

	m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Windows in session `%s`:", b.cfg.DefaultSession))
	m.ParseMode = tgbotapi.ModeMarkdown
	m.ReplyMarkup = BuildTabsKeyboard(windows)
	_, _ = b.api.Send(m)
}
// handleUpdateCommand checks for new releases on GitHub.
func (b *Bot) handleUpdateCommand(chatID int64) {
	info, err := version.CheckLatest()
	if err != nil {
		m := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error checking updates: %v", err))
		_, _ = b.api.Send(m)
		return
	}

	var reply string
	if info.HasUpdate {
		reply = fmt.Sprintf(
			"New version available!\n\n"+
				"Current: `%s`\n"+
				"Latest: `%s`\n\n"+
				"Release link: %s",
			version.CurrentVersion(),
			info.LatestVersion,
			info.ReleaseURL,
		)
	} else {
		reply = fmt.Sprintf("You are on the latest version (%s).", version.CurrentVersion())
	}
	m := tgbotapi.NewMessage(chatID, reply)
	m.ParseMode = tgbotapi.ModeMarkdown
	_, _ = b.api.Send(m)
}

