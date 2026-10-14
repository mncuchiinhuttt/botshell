package render

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	fontInitOnce sync.Once
	regularFace  font.Face
	boldFace     font.Face

	glyphWidth  int
	glyphHeight int
	fontAscent  int
	padX        = 16
	padY        = 16
)

func initFonts() {
	fontInitOnce.Do(func() {
		regFont, err := opentype.Parse(gomono.TTF)
		if err != nil {
			panic(err)
		}
		bldFont, err := opentype.Parse(gomonobold.TTF)
		if err != nil {
			panic(err)
		}

		opts := &opentype.FaceOptions{
			Size:    13,
			DPI:     120, // 1.25x crisp render
			Hinting: font.HintingFull,
		}

		rFace, err := opentype.NewFace(regFont, opts)
		if err != nil {
			panic(err)
		}
		bFace, err := opentype.NewFace(bldFont, opts)
		if err != nil {
			panic(err)
		}

		regularFace = rFace
		boldFace = bFace

		adv, _ := regularFace.GlyphAdvance('M')
		glyphWidth = adv.Ceil()
		if glyphWidth < 8 {
			glyphWidth = 10
		}

		metrics := regularFace.Metrics()
		glyphHeight = metrics.Height.Ceil()
		fontAscent = metrics.Ascent.Ceil()
	})
}

// RenderPNG transforms a 2D cell grid into an in-memory high-res PNG image.
func RenderPNG(grid [][]Cell) ([]byte, error) {
	initFonts()

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

	imgWidth := maxCols*glyphWidth + padX*2
	imgHeight := numRows*glyphHeight + padY*2

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	// Fill background with dark theme color
	draw.Draw(img, img.Bounds(), &image.Uniform{C: ColorDefaultBG}, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst: img,
	}

	for y, row := range grid {
		for x, cell := range row {
			cellLeft := padX + x*glyphWidth
			cellTop := padY + y*glyphHeight

			// Draw cell background if different from default
			if cell.BG != ColorDefaultBG {
				cellRect := image.Rect(cellLeft, cellTop, cellLeft+glyphWidth, cellTop+glyphHeight)
				draw.Draw(img, cellRect, &image.Uniform{C: cell.BG}, image.Point{}, draw.Src)
			}

			if cell.Rune == ' ' || cell.Rune == 0 {
				continue
			}

			if cell.Bold {
				drawer.Face = boldFace
			} else {
				drawer.Face = regularFace
			}

			drawer.Src = &image.Uniform{C: cell.FG}
			drawer.Dot = fixed.P(cellLeft, cellTop+fontAscent)
			drawer.DrawString(string(cell.Rune))
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
