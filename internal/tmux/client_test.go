package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListWindowsParsing(t *testing.T) {
	// Create a mock tmux executable script
	tempDir := t.TempDir()
	mockBin := filepath.Join(tempDir, "tmux")

	script := `#!/bin/sh
if [ "$1" = "list-windows" ]; then
    echo "0:bash:0"
    echo "1:nvim:1"
    echo "2:harness:0"
    exit 0
fi
if [ "$1" = "has-session" ]; then
    exit 0
fi
exit 0
`
	if err := os.WriteFile(mockBin, []byte(script), 0755); err != nil {
		t.Fatalf("failed to write mock script: %v", err)
	}

	client := &Client{binPath: mockBin}

	windows, err := client.ListWindows("test-session")
	if err != nil {
		t.Fatalf("ListWindows failed: %v", err)
	}

	if len(windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(windows))
	}

	if windows[0].Index != 0 || windows[0].Name != "bash" || windows[0].Active {
		t.Errorf("unexpected window 0: %+v", windows[0])
	}
	if windows[1].Index != 1 || windows[1].Name != "nvim" || !windows[1].Active {
		t.Errorf("unexpected window 1: %+v", windows[1])
	}
	if windows[2].Index != 2 || windows[2].Name != "harness" || windows[2].Active {
		t.Errorf("unexpected window 2: %+v", windows[2])
	}
}

func TestHasSessionMock(t *testing.T) {
	tempDir := t.TempDir()
	mockBin := filepath.Join(tempDir, "tmux")
	script := `#!/bin/sh
if [ "$1" = "has-session" ] && [ "$3" = "exists" ]; then
    exit 0
else
    exit 1
fi
`
	if err := os.WriteFile(mockBin, []byte(script), 0755); err != nil {
		t.Fatalf("failed to write mock script: %v", err)
	}

	client := &Client{binPath: mockBin}

	if !client.HasSession("exists") {
		t.Errorf("expected session 'exists' to be true")
	}
	if client.HasSession("missing") {
		t.Errorf("expected session 'missing' to be false")
	}
}
