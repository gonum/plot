package vecgfx

import (
	"fmt"
	"image/color"
	"testing"
)

// DrawFontExtents draws some text and denotes the
// various extents and width with lines. Expects
// about a 4x4 inch canvas.
func DrawFontExtents(t *testing.T, c Canvas) {
	x, y := c.DPI(), 2*c.DPI()
	str := "Eloquent"
	font, err := MakeFont("Times-Roman", 18)
	if err != nil {
		t.Fatal(err)
	}
	width := font.Width(str) / PtInch * c.DPI()
	ext := font.Extents()
	des := ext.Descent / PtInch * c.DPI()
	asc := ext.Ascent / PtInch * c.DPI()

	c.FillText(font, x, y, str)

	// baseline
	path := Path{}
	path.Move(x, y)
	path.Line(x+width, y)
	c.Stroke(path)

	// descent
	c.SetColor(color.RGBA{G: 255, A: 255})
	path = Path{}
	path.Move(x, y+des)
	path.Line(x+width, y+des)
	c.Stroke(path)

	// ascent
	c.SetColor(color.RGBA{B: 255, A: 255})
	path = Path{}
	path.Move(x, y+asc)
	path.Line(x+width, y+asc)
	c.Stroke(path)
}

// DrawFonts draws some text in all of the various
// fonts along with a box to make sure that their
// sizes are computed correctly.
func DrawFonts(t *testing.T, c Canvas) {
	y := 0.0
	for fname := range FontMap {
		font, err := MakeFont(fname, 12)
		if err != nil {
			t.Fatal(err)
		}

		w := font.Width(fname+"Xqg") / PtInch * c.DPI()
		h := font.Extents().Ascent / PtInch * c.DPI()

		// Shift the bottom font up so that its descents
		// aren't clipped.
		if y == 0.0 {
			y -= font.Extents().Descent / PtInch * c.DPI()
		}

		c.FillText(font, 0, y, fname+"Xqg")
		fmt.Println(fname)

		path := Path{}
		path.Move(0, y+h)
		path.Line(w, y+h)
		path.Line(w, y)
		path.Line(0, y)
		path.Close()
		c.Stroke(path)

		y += h
	}
}
