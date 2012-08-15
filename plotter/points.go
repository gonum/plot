// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Points implements the Plotter interface, drawing
// a set of points.
type Points struct {
	XYs

	// LineStyle is the style of the line connecting
	// the points.  If the color is nil then no line is
	// drawn.
	plot.LineStyle

	// GlyphStyle is the style of the glyphs drawn
	// at each point.  If Shape or Color are nil 
	// then no glyphs are drawn.
	plot.GlyphStyle
}

// NewLine returns a Points that uses the
// default line style and does not draw
// glyphs.
func NewLine(xys XYer) *Points {
	return &Points{
		XYs:       CopyXYs(xys),
		LineStyle: DefaultLineStyle,
	}
}

// NewLinePoints returns a Points that uses
// both the default line and glyph styles.
func NewLinePoints(xys XYer) *Points {
	return &Points{
		XYs:        CopyXYs(xys),
		LineStyle:  DefaultLineStyle,
		GlyphStyle: DefaultGlyphStyle,
	}
}

// NewScatter returns a Points that uses the
// default glyph style and does not draw a line.
func NewScatter(xys XYer) *Points {
	return &Points{
		XYs:        CopyXYs(xys),
		GlyphStyle: DefaultGlyphStyle,
	}
}

// Plot draws the Points, implementing the plot.Plotter
// interface.
func (pts *Points) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)
	ps := make([]plot.Point, len(pts.XYs))
	for i, p := range pts.XYs {
		ps[i].X = trX(p.X)
		ps[i].Y = trY(p.Y)
	}
	if pts.LineStyle.Color != nil {
		da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
	}
	if pts.GlyphStyle.Shape != nil && pts.GlyphStyle.Color != nil {
		for _, p := range ps {
			da.DrawGlyph(pts.GlyphStyle, p)
		}
	}
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger
// interface.
func (pts *Points) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(pts)
}

// GlyphBoxes returns a slice of plot.GlyphBoxes.
// If the GlyphStyle.Shape is non-nil then there
// is a plot.GlyphBox for each glyph, otherwise
// the returned slice is empty.  This implements the
// plot.GlyphBoxer interface.
func (pts *Points) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	if pts.GlyphStyle.Shape != nil && pts.GlyphStyle.Color != nil {
		return []plot.GlyphBox{}
	}
	bs := make([]plot.GlyphBox, len(pts.XYs))
	for i, p := range pts.XYs {
		bs[i].X = plt.X.Norm(p.X)
		bs[i].Y = plt.Y.Norm(p.Y)
		bs[i].Rect = pts.GlyphStyle.Rect()
	}
	return bs
}

// Thumbnail the thumbnail for the Points,
// implementing the plot.Thumbnailer interface.
func (pts *Points) Thumbnail(da *plot.DrawArea) {
	if pts.LineStyle.Color != nil {
		y := da.Center().Y
		da.StrokeLine2(pts.LineStyle, da.Min.X, y, da.Max().X, y)
	}
	if pts.GlyphStyle.Shape != nil && pts.GlyphStyle.Color != nil {
		da.DrawGlyph(pts.GlyphStyle, da.Center())
	}
}
