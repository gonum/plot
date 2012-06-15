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

// A Yer wraps methods for getting a set of Y data values.
type Yer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// Y returns a Y value
	Y(int) float64
}

// Line implements the Data interface, drawing a line
// for the Plot method.
type Line struct {
	XYer
	LineStyle
}

// Plot implements the Plot method of the Data interface,
// drawing a line that connects each point in the Line.
func (l Line) Plot(da DrawArea, p *Plot) {
	line := make([]Point, l.Len())
	for i := range line {
		line[i].X = da.X(p.X.Norm(l.X(i)))
		line[i].Y = da.Y(p.Y.Norm(l.Y(i)))
	}
	da.StrokeLines(l.LineStyle, da.ClipLinesXY(line)...)
}

// Extents implemnets the Extents function of the
// Data interface.
func (s Line) Extents() (xmin, ymin, xmax, ymax float64) {
	return xyExtents(s.XYer)
}


// Scatter implements the Data interface, drawing
// glyphs at each of the given points.
type Scatter struct {
	XYer
	GlyphStyle
}

// Plot implements the Plot method of the Data interface,
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

// Extents implemnets the Extents function of the
// Data interface.
func (s Scatter) Extents() (xmin, ymin, xmax, ymax float64) {
	return xyExtents(s.XYer)
}

// xyData wraps an XYer with an Extents method.
type xyData struct {
	XYer
}

// xyExtents returns the minimum and maximum x
// and y values of all points from the XYer.
func xyExtents(xy XYer) (xmin, ymin, xmax, ymax float64) {
	xmin = math.Inf(1)
	ymin = xmin
	xmax = math.Inf(-1)
	ymax = xmax
	for i := 0; i < xy.Len(); i++ {
		x, y := xy.X(i), xy.Y(i)
		xmin = math.Min(xmin, x)
		xmax = math.Max(xmax, x)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
	}
	return
}

// Box implements the Data interface, drawing a boxplot.
type Box struct {
	Yer

	// X is the X value, in data coordinates, at which
	// to draw the box.
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

	// Med, Q1, and Q3 are the median, first, and third
	// quartiles respectively.
	Med, Q1, Q3 float64

	// Points is a slice containing the indices Y values
	// that should be drawn separately as points.
	Points []int
}

// NewBox returns a Data which draws a box plot
// of the given y values at the given x value.
func NewBox(w vg.Length, x float64, ys Yer) *Box {
	sorted := sortedIndices(ys)
	return &Box{
		Yer: ys,
		X: x,
		Width: w,
		BoxStyle: DefaultLineStyle,
		WhiskerStyle: DefaultLineStyle,
		CapWidth: w / 2,
		GlyphStyle: DefaultGlyphStyle,
		Med: median(ys, sorted),
		Q1: percentile(ys, sorted, 0.25),
		Q3: percentile(ys, sorted, 0.75),
		Points: tukeyPoints(ys, sortedIndices(ys)),
	}
}

// Plot implements the Plot function of the Data interface,
// drawing a boxplot.
func (b *Box) Plot(da DrawArea, p *Plot) {
	x := da.X(p.X.Norm(b.X))
	q1y := da.Y(p.Y.Norm(b.Q1))
	q3y := da.Y(p.Y.Norm(b.Q3))
	medy := da.Y(p.Y.Norm(b.Med))
	box := da.ClipLinesY([]Point{
		{ x - b.Width/2, q1y }, { x - b.Width/2, q3y },
		{ x + b.Width/2, q3y }, { x + b.Width/2, q1y },
		{ x - b.Width/2 - b.BoxStyle.Width/2, q1y } },
		[]Point{ { x - b.Width/2, medy }, { x + b.Width/2, medy } })
	da.StrokeLines(b.BoxStyle, box...)

	min, max := b.Q1, b.Q3
	if filtered := filteredIndices(b.Yer, b.Points); len(filtered) > 0 {
		min = b.Y(filtered[0])
		max = b.Y(filtered[len(filtered)-1])
	}
	miny := da.Y(p.Y.Norm(min))
	maxy := da.Y(p.Y.Norm(max))
	whisk := da.ClipLinesY([]Point{{x, q3y}, {x, maxy} },
		[]Point{ {x - b.CapWidth/2, maxy}, {x + b.CapWidth/2, maxy} },
		[]Point{ {x, q1y}, {x, miny} },
		[]Point{ {x - b.CapWidth/2, miny}, {x + b.CapWidth/2, miny} })
	da.StrokeLines(b.WhiskerStyle, whisk...)

	for _, i := range b.Points {
		da.DrawGlyph(b.GlyphStyle,  Point{x, da.Y(p.Y.Norm(b.Y(i)))})
	}
}

// Extents implements the Extents function of the Data
// interface.
func (b *Box) Extents() (xmin, ymin, xmax, ymax float64) {
	xmin = b.X
	ymin = xmin
	xmax = b.X
	ymax = xmax
	for i := 0; i < b.Len(); i++ {
		y := b.Y(i)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
	}
	return
}

// GlyphBoxes returns a slice of GlyphBoxes for the
// points and for the median line of the boxplot.
func (b *Box) GlyphBoxes(p *Plot) (boxes []GlyphBox) {
	x := p.X.Norm(b.X)
	boxes = append(boxes, GlyphBox {
		X: x,
		Y: p.Y.Norm(b.Med),
		Rect: Rect{
			Min: Point{ X: -(b.Width/2 + b.BoxStyle.Width/2)},
			Size: Point{ X: b.Width + b.BoxStyle.Width },
		},
	})

	r := b.GlyphStyle.Radius
	rect := Rect{ Point{-r, -r}, Point{r*2, r*2} }
	for _, i := range b.Points {
		box := GlyphBox{
			X:    x,
			Y:    p.Y.Norm(b.Y(i)),
			Rect: rect,
		}
		boxes = append(boxes, box)
	}
	return
}

// tukeyPoints returns values that are more than ½ of the
// inter-quartile range beyond the 1st and 3rd quartile.
// According to John Tukey, these values are reasonable
// to draw separately as points.
func tukeyPoints(ys Yer, sorted []int) (pts []int) {
	q1 := percentile(ys, sorted, 0.25)
	q3 := percentile(ys, sorted, 0.75)
	min := q1 - 1.5*(q3 - q1)
	max := q3 + 1.5*(q3 - q1)
	for _, i := range sorted {
		if y := ys.Y(i); y > max || y < min {
			pts = append(pts, i)
		}
	}
	return
}

// median returns the median Y value given a sorted
// slice of indices.
func median(ys Yer, sorted []int) float64 {
	med := ys.Y(sorted[len(sorted)/2])
	if len(sorted) % 2 == 0 {
		med += ys.Y(sorted[len(sorted)/2 - 1])
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
	return yk1 + d * (yk - yk1)
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

// Points is a slice of X, Y pairs, implementing the
// XYer interface.
type Points []struct{ X, Y float64 }

// Len returns the number of points.
func (p Points) Len() int {
	return len(p)
}

// X returns the ith X value.
func (p Points) X(i int) float64 {
	return p[i].X
}

// Y returns the ith Y value.
func (p Points) Y(i int) float64 {
	return p[i].Y
}

// Values is a slice of values, implementing the Yer
// interface.
type Values []float64

// Len returns the number of values.
func (v Values) Len() int {
	return len(v)
}

// Y returns the ith Y value.
func (v Values) Y(i int) float64 {
	return v[i]
}