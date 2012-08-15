// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
)

// GlyphBoxes implements the Plotter interface, drawing
// all of the glyph boxes of the plot.  This is intended for
// debugging.
type GlyphBoxes struct {
	plot.LineStyle
}

func NewGlyphBoxes() *GlyphBoxes {
	g := new(GlyphBoxes)
	g.Color = color.RGBA{R: 255, A: 255}
	g.Width = vg.Points(0.25)
	return g
}

func (g GlyphBoxes) Plot(da plot.DrawArea, plt *plot.Plot) {
	for _, b := range plt.GlyphBoxes(plt) {
		x := da.X(b.X) + b.Rect.Min.X
		y := da.Y(b.Y) + b.Rect.Min.Y
		da.StrokeLines(g.LineStyle, []plot.Point{
			{x, y},
			{x + b.Rect.Size.X, y},
			{x + b.Rect.Size.X, y + b.Rect.Size.Y},
			{x, y + b.Rect.Size.Y},
			{x, y},
		})
	}
}
