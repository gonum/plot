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
	return &BoxPlot{
		Yer:          ys,
		X:            x,
		Width:        width,
		BoxStyle:     DefaultLineStyle,
		WhiskerStyle: DefaultLineStyle,
		CapWidth:     width / 2,
		GlyphStyle:   DefaultGlyphStyle,
	}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b *BoxPlot) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	x := trX(b.X)
	q1, med, q3, points := b.Statistics()
	q1y, medy, q3y := trY(q1), trY(med), trY(q3)
	box := da.ClipLinesY([]plot.Point{
		{x - b.Width/2, q1y}, {x - b.Width/2, q3y},
		{x + b.Width/2, q3y}, {x + b.Width/2, q1y},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1y}},
		[]plot.Point{{x - b.Width/2, medy}, {x + b.Width/2, medy}})
	da.StrokeLines(b.BoxStyle, box...)

	min, max := q1, q3
	if filtered := filteredIndices(b.Yer, points); len(filtered) > 0 {
		min = b.Y(filtered[0])
		max = b.Y(filtered[len(filtered)-1])
	}
	miny, maxy := trY(min), trY(max)
	whisk := da.ClipLinesY([]plot.Point{{x, q3y}, {x, maxy}},
		[]plot.Point{{x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy}},
		[]plot.Point{{x, q1y}, {x, miny}},
		[]plot.Point{{x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

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
	_, med, _, pts := b.Statistics()
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

// Statistics returns the `boxplot' statistics: the
// first quartile, the median, the third quartile,
// and a slice of indices to be drawn as separate
// points. This latter slice is computed as
// recommended by John Tukey in his book
// Exploratory Data Analysis: all values that are 1.5x
// the inter-quartile range before the first quartile
// and 1.5x the inter-quartile range after the third
// quartile.
func (b *BoxPlot) Statistics() (q1, med, q3 float64, points []int) {
	sorted := sortedIndices(b)
	q1 = percentile(b, sorted, 0.25)
	med = median(b, sorted)
	q3 = percentile(b, sorted, 0.75)
	points = tukeyPoints(b, sorted)
	return
}

// median returns the median Y value given a sorted
// slice of indices.
func median(ys Yer, sorted []int) float64 {
	med := ys.Y(sorted[len(sorted)/2])
	if len(sorted)%2 == 0 {
		med += ys.Y(sorted[len(sorted)/2-1])
		med /= 2
	}
	return med
}

// percentile returns the given percentile.
// According to Wikipedia, this technique is
// an alternative technique recommended
// by National Institute of Standards and
// Technology (NIST), and is used by MS
// Excel 2007.
func percentile(ys Yer, sorted []int, p float64) float64 {
	n := p*float64(len(sorted)-1) + 1
	k := math.Floor(n)
	d := n - k
	if n <= 1 {
		return ys.Y(sorted[0])
	} else if n >= float64(len(sorted)) {
		return ys.Y(sorted[len(sorted)-1])
	}
	yk := ys.Y(sorted[int(k)])
	yk1 := ys.Y(sorted[int(k)-1])
	return yk1 + d*(yk-yk1)
}

// sortedIndices returns a slice of the indices sorted in
// ascending order of their corresponding Y value.
func sortedIndices(ys Yer) []int {
	data := make([]int, ys.Len())
	for i := range data {
		data[i] = i
	}
	sort.Sort(ySorter{ys, data})
	return data
}

// tukeyPoints returns indices of values that are more than
// 1½ of the inter-quartile range beyond the 1st and 3rd
// quartile. According to John Tukey (Exploratory Data Analysis),
// these values are reasonable to draw separately as points.
func tukeyPoints(ys Yer, sorted []int) (pts []int) {
	q1 := percentile(ys, sorted, 0.25)
	q3 := percentile(ys, sorted, 0.75)
	min := q1 - 1.5*(q3-q1)
	max := q3 + 1.5*(q3-q1)
	for _, i := range sorted {
		if y := ys.Y(i); y > max || y < min {
			pts = append(pts, i)
		}
	}
	return
}

// filteredIndices returns a slice of the indices sorted in
// ascending order of their corresponding Y value, and
// excluding all indices in outList.
func filteredIndices(ys Yer, outList []int) (data []int) {
	out := make([]bool, ys.Len())
	for _, o := range outList {
		out[o] = true
	}
	for i := 0; i < ys.Len(); i++ {
		if !out[i] {
			data = append(data, i)
		}
	}
	sort.Sort(ySorter{ys, data})
	return data
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
	q1, med, q3, points := b.Statistics()
	q1x, medx, q3x := trX(q1), trX(med), trX(q3)
	box := da.ClipLinesX([]plot.Point{
		{q1x, y - b.Width/2}, {q3x, y - b.Width/2},
		{q3x, y + b.Width/2}, {q1x, y + b.Width/2},
		{q1x, y - b.Width/2 - b.BoxStyle.Width/2}},
		[]plot.Point{{medx, y - b.Width/2}, {medx, y + b.Width/2}})
	da.StrokeLines(b.BoxStyle, box...)

	min, max := q1, q3
	if filtered := filteredIndices(b.Yer, points); len(filtered) > 0 {
		min = b.Y(filtered[0])
		max = b.Y(filtered[len(filtered)-1])
	}
	minx, maxx := trX(min), trX(max)
	whisk := da.ClipLinesX([]plot.Point{{q3x, y}, {maxx, y}},
		[]plot.Point{{maxx, y - b.CapWidth/2}, {maxx, y + b.CapWidth/2}},
		[]plot.Point{{q1x, y}, {minx, y}},
		[]plot.Point{{minx, y - b.CapWidth/2}, {minx, y + b.CapWidth/2}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

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
	_, med, _, pts := b.Statistics()
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
