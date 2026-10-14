package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"botshell/internal/tmux"
)

// BuildScreenKeyboard creates the interactive control inline keyboard.
func BuildScreenKeyboard(currentSession string, currentWin int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		// Row 1: Quick Actions
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Refresh", "act:refresh"),
			tgbotapi.NewInlineKeyboardButtonData("Tabs", "act:tabs"),
			tgbotapi.NewInlineKeyboardButtonData("Text Mode", "act:text"),
		),
		// Row 2: Control Keys
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Esc", "key:Escape"),
			tgbotapi.NewInlineKeyboardButtonData("Ctrl+O", "key:C-o"),
			tgbotapi.NewInlineKeyboardButtonData("Ctrl+C", "key:C-c"),
			tgbotapi.NewInlineKeyboardButtonData("Tab", "key:Tab"),
		),
		// Row 3: Navigation & Confirmation
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Up", "key:Up"),
			tgbotapi.NewInlineKeyboardButtonData("Down", "key:Down"),
			tgbotapi.NewInlineKeyboardButtonData("Enter", "key:Enter"),
			tgbotapi.NewInlineKeyboardButtonData("y", "lit:y"),
			tgbotapi.NewInlineKeyboardButtonData("n", "lit:n"),
		),
	)
}

// BuildTabsKeyboard creates buttons for switching tmux windows.
func BuildTabsKeyboard(windows []tmux.Window) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for _, w := range windows {
		label := fmt.Sprintf("[%d] %s", w.Index, w.Name)
		if w.Active {
			label = fmt.Sprintf("* [%d] %s", w.Index, w.Name)
		}
		data := fmt.Sprintf("win:%d", w.Index)
		btn := tgbotapi.NewInlineKeyboardButtonData(label, data)
		currentRow = append(currentRow, btn)

		if len(currentRow) == 2 {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back to Screen", "act:refresh"),
	))

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}
