// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"errors"
	"image/color"
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// A BarChart presents grouped data with rectangular bars
// with lengths proportional to the data values.
type BarChart struct {
	Values

	// Width is the width of the bars.
	Width vg.Length

	// Color is the fill color of the bars.
	Color color.Color

	// LineStyle is the style of the outline of the bars.
	draw.LineStyle

	// Offset is added to the X location of each bar.
	// When the Offset is zero, the bars are drawn
	// centered at their X location.
	Offset vg.Length

	// XMin is the X location of the first bar.  XMin
	// can be changed to move groups of bars
	// down the X axis in order to make grouped
	// bar charts.
	XMin float64

	// Horizontal dictates whether the bars should be in the vertical
	// (default) or horizontal direction. If Horizontal is true, all
	// X locations and distances referred to here will actually be Y
	// locations and distances.
	Horizontal bool

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
func (b *BarChart) Plot(c draw.Canvas, plt *plot.Plot) {
	var trCat, trVal func(float64) vg.Length
	switch b.Horizontal {
	case false:
		trCat, trVal = plt.Transforms(&c)
	case true:
		trVal, trCat = plt.Transforms(&c)
	}

	for i, ht := range b.Values {
		catVal := b.XMin + float64(i)
		catMin := trCat(float64(catVal))
		switch b.Horizontal {
		case false:
			if !c.ContainsX(catMin) {
				continue
			}
		case true:
			if !c.ContainsY(catMin) {
				continue
			}
		}
		catMin = catMin - b.Width/2 + b.Offset
		catMax := catMin + b.Width
		bottom := b.stackedOn.BarHeight(i)
		valMin := trVal(bottom)
		valMax := trVal(bottom + ht)

		var pts []draw.Point
		var poly []draw.Point
		switch b.Horizontal {
		case false:
			pts = []draw.Point{
				{catMin, valMin},
				{catMin, valMax},
				{catMax, valMax},
				{catMax, valMin},
			}
			poly = c.ClipPolygonY(pts)
		case true:
			pts = []draw.Point{
				{valMin, catMin},
				{valMin, catMax},
				{valMax, catMax},
				{valMax, catMin},
			}
			poly = c.ClipPolygonX(pts)
		}
		c.FillPolygon(b.Color, poly)

		var outline [][]draw.Point
		switch b.Horizontal {
		case false:
			pts = append(pts, draw.Point{X: catMin, Y: valMin})
			outline = c.ClipLinesY(pts)
		case true:
			pts = append(pts, draw.Point{X: valMin, Y: catMin})
			outline = c.ClipLinesX(pts)
		}
		c.StrokeLines(b.LineStyle, outline...)
	}
}

// DataRange implements the plot.DataRanger interface.
func (b *BarChart) DataRange() (xmin, xmax, ymin, ymax float64) {
	catMin := b.XMin
	catMax := catMin + float64(len(b.Values)-1)

	valMin := math.Inf(1)
	valMax := math.Inf(-1)
	for i, val := range b.Values {
		valBot := b.stackedOn.BarHeight(i)
		valTop := valBot + val
		valMin = math.Min(valMin, math.Min(valBot, valTop))
		valMax = math.Max(valMax, math.Max(valBot, valTop))
	}
	if !b.Horizontal {
		return catMin, catMax, valMin, valMax
	}
	return valMin, valMax, catMin, catMax
}

// GlyphBoxes implements the GlyphBoxer interface.
func (b *BarChart) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(b.Values))
	for i := range b.Values {
		cat := b.XMin + float64(i)
		switch b.Horizontal {
		case false:
			boxes[i].X = plt.X.Norm(cat)
			boxes[i].Rectangle = draw.Rectangle{
				Min: draw.Point{X: b.Offset - b.Width/2},
				Max: draw.Point{X: b.Offset + b.Width/2},
			}
		case true:
			boxes[i].Y = plt.Y.Norm(cat)
			boxes[i].Rectangle = draw.Rectangle{
				Min: draw.Point{Y: b.Offset - b.Width/2},
				Max: draw.Point{Y: b.Offset + b.Width/2},
			}
		}
	}
	return boxes
}

// Thumbnail fulfills the plot.Thumbnailer interface.
func (b *BarChart) Thumbnail(c *draw.Canvas) {
	pts := []draw.Point{
		{c.Min.X, c.Min.Y},
		{c.Min.X, c.Max.Y},
		{c.Max.X, c.Max.Y},
		{c.Max.X, c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(b.Color, poly)

	pts = append(pts, draw.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	c.StrokeLines(b.LineStyle, outline...)
}
