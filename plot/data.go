// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plot

import (
	"code.google.com/p/plotinum/vg"
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
	line := make([]Point, l.Len())
	for i := range line {
		line[i].X = da.X(p.X.Norm(l.X(i)))
		line[i].Y = da.Y(p.Y.Norm(l.Y(i)))
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
	for i := 0; i < s.Len(); i++ {
		x, y := da.X(p.X.Norm(s.X(i))), da.Y(p.Y.Norm(s.Y(i)))
		da.DrawGlyph(s.GlyphStyle, Point{x, y})
	}
}

// GlyphBoxes returns a slice of GlyphBoxes, one for
// each of the glyphs in the Scatter.
func (s Scatter) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	r := Rect{
		Point{-s.Radius, -s.Radius},
		Point{s.Radius * 2, s.Radius * 2},
	}
	for i := 0; i < s.Len(); i++ {
		box := GlyphBox{
			X:    p.X.Norm(s.X(i)),
			Y:    p.Y.Norm(s.Y(i)),
			Rect: r,
		}
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
	// XYLabeler has a set of labels located in data coordinates.
	XYLabeler

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

// MakeLabels returns a Labels using the default TextStyle,
// with the labels left-aligned above the corresponding
// X, Y point.
func MakeLabels(ls XYLabeler) (Labels, error) {
	labelFont, err := vg.MakeFont(defaultFont, vg.Points(10))
	if err != nil {
		return Labels{}, err
	}
	return Labels{
		XYLabeler: ls,
		TextStyle: TextStyle{Font: labelFont},
	}, nil
}

// Plot implements the Plotter interface for Labels.
func (l Labels) Plot(da DrawArea, p *Plot) {
	for i := 0; i < l.Len(); i++ {
		x, y := da.X(p.X.Norm(l.X(i))), da.Y(p.Y.Norm(l.Y(i)))
		if da.Contains(Point{x, y}) {
			da.FillText(l.TextStyle, x+l.XOffs, y+l.YOffs, l.XAlign, l.YAlign, l.Label(i))
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
		w := l.TextStyle.Width(l.Label(i))
		h := l.TextStyle.Height(l.Label(i))
		rect := Rect{
			Min: Point{
				w*vg.Length(l.XAlign) + l.XOffs,
				h*vg.Length(l.YAlign) + l.YOffs},
			Size: Point{w, h},
		}
		box := GlyphBox{
			X:    p.X.Norm(l.X(i)),
			Y:    p.Y.Norm(l.Y(i)),
			Rect: rect,
		}
		boxes = append(boxes, box)
	}
	return
}

// Box implements the Plotter interface, drawing a boxplot.
type Box struct {
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

// NewBox returns new Box representing a distribution
// of values.   The width parameter is the width of the
// box. The box surrounds the center of the range of
// data values, the middle line is the median, the points
// are as described in the Statistics method, and the
// whiskers extend to the extremes of all data that are
// not drawn as separate points.
func NewBox(width vg.Length, x float64, ys Yer) *Box {
	return &Box{
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
func (b *Box) Plot(da DrawArea, p *Plot) {
	x := da.X(p.X.Norm(b.X))
	q1, med, q3, points := b.Statistics()
	q1y := da.Y(p.Y.Norm(q1))
	q3y := da.Y(p.Y.Norm(q3))
	medy := da.Y(p.Y.Norm(med))
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
	miny := da.Y(p.Y.Norm(min))
	maxy := da.Y(p.Y.Norm(max))
	whisk := da.ClipLinesY([]Point{{x, q3y}, {x, maxy}},
		[]Point{{x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy}},
		[]Point{{x, q1y}, {x, miny}},
		[]Point{{x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, Point{x, da.Y(p.Y.Norm(b.Y(i)))})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b HorizBox) DataRange() (xmin, xmax, ymin, ymax float64) {
	ymin, ymax = b.X, b.X
	xmin, xmax = yDataRange(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b *Box) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	_, med, _, pts := b.Statistics()
	boxes = append(boxes, GlyphBox{
		X: p.X.Norm(b.X),
		Y: p.Y.Norm(med),
		Rect: Rect{
			Min:  Point{X: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: Point{X: b.Width + b.BoxStyle.Width},
		},
	})

	r := b.GlyphStyle.Radius
	rect := Rect{Point{-r, -r}, Point{r * 2, r * 2}}
	for _, i := range pts {
		boxes = append(boxes, GlyphBox{
			X:    p.X.Norm(b.X),
			Y:    p.Y.Norm(b.Y(i)),
			Rect: rect,
		})
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
func (b *Box) Statistics() (q1, med, q3 float64, points []int) {
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

// tukeyPoints returns values that are more than 1½ of the
// inter-quartile range beyond the 1st and 3rd quartile.
// According to John Tukey, these values are reasonable
// to draw separately as points.
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
type HorizBox struct {
	*Box
}

// NewHorizBox returns a HorizBox.  This is the
// same as NewBox except that the box draws
// horizontally instead of vertically.
func MakeHorizBox(width vg.Length, y float64, vals Yer) HorizBox {
	return HorizBox{NewBox(width, y, vals)}
}

// Plot implements the Plot function of the Plotter interface,
// drawing a boxplot.
func (b HorizBox) Plot(da DrawArea, p *Plot) {
	y := da.Y(p.Y.Norm(b.X))
	q1, med, q3, points := b.Statistics()
	q1x := da.X(p.X.Norm(q1))
	q3x := da.X(p.X.Norm(q3))
	medx := da.X(p.X.Norm(med))
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
	minx := da.X(p.X.Norm(min))
	maxx := da.X(p.X.Norm(max))
	whisk := da.ClipLinesX([]Point{{q3x, y}, {maxx, y}},
		[]Point{{maxx, y - b.CapWidth/2}, {maxx, y + b.CapWidth/2}},
		[]Point{{q1x, y}, {minx, y}},
		[]Point{{minx, y - b.CapWidth/2}, {minx, y + b.CapWidth/2}})
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range points {
		da.DrawGlyph(b.GlyphStyle, Point{da.X(p.X.Norm(b.Y(i))), y})
	}
}

// DataRange returns the minimum and maximum X and Y values
func (b *Box) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = b.X, b.X
	ymin, ymax = yDataRange(b)
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b HorizBox) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	_, med, _, pts := b.Statistics()
	boxes = append(boxes, GlyphBox{
		X: p.X.Norm(med),
		Y: p.Y.Norm(b.X),
		Rect: Rect{
			Min:  Point{Y: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: Point{Y: b.Width + b.BoxStyle.Width},
		},
	})

	r := b.GlyphStyle.Radius
	rect := Rect{Point{-r, -r}, Point{r * 2, r * 2}}
	for _, i := range pts {
		boxes = append(boxes, GlyphBox{
			X:    p.X.Norm(b.Y(i)),
			Y:    p.Y.Norm(b.X),
			Rect: rect,
		})
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

// X returns the ith X value.
func (p XYs) X(i int) float64 {
	return p[i].X
}

// Y returns the ith Y value.
func (p XYs) Y(i int) float64 {
	return p[i].Y
}

// An XYLabeler wraps the XYer methods along with
// a Label method that returns a label for the
// corresponding X,Y point.
type XYLabeler interface {
	XYer
	Label(int) string
}

// XYLabels implements the XYLabeler interface.
type XYLabels []struct {
	X, Y  float64
	Label string
}

// Len returns the number of points.
func (p XYLabels) Len() int {
	return len(p)
}

// X returns the ith X value.
func (p XYLabels) X(i int) float64 {
	return p[i].X
}

// Y returns the ith Y value.
func (p XYLabels) Y(i int) float64 {
	return p[i].Y
}

// Label returns the ith Label value.
func (p XYLabels) Label(i int) string {
	return p[i].Label
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
func (v Ys) Len() int {
	return len(v)
}

// Y returns the ith Y value.
func (v Ys) Y(i int) float64 {
	return v[i]
}
