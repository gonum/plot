package plt

import (
	"math"
	"image/color"
	"code.google.com/p/plotinum/vg"
)

var (
	// DefaultLineStyle is the LineStyle that is used if none
	// is specified otherwise.
	DefaultLineStyle = LineStyle{
		Width: vg.Points(0.75),
		Color: color.Black,
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
	points []Point
	LineStyle
}

// MakeLine returns a Data that plots as a line in the given
// style connecting the given points.
func MakeLine(sty LineStyle, pts ...Point) Data {
	ptscopy := make([]Point, len(pts))
	copy(ptscopy, pts)
	return lineData{ points: ptscopy, LineStyle: sty }
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

func (l lineData) extents() (xmin, ymin, xmax, ymax float64) {
	xmin = math.Inf(1)
	ymin = xmin
	xmax = math.Inf(-1)
	ymax = xmax
	for _, pt := range l.points {
		xmin = math.Min(xmin, pt.X)
		xmax = math.Max(xmax, pt.X)
		ymin = math.Min(ymin, pt.Y)
		ymax = math.Max(ymax, pt.Y)
	}
	return
}

// Point is a structure that represents a point in the
// data coordinate system.
type Point struct {
	X, Y float64
}