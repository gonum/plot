// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plot

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
	"sort"
)

var (
	// DefaultLineStyle is a reasonable default LineStyle
	// for drawing most lines in a plot.
	DefaultLineStyle = LineStyle{
		Width: vg.Points(0.5),
		Color: color.Black,
	}

	// DefaultGlyhpStyle is a reasonable default GlyphStyle
	// for drawing points on a plot.
	DefaultGlyphStyle = GlyphStyle{
		Radius: vg.Points(2),
		Color:  color.Black,
	}
)

// Line implements the Plotter interface, drawing a line
// for the Plot method.
type Line struct {
	XYer
	LineStyle
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (l Line) Plot(da DrawArea, p *Plot) {
	trX, trY := p.Transforms(&da)
	line := make([]Point, l.Len())
	for i := range line {
		line[i].X = trX(l.X(i))
		line[i].Y = trY(l.Y(i))
	}
	da.StrokeLines(l.LineStyle, da.ClipLinesXY(line)...)
}

// DataRange returns the minimum and maximum X and Y values
func (l Line) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = xDataRange(l)
	ymin, ymax = yDataRange(l)
	return
}

// Scatter implements the Plotter interface, drawing
// glyphs at each of the given points.
type Scatter struct {
	XYer
	GlyphStyle
}

// Plot implements the Plot method of the Plotter interface,
// drawing a glyph for each point in the Scatter.
func (s Scatter) Plot(da DrawArea, p *Plot) {
	trX, trY := p.Transforms(&da)
	for i := 0; i < s.Len(); i++ {
		da.DrawGlyph(s.GlyphStyle, Point{trX(s.X(i)), trY(s.Y(i))})
	}
}

// GlyphBoxes returns a slice of GlyphBoxes, one for
// each of the glyphs in the Scatter.
func (s Scatter) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	for i := 0; i < s.Len(); i++ {
		x, y := p.X.Norm(s.X(i)), p.Y.Norm(s.Y(i))
		box := GlyphBox{X: x, Y: y, Rect: s.Rect()}
		boxes = append(boxes, box)
	}
	return
}

// DataRange returns the minimum and maximum X and Y values
func (s Scatter) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = xDataRange(s)
	ymin, ymax = yDataRange(s)
	return
}

// Labels implements the Plotter interface, drawing
// a set of labels on the plot.
type Labels struct {
	// XYLabeller has a set of labels located in data coordinates.
	XYLabeller

	// TextStyle gives the style of the labels.
	TextStyle

	// XAlign and YAlign are multiplied by the width
	// and height of each label respectively and the
	// added to the final location.  E.g., XAlign=-0.5
	// and YAlign=-0.5 centers the label at the given
	// X, Y location, and XAlign=0, YAlign=0 aligns
	// the text to the left of the point, and XAlign=-1,
	// YAlign=0 aligns the text to the right of the point.
	XAlign, YAlign float64

	// XOffs and YOffs are added directly to the final
	// label X and Y location respectively.
	XOffs, YOffs vg.Length
}

// Labels returns a Labels using the default TextStyle,
// with the labels left-aligned above the corresponding
// X, Y point.
func MakeLabels(ls XYLabeller) (Labels, error) {
	labelFont, err := vg.MakeFont(defaultFont, vg.Points(10))
	if err != nil {
		return Labels{}, err
	}
	return Labels{XYLabeller: ls, TextStyle: TextStyle{Font: labelFont}}, nil
}

// Plot implements the Plotter interface for Labels.
func (l Labels) Plot(da DrawArea, p *Plot) {
	trX, trY := p.Transforms(&da)
	for i := 0; i < l.Len(); i++ {
		x, y := trX(l.X(i))+l.XOffs, trY(l.Y(i))+l.YOffs
		if da.Contains(Point{x, y}) {
			da.FillText(l.TextStyle, x, y, l.XAlign, l.YAlign, l.Label(i))
		}
	}
}

// DataRange returns the minimum and maximum X and Y values
func (l Labels) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = xDataRange(l)
	ymin, ymax = yDataRange(l)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes, one for
