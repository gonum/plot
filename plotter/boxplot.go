// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"math"
	"sort"
)

// BoxPlot implements the Plotter interface, drawing
// a boxplot to represent the distribution of values.
type BoxPlot struct {
	// Values is a copy of the values of the values used to
	// create this box plot.
	Values

	// Location is the location of the box along its axis.
	Location float64

	// Width is the width used to draw the box.
	Width vg.Length

	// CapWidth is the width of the cap used to top
	// off a whisker.
	CapWidth vg.Length

	// Median is the median value of the data.
	Median float64

	// Quartile1 and Quartile3 are the first and
	// third quartiles of the data respectively.
	Quartile1, Quartile3 float64

	// AdjLow and AdjHigh are the `adjacent' values
	// on the low and high ends of the data.  The
	// adjacent values are the points to which the
	// whiskers are drawn.
	AdjLow, AdjHigh float64

	// Min and Max are the extreme values of the data.
	Min, Max float64

	// Outside are the indices of Vs for the outside points.
	Outside []int

	// GlyphStyle is the style of the outside point glyphs.
	plot.GlyphStyle

	// BoxStyle is the line style for the box.
	BoxStyle plot.LineStyle

	// MedianStyle is the line style for the median line.
	MedianStyle plot.LineStyle

	// WhiskerStyle is the line style used to draw the
	// whiskers.
	WhiskerStyle plot.LineStyle
}

// NewBoxPlot returns a new BoxPlot that represents
// the distribution of the given values.  The style of
// the box plot is that used for Tukey's schematic
// plots is ``Exploratory Data Analysis.''
//
// An error is returned if the boxplot is created with
// no values.
//
// The fence values are 1.5x the interquartile before
// the first quartile and after the third quartile.  Any
// value that is outside of the fences are drawn as
// Outside points.  The adjacent values (to which the
// whiskers stretch) are the minimum and maximum
// values that are not outside the fences.
func NewBoxPlot(w vg.Length, loc float64, values Valuer) *BoxPlot {
	b := new(BoxPlot)
	b.Location = loc
	b.Width = w
	b.CapWidth = 3 * w / 4
	b.GlyphStyle = DefaultGlyphStyle
	b.BoxStyle = DefaultLineStyle
	b.MedianStyle = DefaultLineStyle
	b.WhiskerStyle = plot.LineStyle{
		Width:  vg.Points(0.5),
		Dashes: []vg.Length{vg.Points(4), vg.Points(2)},
	}

	b.Values = CopyValues(values)
	sorted := CopyValues(values)
	sort.Float64s(sorted)
	if len(sorted) == 0 {
		b.Width = 0
		b.GlyphStyle.Radius = 0
		b.BoxStyle.Width = 0
		b.MedianStyle.Width = 0
		b.WhiskerStyle.Width = 0
		return b
	} else if len(sorted) == 1 {
		b.Median = sorted[0]
		b.Quartile1 = sorted[0]
		b.Quartile3 = sorted[0]
	} else {
		b.Median = median(sorted)
		b.Quartile1 = median(sorted[:len(sorted)/2])
		b.Quartile3 = median(sorted[len(sorted)/2:])
	}
	b.Min = sorted[0]
	b.Max = sorted[len(sorted)-1]

	low := b.Quartile1 - 1.5*(b.Quartile3-b.Quartile1)
	high := b.Quartile3 + 1.5*(b.Quartile3-b.Quartile1)
	b.AdjLow = math.Inf(1)
	b.AdjHigh = math.Inf(-1)
	for i, v := range b.Values {
		if v > high || v < low {
			b.Outside = append(b.Outside, i)
			continue
		}
		if v < b.AdjLow {
			b.AdjLow = v
		}
		if v > b.AdjHigh {
			b.AdjHigh = v
		}
	}

	return b
}

// median returns the median value from a
// sorted Values.
func median(vs Values) float64 {
	if len(vs) == 1 {
		return vs[0]
	}
	med := vs[len(vs)/2]
	if len(vs)%2 == 0 {
		med += vs[len(vs)/2-1]
		med /= 2
	}
	return med
}

func (b *BoxPlot) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)
	x := trX(b.Location)
	if !da.ContainsX(x) {
		return
	}

	med := trY(b.Median)
	q1 := trY(b.Quartile1)
	q3 := trY(b.Quartile3)
	aLow := trY(b.AdjLow)
	aHigh := trY(b.AdjHigh)

	box := da.ClipLinesY([]plot.Point{
		{x - b.Width/2, q1},
		{x - b.Width/2, q3},
		{x + b.Width/2, q3},
		{x + b.Width/2, q1},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1},
	})
	da.StrokeLines(b.BoxStyle, box...)

	medLine := da.ClipLinesY([]plot.Point{
		{x - b.Width/2, med},
		{x + b.Width/2, med},
	})
	da.StrokeLines(b.MedianStyle, medLine...)

	cap := b.CapWidth / 2
	whisks := da.ClipLinesY([]plot.Point{{x, q3}, {x, aHigh}},
		[]plot.Point{{x - cap, aHigh}, {x + cap, aHigh}},
		[]plot.Point{{x, q1}, {x, aLow}},
		[]plot.Point{{x - cap, aLow}, {x + cap, aLow}})
	da.StrokeLines(b.WhiskerStyle, whisks...)

	for _, out := range b.Outside {
		y := trY(b.Value(out))
		da.DrawGlyph(b.GlyphStyle, plot.Point{x, y})
	}
}

