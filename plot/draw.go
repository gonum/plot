package plot

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
	"strings"
)

// A DrawArea is a vector graphics canvas along with
// an associated Rect defining a section of the canvas
// to which drawing should take place.
type DrawArea struct {
	vg.Canvas
	Rect
}

// TextStyle describes what text will look like.
type TextStyle struct {
	// Color is the text color.
	Color color.Color

	// Font is the font description.
	Font vg.Font
}

// LineStyle describes what a line will look like.
type LineStyle struct {
	// Color is the color of the line.
	Color color.Color

	// Width is the width of the line.
	Width vg.Length

	Dashes   []vg.Length
	DashOffs vg.Length
}

// A GlyphShape is a lable representing a shape for drawing
// a glyph that represents a point.
//
// GlyphShape values that corresponding to uppercase ASCII
// letters (A'–'Z'), represent the shape of the corresponding
// character.  A handful of other GlyphShape values are
// defined as constants, all other GlyphShape values are
// invalid.
type GlyphShape uint8

const (
	// RingGlyph is an outlined circle
	RingGlyph GlyphShape = iota

	// CircleGlyph is a filled circle
	CircleGlyph
)

// A GlyphStyle specifies the look of a glyph used to draw
// a point on a plot.
type GlyphStyle struct {
	// Color is the color used to draw the glyph.
	color.Color

	// Shape is the shape of the glyph.
	Shape GlyphShape

	// Radius specifies the size of the glyph's radius.
	Radius vg.Length
}

// NewDrawArea returns a new DrawArea of a specified
// size using the given canvas.
func NewDrawArea(c vg.Canvas, w, h vg.Length) *DrawArea {
	return &DrawArea{Canvas: c, Rect: Rect{Min: Point{}, Size: Point{w, h}}}
}

// Center returns the center point of the area
func (da *DrawArea) Center() Point {
	return Point{
		X: (da.Max().X-da.Min.X)/2 + da.Min.X,
		Y: (da.Max().Y-da.Min.Y)/2 + da.Min.Y,
	}
}

// Contains returns true if the DrawArea contains the point.
func (da *DrawArea) Contains(p Point) bool {
	return p.X <= da.Max().X && p.X >= da.Min.X &&
		p.Y <= da.Max().Y && p.Y >= da.Min.Y
}

// X returns the value of x, given in the unit range,
// in the drawing coordinates of this draw area.
// A value of 0, for example, will return the minimum
// x value of the draw area and a value of 1 will
// return the maximum.
func (da *DrawArea) X(x float64) vg.Length {
	return vg.Length(x)*(da.Max().X-da.Min.X) + da.Min.X
}

// Y returns the value of x, given in the unit range,
// in the drawing coordinates of this draw area.
// A value of 0, for example, will return the minimum
// y value of the draw area and a value of 1 will
// return the maximum.
func (da *DrawArea) Y(y float64) vg.Length {
	return vg.Length(y)*(da.Max().Y-da.Min.Y) + da.Min.Y
}

// crop returns a new DrawArea corresponding to the receiver
// area with the given number of inches added to the minimum
// and maximum x and y values of the DrawArea's Rect.
func (da *DrawArea) crop(minx, miny, maxx, maxy vg.Length) *DrawArea {
	minpt := Point{
		X: da.Min.X + minx,
		Y: da.Min.Y + miny,
	}
	sz := Point{
		X: da.Max().X + maxx - minpt.X,
		Y: da.Max().Y + maxy - minpt.Y,
	}
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		Rect:      Rect{Min: minpt, Size: sz},
	}
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

// DrawGlyph draws a glyph at a specified location.  If
// the location is outside of the DrawArea then it is
// not drawn.
func (da *DrawArea) DrawGlyph(sty GlyphStyle, pt Point) {
	if !da.Contains(pt) {
		return
	}

	da.setLineStyle(LineStyle{Width: vg.Points(0.5)})
	da.SetColor(sty.Color)

	switch {
	case sty.Shape == CircleGlyph:
		var p vg.Path
		p.Move(pt.X+sty.Radius, pt.Y)
		p.Arc(pt.X, pt.Y, sty.Radius, 0, 2*math.Pi)
		p.Close()
		da.Fill(p)

	case sty.Shape == RingGlyph:
		var p vg.Path
		p.Move(pt.X+sty.Radius, pt.Y)
		p.Arc(pt.X, pt.Y, sty.Radius, 0, 2*math.Pi)
		p.Close()
		da.Stroke(p)

	case sty.Shape >= 'A' && sty.Shape <= 'Z':
		font, err := vg.MakeFont(defaultFont, sty.Radius*2)
		if err != nil {
			panic(err)
		}
		str := string([]byte{byte(sty.Shape)})
		x := pt.X - font.Width(str)/2
		y := pt.Y + font.Extents().Descent
		da.FillString(font, x, y, str)

	default:
		panic(fmt.Sprintf("Invalid GlyphShape: %d", sty.Shape))
	}
}

// StrokeLine draws a line connecting a set of points
// in the given DrawArea.
func (da *DrawArea) StrokeLines(sty LineStyle, lines ...[]Point) {
	if len(lines) == 0 {
		return
	}

	da.setLineStyle(sty)

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		var p vg.Path
		p.Move(l[0].X, l[0].Y)
		for _, pt := range l {
			p.Line(pt.X, pt.Y)
		}
		da.Stroke(p)
	}
}

