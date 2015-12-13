// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"errors"
	"math"
	"sort"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// fiveStatPlot contains the shared fields for quartile
// and box-whisker plots.
type fiveStatPlot struct {
	// Values is a copy of the values of the values used to
	// create this box plot.
	Values

	// Location is the location of the box along its axis.
	Location float64

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
}

// BoxPlot implements the Plotter interface, drawing
// a boxplot to represent the distribution of values.
type BoxPlot struct {
	fiveStatPlot

	// Offset is added to the x location of each box.
	// When the Offset is zero, the boxes are drawn
	// centered at their x location.
	Offset vg.Length

	// Width is the width used to draw the box.
	Width vg.Length

	// CapWidth is the width of the cap used to top
	// off a whisker.
	CapWidth vg.Length

	// GlyphStyle is the style of the outside point glyphs.
	GlyphStyle draw.GlyphStyle

	// BoxStyle is the line style for the box.
	BoxStyle draw.LineStyle

	// MedianStyle is the line style for the median line.
	MedianStyle draw.LineStyle

	// WhiskerStyle is the line style used to draw the
	// whiskers.
	WhiskerStyle draw.LineStyle
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
func NewBoxPlot(w vg.Length, loc float64, values Valuer) (*BoxPlot, error) {
	if w < 0 {
		return nil, errors.New("Negative boxplot width")
	}

	b := new(BoxPlot)
	var err error
	if b.fiveStatPlot, err = newFiveStat(w, loc, values); err != nil {
		return nil, err
	}

	b.Width = w
	b.CapWidth = 3 * w / 4

	b.GlyphStyle = DefaultGlyphStyle
	b.BoxStyle = DefaultLineStyle
	b.MedianStyle = DefaultLineStyle
	b.WhiskerStyle = draw.LineStyle{
		Width:  vg.Points(0.5),
		Dashes: []vg.Length{vg.Points(4), vg.Points(2)},
	}

	if len(b.Values) == 0 {
		b.Width = 0
		b.GlyphStyle.Radius = 0
		b.BoxStyle.Width = 0
		b.MedianStyle.Width = 0
		b.WhiskerStyle.Width = 0
	}

	return b, nil
}

func newFiveStat(w vg.Length, loc float64, values Valuer) (fiveStatPlot, error) {
	var b fiveStatPlot
	b.Location = loc

	var err error
	if b.Values, err = CopyValues(values); err != nil {
		return fiveStatPlot{}, err
	}

	sorted := make(Values, len(b.Values))
	copy(sorted, b.Values)
	sort.Float64s(sorted)

	if len(sorted) == 1 {
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

	return b, nil
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

func (b *BoxPlot) Plot(c draw.Canvas, plt *plot.Plot, xAxis, yAxis *plot.Axis) {
	trX, trY := plt.Transforms(&c, xAxis, yAxis)
	x := trX(b.Location)
	if !c.ContainsX(x) {
		return
	}
	x += b.Offset

	med := trY(b.Median)
	q1 := trY(b.Quartile1)
	q3 := trY(b.Quartile3)
	aLow := trY(b.AdjLow)
	aHigh := trY(b.AdjHigh)

	box := c.ClipLinesY([]draw.Point{
		{x - b.Width/2, q1},
		{x - b.Width/2, q3},
		{x + b.Width/2, q3},
		{x + b.Width/2, q1},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1},
	})
	c.StrokeLines(b.BoxStyle, box...)

	medLine := c.ClipLinesY([]draw.Point{
		{x - b.Width/2, med},
		{x + b.Width/2, med},
	})
	c.StrokeLines(b.MedianStyle, medLine...)

	cap := b.CapWidth / 2
	whisks := c.ClipLinesY([]draw.Point{{x, q3}, {x, aHigh}},
		[]draw.Point{{x - cap, aHigh}, {x + cap, aHigh}},
		[]draw.Point{{x, q1}, {x, aLow}},
		[]draw.Point{{x - cap, aLow}, {x + cap, aLow}})
	c.StrokeLines(b.WhiskerStyle, whisks...)

	for _, out := range b.Outside {
		y := trY(b.Value(out))
		if c.ContainsY(y) {
			c.DrawGlyphNoClip(b.GlyphStyle, draw.Point{x, y})
		}
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
func (b *BoxPlot) GlyphBoxes(plt *plot.Plot, x, y *plot.Axis) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(b.Outside)+1)
	for i, out := range b.Outside {
		bs[i].X = plt.X.Norm(b.Location)
		bs[i].Y = plt.Ys[0].Norm(b.Value(out))
		bs[i].Rectangle = b.GlyphStyle.Rectangle()
	}
	bs[len(bs)-1].X = x.Norm(b.Location)
	bs[len(bs)-1].Y = y.Norm(b.Median)
	bs[len(bs)-1].Rectangle = draw.Rectangle{
		Min: draw.Point{X: b.Offset - (b.Width/2 + b.BoxStyle.Width/2)},
		Max: draw.Point{X: b.Offset + (b.Width/2 + b.BoxStyle.Width/2)},
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
func MakeHorizBoxPlot(w vg.Length, loc float64, vs Valuer) (HorizBoxPlot, error) {
	b, err := NewBoxPlot(w, loc, vs)
	return HorizBoxPlot{b}, err
}

func (b HorizBoxPlot) Plot(c draw.Canvas, plt *plot.Plot, xAxis, yAxis *plot.Axis) {
	trX, trY := plt.Transforms(&c, xAxis, yAxis)
	y := trY(b.Location)
	if !c.ContainsY(y) {
		return
	}
	y += b.Offset

	med := trX(b.Median)
	q1 := trX(b.Quartile1)
	q3 := trX(b.Quartile3)
	aLow := trX(b.AdjLow)
	aHigh := trX(b.AdjHigh)

	box := c.ClipLinesX([]draw.Point{
		{q1, y - b.Width/2},
		{q3, y - b.Width/2},
		{q3, y + b.Width/2},
		{q1, y + b.Width/2},
		{q1, y - b.Width/2 - b.BoxStyle.Width/2},
	})
	c.StrokeLines(b.BoxStyle, box...)

	medLine := c.ClipLinesX([]draw.Point{
		{med, y - b.Width/2},
		{med, y + b.Width/2},
	})
	c.StrokeLines(b.MedianStyle, medLine...)

	cap := b.CapWidth / 2
	whisks := c.ClipLinesX([]draw.Point{{q3, y}, {aHigh, y}},
		[]draw.Point{{aHigh, y - cap}, {aHigh, y + cap}},
		[]draw.Point{{q1, y}, {aLow, y}},
		[]draw.Point{{aLow, y - cap}, {aLow, y + cap}})
	c.StrokeLines(b.WhiskerStyle, whisks...)

	for _, out := range b.Outside {
		x := trX(b.Value(out))
		if c.ContainsX(x) {
			c.DrawGlyphNoClip(b.GlyphStyle, draw.Point{x, y})
		}
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
func (b HorizBoxPlot) GlyphBoxes(plt *plot.Plot, x, y *plot.Axis) []plot.GlyphBox {
	bs := make([]plot.GlyphBox, len(b.Outside)+1)
	for i, out := range b.Outside {
		bs[i].X = plt.X.Norm(b.Value(out))
		bs[i].Y = plt.Ys[0].Norm(b.Location)
		bs[i].Rectangle = b.GlyphStyle.Rectangle()
	}
	bs[len(bs)-1].X = x.Norm(b.Median)
	bs[len(bs)-1].Y = y.Norm(b.Location)
	bs[len(bs)-1].Rectangle = draw.Rectangle{
		Min: draw.Point{Y: b.Offset - (b.Width/2 + b.BoxStyle.Width/2)},
		Max: draw.Point{Y: b.Offset + (b.Width/2 + b.BoxStyle.Width/2)},
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
