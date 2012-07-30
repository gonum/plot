// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vg

import (
	"fmt"
	"image/color"
	"testing"
	"sort"
)

// DrawFontExtents draws some text and denotes the
// various extents and width with lines. Expects
// about a 4x4 inch canvas.
func DrawFontExtents(t *testing.T, c Canvas) {
	x, y := Inches(1), Inches(2)
	str := "Eloquent"
	font, err := MakeFont("Times-Roman", 18)
	if err != nil {
		t.Fatal(err)
	}
	width := font.Width(str)
	ext := font.Extents()
	des := ext.Descent
	asc := ext.Ascent

	c.FillString(font, x, y, str)

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
	y := Points(0)
	var fonts []string
	for fname := range FontMap {
		fonts = append(fonts, fname)
	}
	sort.Strings(fonts)
	for _, fname := range fonts {
		font, err := MakeFont(fname, 12)
		if err != nil {
			t.Fatal(err)
		}

		w := font.Width(fname + "Xqg")
		h := font.Extents().Ascent

		// Shift the bottom font up so that its descents
		// aren't clipped.
		if y == 0.0 {
			y -= font.Extents().Descent
		}

		c.FillString(font, 0, y, fname+"Xqg")
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
