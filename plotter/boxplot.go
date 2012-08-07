// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"sort"
)

// BoxPlot implements the Plotter interface, drawing a box plot.
type BoxPlot struct {
	Yer

	// X is the X or Y value, in data coordinates, around
	// which the box is centered.
	X float64

	// Width is the width of the box.
	Width vg.Length

	// BoxStyle is the style used to draw the line
	// around the box, the median line.
	BoxStyle plot.LineStyle

	// MedStyle is the style of the median line.
	MedStyle plot.LineStyle

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
func NewBoxPlot(width vg.Length, x float64, ys Yer) *BoxPlot {
	ws := DefaultLineStyle
	ws.Dashes = []vg.Length{ vg.Points(4), vg.Points(2) }
	return &BoxPlot{
		Yer:          ys,
		X:            x,
		Width:        width,
		BoxStyle:     DefaultLineStyle,
		MedStyle:     DefaultLineStyle,
		WhiskerStyle: ws,
		CapWidth:     3*width / 4,
		GlyphStyle:   DefaultGlyphStyle,
	}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b *BoxPlot) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	x := trX(b.X)
	min, q1, med, q3, max, points := b.Statistics()

	q1y, medy, q3y := trY(q1), trY(med), trY(q3)
	da.StrokeLines(b.BoxStyle, da.ClipLinesY([]plot.Point{
		{x - b.Width/2, q1y}, {x - b.Width/2, q3y},
		{x + b.Width/2, q3y}, {x + b.Width/2, q1y},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1y}})...)

	da.StrokeLines(b.MedStyle,  da.ClipLinesY([]plot.Point{
		{x - b.Width/2, medy}, {x + b.Width/2, medy},
	})...)

	miny, maxy := trY(min), trY(max)
	da.StrokeLines(b.WhiskerStyle, da.ClipLinesY([]plot.Point{{x, q3y}, {x, maxy}},
		[]plot.Point{{x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy}},
		[]plot.Point{{x, q1y}, {x, miny}},
		[]plot.Point{{x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny}})...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, plot.Point{x, trY(b.Y(i))})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b *BoxPlot) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = b.X, b.X
	ymin, ymax = yDataRange(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b *BoxPlot) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	_, _, med, _, _, pts := b.Statistics()
	boxes = append(boxes, plot.GlyphBox{
		X: p.X.Norm(b.X),
		Y: p.Y.Norm(med),
		Rect: plot.Rect{
			Min:  plot.Point{X: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: plot.Point{X: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.X), p.Y.Norm(b.Y(i))
		box := plot.GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
		boxes = append(boxes, box)
	}
	return
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
		y0 := b.Y(0)
		return y0, y0, y0, y0, y0, []int{}
	}

	sorted := make([]int, b.Len())
	for i := range sorted {
		sorted[i] = i
	}
	sort.Sort(ySorter{b, sorted})

	q1 = median(b, sorted[:len(sorted)/2])
	med = median(b, sorted)
	q3 = median(b, sorted[len(sorted)/2:])

	hlow := q1 - 1.5*(q3-q1)
	hhigh := q3 + 1.5*(q3-q1)
	min = b.Y(sorted[b.Len()-1])
	max = b.Y(sorted[0])
	for _, i := range sorted {
		y := b.Y(i)
		if y > hhigh || y < hlow {
			outside = append(outside, i)
		} else if y > max {
			max = y			
		} else if y < min {
			min = y
		}
	}
	return
}

// median returns the median Y value given a sorted
// slice of indices.
func median(ys Yer, sorted []int) float64 {
	if len(sorted) == 1 {
		return ys.Y(sorted[0])
	}
	med := ys.Y(sorted[len(sorted)/2])
	if len(sorted)%2 == 0 {
		med += ys.Y(sorted[len(sorted)/2-1])
		med /= 2
	}
	return med
}

// HorizBox is a boxplot that draws horizontally
// instead of vertically.  The distribution of Y values
// are shown along the X axis.  The box is centered
// around the Y value that corresponds to the X
// value of the box.
type HorizBoxPlot struct {
	*BoxPlot
}

// NewHorizBoxPlot returns a HorizBox.  This is the
// same as NewBoxPlot except that the box draws
// horizontally instead of vertically.
func MakeHorizBoxPlot(width vg.Length, y float64, vals Yer) HorizBoxPlot {
	return HorizBoxPlot{NewBoxPlot(width, y, vals)}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b HorizBoxPlot) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	y := trY(b.X)
	min, q1, med, q3, max, points := b.Statistics()

	q1x, medx, q3x := trX(q1), trX(med), trX(q3)
	da.StrokeLines(b.BoxStyle, da.ClipLinesX([]plot.Point{
		{q1x, y - b.Width/2}, {q3x, y - b.Width/2},
		{q3x, y + b.Width/2}, {q1x, y + b.Width/2},
		{q1x, y - b.Width/2 - b.BoxStyle.Width/2}})...)

	da.StrokeLines(b.MedStyle, da.ClipLinesX([]plot.Point{
		{medx, y - b.Width/2}, {medx, y + b.Width/2},
	})...)

	minx, maxx := trX(min), trX(max)
	da.StrokeLines(b.WhiskerStyle, da.ClipLinesX([]plot.Point{{q3x, y}, {maxx, y}},
		[]plot.Point{{maxx, y - b.CapWidth/2}, {maxx, y + b.CapWidth/2}},
		[]plot.Point{{q1x, y}, {minx, y}},
		[]plot.Point{{minx, y - b.CapWidth/2}, {minx, y + b.CapWidth/2}})...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, plot.Point{trX(b.Y(i)), y})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b HorizBoxPlot) DataRange() (xmin, xmax, ymin, ymax float64) {
	ymin, ymax = b.X, b.X
	xmin, xmax = yDataRange(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b HorizBoxPlot) GlyphBoxes(p *plot.Plot) (boxes []plot.GlyphBox) {
	_, _, med, _, _, pts := b.Statistics()
	boxes = append(boxes, plot.GlyphBox{
		X: p.X.Norm(med),
		Y: p.Y.Norm(b.X),
		Rect: plot.Rect{
			Min:  plot.Point{Y: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: plot.Point{Y: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.Y(i)), p.Y.Norm(b.X)
		box := plot.GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
		boxes = append(boxes, box)
	}
	return
}
