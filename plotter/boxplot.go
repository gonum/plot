// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
)

// BoxPlot implements the Plotter interface, drawing a box plot.
type BoxPlot struct {
	Valuer

	// Loc is the X or Y location, in data coordinates,
	// around which the box is centered.
	Loc float64

	// Width is the width of the box.
	Width vg.Length

	// BoxStyle is the style used to draw the line
	// around the box, the median line.
	BoxStyle plot.LineStyle

	// MedianStyle is the style of the median line.
	MedianStyle plot.LineStyle

	// WhiskerStyle is the style used to draw the
	// whiskers.
	WhiskerStyle plot.LineStyle

	// CapWidth is the width of the cap on the whiskers.
	CapWidth vg.Length

	// GlyphStyle is the style of the points.
	GlyphStyle plot.GlyphStyle
}

// NewBoxPlot returns new Box representing a distribution
// of values.   The width parameter is the width of the
// box. The box surrounds the center of the range of
// data values, the middle line is the median, the points
// are as described in the Statistics method, and the
// whiskers extend to the extremes of all data that are
// not drawn as separate points.
func NewBoxPlot(width vg.Length, loc float64, vs Valuer) *BoxPlot {
	ws := DefaultLineStyle
	ws.Dashes = []vg.Length{ vg.Points(4), vg.Points(2) }
	return &BoxPlot{
		Valuer:          vs,
		Loc:            loc,
		Width:        width,
		BoxStyle:     DefaultLineStyle,
		MedianStyle:     DefaultLineStyle,
		WhiskerStyle: ws,
		CapWidth:     3*width / 4,
		GlyphStyle:   DefaultGlyphStyle,
	}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b *BoxPlot) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	x := trX(b.Loc)
	min, q1, med, q3, max, points := b.Statistics()

	q1y, medy, q3y := trY(q1), trY(med), trY(q3)
	da.StrokeLines(b.BoxStyle, da.ClipLinesY([]plot.Point{
		{x - b.Width/2, q1y}, {x - b.Width/2, q3y},
		{x + b.Width/2, q3y}, {x + b.Width/2, q1y},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1y}})...)

	da.StrokeLines(b.MedianStyle,  da.ClipLinesY([]plot.Point{
		{x - b.Width/2, medy}, {x + b.Width/2, medy},
	})...)

	miny, maxy := trY(min), trY(max)
	whisk := da.ClipLinesY([]plot.Point{{x, q3y}, {x, maxy}},
		[]plot.Point{{x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy}},
		[]plot.Point{{x, q1y}, {x, miny}},
		[]plot.Point{{x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, plot.Point{x, trY(b.Value(i))})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b *BoxPlot) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = b.Loc, b.Loc
	ymin, ymax = Range(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b *BoxPlot) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	_, _, med, _, _, pts := b.Statistics()
	boxes = append(boxes, plot.GlyphBox{
		X: p.X.Norm(b.Loc),
		Y: p.Y.Norm(med),
		Rect: plot.Rect{
			Min:  plot.Point{X: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: plot.Point{X: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.Loc), p.Y.Norm(b.Value(i))
		box := plot.GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
		boxes = append(boxes, box)
	}
	return
}

// X returns the Loc value of the box plot, i.e.,
// its X value.
func (b *BoxPlot) X(i int) float64 {
	return b.Loc
}

// Y returns the ith value of the box plot.
func (b *BoxPlot) Y(i int) float64 {
	return b.Value(i)
}

// HorizBox is a boxplot that draws horizontally
// instead of vertically.  The distribution of Y values
// are shown along the X axis.  The box is centered
// around the Y value that corresponds to the Loc
// field of the box.
type HorizBoxPlot struct {
	*BoxPlot
}

// NewHorizBoxPlot returns a HorizBox.  This is the
// same as NewBoxPlot except that the box draws
// horizontally instead of vertically.
func MakeHorizBoxPlot(width vg.Length, loc float64, vs Valuer) HorizBoxPlot {
	return HorizBoxPlot{NewBoxPlot(width, loc, vs)}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b HorizBoxPlot) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	y := trY(b.Loc)
	min, q1, med, q3, max, points := b.Statistics()

	q1x, medx, q3x := trX(q1), trX(med), trX(q3)
	da.StrokeLines(b.BoxStyle, da.ClipLinesX([]plot.Point{
		{q1x, y - b.Width/2}, {q3x, y - b.Width/2},
		{q3x, y + b.Width/2}, {q1x, y + b.Width/2},
		{q1x, y - b.Width/2 - b.BoxStyle.Width/2}})...)

	da.StrokeLines(b.MedianStyle, da.ClipLinesX([]plot.Point{
		{medx, y - b.Width/2}, {medx, y + b.Width/2},
	})...)

	minx, maxx := trX(min), trX(max)
	whisk := da.ClipLinesX([]plot.Point{{q3x, y}, {maxx, y}},
		[]plot.Point{{maxx, y - b.CapWidth/2}, {maxx, y + b.CapWidth/2}},
		[]plot.Point{{q1x, y}, {minx, y}},
		[]plot.Point{{minx, y - b.CapWidth/2}, {minx, y + b.CapWidth/2}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, plot.Point{trX(b.Value(i)), y})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b HorizBoxPlot) DataRange() (xmin, xmax, ymin, ymax float64) {
	ymin, ymax = b.Loc, b.Loc
	xmin, xmax = Range(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b HorizBoxPlot) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	_, _, med, _, _, pts := b.Statistics()
	boxes = append(boxes, plot.GlyphBox{
		X: p.X.Norm(med),
		Y: p.Y.Norm(b.Loc),
		Rect: plot.Rect{
			Min:  plot.Point{Y: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: plot.Point{Y: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.Value(i)), p.Y.Norm(b.Loc)
		box := plot.GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
		boxes = append(boxes, box)
	}
	return
}

// X returns the ith value of the box plot.
func (b *HorizBoxPlot) X(i int) float64 {
	return b.Value(i)
}

// Y returns the location of the box plot.
func (b *HorizBoxPlot) Y(i int) float64 {
	return b.Loc
}

// Statistics returns the 5 `boxplot' statistics, and
// the indices of `outside' points that should be drawn
// separately.
//
// The outside points are chosen as recommended
// by John Tukey in ``Exploratory Data Analysis'':
// Points that are more than 1.5x the inter-quartile
// range before the 1st and after the 3rd quartile.
func (b *BoxPlot) Statistics() (min, q1, med, q3, max float64, outside []int) {
	if b.Len() <= 1 {
		v := b.Value(0)
		return v, v, v, v, v, []int{}
	}

	sorted := SortedIndices(b)
	q1 = median(b, sorted[:len(sorted)/2])
	med = median(b, sorted)
	q3 = median(b, sorted[len(sorted)/2:])

	hlow := q1 - 1.5*(q3-q1)
	hhigh := q3 + 1.5*(q3-q1)
	min = b.Value(sorted[b.Len()-1])
	max = b.Value(sorted[0])

	for _, i := range sorted {
		v := b.Value(i)
		if v > hhigh || v < hlow {
			outside = append(outside, i)
		} else if v > max {
			max = v			
		} else if v < min {
			min = v
		}
	}
	return
}

// median returns the median value given a sorted
// slice of indices.
func median(vs Valuer, sorted []int) float64 {
	if len(sorted) == 1 {
		return vs.Value(sorted[0])
	}
	med := vs.Value(sorted[len(sorted)/2])
	if len(sorted)%2 == 0 {
		med += vs.Value(sorted[len(sorted)/2-1])
		med /= 2
	}
	return med
}

type boxPlotPoints struct {
	xys XYer
	labels Labeller
	inds []int
}

func (b boxPlotPoints) Len() int {
	return len(b.inds)
}

func (b boxPlotPoints) X(i int) float64 {
	return b.xys.X(b.inds[i])
}

func (b boxPlotPoints) Y(i int) float64 {
	return b.xys.Y(b.inds[i])
}

func (b boxPlotPoints) Label(i int) string {
	return b.labels.Label(b.inds[i])
}