// each of the labels.
func (l Labels) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	for i := 0; i < l.Len(); i++ {
		x, y := p.X.Norm(l.X(i)), p.Y.Norm(l.Y(i))
		txt := l.Label(i)
		rect := l.Rect(txt)
		rect.Min.X += l.Width(txt)*vg.Length(l.XAlign) + l.XOffs
		rect.Min.Y += l.Height(txt)*vg.Length(l.YAlign) + l.YOffs
		box := GlyphBox{X: x, Y: y, Rect: rect}
		boxes = append(boxes, box)
	}
	return
}

// ErrorBars implement the Plotter interface, drawing
// either vertical, horizontal, or both vertical and
// horizontal error bars.
type ErrorBars struct {
	// XYer gives the center location of the error
	// bar.  If the XYer also implements the XErrorer
	// interface then a horizontal error bar is plotted.
	// If the XYer also implements YErrorer then
	// a vertical error bar is drawn.
	XYer

	// LineStyle is the style of the errorbar and cap
	// lines.
	LineStyle

	// CapWidth is the width of the error bar caps.
	CapWidth vg.Length
}

// MakeErrorBars returns ErrorBars using the default
// line style and cap width.  If the XYer doesn't implement
// either XErrorer or YErrorer then a non-nil error
// is returned, and while the resulting ErrorBars can be
// successfully added to a plot, they will not draw
// anything.
func MakeErrorBars(xy XYer) (bars ErrorBars, err error) {
	_, xerr := xy.(XErrorer)
	_, yerr := xy.(YErrorer)
	if !xerr && !yerr {
		err = fmt.Errorf("Type %T doesn't implement XErrorer or YErrorer")
	}
	bars.XYer = xy
	bars.CapWidth = vg.Points(6)
	bars.LineStyle = DefaultLineStyle
	return
}

// Plot implements the Plotter interface, drawing
// error bars in either the X or Y directions or both.
func (e ErrorBars) Plot(da DrawArea, p *Plot) {
	e.plotVerticalBars(&da, p)
	e.plotHorizontalBars(&da, p)
}

// plotVerticalBars plots the vertical error bars
// if this ErrorBars implements the YErrorer interface.
func (e ErrorBars) plotVerticalBars(da *DrawArea, p *Plot) {
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	trX, trY := p.Transforms(da)
	for i := 0; i < e.Len(); i++ {
		errlow, errhigh := yerr.YError(i)
		y := e.Y(i)
		min, max := trY(y+errlow), trY(y+errhigh)
		x := trX(e.X(i))
		da.StrokeLines(e.LineStyle, da.ClipLinesXY([]Point{{x, min}, {x, max}})...)
		e.plotVerticalCap(da, Point{x, min})
		e.plotVerticalCap(da, Point{x, max})
	}
}

// plotVerticalCap plots a horizontal line, centered
// at the given point, capping a vertical errorbar.
func (e ErrorBars) plotVerticalCap(da *DrawArea, pt Point) {
	w := e.CapWidth / 2
	if da.Contains(pt) {
		da.StrokeLine2(e.LineStyle, pt.X-w, pt.Y, pt.X+w, pt.Y)
	}
}

// plotHorizontalBars plots the horizontal error bars
// if this ErrorBars implements the XErrorer interface.
func (e ErrorBars) plotHorizontalBars(da *DrawArea, p *Plot) {
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	trX, trY := p.Transforms(da)
	for i := 0; i < e.Len(); i++ {
		errlow, errhigh := xerr.XError(i)
		x := e.X(i)
		min, max := trX(x+errlow), trX(x+errhigh)
		y := trY(e.Y(i))
		da.StrokeLines(e.LineStyle, da.ClipLinesXY([]Point{{min, y}, {max, y}})...)
		e.plotHorizontalCap(da, Point{min, y})
		e.plotHorizontalCap(da, Point{max, y})
	}
}

