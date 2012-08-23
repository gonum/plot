// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Line implements the Plotter interface, drawing a line.
type Line struct {
	// XYs is a copy of the points for this line.
	XYs

	// LineStyle is the style of the line connecting
	// the points.
	plot.LineStyle
}

// NewLine returns a Line that uses the
// default line style and does not draw
// glyphs.
func NewLine(xys XYer) *Line {
	return &Line{
		XYs:       CopyXYs(xys),
		LineStyle: DefaultLineStyle,
	}
}

// Plot draws the Line, implementing the plot.Plotter
// interface.
func (pts *Line) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)
	ps := make([]plot.Point, len(pts.XYs))
	for i, p := range pts.XYs {
		ps[i].X = trX(p.X)
		ps[i].Y = trY(p.Y)
	}
	da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger
// interface.
func (pts *Line) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(pts)
}

// Thumbnail the thumbnail for the Line,
// implementing the plot.Thumbnailer interface.
func (pts *Line) Thumbnail(da *plot.DrawArea) {
	y := da.Center().Y
	da.StrokeLine2(pts.LineStyle, da.Min.X, y, da.Max().X, y)
}

// LinePoints is a line with glyphs at each point.
type LinePoints struct {
	*Line
	Points *Scatter
}

// MakeLinePoints returns both a Line and a
// Points for the given point data.
func MakeLinePoints(xys XYer) LinePoints {
	s := NewScatter(xys)
	l := &Line{
		XYs:       s.XYs,
		LineStyle: DefaultLineStyle,
	}
	return LinePoints{ Line: l, Points: s }
}

// Plot implements the plot.Plotter interface.
func (l LinePoints) Plot(da plot.DrawArea, plt *plot.Plot) {
	l.Line.Plot(da, plt)
	l.Points.Plot(da, plt)
}

// DataRange implements the plot.DataRanger interface.
func (l LinePoints) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(l.Line)
}

// Thumbnail implements the plot.Thumbnailer interface.
func (l LinePoints) Thumbnail(da *plot.DrawArea) {
	l.Line.Thumbnail(da)
	l.Points.Thumbnail(da)
}