// DataRange returns the minimum and maximum x
// and y values, implementing the plot.DataRanger
// interface.
func (b *BoxPlot) DataRange() (float64, float64, float64, float64) {
	return b.Location, b.Location, b.Min, b.Max
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot,
// implementing the plot.GlyphBoxer interface
func (b *BoxPlot) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(b.Outside)+1)
	for i, out := range b.Outside {
		bs[i].X = plt.X.Norm(b.Location)
		bs[i].Y = plt.Y.Norm(b.Value(out))
		bs[i].Rect = b.GlyphStyle.Rect()
	}
	bs[len(bs)-1].X = plt.X.Norm(b.Location)
	bs[len(bs)-1].Y = plt.Y.Norm(b.Median)
	bs[len(bs)-1].Rect = plot.Rect{
		Min:  plot.Point{X: -(b.Width/2 + b.BoxStyle.Width/2)},
		Size: plot.Point{X: b.Width + b.BoxStyle.Width},
	}
	return bs
}

// OutsideLabels returns a *Labels that will plot
// a label for each of the outside points.  The
// labels are assumed to correspond to the
// points used to create the box plot.
func (b *BoxPlot) OutsideLabels(labels Labeller) (*Labels, error) {
	strs := make([]string, len(b.Outside))
	for i, out := range b.Outside {
		strs[i] = labels.Label(out)
	}
	o := boxPlotOutsideLabels{b, strs}
	ls, err := NewLabels(o)
	if err != nil {
		return nil, err
	}
	ls.XOffset += b.GlyphStyle.Radius / 2
	ls.YOffset += b.GlyphStyle.Radius / 2
	return ls, nil
}

type boxPlotOutsideLabels struct {
	box    *BoxPlot
	labels []string
}

func (o boxPlotOutsideLabels) Len() int {
	return len(o.box.Outside)
}

func (o boxPlotOutsideLabels) XY(i int) (float64, float64) {
	return o.box.Location, o.box.Value(o.box.Outside[i])
}

func (o boxPlotOutsideLabels) Label(i int) string {
	return o.labels[i]
}

// HorizBoxPlot is like a regular BoxPlot, however,
// it draws horizontally instead of Vertically.
type HorizBoxPlot struct{ *BoxPlot }

// MakeHorizBoxPlot returns a HorizBoxPlot,
// plotting the values in a horizontal box plot
// centered along a fixed location of the y axis.
func MakeHorizBoxPlot(w vg.Length, loc float64, vs Values) HorizBoxPlot {
	return HorizBoxPlot{NewBoxPlot(w, loc, vs)}
}

func (b HorizBoxPlot) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, trY := plt.Transforms(&da)
	y := trY(b.Location)
	if !da.ContainsY(y) {
		return
	}

	med := trX(b.Median)
	q1 := trX(b.Quartile1)
	q3 := trX(b.Quartile3)
	aLow := trX(b.AdjLow)
	aHigh := trX(b.AdjHigh)

	box := da.ClipLinesX([]plot.Point{
		{q1, y - b.Width/2},
		{q3, y - b.Width/2},
		{q3, y + b.Width/2},
		{q1, y + b.Width/2},
		{q1, y - b.Width/2 - b.BoxStyle.Width/2},
	})
	da.StrokeLines(b.BoxStyle, box...)

	medLine := da.ClipLinesX([]plot.Point{
		{med, y - b.Width/2},
		{med, y + b.Width/2},
	})
	da.StrokeLines(b.MedianStyle, medLine...)

	cap := b.CapWidth / 2
	whisks := da.ClipLinesX([]plot.Point{{q3, y}, {aHigh, y}},
		[]plot.Point{{aHigh, y - cap}, {aHigh, y + cap}},
		[]plot.Point{{q1, y}, {aLow, y}},
		[]plot.Point{{aLow, y - cap}, {aLow, y + cap}})
	da.StrokeLines(b.WhiskerStyle, whisks...)

	for _, out := range b.Outside {
		x := trX(b.Value(out))
		da.DrawGlyph(b.GlyphStyle, plot.Point{x, y})
	}
}

// DataRange returns the minimum and maximum x
// and y values, implementing the plot.DataRanger
// interface.
func (b HorizBoxPlot) DataRange() (float64, float64, float64, float64) {
	return b.Min, b.Max, b.Location, b.Location
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot,
// implementing the plot.GlyphBoxer interface
func (b HorizBoxPlot) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(b.Outside)+1)
	for i, out := range b.Outside {
		bs[i].X = plt.X.Norm(b.Value(out))
		bs[i].Y = plt.Y.Norm(b.Location)
		bs[i].Rect = b.GlyphStyle.Rect()
	}
	bs[len(bs)-1].X = plt.X.Norm(b.Median)
	bs[len(bs)-1].Y = plt.Y.Norm(b.Location)
	bs[len(bs)-1].Rect = plot.Rect{
		Min:  plot.Point{Y: -(b.Width/2 + b.BoxStyle.Width/2)},
		Size: plot.Point{Y: b.Width + b.BoxStyle.Width},
	}
	return bs
}

// OutsideLabels returns a *Labels that will plot
// a label for each of the outside points.  The
// labels are assumed to correspond to the
// points used to create the box plot.
func (b *HorizBoxPlot) OutsideLabels(labels Labeller) (*Labels, error) {
	strs := make([]string, len(b.Outside))
	for i, out := range b.Outside {
		strs[i] = labels.Label(out)
	}
	o := horizBoxPlotOutsideLabels{
		boxPlotOutsideLabels{b.BoxPlot, strs},
	}
	ls, err := NewLabels(o)
	if err != nil {
		return nil, err
	}
	ls.XOffset += b.GlyphStyle.Radius / 2
	ls.YOffset += b.GlyphStyle.Radius / 2
	return ls, nil
}

type horizBoxPlotOutsideLabels struct {
	boxPlotOutsideLabels
}

func (o horizBoxPlotOutsideLabels) XY(i int) (float64, float64) {
	return o.box.Value(o.box.Outside[i]), o.box.Location
}
