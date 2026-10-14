package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Window represents a tmux window.
type Window struct {
	Index  int    `json:"index"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// Client wraps operations on the tmux CLI.
type Client struct {
	binPath string
}

// NewClient creates a new tmux client instance, checking path.
func NewClient() (*Client, error) {
	path, err := exec.LookPath("tmux")
	if err != nil {
		// Fallback to common Homebrew / Mac paths
		candidates := []string{
			"/opt/homebrew/bin/tmux",
			"/usr/local/bin/tmux",
			"/usr/bin/tmux",
		}
		for _, c := range candidates {
			if _, lookErr := exec.LookPath(c); lookErr == nil {
				path = c
				err = nil
				break
			}
		}
	}
	if err != nil {
		return nil, errors.New("tmux is not installed or not in PATH. Please install with 'brew install tmux'")
	}
	return &Client{binPath: path}, nil
}

// Run executes a tmux command with provided args.
func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command(c.binPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("tmux %s failed: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

// HasSession checks if the session exists.
func (c *Client) HasSession(session string) bool {
	_, err := c.Run("has-session", "-t", session)
	return err == nil
}

// EnsureSession verifies the session exists or creates a new detached session.
func (c *Client) EnsureSession(session string) error {
	if c.HasSession(session) {
		return nil
	}
	_, err := c.Run("new-session", "-d", "-s", session)
	return err
}

// CapturePane extracts the visible pane content.
// If withANSI is true, passes -e to retain color escapes.
// If startLine < 0, passes -S <startLine> (e.g. -40 for last 40 lines).
func (c *Client) CapturePane(target string, withANSI bool, startLine int) (string, error) {
	args := []string{"capture-pane", "-t", target, "-p"}
	if withANSI {
		args = append(args, "-e")
	}
	if startLine < 0 {
		args = append(args, "-S", strconv.Itoa(startLine))
	}
	return c.Run(args...)
}

// SendLiteral sends verbatim text into the target pane (no key expansion).
func (c *Client) SendLiteral(target string, text string) error {
	_, err := c.Run("send-keys", "-t", target, "-l", text)
	return err
}

// SendKeys sends named keys into the target pane (e.g. Enter, Escape, C-c, C-o, Up, Down).
func (c *Client) SendKeys(target string, keys ...string) error {
	args := append([]string{"send-keys", "-t", target}, keys...)
	_, err := c.Run(args...)
	return err
}

// ListWindows returns all windows in the given session.
func (c *Client) ListWindows(session string) ([]Window, error) {
	// Format: index:name:active
	out, err := c.Run("list-windows", "-t", session, "-F", "#{window_index}:#{window_name}:#{window_active}")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	windows := make([]Window, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		idx, _ := strconv.Atoi(parts[0])
		active := parts[2] == "1"
		windows = append(windows, Window{
			Index:  idx,
			Name:   parts[1],
			Active: active,
		})
	}
	return windows, nil
}

// SelectWindow switches the active window in a session.
func (c *Client) SelectWindow(session string, index int) error {
	target := fmt.Sprintf("%s:%d", session, index)
	_, err := c.Run("select-window", "-t", target)
	return err
}

// NewWindow creates a new window in the specified session.
func (c *Client) NewWindow(session string, name string) error {
	args := []string{"new-window", "-t", session}
	if name != "" {
		args = append(args, "-n", name)
	}
	_, err := c.Run(args...)
	return err
}
