package render

import (
	"image/color"
	"strconv"
	"strings"
)

// Cell represents a single character on the terminal screen with styling.
type Cell struct {
	Rune rune
	FG   color.RGBA
	BG   color.RGBA
	Bold bool
}

// Default colors: Modern Catppuccin Macchiato dark theme
var (
	ColorDefaultBG = color.RGBA{R: 0x1e, G: 0x1e, B: 0x2e, A: 0xff} // #1e1e2e
	ColorDefaultFG = color.RGBA{R: 0xcd, G: 0xd6, B: 0xf4, A: 0xff} // #cdd6f4
)

// Standard 16 ANSI colors
var ANSI16 = []color.RGBA{
	{R: 0x45, G: 0x47, B: 0x5a, A: 0xff}, // 0: Black
	{R: 0xf3, G: 0x8b, B: 0xa8, A: 0xff}, // 1: Red
	{R: 0xa6, G: 0xe3, B: 0xa1, A: 0xff}, // 2: Green
	{R: 0xf9, G: 0xe2, B: 0xaf, A: 0xff}, // 3: Yellow
	{R: 0x89, G: 0xb4, B: 0xfa, A: 0xff}, // 4: Blue
	{R: 0xcb, G: 0xa6, B: 0xf7, A: 0xff}, // 5: Magenta
	{R: 0x94, G: 0xe2, B: 0xd5, A: 0xff}, // 6: Cyan
	{R: 0xba, G: 0xc2, B: 0xde, A: 0xff}, // 7: White
	{R: 0x58, G: 0x5b, B: 0x70, A: 0xff}, // 8: Bright Black
	{R: 0xf3, G: 0x8b, B: 0xa8, A: 0xff}, // 9: Bright Red
	{R: 0xa6, G: 0xe3, B: 0xa1, A: 0xff}, // 10: Bright Green
	{R: 0xf9, G: 0xe2, B: 0xaf, A: 0xff}, // 11: Bright Yellow
	{R: 0x89, G: 0xb4, B: 0xfa, A: 0xff}, // 12: Bright Blue
	{R: 0xf5, G: 0xc2, B: 0xe7, A: 0xff}, // 13: Bright Magenta
	{R: 0x94, G: 0xe2, B: 0xd5, A: 0xff}, // 14: Bright Cyan
	{R: 0xa6, G: 0xad, B: 0xc8, A: 0xff}, // 15: Bright White
}

// ParseANSI processes text containing ANSI SGR escape sequences into a 2D grid of Cells.
func ParseANSI(raw string) [][]Cell {
	rawLines := strings.Split(raw, "\n")
	grid := make([][]Cell, 0, len(rawLines))

	for _, line := range rawLines {
		var row []Cell
		currentFG := ColorDefaultFG
		currentBG := ColorDefaultBG
		bold := false

		runes := []rune(line)
		n := len(runes)

		for i := 0; i < n; i++ {
			r := runes[i]
			if r == 0x1b && i+1 < n {
				next := runes[i+1]

				// 1. CSI sequence: ESC [
				if next == '[' {
					j := i + 2
					// Collect parameter bytes
					for j < n && ((runes[j] >= '0' && runes[j] <= '9') || runes[j] == ';' || runes[j] == '?') {
						j++
					}
					if j < n {
						finalChar := runes[j]
						if finalChar == 'm' {
							// SGR styling command
							sgrStr := string(runes[i+2 : j])
							currentFG, currentBG, bold = applySGR(sgrStr, currentFG, currentBG, bold)
						}
						// For all CSI commands (m, H, J, K, etc.), skip the whole sequence
						i = j
						continue
					}
				}

				// 2. OSC sequence: ESC ] (Operating System Command, e.g. OSC 8 hyperlinks, OSC 0 window title)
				// Format: ESC ] ... BEL (0x07) or ESC ] ... ESC \ (ST)
				if next == ']' {
					j := i + 2
					for j < n {
						if runes[j] == 0x07 { // BEL
							break
						}
						if runes[j] == 0x1b && j+1 < n && runes[j+1] == '\\' { // ST (ESC \)
							j++
							break
						}
						j++
					}
					if j < n {
						i = j
						continue
					}
				}

				// 3. Other escape sequences: ESC ( B, ESC = , ESC >, etc.
				if i+2 < n && (next == '(' || next == ')' || next == '*' || next == '+') {
					i += 2
					continue
				}
				if next >= ' ' && next <= '~' {
					i++
					continue
				}
			}

			// Handle carriage return and tabs
			if r == '\r' || r == '\t' {
				if r == '\t' {
					// Tab stops (4 spaces)
					for range 4 {
						row = append(row, Cell{Rune: ' ', FG: currentFG, BG: currentBG, Bold: bold})
					}
				}
				continue
			}

			// Filter unprintable control characters except standard space
			if r < 32 && r != ' ' {
				continue
			}

			row = append(row, Cell{Rune: r, FG: currentFG, BG: currentBG, Bold: bold})
		}

		grid = append(grid, row)
	}

	return grid
}

func applySGR(sgr string, fg, bg color.RGBA, bold bool) (color.RGBA, color.RGBA, bool) {
	if sgr == "" || sgr == "0" {
		return ColorDefaultFG, ColorDefaultBG, false
	}

	parts := strings.Split(sgr, ";")
	for k := 0; k < len(parts); k++ {
		code, err := strconv.Atoi(parts[k])
		if err != nil {
			continue
		}

		switch {
		case code == 0:
			fg = ColorDefaultFG
			bg = ColorDefaultBG
			bold = false
		case code == 1:
			bold = true
		case code == 22:
			bold = false
		case code >= 30 && code <= 37:
			fg = ANSI16[code-30]
		case code == 39:
			fg = ColorDefaultFG
		case code >= 40 && code <= 47:
			bg = ANSI16[code-40]
		case code == 49:
			bg = ColorDefaultBG
		case code >= 90 && code <= 97:
			fg = ANSI16[code-90+8]
		case code >= 100 && code <= 107:
			bg = ANSI16[code-100+8]
		case code == 38 && k+2 < len(parts) && parts[k+1] == "5":
			// 256 color FG: 38;5;idx
			idx, _ := strconv.Atoi(parts[k+2])
			if idx >= 0 && idx < 16 {
				fg = ANSI16[idx]
			}
			k += 2
		case code == 48 && k+2 < len(parts) && parts[k+1] == "5":
			// 256 color BG: 48;5;idx
			idx, _ := strconv.Atoi(parts[k+2])
			if idx >= 0 && idx < 16 {
				bg = ANSI16[idx]
			}
			k += 2
		case code == 38 && k+4 < len(parts) && parts[k+1] == "2":
			// RGB FG: 38;2;r;g;b
			r, _ := strconv.Atoi(parts[k+2])
			g, _ := strconv.Atoi(parts[k+3])
			b, _ := strconv.Atoi(parts[k+4])
			fg = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
			k += 4
		case code == 48 && k+4 < len(parts) && parts[k+1] == "2":
			// RGB BG: 48;2;r;g;b
			r, _ := strconv.Atoi(parts[k+2])
			g, _ := strconv.Atoi(parts[k+3])
			b, _ := strconv.Atoi(parts[k+4])
			bg = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}
			k += 4
		}
	}

	return fg, bg, bold
}
