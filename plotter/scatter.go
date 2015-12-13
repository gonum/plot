// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg/draw"
)

// Scatter implements the Plotter interface, drawing
// a glyph for each of a set of points.
type Scatter struct {
	// XYs is a copy of the points for this scatter.
	XYs

	// GlyphStyle is the style of the glyphs drawn
	// at each point.
	draw.GlyphStyle
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
func (pts *Scatter) Plot(c draw.Canvas, plt *plot.Plot, xAxis, yAxis *plot.Axis) {
	trX, trY := plt.Transforms(&c, xAxis, yAxis)
	for _, p := range pts.XYs {
		c.DrawGlyph(pts.GlyphStyle, draw.Point{trX(p.X), trY(p.Y)})
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
func (pts *Scatter) GlyphBoxes(plt *plot.Plot, x, y *plot.Axis) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(pts.XYs))
	for i, p := range pts.XYs {
		bs[i].X = x.Norm(p.X)
		bs[i].Y = y.Norm(p.Y)
		bs[i].Rectangle = pts.GlyphStyle.Rectangle()
	}
	return bs
}

// Thumbnail the thumbnail for the Scatter,
// implementing the plot.Thumbnailer interface.
func (pts *Scatter) Thumbnail(c *draw.Canvas) {
	c.DrawGlyph(pts.GlyphStyle, c.Center())
}
