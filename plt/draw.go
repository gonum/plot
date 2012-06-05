package plt

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
)

var (
	Black = color.RGBA{A: 255}
	White = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Red   = color.RGBA{R: 255, A: 255}
	Green = color.RGBA{G: 255, A: 255}
	Blue  = color.RGBA{B: 255, A: 255}
)

// A DrawArea is a vector graphics canvas along with
// an associated rect defining a section of the canvas
// to which drawing should take place.
type DrawArea struct {
	vg.Canvas
	font vg.Font
	rect
}

// center returns the center point of the area
func (da *DrawArea) center() point {
	return point{
		X: (da.Max().X-da.Min.X)/2 + da.Min.X,
		Y: (da.Max().Y-da.Min.Y)/2 + da.Min.Y,
	}
}

// x returns the value of x, given in the unit range,
// in the drawing coordinates of this draw area.
// A value of 0, for example, will return the minimum
// x value of the draw area and a value of 1 will
// return the maximum.
func (da *DrawArea) x(x float64) vg.Length {
	return vg.Length(x)*(da.Max().X-da.Min.X) + da.Min.X
}

// y returns the value of x, given in the unit range,
// in the drawing coordinates of this draw area.
// A value of 0, for example, will return the minimum
// y value of the draw area and a value of 1 will
// return the maximum.
func (da *DrawArea) y(y float64) vg.Length {
	return vg.Length(y)*(da.Max().Y-da.Min.Y) + da.Min.Y
}

// crop returns a new DrawArea corresponding to the receiver
// area with the given number of inches added to the minimum
// and maximum x and y values of the DrawArea's rect.
func (da *DrawArea) crop(minx, miny, maxx, maxy vg.Length) *DrawArea {
	minpt := point{
		X: da.rect.Min.X + minx,
		Y: da.rect.Min.Y + miny,
	}
	sz := point{
		X: da.Max().X + maxx - minpt.X,
		Y: da.Max().Y + maxy - minpt.Y,
	}
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		font:      da.font,
		rect:      rect{Min: minpt, Size: sz},
	}
}

// squishX returns a new DrawArea with a squished width such
// that any of the given set of glyphs will draw within the original
// draw area when they are mapped to the coordinate system
// of the returned DrawArea.
//
// The location of the glyphs that are given as a parameter are
// assumed to be on the unit interval, with 0 meaning the left-most
// side of the draw area and 1 meaning the right-most side.
func (da *DrawArea) squishX(boxes []glyphBox) *DrawArea {
	if len(boxes) == 0 {
		return da
	}

	var left, right vg.Length
	minx, maxx := vg.Length(math.Inf(1)), vg.Length(math.Inf(-1))
	for _, b := range boxes {
		if x := da.x(b.X) + b.rect.Min.X; x < minx {
			left = vg.Length(b.X)
			minx = x
		}
		if x := da.x(b.X) + b.Min.X + b.Size.X; x > maxx {
			right = vg.Length(b.X)
			maxx = x
		}
	}

	if minx >= da.Min.X {
		minx = da.Min.X
	}
	if maxx <= da.Max().X {
		maxx = da.Max().X
	}

	// where we want the left and right points to end up
	l := da.Min.X + (da.Min.X - minx)
	r := da.Max().X - (maxx - da.Max().X)
	n := (left*r - right*l) / (left - right)
	m := ((left-1)*r - right*l + l) / (left - right)
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		font:      da.font,
		rect: rect{
			Min:  point{X: n, Y: da.Min.Y},
			Size: point{X: m - n, Y: da.Size.Y},
		},
	}
}

