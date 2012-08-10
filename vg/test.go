// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vg

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"testing"
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
		font, err := MakeFont(fname, 20)
		if err != nil {
			t.Fatal(err)
		}

		w := font.Width(fname + "Xqg")
		h := font.Extents().Ascent

		c.FillString(font, 0, y-font.Extents().Descent, fname+"Xqg")
		fmt.Println(fname)

		var path Path
		path.Move(0, y+h)
		path.Line(w, y+h)
		path.Line(w, y)
		path.Line(0, y)
		path.Close()
		c.Stroke(path)

		path = Path{}
		c.SetColor(color.RGBA{B:255, A:255})
		c.SetLineDash([]Length{ Points(5), Points(3) }, 0)
		path.Move(0, y-font.Extents().Descent)
		path.Line(w, y-font.Extents().Descent)
		c.Stroke(path)
		c.SetColor(color.Black)
		c.SetLineDash([]Length{}, 0)

		y += h
	}
}

// DrawArcs draws some arcs to the canvas.
// The canvas is assumed to be 4 inches square.
func DrawArcs(t *testing.T, c Canvas) {
	green := color.RGBA{G: 255, A: 255}

	var p Path
	p.Move(Inches(3), Inches(2))
	p.Arc(Inches(2), Inches(2), Inches(1), 0, 2*math.Pi)
	c.SetColor(color.RGBA{B: 255, A: 255})
	c.Fill(p)

	p = Path{}
	p.Move(Inches(4), Inches(2))
	p.Line(Inches(3), Inches(2))
	p.Arc(Inches(2), Inches(2), Inches(1), 0, 5*math.Pi/2)
	p.Line(Inches(2), Inches(4))
	c.SetColor(color.RGBA{R: 255, A: 255})
	c.SetLineWidth(Points(3))
	c.Stroke(p)

	p = Path{}
	p.Move(Inches(0), Inches(2))
	p.Line(Inches(1), Inches(2))
	p.Arc(Inches(2), Inches(2), Inches(1), math.Pi, -7*math.Pi/2)
	p.Line(Inches(2), Inches(0))
	c.SetColor(color.Black)
	c.SetLineWidth(Points(1))
	c.Stroke(p)

	p = Path{}
	p.Move(Inches(0), Inches(1))
	p.Arc(Inches(1), Inches(1), Inches(1), math.Pi, math.Pi/2)
	c.SetLineWidth(Points(3))
	c.SetColor(green)
	c.Stroke(p)

	p = Path{}
	p.Move(Inches(1), Inches(0))
	p.Arc(Inches(1), Inches(1), Inches(1), 3*math.Pi/2, -math.Pi/2)
	c.SetLineWidth(Points(1))
	c.SetColor(color.Black)
	c.Stroke(p)

	p = Path{}
	p.Move(Inches(3), Inches(2))
	p.Arc(Inches(3), Inches(3), Inches(1), 3*math.Pi/2, 3*math.Pi/2)
	c.SetLineWidth(Points(3))
	c.SetColor(green)
	c.Stroke(p)

	p = Path{}
	p.Move(Inches(2), Inches(3))
	p.Arc(Inches(3), Inches(3), Inches(1), math.Pi, -3*math.Pi/2)
	c.SetLineWidth(Points(1))
	c.SetColor(color.Black)
	c.Stroke(p)
}
