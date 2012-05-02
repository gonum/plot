package plotinum

import (
	"fmt"
	"math"
	"plotinum/vecgfx"
	"image/color"
)

var (
	Black = color.RGBA{A: 255}
	White = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Red = color.RGBA{R: 255, A: 255}
	Green = color.RGBA{G: 255, A: 255}
	Blue = color.RGBA{B: 255, A: 255}
)

// A DrawArea is a vector graphics canvas along with
// an associated rectangle defining a section of the canvas
// to which drawing should take place.
type DrawArea struct {
	vecgfx.Canvas
	font vecgfx.Font
	Rect
}

// Center returns the center point of the area
func (da *DrawArea) Center() Point {
	return Point{
		X: (da.Max.X - da.Min.X)/2 + da.Min.X,
		Y: (da.Max.Y - da.Min.Y)/2 + da.Min.Y,
	}
}

// SetTextStyle sets the current text style
func (da *DrawArea) SetTextStyle(sty TextStyle) {
	da.SetColor(sty.Color)
	da.font = sty.Font
}

func (da *DrawArea) Pt2Dot(pt float64) float64 {
	return pt/vecgfx.PtInch * da.DPI()
}

// TextStyle describes what text will look like.
type TextStyle struct {
	// Color is the text color.
	Color color.Color

	// Font is the font description.
	Font vecgfx.Font
}

// MakeFont returns a font object.
// This function is merely included for convenience so that
// the user doesn't have to import the vecgfx package.
func MakeFont(name string, size float64) (vecgfx.Font, error) {
	return vecgfx.MakeFont(name, size)
}

// FontAscent returns the ascent of the current font in dots.
func (da *DrawArea) FontAscent() float64 {
	return da.font.Extents().Ascent/vecgfx.PtInch * da.DPI()
}

// FontDescent returns the descent of the current font in dots.
func (da *DrawArea) FontDescent() float64 {
	return da.font.Extents().Descent/vecgfx.PtInch * da.DPI()
}

// TextWidth returns the width of the text in the current font in dots.
func (da *DrawArea) TextWidth(str string) float64 {
	return da.font.Width(str)/vecgfx.PtInch * da.DPI()
}

// Text fills the text to the drawing area.  The string is created
// using the printf-style format specification and the text is
// located at x + width*fx, y + height*fy, where width and height
// are the width and height of the rendered string.
func (da *DrawArea) Text(x, y, fx, fy float64, f string, v ...interface{}) {
	if da.font.Font() == nil {
		panic("Drawing text without a current font set");
	}
	str := fmt.Sprintf(f, v...)
	w := da.font.Width(str)/vecgfx.PtInch * da.DPI()
	h := da.font.Extents().Ascent/vecgfx.PtInch * da.DPI()
	da.FillText(da.font, x+w*fx, y+h*fy, str)
}

// SetLineStyle sets the current line style
func (da *DrawArea) SetLineStyle(sty LineStyle) {
	da.SetColor(sty.Color)
	da.SetLineWidth(sty.Width*da.DPI())
	var dashDots []float64
	for _, dash := range sty.Dashes {
		dashDots = append(dashDots, dash*da.DPI())
	}
	da.SetLineDash(dashDots, sty.DashOffs*da.DPI())

}

// LineStyle describes what a line will look like.
type LineStyle struct {
	// Color is the color of the line.
	Color color.Color

	// Width is the width of the line in inches.
	Width float64

	Dashes []float64
	DashOffs float64
}

// Line draws a line connecting the given points.
func (da *DrawArea) Line(pts XYer) {
	l := pts.Len()
	if l == 0 {
		return
	}

	var p vecgfx.Path
	p.Move(pts.X(0), pts.Y(0))
	for i := 1; i < l; i++ {
		p.Line(pts.X(i), pts.Y(i))
	}
	da.Stroke(p)
}

// ClippedLine draws a line that is clipped at the bounds
// the DrawArea.
func (da *DrawArea) ClippedLine(pts XYer) {
	// clip right
	lines0 := clip(isLeft, Point{da.Max.X, da.Min.Y}, Point{-1, 0}, pts)

	// clip bottom
	var lines1 []Line
	for _, line := range lines0 {
		ls := clip(isAbove, Point{da.Min.X, da.Min.Y}, Point{0, -1}, line)
		lines1 = append(lines1, ls...)
	}

	// clip left
	lines0 = lines0[:0]
	for _, line := range lines1 {
		ls := clip(isRight, Point{da.Min.X, da.Min.Y}, Point{1, 0}, line)
		lines0 = append(lines0, ls...)
	}

	// clip top
	lines1 = lines1[:0]
	for _, line := range lines0 {
		ls := clip(isBelow, Point{da.Min.X, da.Max.Y}, Point{0, 1}, line)
		lines1 = append(lines1, ls...)
	}

	for _, l := range lines1 {
		da.Line(l)
	}
	return
}

// clip performs clipping in a single clipping line specified
// by the norm, clip point, and in function.
func clip(in func(Point, Point) bool, clip, norm Point, pts XYer) (lines []Line) {
	var l Line
	length := pts.Len()
	for i := 1; i < length; i++ {
		cur, next := Pt(pts, i-1), Pt(pts, i)
		curIn, nextIn := in(cur, clip), in(next, clip)
		switch {
		case curIn && nextIn:
			l = append(l, cur)

		case curIn && !nextIn:
			l = append(l, cur, isect(cur, next, clip, norm))
			lines = append(lines, l)
			l = Line{}

		case !curIn && !nextIn:
			// do nothing

		default: // !curIn && nextIn
			l = append(l, isect(cur, next, clip, norm))
		}
		if nextIn && i == length-1 {
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

func isLeft(p, clip Point) bool {
	return p.X <= clip.X+slop
}

func isRight(p, clip Point) bool {
	return p.X >= clip.X-slop
}

func isBelow(p, clip Point) bool {
	return p.Y <= clip.Y+slop
}

func isAbove(p, clip Point) bool {
	return p.Y >= clip.Y-slop
}

// isect returns the intersection of a line p0→p1 with the
// clipping line specified by the clip point and normal.
func isect(p0, p1, clip, norm Point) Point {
	// t = (norm · (p0 - clip)) / (norm · (p0 - p1))
	t := p0.Minus(clip).Dot(norm) / p0.Minus(p1).Dot(norm)

	// p = p0 + t*(p1 - p0)
	return p1.Minus(p0).Scale(t).Plus(p0)
}

// CirclePath returns the path of a circle centered at x,y with
// radius r.
func CirclePath(x, y, r float64) (p vecgfx.Path) {
	p.Move(x+r, y)
	p.Arc(x, y, r, 0, 2*math.Pi)
	p.Close()
	return
}

// EqTrianglePath returns the path for an equilateral triangle
// that is circumscribed by a circle centered at x,y with
// radius r.  One point of the triangle is directly above the
// center point of the circle.
func EqTrianglePath(x, y, r float64) (p vecgfx.Path) {
	p.Move(x, y+r)
	p.Line(x+r*math.Cos(math.Pi/6), y-r*math.Sin(math.Pi/6))
	p.Line(x-r*math.Cos(math.Pi/6), y-r*math.Sin(math.Pi/6))
	p.Close()
	return
}

// RectPath returns the path of a rectangle specified by its
// upper left corner, width and height.
func RectPath(r Rect) (p vecgfx.Path) {
	p.Move(r.Min.X, r.Min.Y)
	p.Line(r.Max.X, r.Min.Y)
	p.Line(r.Max.X, r.Max.Y)
	p.Line(r.Min.X, r.Max.Y)
	p.Close()
	return
}
