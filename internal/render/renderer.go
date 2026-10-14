package render

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	charWidth  = 7
	charHeight = 13
	padX       = 12
	padY       = 12
)

// RenderPNG transforms a 2D cell grid into an in-memory PNG image.
func RenderPNG(grid [][]Cell) ([]byte, error) {
	numRows := len(grid)
	if numRows == 0 {
		numRows = 1
	}

	maxCols := 1
	for _, row := range grid {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	imgWidth := maxCols*charWidth + padX*2
	imgHeight := numRows*charHeight + padY*2

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	// Fill whole background with dark theme color
	draw.Draw(img, img.Bounds(), &image.Uniform{C: ColorDefaultBG}, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  img,
		Face: basicfont.Face7x13,
	}

	for y, row := range grid {
		for x, cell := range row {
			cellLeft := padX + x*charWidth
			cellTop := padY + y*charHeight

			// Draw cell background if different from default
			if cell.BG != ColorDefaultBG {
				cellRect := image.Rect(cellLeft, cellTop, cellLeft+charWidth, cellTop+charHeight)
				draw.Draw(img, cellRect, &image.Uniform{C: cell.BG}, image.Point{}, draw.Src)
			}

			if cell.Rune == ' ' || cell.Rune == 0 {
				continue
			}

			// basicfont baseline is at +10
			drawer.Src = &image.Uniform{C: cell.FG}
			drawer.Dot = fixed.P(cellLeft, cellTop+10)
			drawer.DrawString(string(cell.Rune))

			// If bold, draw again offset by 1px
			if cell.Bold {
				drawer.Dot = fixed.P(cellLeft+1, cellTop+10)
				drawer.DrawString(string(cell.Rune))
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RenderANSIToPNG convenience function parsing ANSI text directly to PNG bytes.
func RenderANSIToPNG(ansiText string) ([]byte, error) {
	grid := ParseANSI(ansiText)
	return RenderPNG(grid)
}