// squishY returns a new DrawArea with a squished height such
// that any of the given set of glyphs will draw within the original
// draw area when they are mapped to the coordinate system
// of the returned DrawArea.
//
// The location of the glyphs that are given as a parameter are
// assumed to be on the unit interval, with 0 meaning the
// bottom-most side of the draw area and 1 meaning the
// top-most side.
func (da *DrawArea) squishY(boxes []glyphBox) *DrawArea {
	if len(boxes) == 0 {
		return da
	}

	var top, bot vg.Length
	miny, maxy := vg.Length(math.Inf(1)), vg.Length(math.Inf(-1))
	for _, b := range boxes {
		if y := da.y(b.Y) + b.rect.Min.Y; y < miny {
			bot = vg.Length(b.Y)
			miny = y
		}
		if y := da.y(b.Y) + b.Min.Y + b.Size.Y; y > maxy {
			top = vg.Length(b.Y)
			maxy = y
		}
	}

	if miny >= da.Min.Y {
		miny = da.Min.Y
	}
	if maxy <= da.Max().Y {
		maxy = da.Max().Y
	}

	// where we want the top and bottom points to end up
	b := da.Min.Y + (da.Min.Y - miny)
	t := da.Max().Y - (maxy - da.Max().Y)
	n := (bot*t - top*b) / (bot - top)
	m := ((bot-1)*t - top*b + b) / (bot - top)
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		font:      da.font,
		rect: rect{
			Min:  point{X: da.Min.X, Y: n},
			Size: point{X: da.Size.X, Y: m - n, },
		},
	}
}

// A glyphBox describes the location of a glyph
// and the offset/size of its bounding box.
type glyphBox struct {
	// The glyph location in normalized coordinates.
	X, Y float64
	// rect is the offset of the glyph's minimum drawing
	// point relative to the glyph location and its size.
	rect
}

// setTextStyle sets the current text style
func (da *DrawArea) setTextStyle(sty TextStyle) {
	da.SetColor(sty.Color)
	da.font = sty.Font
}

// TextStyle describes what text will look like.
type TextStyle struct {
	// Color is the text color.
	Color color.Color

	// Font is the font description.
	Font vg.Font
}

// MakeFont returns a font object.
// This function is merely included for convenience so that
// the user doesn't have to import the vg package.
func MakeFont(name string, size vg.Length) (vg.Font, error) {
	return vg.MakeFont(name, size)
}

// text fills the text to the drawing area.  The string is created
// using the printf-style format specification and the text is
// located at x + width*fx, y + height*fy, where width and height
// are the width and height of the rendered string.
func (da *DrawArea) text(x, y vg.Length, fx, fy float64, f string, v ...interface{}) {
	if da.font.Font() == nil {
		panic("Drawing text without a current font set")
	}
	str := fmt.Sprintf(f, v...)
	w := da.font.Width(str)
	h := da.font.Extents().Ascent
	da.FillText(da.font, x+w*vg.Length(fx), y+h*vg.Length(fy), str)
}

// setLineStyle sets the current line style
func (da *DrawArea) setLineStyle(sty LineStyle) {
	da.SetColor(sty.Color)
	da.SetLineWidth(sty.Width)
	var dashDots []vg.Length
	for _, dash := range sty.Dashes {
		dashDots = append(dashDots, dash)
	}
	da.SetLineDash(dashDots, sty.DashOffs)

}

// LineStyle describes what a line will look like.
type LineStyle struct {
	// Color is the color of the line.
	Color color.Color

	// Width is the width of the line in inches.
	Width vg.Length

	Dashes   []vg.Length
	DashOffs vg.Length
}

// line draws a line connecting the given points.
func (da *DrawArea) line(pts []point) {
	if len(pts) == 0 {
		return
	}

	var p vg.Path
	p.Move(pts[0].X, pts[0].Y)
	for _, pt := range pts {
		p.Line(pt.X, pt.Y)
	}
	da.Stroke(p)
}

