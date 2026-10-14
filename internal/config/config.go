package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds user and runtime settings for botshell.
type Config struct {
	BotToken            string  `json:"bot_token"`
	AllowedChatIDs      []int64 `json:"allowed_chat_ids"`
	DefaultSession      string  `json:"default_session"`
	AlwaysAwake         bool    `json:"always_awake"`
	CaptureHistoryLines int     `json:"capture_history_lines"`
}

// DefaultConfig returns baseline options.
func DefaultConfig() *Config {
	return &Config{
		DefaultSession:      "botshell",
		AlwaysAwake:         true,
		CaptureHistoryLines: 40,
		AllowedChatIDs:      []int64{},
	}
}

// DefaultPath returns the standard configuration file path (~/.config/botshell/config.json).
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}
	return filepath.Join(home, ".config", "botshell", "config.json"), nil
}

// Load loads configuration from path. If file does not exist, returns os.ErrNotExist.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return cfg, nil
}

// Save writes configuration to disk with 0600 permissions.
func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", path, err)
	}

	return nil
}

// Validate checks whether essential parameters are provided.
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return errors.New("bot_token must not be empty")
	}
	if len(c.AllowedChatIDs) == 0 {
		return errors.New("allowed_chat_ids must contain at least one ID")
	}
	if c.DefaultSession == "" {
		c.DefaultSession = "botshell"
	}
	if c.CaptureHistoryLines <= 0 {
		c.CaptureHistoryLines = 40
	}
	return nil
}

// IsChatAllowed checks if a given Telegram chat ID or user ID is permitted.
func (c *Config) IsChatAllowed(chatID int64) bool {
	for _, id := range c.AllowedChatIDs {
		if id == chatID {
			return true
		}
	}
	return false
}
