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

// NewLine returns a Line that uses the default line style and
// does not draw glyphs.
func NewLine(xys XYer) (*Line, error) {
	data, err := CopyXYs(xys)
	if err != nil {
		return nil, err
	}
	return &Line{
		XYs:       data,
		LineStyle: DefaultLineStyle,
	}, nil
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

// NewLinePoints returns both a Line and a
// Points for the given point data.
func NewLinePoints(xys XYer) (*Line, *Scatter, error) {
	s, err := NewScatter(xys)
	if err != nil {
		return nil, nil, err
	}
	l := &Line{
		XYs:       s.XYs,
		LineStyle: DefaultLineStyle,
	}
	return l, s, nil
}
