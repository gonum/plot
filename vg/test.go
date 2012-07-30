// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vg

import (
	"fmt"
	"image/color"
	"testing"
	"sort"
	"math"
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
	c.SetColor(color.RGBA{G:255, A:255})
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

// DrawArcs draws some arcs to the canvas.
// The canvas is assumed to be 4 inches square.
func DrawArcs(t *testing.T, c Canvas) {
	green := color.RGBA{G: 255, A: 255}

	var p0 Path
	p0.Move(Inches(0), Inches(1))
	p0.Arc(Inches(1), Inches(1), Inches(1), math.Pi, 3*math.Pi/2)
	c.SetLineWidth(Points(3))
	c.SetColor(green)
	c.Stroke(p0)

	var p05 Path
	p05.Move(Inches(1), Inches(0))
	p05.Arc(Inches(1), Inches(1), Inches(1), 3*math.Pi/2, math.Pi)
	c.SetLineWidth(Points(1))
	c.SetColor(color.Black)
	c.Stroke(p0)

	var p1 Path
	p1.Move(Inches(3), Inches(2))
	p1.Arc(Inches(3), Inches(3), Inches(1), 3*math.Pi/2, math.Pi)
	c.SetLineWidth(Points(3))
	c.SetColor(green)
	c.Stroke(p1)

	var p15 Path
	p15.Move(Inches(2), Inches(3))
	p15.Arc(Inches(3), Inches(3), Inches(1), math.Pi, 3*math.Pi/2)
	c.SetLineWidth(Points(1))
	c.SetColor(color.Black)
	c.Stroke(p1)

	var p2 Path
	p2.Move(Inches(3), Inches(2))
	p2.Arc(Inches(2), Inches(2), Inches(1), 0, 2*math.Pi)
	c.Stroke(p2)
}