// plotHorizontalCap plots a vertical line, centered
// at the given point, capping a horizontal errorbar.
func (e ErrorBars) plotHorizontalCap(da *DrawArea, pt Point) {
	w := e.CapWidth / 2
	if da.Contains(pt) {
		da.StrokeLine2(e.LineStyle, pt.X, pt.Y-w, pt.X, pt.Y+w)
	}
}

// GlyphBoxes implements the GlyphBoxer interface,
// ensuring that the caps of the error bars are not
// clipped by the edge of the plot.
func (e ErrorBars) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	boxes = append(boxes, e.verticalGlyphBoxes(p)...)
	boxes = append(boxes, e.horizontalGlyphBoxes(p)...)
	return
}

// verticalGlyphBoxes returns the GlyphBoxes
// for the vertical error bar caps.
func (e ErrorBars) verticalGlyphBoxes(p *Plot) (boxes []GlyphBox) {
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	vertRect := Rect{Min: Point{X: -e.CapWidth / 2}, Size: Point{X: e.CapWidth}}
	for i := 0; i < e.Len(); i++ {
		x, y := e.X(i), e.Y(i)
		errlow, errhigh := yerr.YError(i)
		min, max := p.Y.Norm(y+errlow), p.Y.Norm(y+errhigh)
		boxes = append(boxes,
			GlyphBox{X: p.X.Norm(x), Y: min, Rect: vertRect},
			GlyphBox{X: p.X.Norm(x), Y: max, Rect: vertRect})
	}
	return
}

// horizontalGlyphBoxes returns the GlyphBoxes
// for the horizontal error bar caps.
func (e ErrorBars) horizontalGlyphBoxes(p *Plot) (boxes []GlyphBox) {
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	horzRect := Rect{Min: Point{Y: -e.CapWidth / 2}, Size: Point{Y: e.CapWidth}}
	for i := 0; i < e.Len(); i++ {
		x, y := e.X(i), e.Y(i)
		errlow, errhigh := xerr.XError(i)
		min, max := p.X.Norm(x+errlow), p.X.Norm(x+errhigh)
		boxes = append(boxes,
			GlyphBox{X: min, Y: p.Y.Norm(y), Rect: horzRect},
			GlyphBox{X: max, Y: p.Y.Norm(y), Rect: horzRect})
	}
	return

}

// DataRange implements the DataRanger interface,
// returning the minimum and maximum X and Y
// values of the error bars.
func (e ErrorBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = e.xDataRange()
	ymin, ymax = e.yDataRange()
	return
}

// xDataRange returns the range of x values
// for the error bars.
func (e ErrorBars) xDataRange() (xmin, xmax float64) {
	xmin, xmax = xDataRange(e)
	xerr, ok := e.XYer.(XErrorer)
	if !ok {
		return
	}
	for i := 0; i < e.Len(); i++ {
		x := e.X(i)
		errlow, errhigh := xerr.XError(i)
		xmin = math.Min(xmin, x+errlow)
		xmax = math.Max(xmax, x+errlow)
		xmin = math.Min(xmin, x+errhigh)
		xmax = math.Max(xmax, x+errhigh)
	}
	return
}