// clippedLine draws a line that is clipped at the bounds
// the DrawArea.
func (da *DrawArea) clippedLine(pts []point) {
	// clip right
	lines0 := clip(isLeft, point{da.Max().X, da.Min.Y}, point{-1, 0}, pts)

	// clip bottom
	var lines1 [][]point
	for _, line := range lines0 {
		ls := clip(isAbove, point{da.Min.X, da.Min.Y}, point{0, -1}, line)
		lines1 = append(lines1, ls...)
	}

	// clip left
	lines0 = lines0[:0]
	for _, line := range lines1 {
		ls := clip(isRight, point{da.Min.X, da.Min.Y}, point{1, 0}, line)
		lines0 = append(lines0, ls...)
	}

	// clip top
	lines1 = lines1[:0]
	for _, line := range lines0 {
		ls := clip(isBelow, point{da.Min.X, da.Max().Y}, point{0, 1}, line)
		lines1 = append(lines1, ls...)
	}

	for _, l := range lines1 {
		da.line(l)
	}
	return
}

// clip performs clipping in a single clipping line specified
// by the norm, clip point, and in function.
func clip(in func(point, point) bool, clip, norm point, pts []point) (lines [][]point) {
	var l []point
	for i := 0; i < len(pts); i++ {
		cur, next := pts[i-1], pts[i]
		curIn, nextIn := in(cur, clip), in(next, clip)
		switch {
		case curIn && nextIn:
			l = append(l, cur)

		case curIn && !nextIn:
			l = append(l, cur, isect(cur, next, clip, norm))
			lines = append(lines, l)
			l = []point{}

		case !curIn && !nextIn:
			// do nothing

		default: // !curIn && nextIn
			l = append(l, isect(cur, next, clip, norm))
		}
		if nextIn && i == len(pts)-1 {
			l = append(l, next)
		}
	}
	if len(l) > 1 {
		lines = append(lines, l)
	}
	return
}

// slop is some slop for floating point equality
const slop = 3e-8 // ≈ √1⁻¹⁵

func isLeft(p, clip point) bool {
	return p.X <= clip.X+slop
}

func isRight(p, clip point) bool {
	return p.X >= clip.X-slop
}

func isBelow(p, clip point) bool {
	return p.Y <= clip.Y+slop
}

func isAbove(p, clip point) bool {
	return p.Y >= clip.Y-slop
}

// isect returns the intersection of a line p0→p1 with the
// clipping line specified by the clip point and normal.
func isect(p0, p1, clip, norm point) point {
	// t = (norm · (p0 - clip)) / (norm · (p0 - p1))
	t := p0.minus(clip).dot(norm) / p0.minus(p1).dot(norm)

	// p = p0 + t*(p1 - p0)
	return p1.minus(p0).scale(t).plus(p0)
}

// rectPath returns the path of a rectangle specified by its
// upper left corner, width and height.
func rectPath(r rect) (p vg.Path) {
	p.Move(r.Min.X, r.Min.Y)
	p.Line(r.Max().X, r.Min.Y)
	p.Line(r.Max().X, r.Max().Y)
	p.Line(r.Min.X, r.Max().Y)
	p.Close()
	return
}

// A rect represents a rectangular region of 2d space.
type rect struct {
	Min, Size point
}

// Max returns the maxmium x and y values of a rect.
func (r rect) Max() point {
	return point{
		X: r.Min.X + r.Size.X,
		Y: r.Min.Y + r.Size.Y,
	}
}

// A point is a location in 2d space.
type point struct {
	X, Y vg.Length
}

// Pt returns a point with the given x and y values.
func Pt(x, y vg.Length) point {
	return point{x, y}
}

// dot returns the dot product of two points.
func (p point) dot(q point) vg.Length {
	return p.X*q.X + p.Y*q.Y
}

// plus returns the component-wise sum of two points.
func (p point) plus(q point) point {
	return point{p.X + q.X, p.Y + q.Y}
}

// minus returns the component-wise difference of two points.
func (p point) minus(q point) point {
	return point{p.X - q.X, p.Y - q.Y}
}

// scale returns the component-wise product of a point and a scalar.
func (p point) scale(s vg.Length) point {
	return point{p.X * s, p.Y * s}
}
