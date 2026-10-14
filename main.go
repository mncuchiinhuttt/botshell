package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"botshell/internal/bot"
	"botshell/internal/config"
	"botshell/internal/logger"
	"botshell/internal/power"
	"botshell/internal/tmux"
	"botshell/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPathFlag := flag.String("config", "", "path to config JSON file")
	setupFlag := flag.Bool("setup", false, "run interactive setup wizard")
	daemonFlag := flag.Bool("daemon", false, "run silently in background without interactive TUI")
	flag.Parse()
	cfgPath := *configPathFlag
	if cfgPath == "" {
		var err error
		cfgPath, err = config.DefaultPath()
		if err != nil {
			log.Fatalf("Error determining default config path: %v", err)
		}
	}

	// Run interactive setup if requested or if config does not exist
	if *setupFlag {
		if _, err := config.RunSetup(os.Stdin, os.Stdout, cfgPath); err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
		fmt.Printf("Setup completed. Run 'botshell' to start.\n")
		return
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No configuration found at %s. Launching setup wizard...\n", cfgPath)
			cfg, err = config.RunSetup(os.Stdin, os.Stdout, cfgPath)
			if err != nil {
				log.Fatalf("Setup wizard failed: %v", err)
			}
		} else {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v. Please run 'botshell --setup'", err)
	}

	// 1. Initialize tmux client and ensure target session
	tmuxClient, err := tmux.NewClient()
	if err != nil {
		log.Fatalf("tmux error: %v", err)
	}

	if err := tmuxClient.EnsureSession(cfg.DefaultSession); err != nil {
		log.Fatalf("Failed to ensure tmux session %s: %v", cfg.DefaultSession, err)
	}
	log.Printf("tmux session '%s' ready.", cfg.DefaultSession)

	// 2. Start macOS Always-Awake Power Guard if configured
	var powerGuard *power.SleepGuard
	if cfg.AlwaysAwake {
		var err error
		powerGuard, err = power.StartSleepGuard(os.Getpid())
		if err != nil {
			log.Printf("Warning: failed to start caffeinate sleep guard: %v", err)
		} else {
			log.Println("macOS Always-Awake (caffeinate) active.")
		}
	}

	// 3. Setup context and signal listener for clean exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		if powerGuard != nil {
			log.Println("Stopping Sleep Guard...")
			_ = powerGuard.Stop()
		}
		log.Printf("Cleaning up tmux session '%s'...", cfg.DefaultSession)
		if err := tmuxClient.KillSession(cfg.DefaultSession); err != nil {
			log.Printf("Notice: failed to kill tmux session: %v", err)
		}
	}()

	// 4. Initialize and start Telegram Bot
	telegramBot, err := bot.New(cfg, tmuxClient, powerGuard)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram Bot: %v", err)
	}
	if err := telegramBot.RegisterCommands(); err != nil {
		log.Printf("Warning: failed to register slash commands menu: %v", err)
	}

	// Start Bot polling loop in background goroutine
	go func() {
		if err := telegramBot.Start(ctx); err != nil {
			logger.Global().Add("Bot error: %v", err)
		}
	}()

	if *daemonFlag {
		// Silent mode: block on context cancellation
		<-ctx.Done()
	} else {
		// Interactive TUI mode
		m := tui.NewModel(cfg, tmuxClient, powerGuard, os.Getpid())
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Printf("TUI error: %v", err)
		}
		stop()
	}

	log.Println("BotShell shutdown cleanly.")
}
