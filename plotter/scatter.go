// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Scatter implements the Plotter interface, drawing
// glyphs at each of the given points.
type Scatter struct {
	XYer
	plot.GlyphStyle
}

// Plot implements the Plot method of the Plotter interface,
// drawing a glyph for each point in the Scatter.
func (s Scatter) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	for i := 0; i < s.Len(); i++ {
		da.DrawGlyph(s.GlyphStyle, plot.Point{trX(s.X(i)), trY(s.Y(i))})
	}
}

// GlyphBoxes returns a slice of GlyphBoxes, one for
// each of the glyphs in the Scatter.
func (s Scatter) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	for i := 0; i < s.Len(); i++ {
		x, y := p.X.Norm(s.X(i)), p.Y.Norm(s.Y(i))
		box := plot.GlyphBox{X: x, Y: y, Rect: s.Rect()}
		boxes = append(boxes, box)
	}
	return
}

// DataRange returns the minimum and maximum X and Y values
func (s Scatter) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = xDataRange(s)
	ymin, ymax = yDataRange(s)
	return
}

// Thumbnail draws a glyph in the center of a DrawArea
// as a thumbnail image representing the GlyhpStyle of
// the Scatter.
func (s Scatter) Thumbnail(da *plot.DrawArea) {
	da.DrawGlyph(s.GlyphStyle, da.Center())
}
