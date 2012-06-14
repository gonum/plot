package plt

import (
	"code.google.com/p/plotinum/vg"
	"image/color"
	"math"
)

var (
	DefaultLineStyle = LineStyle{
		Width: vg.Points(0.75),
		Color: color.Black,
	}

	DefaultGlyphStyle = GlyphStyle{
		Radius: vg.Points(2),
		Color:  color.Black,
	}
)

// Data is an interface that wraps all of the methods required
// to add data elements to a plot.
type Data interface {
	// Plot draws the data to the given DrawArea using
	// the axes from the given Plot.
	Plot(DrawArea, *Plot)

	// Extents returns the minimum and maximum
	// values of the data.
	Extents() (xmin, ymin, xmax, ymax float64)
}

// An XYer wraps methods for a type that can return
// multiple X and Y data values.
type XYer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// X returns an X value
	X(int) float64

	// Y returns a Y value
	Y(int) float64
}

// lineData implements the Data interface, drawing a line
// for the Plot method.
type lineData struct {
	xyData
	LineStyle
}

// MakeLine returns a Data that plots as a line in the given
// style connecting the given points.
func MakeLine(sty LineStyle, pts XYer) Data {
	return lineData{xyData: xyData{pts}, LineStyle: sty}
}

func (l lineData) Plot(da DrawArea, plt *Plot) {
	pts := make([]Point, l.pts.Len())
	for i := range pts {
		pts[i].X = da.X(plt.X.Norm(l.pts.X(i)))
		pts[i].Y = da.Y(plt.Y.Norm(l.pts.Y(i)))
	}
	da.StrokeClippedLine(l.LineStyle, pts...)
}

// scatterData implements the Data interface, drawing
// glyphs at each of the given points.
type scatterData struct {
	xyData
	GlyphStyle
}

// MakeScatter returns a Data interface, drawing the given
// points as glyphs for the Plot method.
func MakeScatter(sty GlyphStyle, pts XYer) Data {
	return scatterData{xyData: xyData{pts}, GlyphStyle: sty}
}

func (s scatterData) Plot(da DrawArea, plt *Plot) {
	for i := 0; i < s.pts.Len(); i++ {
		x, y := da.X(plt.X.Norm(s.pts.X(i))), da.Y(plt.Y.Norm(s.pts.Y(i)))
		da.DrawGlyph(s.GlyphStyle, Point{x, y})
	}
}

func (s scatterData) GlyphBoxes(plt *Plot) (boxes []GlyphBox) {
	r := Rect{ Point{ -s.Radius, -s.Radius }, Point{ s.Radius*2, s.Radius*2 } }
	for i := 0; i < s.pts.Len(); i++ {
		box := GlyphBox{
			X: plt.X.Norm(s.pts.X(i)),
			Y: plt.Y.Norm(s.pts.Y(i)),
			Rect: r,
		}
		boxes = append(boxes, box)
	}
	return
}

// xyData wraps an XYer with an Extents method.
type xyData struct {
	pts XYer
}

// extents returns the minimum and maximum x
// and y values of all points from the XYer.
func (xy xyData) Extents() (xmin, ymin, xmax, ymax float64) {
	xmin = math.Inf(1)
	ymin = xmin
	xmax = math.Inf(-1)
	ymax = xmax
	for i := 0; i < xy.pts.Len(); i++ {
		x, y := xy.pts.X(i), xy.pts.Y(i)
		xmin = math.Min(xmin, x)
		xmax = math.Max(xmax, x)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
	}
	return
}

// DataPoints is a slice of X, Y pairs, implementing the
// XYer interface.
type DataPoints []struct{ X, Y float64 }

func (p DataPoints) Len() int {
	return len(p)
}

func (p DataPoints) X(i int) float64 {
	return p[i].X
}

func (p DataPoints) Y(i int) float64 {
	return p[i].Y
}