// StrokeLine2 draws a line between two points in the given
// DrawArea.
func (da *DrawArea) StrokeLine2(sty LineStyle, x0, y0, x1, y1 vg.Length) {
	da.StrokeLines(sty, []Point{{x0, y0}, {x1, y1}})
}

// ClipLineXY returns a slice of lines that
// represent the given line clipped in both
// X and Y directions.
func (da *DrawArea) ClipLinesXY(lines ...[]Point) [][]Point {
	return da.ClipLinesY(da.ClipLinesX(lines...)...)
}

// ClipLineX returns a slice of lines that
// represent the given line clipped in the
// X direction.
func (da *DrawArea) ClipLinesX(lines ...[]Point) (clipped [][]Point) {
	var lines1 [][]Point
	for _, line := range lines {
		ls := clip(isLeft, Point{da.Max().X, da.Min.Y}, Point{-1, 0}, line)
		lines1 = append(lines1, ls...)
	}
	for _, line := range lines1 {
		ls := clip(isRight, Point{da.Min.X, da.Min.Y}, Point{1, 0}, line)
		clipped = append(clipped, ls...)
	}
	return
}

// ClipLineY returns a slice of lines that
// represent the given line clipped in the
// Y direction.
func (da *DrawArea) ClipLinesY(lines ...[]Point) (clipped [][]Point) {
	var lines1 [][]Point
	for _, line := range lines {
		ls := clip(isAbove, Point{da.Min.X, da.Min.Y}, Point{0, -1}, line)
		lines1 = append(lines1, ls...)
	}
	for _, line := range lines1 {
		ls := clip(isBelow, Point{da.Min.X, da.Max().Y}, Point{0, 1}, line)
		clipped = append(clipped, ls...)
	}
	return
}

// clip performs clipping in a single clipping line specified
// by the norm, clip point, and in function.
func clip(in func(Point, Point) bool, clip, norm Point, pts []Point) (lines [][]Point) {
	var l []Point
	for i := 1; i < len(pts); i++ {
		cur, next := pts[i-1], pts[i]
		curIn, nextIn := in(cur, clip), in(next, clip)
		switch {
		case curIn && nextIn:
			l = append(l, cur)

		case curIn && !nextIn:
			l = append(l, cur, isect(cur, next, clip, norm))
			lines = append(lines, l)
			l = []Point{}

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
	t := p0.minus(clip).dot(norm) / p0.minus(p1).dot(norm)

	// p = p0 + t*(p1 - p0)
	return p1.minus(p0).scale(t).plus(p0)
}

// FillText fills lines of text in the draw area.
// The text is offset by its width times xalign and
// its height times yalign.  x and y give the bottom
// left corner of the text befor e it is offset.
func (da *DrawArea) FillText(sty TextStyle, x, y vg.Length, xalign, yalign float64, txt string) {
	txt = strings.TrimRight(txt, "\n")
	if len(txt) == 0 {
		return
	}

	da.SetColor(sty.Color)

	ht := sty.Height(txt)
	y += ht*vg.Length(yalign) - sty.Font.Extents().Ascent
	nl := textNLines(txt)
	for i, line := range strings.Split(txt, "\n") {
		xoffs := vg.Length(xalign) * sty.Font.Width(line)
		n := vg.Length(nl - i)
		da.FillString(sty.Font, x+xoffs, y+n*sty.Font.Size, line)
	}
}

// Width returns the width of lines of text
// when using the given font.
func (sty TextStyle) Width(txt string) (max vg.Length) {
	txt = strings.TrimRight(txt, "\n")
	for _, line := range strings.Split(txt, "\n") {
		if w := sty.Font.Width(line); w > max {
			max = w
		}
	}
	return
}

// Height returns the height of the text when using
// the given font.
func (sty TextStyle) Height(txt string) vg.Length {
	nl := textNLines(txt)
	if nl == 0 {
		return vg.Length(0)
	}
	e := sty.Font.Extents()
	return e.Height*vg.Length(nl-1) + e.Ascent
}

// textNLines returns the number of lines in the text.
func textNLines(txt string) int {
	txt = strings.TrimRight(txt, "\n")
	if len(txt) == 0 {
		return 0
	}
	n := 1
	for _, r := range txt {
		if r == '\n' {
			n++
		}
	}
	return n
}

// rectPath returns the path of a Rectangle specified by its
// upper left corner, width and height.
func rectPath(r Rect) (p vg.Path) {
	p.Move(r.Min.X, r.Min.Y)
	p.Line(r.Max().X, r.Min.Y)
	p.Line(r.Max().X, r.Max().Y)
	p.Line(r.Min.X, r.Max().Y)
	p.Close()
	return
}

// A Rect represents a Rectangular region of 2d space.
type Rect struct {
	Min, Size Point
}

// Max returns the maxmium x and y values of a Rect.
func (r Rect) Max() Point {
	return Point{
		X: r.Min.X + r.Size.X,
		Y: r.Min.Y + r.Size.Y,
	}
}

// A point is a location in 2d space.
type Point struct {
	X, Y vg.Length
}

// dot returns the dot product of two points.
func (p Point) dot(q Point) vg.Length {
	return p.X*q.X + p.Y*q.Y
}

// plus returns the component-wise sum of two points.
func (p Point) plus(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// minus returns the component-wise difference of two points.
func (p Point) minus(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// scale returns the component-wise product of a point and a scalar.
func (p Point) scale(s vg.Length) Point {
	return Point{p.X * s, p.Y * s}
}
