package render

import (
	"bytes"
	"image/png"
	"testing"
)

func TestParseANSI(t *testing.T) {
	raw := "\x1b[31mRed\x1b[0m \x1b[1;32mBoldGreen\x1b[0m"
	grid := ParseANSI(raw)

	if len(grid) != 1 {
		t.Fatalf("expected 1 row, got %d", len(grid))
	}

	row := grid[0]
	// "Red BoldGreen" has 13 characters
	if len(row) != 13 {
		t.Fatalf("expected 13 cells, got %d", len(row))
	}

	// First 3 chars ('R', 'e', 'd') should be Red (ANSI16[1])
	if row[0].Rune != 'R' || row[0].FG != ANSI16[1] {
		t.Errorf("expected 'R' to be red, got rune %c, fg %+v", row[0].Rune, row[0].FG)
	}

	// 4th char is space ' '
	if row[3].Rune != ' ' {
		t.Errorf("expected space, got %c", row[3].Rune)
	}

	// 'B' at index 4 should be Bold and Green (ANSI16[2])
	if row[4].Rune != 'B' || !row[4].Bold || row[4].FG != ANSI16[2] {
		t.Errorf("expected 'B' to be bold green, got: %+v", row[4])
	}
}

func TestRenderPNG(t *testing.T) {
	sampleANSI := "+---------------------------------+\n" +
		"| \x1b[1;34m[high]\x1b[0m Gemini 3.8 Flash (AG)  |\n" +
		"| \x1b[32mTODO:\x1b[0m Finish task 1             |\n" +
		"+---------------------------------+"

	pngBytes, err := RenderANSIToPNG(sampleANSI)
	if err != nil {
		t.Fatalf("RenderANSIToPNG failed: %v", err)
	}

	if len(pngBytes) == 0 {
		t.Fatal("expected non-empty png output")
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("failed to decode generated PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Errorf("unexpected image bounds: %+v", bounds)
	}
}