// yDataRange returns the range of y values
// for the error bars.
func (e ErrorBars) yDataRange() (ymin, ymax float64) {
	ymin, ymax = yDataRange(e)
	yerr, ok := e.XYer.(YErrorer)
	if !ok {
		return
	}
	for i := 0; i < e.Len(); i++ {
		y := e.Y(i)
		errlow, errhigh := yerr.YError(i)
		ymin = math.Min(ymin, y+errlow)
		ymax = math.Max(ymax, y+errlow)
		ymin = math.Min(ymin, y+errhigh)
		ymax = math.Max(ymax, y+errhigh)
	}
	return
}

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
	BoxStyle LineStyle

	// WhiskerStyle is the style used to draw the
	// whiskers.
	WhiskerStyle LineStyle

	// CapWidth is the width of the cap on the whiskers.
	CapWidth vg.Length

	// GlyphStyle is the style of the points.
	GlyphStyle GlyphStyle
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
func (b *BoxPlot) Plot(da DrawArea, p *Plot) {
	trX, trY := p.Transforms(&da)
	x := trX(b.X)
	q1, med, q3, points := b.Statistics()
	q1y, medy, q3y := trY(q1), trY(med), trY(q3)
	box := da.ClipLinesY([]Point{
		{x - b.Width/2, q1y}, {x - b.Width/2, q3y},
		{x + b.Width/2, q3y}, {x + b.Width/2, q1y},
		{x - b.Width/2 - b.BoxStyle.Width/2, q1y}},
		[]Point{{x - b.Width/2, medy}, {x + b.Width/2, medy}})
	da.StrokeLines(b.BoxStyle, box...)

	min, max := q1, q3
	if filtered := filteredIndices(b.Yer, points); len(filtered) > 0 {
		min = b.Y(filtered[0])
		max = b.Y(filtered[len(filtered)-1])
	}
	miny, maxy := trY(min), trY(max)
	whisk := da.ClipLinesY([]Point{{x, q3y}, {x, maxy}},
		[]Point{{x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy}},
		[]Point{{x, q1y}, {x, miny}},
		[]Point{{x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, Point{x, trY(b.Y(i))})
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
func (b *BoxPlot) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	_, med, _, pts := b.Statistics()
	boxes = append(boxes, GlyphBox{
		X: p.X.Norm(b.X),
		Y: p.Y.Norm(med),
		Rect: Rect{
			Min:  Point{X: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: Point{X: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.X), p.Y.Norm(b.Y(i))
		box := GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
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
func (b HorizBoxPlot) Plot(da DrawArea, p *Plot) {
	trX, trY := p.Transforms(&da)
	y := trY(b.X)
	q1, med, q3, points := b.Statistics()
	q1x, medx, q3x := trX(q1), trX(med), trX(q3)
	box := da.ClipLinesX([]Point{
		{q1x, y - b.Width/2}, {q3x, y - b.Width/2},
		{q3x, y + b.Width/2}, {q1x, y + b.Width/2},
		{q1x, y - b.Width/2 - b.BoxStyle.Width/2}},
		[]Point{{medx, y - b.Width/2}, {medx, y + b.Width/2}})
	da.StrokeLines(b.BoxStyle, box...)

	min, max := q1, q3
	if filtered := filteredIndices(b.Yer, points); len(filtered) > 0 {
		min = b.Y(filtered[0])
		max = b.Y(filtered[len(filtered)-1])
	}
	minx, maxx := trX(min), trX(max)
	whisk := da.ClipLinesX([]Point{{q3x, y}, {maxx, y}},
		[]Point{{maxx, y - b.CapWidth/2}, {maxx, y + b.CapWidth/2}},
		[]Point{{q1x, y}, {minx, y}},
		[]Point{{minx, y - b.CapWidth/2}, {minx, y + b.CapWidth/2}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, Point{trX(b.Y(i)), y})
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
func (b HorizBoxPlot) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	_, med, _, pts := b.Statistics()
	boxes = append(boxes, GlyphBox{
		X: p.X.Norm(med),
		Y: p.Y.Norm(b.X),
		Rect: Rect{
			Min:  Point{Y: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: Point{Y: b.Width + b.BoxStyle.Width},
		},
	})
	for _, i := range pts {
		x, y := p.X.Norm(b.Y(i)), p.Y.Norm(b.X)
		box := GlyphBox{X: x, Y: y, Rect: b.GlyphStyle.Rect()}
		boxes = append(boxes, box)
	}
	return
}

// ySorted implements sort.Interface, sorting a slice
// of indices for the given Yer.
type ySorter struct {
	Yer
	inds []int
}

// Len returns the number of indices.
func (y ySorter) Len() int {
	return len(y.inds)
}

// Less returns true if the Y value at index i
// is less than the Y value at index j.
func (y ySorter) Less(i, j int) bool {
	return y.Y(y.inds[i]) < y.Y(y.inds[j])
}

// Swap swaps the ith and jth indices.
func (y ySorter) Swap(i, j int) {
	y.inds[i], y.inds[j] = y.inds[j], y.inds[i]
}

// xDataRange returns the minimum and maximum x
// values of all points from the XYer.
func xDataRange(xys XYer) (xmin, xmax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	for i := 0; i < xys.Len(); i++ {
		x := xys.X(i)
		xmin = math.Min(xmin, x)
		xmax = math.Max(xmax, x)
	}
	return
}

// yDataRange returns the minimum and maximum x
// values of all points from the XYer.
func yDataRange(ys Yer) (ymin, ymax float64) {
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for i := 0; i < ys.Len(); i++ {
		y := ys.Y(i)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
	}
	return
}

// A Yer wraps methods for getting a set of Y data values.
type Yer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// Y returns a Y value
	Y(int) float64
}

// Ys is a slice of values, implementing the Yer
// interface.
type Ys []float64

// Len returns the number of values.
func (ys Ys) Len() int {
	return len(ys)
}

// Less returns true if the ith Y value is less than
// the jth Y value.
func (ys Ys) Less(i, j int) bool {
	return ys[i] < ys[j]
}

// Swap swaps the ith and jth values.
func (ys Ys) Swap(i, j int) {
	ys[i], ys[j] = ys[j], ys[i]
}

// Y returns the ith Y value.
func (ys Ys) Y(i int) float64 {
	return ys[i]
}

// An XYer wraps methods for getting a set of
// X and Y data values.
type XYer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// X returns an X value
	X(int) float64

	// Y returns a Y value
	Y(int) float64
}

// XYs is a slice of X, Y pairs, implementing the
// XYer interface.
type XYs []struct {
	X, Y float64
}

// Len returns the number of points.
func (p XYs) Len() int {
	return len(p)
}

// Less returns true if the ith X value is less than
// the jth X value.  This implements the Less
// method of sort.Interface, for sorting points by
// increasing X.
func (p XYs) Less(i, j int) bool {
	return p[i].X < p[j].X
}

// Swap swaps the ith and jth points.
func (p XYs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// X returns the ith X value.
func (p XYs) X(i int) float64 {
	return p[i].X
}

// Y returns the ith Y value.
func (p XYs) Y(i int) float64 {
	return p[i].Y
}

// XYLabeller wraps both XYer and Labeller.
type XYLabeller interface {
	// XYer returns the XY point that is being labelled.
	XYer

	// Label returns the ith label text.
	Label(int) string
}

// XErrorer wraps the XError method.
type XErrorer interface {
	// XError returns the low and high X errors.
	// Both values are added to the corresponding
	// X value to compute the range of error
	// of the X value of the point, so most likely
	// the low value will be negative.
	XError(int) (float64, float64)
}

// YErrorer wraps the YError method.
type YErrorer interface {
	// YError is the same as the XError method
	// of the XErrorer interface, however it
	// applies to the Y values of points instead
	// of the X values.
	YError(int) (float64, float64)
}

// XYLabelErrors implements the XYer, XYLabeller, XErrorer,
// and YErrorer interfaces.
type XYLabelErrors struct {
	XYs
	Labels  []string
	XErrors []struct{ Low, High float64 }
	YErrors []struct{ Low, High float64 }
}

// MakeYXYLabelErrors returns a new XYLabelErrors
// of the given length.
func MakeXYLabelErrors(l int) XYLabelErrors {
	return XYLabelErrors{
		XYs:     make(XYs, l),
		Labels:  make([]string, l),
		XErrors: make([]struct{ Low, High float64 }, l),
		YErrors: make([]struct{ Low, High float64 }, l),
	}
}

// Label implements the XYLabeller interface.
func (xy XYLabelErrors) Label(i int) string {
	return xy.Labels[i]
}

// XError implements the XErrorer interface.
func (xy XYLabelErrors) XError(i int) (float64, float64) {
	return xy.XErrors[i].Low, xy.XErrors[i].High
}

// YError implements the YErrorer interface.
func (xy XYLabelErrors) YError(i int) (float64, float64) {
	return xy.YErrors[i].Low, xy.YErrors[i].High
}
