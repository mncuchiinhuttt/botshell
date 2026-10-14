package bot

import (
	"testing"

	"botshell/internal/tmux"
)

func TestBuildScreenKeyboard(t *testing.T) {
	kb := BuildScreenKeyboard("test-session", 0)
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows in screen keyboard, got %d", len(kb.InlineKeyboard))
	}

	// Verify Refresh button
	row0 := kb.InlineKeyboard[0]
	if len(row0) < 3 || *row0[0].CallbackData != "act:refresh" {
		t.Errorf("expected refresh button in row 0")
	}

	// Verify Escape button in row 1
	row1 := kb.InlineKeyboard[1]
	if len(row1) < 4 || *row1[0].CallbackData != "key:Escape" {
		t.Errorf("expected Escape key button in row 1")
	}
}

func TestBuildTabsKeyboard(t *testing.T) {
	windows := []tmux.Window{
		{Index: 0, Name: "main", Active: false},
		{Index: 1, Name: "nvim", Active: true},
		{Index: 2, Name: "harness", Active: false},
	}

	kb := BuildTabsKeyboard(windows)
	// 3 windows in 2-column format -> 2 rows for windows + 1 row for Back button = 3 rows
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows in tabs keyboard, got %d", len(kb.InlineKeyboard))
	}

	firstBtn := kb.InlineKeyboard[0][0]
	if *firstBtn.CallbackData != "win:0" {
		t.Errorf("expected first button to target win:0, got: %s", *firstBtn.CallbackData)
	}

	activeBtn := kb.InlineKeyboard[0][1]
	if activeBtn.Text != "* [1] nvim" {
		t.Errorf("expected active marker on active window, got: %s", activeBtn.Text)
	}
}
