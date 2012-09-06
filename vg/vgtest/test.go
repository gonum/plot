// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// vgtest contains some functions to "test"
// vg implementations.
package vgtest

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
	"sort"
	"testing"
)

// DrawFontExtents draws some text and denotes the
// various extents and width with lines. Expects
// about a 4x4 inch canvas.
func DrawFontExtents(t *testing.T, c vg.Canvas) {
	x, y := vg.Inches(1), vg.Inches(2)
	str := "Eloquent"
	font, err := vg.MakeFont("Times-Roman", 18)
	if err != nil {
		t.Fatal(err)
	}
	width := font.Width(str)
	ext := font.Extents()
	des := ext.Descent
	asc := ext.Ascent

	c.FillString(font, x, y, str)

	// baseline
	path := vg.Path{}
	path.Move(x, y)
	path.Line(x+width, y)
	c.Stroke(path)

	// descent
	c.SetColor(color.RGBA{G: 255, A: 255})
	path = vg.Path{}
	path.Move(x, y+des)
	path.Line(x+width, y+des)
	c.Stroke(path)

	// ascent
	c.SetColor(color.RGBA{B: 255, A: 255})
	path = vg.Path{}
	path.Move(x, y+asc)
	path.Line(x+width, y+asc)
	c.Stroke(path)
}

// DrawFonts draws some text in all of the various
// fonts along with a box to make sure that their
// sizes are computed correctly.
func DrawFonts(t *testing.T, c vg.Canvas) {
	y := vg.Points(0)
	var fonts []string
	for fname := range vg.FontMap {
		fonts = append(fonts, fname)
	}
	sort.Strings(fonts)
	for _, fname := range fonts {
		font, err := vg.MakeFont(fname, 20)
		if err != nil {
			t.Fatal(err)
		}

		w := font.Width(fname + "Xqg")
		h := font.Extents().Ascent

		c.FillString(font, 0, y-font.Extents().Descent, fname+"Xqg")
		fmt.Println(fname)

		var path vg.Path
		path.Move(0, y+h)
		path.Line(w, y+h)
		path.Line(w, y)
		path.Line(0, y)
		path.Close()
		c.Stroke(path)

		path = vg.Path{}
		c.SetColor(color.RGBA{B: 255, A: 255})
		c.SetLineDash([]vg.Length{vg.Points(5), vg.Points(3)}, 0)
		path.Move(0, y-font.Extents().Descent)
		path.Line(w, y-font.Extents().Descent)
		c.Stroke(path)
		c.SetColor(color.Black)
		c.SetLineDash([]vg.Length{}, 0)

		y += h
	}
}

// DrawArcs draws some arcs to the canvas.
// The canvas is assumed to be 4 inches square.
func DrawArcs(t *testing.T, c vg.Canvas) {
	green := color.RGBA{G: 255, A: 255}

	var p vg.Path
	p.Move(vg.Inches(3), vg.Inches(2))
	p.Arc(vg.Inches(2), vg.Inches(2), vg.Inches(1), 0, 2*math.Pi)
	c.SetColor(color.RGBA{B: 255, A: 255})
	c.Fill(p)

	p = vg.Path{}
	p.Move(vg.Inches(4), vg.Inches(2))
	p.Line(vg.Inches(3), vg.Inches(2))
	p.Arc(vg.Inches(2), vg.Inches(2), vg.Inches(1), 0, 5*math.Pi/2)
	p.Line(vg.Inches(2), vg.Inches(4))
	c.SetColor(color.RGBA{R: 255, A: 255})
	c.SetLineWidth(vg.Points(3))
	c.Stroke(p)

	p = vg.Path{}
	p.Move(vg.Inches(0), vg.Inches(2))
	p.Line(vg.Inches(1), vg.Inches(2))
	p.Arc(vg.Inches(2), vg.Inches(2), vg.Inches(1), math.Pi, -7*math.Pi/2)
	p.Line(vg.Inches(2), vg.Inches(0))
	c.SetColor(color.Black)
	c.SetLineWidth(vg.Points(1))
	c.Stroke(p)

	p = vg.Path{}
	p.Move(vg.Inches(0), vg.Inches(1))
	p.Arc(vg.Inches(1), vg.Inches(1), vg.Inches(1), math.Pi, math.Pi/2)
	c.SetLineWidth(vg.Points(3))
	c.SetColor(green)
	c.Stroke(p)

	p = vg.Path{}
	p.Move(vg.Inches(1), vg.Inches(0))
	p.Arc(vg.Inches(1), vg.Inches(1), vg.Inches(1), 3*math.Pi/2, -math.Pi/2)
	c.SetLineWidth(vg.Points(1))
	c.SetColor(color.Black)
	c.Stroke(p)

	p = vg.Path{}
	p.Move(vg.Inches(3), vg.Inches(2))
	p.Arc(vg.Inches(3), vg.Inches(3), vg.Inches(1), 3*math.Pi/2, 3*math.Pi/2)
	c.SetLineWidth(vg.Points(3))
	c.SetColor(green)
	c.Stroke(p)

	p = vg.Path{}
	p.Move(vg.Inches(2), vg.Inches(3))
	p.Arc(vg.Inches(3), vg.Inches(3), vg.Inches(1), math.Pi, -3*math.Pi/2)
	c.SetLineWidth(vg.Points(1))
	c.SetColor(color.Black)
	c.Stroke(p)
}
