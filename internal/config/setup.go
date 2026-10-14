package config

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RunSetup conducts an interactive setup wizard to configure botshell.
func RunSetup(r io.Reader, w io.Writer, targetPath string) (*Config, error) {
	reader := bufio.NewReader(r)
	cfg := DefaultConfig()

	fmt.Fprintln(w, "=== BotShell Configuration Setup ===")
	fmt.Fprintln(w, "Setting up Telegram Bot and tmux environment.")
	fmt.Fprintln(w)

	// 1. Bot Token
	for {
		fmt.Fprint(w, "Enter Telegram Bot Token: ")
		token, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		token = strings.TrimSpace(token)
		if token != "" {
			cfg.BotToken = token
			break
		}
		fmt.Fprintln(w, "Bot Token cannot be empty. Please retry.")
	}

	// 2. Allowed Chat IDs
	for {
		fmt.Fprint(w, "Enter Allowed Chat/User IDs (comma-separated, e.g. 12345678,98765432): ")
		idsInput, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		idsInput = strings.TrimSpace(idsInput)
		if idsInput != "" {
			parts := strings.Split(idsInput, ",")
			validIDs := []int64{}
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if id, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
					validIDs = append(validIDs, id)
				}
			}
			if len(validIDs) > 0 {
				cfg.AllowedChatIDs = validIDs
				break
			}
		}
		fmt.Fprintln(w, "Please enter at least one valid numeric chat ID.")
	}

	// 3. Default Session Name
	fmt.Fprintf(w, "Enter default tmux session name [%s]: ", cfg.DefaultSession)
	sessionInput, err := reader.ReadString('\n')
	if err == nil {
		sessionInput = strings.TrimSpace(sessionInput)
		if sessionInput != "" {
			cfg.DefaultSession = sessionInput
		}
	}

	// 4. Always Awake
	fmt.Fprint(w, "Enable macOS Always-Awake while running? (Y/n) [Y]: ")
	awakeInput, err := reader.ReadString('\n')
	if err == nil {
		awakeInput = strings.TrimSpace(strings.ToLower(awakeInput))
		if awakeInput == "n" || awakeInput == "no" {
			cfg.AlwaysAwake = false
		} else {
			cfg.AlwaysAwake = true
		}
	}

	// Save to targetPath
	if targetPath != "" {
		if err := cfg.Save(targetPath); err != nil {
			return nil, fmt.Errorf("failed to save configuration: %w", err)
		}
		fmt.Fprintf(w, "\nConfiguration saved to %s\n", targetPath)
	}

	return cfg, nil
}
