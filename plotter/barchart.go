// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"errors"
	"image/color"
	"math"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
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

	// XMin is the X location of the first bar.  XMin
	// can be changed to move groups of bars
	// down the X axis in order to make grouped
	// bar charts.
	XMin float64

	// stackedOn is the bar chart upon which
	// this bar chart is stacked.
	stackedOn *BarChart
}

// NewBarChart returns a new bar chart with a single bar for each value.
// The bars heights correspond to the values and their x locations correspond
// to the index of their value in the Valuer.
func NewBarChart(vs Valuer, width vg.Length) (*BarChart, error) {
	if width <= 0 {
		return nil, errors.New("Width parameter was not positive")
	}
	values, err := CopyValues(vs)
	if err != nil {
		return nil, err
	}
	return &BarChart{
		Values:    values,
		Width:     width,
		Color:     color.Black,
		LineStyle: DefaultLineStyle,
	}, nil
}

// BarHeight returns the maximum y value of the
// ith bar, taking into account any bars upon
// which it is stacked.
func (b *BarChart) BarHeight(i int) float64 {
	ht := 0.0
	if b == nil {
		return 0
	}
	if i >= 0 && i < len(b.Values) {
		ht += b.Values[i]
	}
	if b.stackedOn != nil {
		ht += b.stackedOn.BarHeight(i)
	}
	return ht
}

// StackOn stacks a bar chart on top of another,
// and sets the XMin and Offset to that of the
// chart upon which it is being stacked.
func (b *BarChart) StackOn(on *BarChart) {
	b.XMin = on.XMin
	b.Offset = on.Offset
	b.stackedOn = on
}

// Plot implements the plot.Plotter interface.
func (b *BarChart) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)

	for i, ht := range b.Values {
		x := b.XMin + float64(i)
		xmin := trX(float64(x))
		if !da.ContainsX(xmin) {
			continue
		}
		xmin = xmin - b.Width/2 + b.Offset
		xmax := xmin + b.Width
		bottom := b.stackedOn.BarHeight(i)
		ymin := trY(bottom)
		ymax := trY(bottom + ht)

		pts := []plot.Point{
			{xmin, ymin},
			{xmin, ymax},
			{xmax, ymax},
			{xmax, ymin},
		}
		poly := da.ClipPolygonY(pts)
		da.FillPolygon(b.Color, poly)

		pts = append(pts, plot.Pt(xmin, ymin))
		outline := da.ClipLinesY(pts)
		da.StrokeLines(b.LineStyle, outline...)
	}
}

// DataRange implements the plot.DataRanger interface.
func (b *BarChart) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = b.XMin
	xmax = xmin + float64(len(b.Values)-1)

	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for i, y := range b.Values {
		ybot := b.stackedOn.BarHeight(i)
		ytop := ybot + y
		ymin = math.Min(ymin, math.Min(ybot, ytop))
		ymax = math.Max(ymax, math.Max(ybot, ytop))
	}
	return
}

// GlyphBoxes implements the GlyphBoxer interface.
func (b *BarChart) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(b.Values))
	for i := range b.Values {
		x := b.XMin + float64(i)
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

	pts = append(pts, plot.Pt(da.Min.X, da.Min.Y))
	outline := da.ClipLinesY(pts)
	da.StrokeLines(b.LineStyle, outline...)
}
