// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Scatter implements the Plotter interface, drawing
// a glyph for each of a set of points.
type Scatter struct {
	// XYs is a copy of the points for this scatter.
	XYs

	// GlyphStyle is the style of the glyphs drawn
	// at each point.
	plot.GlyphStyle
}

// NewScatter returns a Scatter that uses the
// default glyph style.
func NewScatter(xys XYer) (*Scatter, error) {
	data, err := CopyXYs(xys)
	if err != nil {
		return nil, err
	}
	return &Scatter{
		XYs:        data,
		GlyphStyle: DefaultGlyphStyle,
	}, err
}

// Plot draws the Scatter, implementing the plot.Plotter
// interface.
func (pts *Scatter) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)
	for _, p := range pts.XYs {
		da.DrawGlyph(pts.GlyphStyle, plot.Pt(trX(p.X), trY(p.Y)))
	}
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger
// interface.
func (pts *Scatter) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(pts)
}

// GlyphBoxes returns a slice of plot.GlyphBoxes,
// implementing the plot.GlyphBoxer interface.
func (pts *Scatter) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(pts.XYs))
	for i, p := range pts.XYs {
		bs[i].X = plt.X.Norm(p.X)
		bs[i].Y = plt.Y.Norm(p.Y)
		bs[i].Rect = pts.GlyphStyle.Rect()
	}
	return bs
}

// Thumbnail the thumbnail for the Scatter,
// implementing the plot.Thumbnailer interface.
func (pts *Scatter) Thumbnail(da *plot.DrawArea) {
	da.DrawGlyph(pts.GlyphStyle, da.Center())
}
