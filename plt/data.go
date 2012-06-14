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
	// plot draws the data to the given drawArea using
	// the axes from the given Plot.
	plot(*drawArea, *Plot)

	// extents returns the minimum and maximum
	// values of the data.
	extents() (xmin, ymin, xmax, ymax float64)
}

// lineData implements the Data interface, drawing a line
// for the plot method.
type lineData struct {
	points
	LineStyle
}

// MakeLine returns a Data that plots as a line in the given
// style connecting the given points.
func MakeLine(sty LineStyle, pts ...Point) Data {
	ptscopy := make([]Point, len(pts))
	copy(ptscopy, pts)
	return lineData{points: ptscopy, LineStyle: sty}
}

func (l lineData) plot(da *drawArea, plt *Plot) {
	da.setLineStyle(l.LineStyle)
	pts := make([]point, len(l.points))
	for i, pt := range l.points {
		pts[i].x = plt.X.x(da, pt.X)
		pts[i].y = plt.Y.y(da, pt.Y)
	}
	strokeClippedLine(da, pts...)
}

// scatterData implements the Data interface, drawing
// glyphs at each of the given points.
type scatterData struct {
	points
	GlyphStyle
}

// MakeScatter returns a Data interface, drawing the given
// points as glyphs for the plot method.
func MakeScatter(sty GlyphStyle, pts ...Point) Data {
	ptscopy := make([]Point, len(pts))
	copy(ptscopy, pts)
	return scatterData{points: ptscopy, GlyphStyle: sty}
}

func (s scatterData) plot(da *drawArea, plt *Plot) {
	for _, pt := range s.points {
		p := point{plt.X.x(da, pt.X), plt.Y.y(da, pt.Y)}
		if da.contains(p) {
			drawGlyph(da, s.GlyphStyle, p)
		}
	}
}

// points is a slice of points.
type points []Point

// extents returns the minimum and maximum x
// and y values of all points.
func (ps points) extents() (xmin, ymin, xmax, ymax float64) {
	xmin = math.Inf(1)
	ymin = xmin
	xmax = math.Inf(-1)
	ymax = xmax
	for _, pt := range ps {
		xmin = math.Min(xmin, pt.X)
		xmax = math.Max(xmax, pt.X)
		ymin = math.Min(ymin, pt.Y)
		ymax = math.Max(ymax, pt.Y)
	}
	return
}

// Point is a point in the 2D data coordinate system.
type Point struct {
	X, Y float64
}
