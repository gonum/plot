// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Line implements the Plotter interface, drawing a line
// for the Plot method.
type Line struct {
	XYer
	plot.LineStyle
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (l Line) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	line := make([]plot.Point, l.Len())
	for i := range line {
		line[i].X = trX(l.X(i))
		line[i].Y = trY(l.Y(i))
	}
	da.StrokeLines(l.LineStyle, da.ClipLinesXY(line)...)
}

// DataRange returns the minimum and maximum X and Y values
func (l Line) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = xDataRange(l)
	ymin, ymax = yDataRange(l)
	return
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the Line.
func (l Line) Thumbnail(da *plot.DrawArea) {
	da.StrokeLine2(l.LineStyle, da.Min.X, da.Center().Y, da.Max().X, da.Center().Y)
}
