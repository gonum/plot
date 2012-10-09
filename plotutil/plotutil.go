// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// plotutil contains a small number of utilites for creating plots.
// This package is under active development so portions of
// it may change.
package plotutil

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
)

// DefaultColors is a set of colors used by the Color funciton.
var DefaultColors = []color.Color{
	color.RGBA{A: 255, R: 190},
	color.RGBA{A: 255, G: 190},
	color.RGBA{A: 255, B: 190},
	color.RGBA{A: 255, R: 190, B: 190},
	color.RGBA{A: 255, G: 190, B: 190},
	color.RGBA{A: 255, R: 250, G: 190, B: 10},
}

// Color returns the ith default color, wrapping
// if i is less than zero or greater than the max
// number of colors in the DefaultColors slice.
func Color(i int) color.Color {
	n := len(DefaultColors)
	if i < 0 {
		return DefaultColors[i%n+n]
	}
	return DefaultColors[i%n]
}

// DefaultGlyphShapes is a set of GlyphDrawers used by
// the Shape function.
var DefaultGlyphShapes = []plot.GlyphDrawer{
	plot.RingGlyph{},
	plot.SquareGlyph{},
	plot.TriangleGlyph{},
	plot.CrossGlyph{},
	plot.PlusGlyph{},
	plot.CircleGlyph{},
	plot.BoxGlyph{},
	plot.PyramidGlyph{},
}

// Shape returns the ith default glyph shape,
// wrapping if i is less than zero or greater
// than the max number of GlyphDrawers
// in the DefaultGlyphShapes slice.
func Shape(i int) plot.GlyphDrawer {
	n := len(DefaultGlyphShapes)
	if i < 0 {
		return DefaultGlyphShapes[i%n+n]
	}
	return DefaultGlyphShapes[i%n]
}

// DefaultDashes is a set of dash patterns used by
// the Dashes function.
var DefaultDashes = [][]vg.Length{
	{},

	{vg.Points(6), vg.Points(2)},

	{vg.Points(2), vg.Points(2)},

	{vg.Points(1), vg.Points(1)},

	{vg.Points(5), vg.Points(2), vg.Points(1), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(10), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(5), vg.Points(2), vg.Points(5), vg.Points(2),
		vg.Points(2), vg.Points(2), vg.Points(2), vg.Points(2)},

	{vg.Points(4), vg.Points(2), vg.Points(4), vg.Points(1),
		vg.Points(1), vg.Points(1), vg.Points(1), vg.Points(1),
		vg.Points(1), vg.Points(1)},
}

// Dashes returns the ith default dash pattern,
// wrapping if i is less than zero or greater
// than the max number of dash patters
// in the DefaultDashes slice.
func Dashes(i int) []vg.Length {
	n := len(DefaultDashes)
	if i < 0 {
		return DefaultDashes[i%n+n]
	}
	return DefaultDashes[i%n]
}
