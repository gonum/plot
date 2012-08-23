// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
)

type BarChart struct {
	Values

	// Width is the width of the bars.
	Width vg.Length

	// Color is the fill color of the bars.
	Color color.Color

	// LineStyle is the style of the outline of the bars.
	plot.LineStyle

	// Offset is added to the x location of each bar.
	// When the Offset is zero, the bars are drawn
	// centered at their x location.
	Offset vg.Length
}

// NewBarChart returns a new bar chart with
// a single bar for each value.  The bars heights
// correspond to the values and their x locations
// correspond to the index of their value in the
// Valuer.
func NewBarChart(vs Valuer, width vg.Length) *BarChart {
	return &BarChart{
		Values:    CopyValues(vs),
		Width:     width,
		Color:     color.Black,
		LineStyle: DefaultLineStyle,
	}
}

// Plot implements the plot.Plotter interface.
func (b *BarChart) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	for x, ht := range b.Values {
		xmin := trX(float64(x))
		if !da.ContainsX(xmin) {
			continue
		}
		xmin = xmin - b.Width/2 + b.Offset
		xmax := xmin + b.Width
		ymin := trY(0)
		ymax := trY(ht)

		pts := []plot.Point{
			{xmin, ymin},
			{xmin, ymax},
			{xmax, ymax},
			{xmax, ymin},
		}
		poly := da.ClipPolygonY(pts)
		da.FillPolygon(b.Color, poly)

		pts = append(pts, plot.Point{xmin, ymin})
		outline := da.ClipLinesY(pts)
		da.StrokeLines(b.LineStyle, outline...)
	}
}

// DataRange implements the plot.DataRanger interface.
func (b *BarChart) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = 0
	xmax = float64(len(b.Values) - 1)
	_, ymax = Range(b)
	ymin = 0
	return
}

// GlyphBoxes implements the GlyphBoxer interface.
func (b *BarChart) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(b.Values))
	for i := range b.Values {
		x := float64(i)
		boxes[i].X = plt.X.Norm(x)
		boxes[i].Rect = plot.Rect{
			Min:  plot.Point{X: b.Offset - b.Width/2},
			Size: plot.Point{X: b.Width},
		}
	}
	return boxes
}

func (b *BarChart) Thumbnail(da *plot.DrawArea) {
	pts := []plot.Point{
		{da.Min.X, da.Min.Y},
		{da.Min.X, da.Max().Y},
		{da.Max().X, da.Max().Y},
		{da.Max().X, da.Min.Y},
	}
	poly := da.ClipPolygonY(pts)
	da.FillPolygon(b.Color, poly)

	pts = append(pts, plot.Point{da.Min.X, da.Min.Y})
	outline := da.ClipLinesY(pts)
	da.StrokeLines(b.LineStyle, outline...)
}