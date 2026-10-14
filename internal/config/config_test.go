package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSaveLoad(t *testing.T) {
	tempDir := t.TempDir()
	cfgPath := filepath.Join(tempDir, "config.json")

	cfg := &Config{
		BotToken:            "test_token_123",
		AllowedChatIDs:      []int64{111222, 333444},
		DefaultSession:      "my-session",
		AlwaysAwake:         true,
		CaptureHistoryLines: 50,
	}

	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.BotToken != cfg.BotToken {
		t.Errorf("expected BotToken %s, got %s", cfg.BotToken, loaded.BotToken)
	}
	if len(loaded.AllowedChatIDs) != 2 || loaded.AllowedChatIDs[0] != 111222 {
		t.Errorf("unexpected AllowedChatIDs: %v", loaded.AllowedChatIDs)
	}
	if !loaded.IsChatAllowed(111222) {
		t.Errorf("expected chat 111222 to be allowed")
	}
	if loaded.IsChatAllowed(999999) {
		t.Errorf("expected chat 999999 to not be allowed")
	}
}

func TestRunSetup(t *testing.T) {
	tempDir := t.TempDir()
	cfgPath := filepath.Join(tempDir, "setup_config.json")

	input := "my-secret-token\n12345, 67890\ncustom-tmux\ny\n"
	r := bytes.NewBufferString(input)
	w := &bytes.Buffer{}

	cfg, err := RunSetup(r, w, cfgPath)
	if err != nil {
		t.Fatalf("RunSetup failed: %v", err)
	}

	if cfg.BotToken != "my-secret-token" {
		t.Errorf("expected token my-secret-token, got %s", cfg.BotToken)
	}
	if len(cfg.AllowedChatIDs) != 2 || cfg.AllowedChatIDs[1] != 67890 {
		t.Errorf("unexpected chat ids: %v", cfg.AllowedChatIDs)
	}
	if cfg.DefaultSession != "custom-tmux" {
		t.Errorf("expected custom-tmux, got %s", cfg.DefaultSession)
	}
	if !cfg.AlwaysAwake {
		t.Errorf("expected AlwaysAwake true")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		t.Errorf("expected config file to be created: %v", err)
	}
}
