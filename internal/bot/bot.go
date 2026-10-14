package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"botshell/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botshell/internal/config"
	"botshell/internal/power"
	"botshell/internal/tmux"
)

// Bot orchestrates the Telegram listener and commands.
type Bot struct {
	api        *tgbotapi.BotAPI
	cfg        *config.Config
	tmux       *tmux.Client
	powerGuard *power.SleepGuard
}

// New creates and initializes a Bot instance.
func New(cfg *config.Config, t *tmux.Client, p *power.SleepGuard) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot api: %w", err)
	}

	return &Bot{
		api:        api,
		cfg:        cfg,
		tmux:       t,
		powerGuard: p,
	}, nil
}
// RegisterCommands registers the bot's slash command menu in Telegram.
func (b *Bot) RegisterCommands() error {
	cmds := []tgbotapi.BotCommand{
		{Command: "screen", Description: "Capture terminal screen as image (PNG)"},
		{Command: "text", Description: "Capture terminal screen as text code block"},
		{Command: "tabs", Description: "List & switch tmux windows"},
		{Command: "newtab", Description: "Create a new window (/newtab <name>)"},
		{Command: "keys", Description: "Send special key (/keys Escape, /keys C-c)"},
		{Command: "update", Description: "Check for newer versions on GitHub"},
		{Command: "start", Description: "Show welcome status & help info"},
	}
	cfg := tgbotapi.NewSetMyCommands(cmds...)
	_, err := b.api.Request(cfg)
	return err
}


// Start runs the update polling loop until context cancellation.
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)
	log.Printf("BotShell started as @%s. Listening for updates...", b.api.Self.UserName)
	logger.Global().Add("BotShell started as @%s. Listening for updates...", b.api.Self.UserName)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			b.handleUpdate(&update)
		}
	}
}

func (b *Bot) handleUpdate(update *tgbotapi.Update) {
	if update.Message != nil {
		if !b.cfg.IsChatAllowed(update.Message.Chat.ID) {
			log.Printf("Unauthorized message from chat ID %d", update.Message.Chat.ID)
			logger.Global().Add("Unauthorized message from chat ID %d", update.Message.Chat.ID)
			return
		}
		logger.Global().Add("Message from @%s: %s", update.Message.From.UserName, update.Message.Text)
		b.processMessage(update.Message)
	} else if update.CallbackQuery != nil {
		if !b.cfg.IsChatAllowed(update.CallbackQuery.Message.Chat.ID) {
			log.Printf("Unauthorized callback from chat ID %d", update.CallbackQuery.Message.Chat.ID)
			return
		}
		b.processCallback(update.CallbackQuery)
	}
}

func (b *Bot) processMessage(msg *tgbotapi.Message) {
	// 1. Check if user explicitly wants to send text to terminal via // (e.g. //exit or //changelog)
	if strings.HasPrefix(msg.Text, "//") {
		termCmd := strings.TrimPrefix(msg.Text, "/") // leaves single slash, e.g. "/exit"
		_ = b.tmux.SendLiteral(b.cfg.DefaultSession, termCmd)
		_ = b.tmux.SendKeys(b.cfg.DefaultSession, "Enter")
		time.Sleep(250 * time.Millisecond)
		b.sendScreen(msg.Chat.ID)
		return
	}

	// 2. Check if it's a known BotShell command
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			b.handleStart(msg)
			return
		case "screen", "s":
			b.sendScreen(msg.Chat.ID)
			return
		case "text", "t":
			b.sendTextMode(msg.Chat.ID)
			return
		case "tabs":
			b.sendTabsList(msg.Chat.ID)
			return
		case "newtab":
			name := strings.TrimSpace(msg.CommandArguments())
			_ = b.tmux.NewWindow(b.cfg.DefaultSession, name)
			time.Sleep(150 * time.Millisecond)
			b.sendScreen(msg.Chat.ID)
			return
		case "keys":
			key := strings.TrimSpace(msg.CommandArguments())
			if key != "" {
				_ = b.tmux.SendKeys(b.cfg.DefaultSession, key)
				time.Sleep(200 * time.Millisecond)
				b.sendScreen(msg.Chat.ID)
			}
			return
		case "update":
			b.handleUpdateCommand(msg.Chat.ID)
			return
		default:
			// Unknown command to BotShell (e.g. /exit, /help, /model, /clear, /compact, /plan)
			// Forward verbatim directly to the terminal!
			logger.Global().Add("Forwarding terminal command: %s", msg.Text)
			_ = b.tmux.SendLiteral(b.cfg.DefaultSession, msg.Text)
			_ = b.tmux.SendKeys(b.cfg.DefaultSession, "Enter")
			time.Sleep(250 * time.Millisecond)
			b.sendScreen(msg.Chat.ID)
			return
		}
	}

	// Normal text input typed by user: forward literal text and press Enter
	if msg.Text != "" {
		_ = b.tmux.SendLiteral(b.cfg.DefaultSession, msg.Text)
		_ = b.tmux.SendKeys(b.cfg.DefaultSession, "Enter")
		time.Sleep(250 * time.Millisecond)
		b.sendScreen(msg.Chat.ID)
	}
}

func (b *Bot) processCallback(cb *tgbotapi.CallbackQuery) {
	// Acknowledge callback immediately to remove spinner
	resp := tgbotapi.NewCallback(cb.ID, "")
	_, _ = b.api.Request(resp)

	data := cb.Data
	chatID := cb.Message.Chat.ID

	switch {
	case data == "act:refresh":
		b.sendScreen(chatID)
	case data == "act:tabs":
		b.sendTabsList(chatID)
	case data == "act:text":
		b.sendTextMode(chatID)
	case strings.HasPrefix(data, "key:"):
		key := strings.TrimPrefix(data, "key:")
		_ = b.tmux.SendKeys(b.cfg.DefaultSession, key)
		time.Sleep(200 * time.Millisecond)
		b.sendScreen(chatID)
	case strings.HasPrefix(data, "lit:"):
		lit := strings.TrimPrefix(data, "lit:")
		_ = b.tmux.SendLiteral(b.cfg.DefaultSession, lit)
		time.Sleep(200 * time.Millisecond)
		b.sendScreen(chatID)
	case strings.HasPrefix(data, "win:"):
		winStr := strings.TrimPrefix(data, "win:")
		if idx, err := strconv.Atoi(winStr); err == nil {
			_ = b.tmux.SelectWindow(b.cfg.DefaultSession, idx)
			time.Sleep(200 * time.Millisecond)
			b.sendScreen(chatID)
		}
	}
}
