// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plot

import (
	"code.google.com/p/plotinum/vg"
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

// A GlyphStyle specifies the look of a glyph used to draw
// a point on a plot.
type GlyphStyle struct {
	// Color is the color used to draw the glyph.
	color.Color

	// Radius specifies the size of the glyph's radius.
	Radius vg.Length

	// Shape draws the shape of the glyph.
	Shape GlyphDrawer
}

// A GlyphDrawer wraps the DrawGlyph function.
type GlyphDrawer interface {
	// DrawGlyph draws the glyph at the given
	// point, with the given color and radius.
	DrawGlyph(*DrawArea, GlyphStyle, Point)
}

// DrawGlyph draws the given glyph to the draw
// area.  If the point is not within the DrawArea
// or the sty.Shape is nil then nothing is drawn.
func (da *DrawArea) DrawGlyph(sty GlyphStyle, pt Point) {
	if sty.Shape == nil || !da.Contains(pt) {
		return
	}
	da.SetColor(sty.Color)
	sty.Shape.DrawGlyph(da, sty, pt)
}

// DrawGlyphNoClip draws the given glyph to the draw
// area.  If the sty.Shape is nil then nothing is drawn.
func (da *DrawArea) DrawGlyphNoClip(sty GlyphStyle, pt Point) {
	if sty.Shape == nil {
		return
	}
	da.SetColor(sty.Color)
	sty.Shape.DrawGlyph(da, sty, pt)
}

// Rect returns the rectangle surrounding this glyph,
// assuming that it is drawn centered at 0,0
func (g GlyphStyle) Rect() Rect {
	return Rect{Point{-g.Radius, -g.Radius}, Point{g.Radius * 2, g.Radius * 2}}
}

// CircleGlyph is a glyph that draws a solid circle.
type CircleGlyph struct{}

// DrawGlyph implements the GlyphDrawer interface.
func (c CircleGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	var p vg.Path
	p.Move(pt.X+sty.Radius, pt.Y)
	p.Arc(pt.X, pt.Y, sty.Radius, 0, 2*math.Pi)
	p.Close()
	da.Fill(p)
}

// RingGlyph is a glyph that draws the outline of a circle.
type RingGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (RingGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	da.SetLineStyle(LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	var p vg.Path
	p.Move(pt.X+sty.Radius, pt.Y)
	p.Arc(pt.X, pt.Y, sty.Radius, 0, 2*math.Pi)
	p.Close()
	da.Stroke(p)
}

const (
	cosπover4 = vg.Length(.707106781202420)
	sinπover6 = vg.Length(.500000000025921)
	cosπover6 = vg.Length(.866025403769473)
)

// SquareGlyph is a glyph that draws the outline of a square.
type SquareGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (SquareGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	da.SetLineStyle(LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	x := (sty.Radius-sty.Radius*cosπover4)/2 + sty.Radius*cosπover4
	var p vg.Path
	p.Move(pt.X-x, pt.Y-x)
	p.Line(pt.X+x, pt.Y-x)
	p.Line(pt.X+x, pt.Y+x)
	p.Line(pt.X-x, pt.Y+x)
	p.Close()
	da.Stroke(p)
}

// BoxGlyph is a glyph that draws a filled square.
type BoxGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (BoxGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	x := (sty.Radius-sty.Radius*cosπover4)/2 + sty.Radius*cosπover4
	var p vg.Path
	p.Move(pt.X-x, pt.Y-x)
	p.Line(pt.X+x, pt.Y-x)
	p.Line(pt.X+x, pt.Y+x)
	p.Line(pt.X-x, pt.Y+x)
	p.Close()
	da.Fill(p)
}

// TriangleGlyph is a glyph that draws the outline of a triangle.
type TriangleGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (TriangleGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	da.SetLineStyle(LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	r := sty.Radius + (sty.Radius-sty.Radius*sinπover6)/2
	var p vg.Path
	p.Move(pt.X, pt.Y+r)
	p.Line(pt.X-r*cosπover6, pt.Y-r*sinπover6)
	p.Line(pt.X+r*cosπover6, pt.Y-r*sinπover6)
	p.Close()
	da.Stroke(p)
}

// PyramidGlyph is a glyph that draws a filled triangle.
type PyramidGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (PyramidGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	r := sty.Radius + (sty.Radius-sty.Radius*sinπover6)/2
	var p vg.Path
	p.Move(pt.X, pt.Y+r)
	p.Line(pt.X-r*cosπover6, pt.Y-r*sinπover6)
	p.Line(pt.X+r*cosπover6, pt.Y-r*sinπover6)
	p.Close()
	da.Fill(p)
}

// PlusGlyph is a glyph that draws a plus sign
type PlusGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (PlusGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	da.SetLineStyle(LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	r := sty.Radius
	var p vg.Path
	p.Move(pt.X, pt.Y+r)
	p.Line(pt.X, pt.Y-r)
	da.Stroke(p)
	p = vg.Path{}
	p.Move(pt.X-r, pt.Y)
	p.Line(pt.X+r, pt.Y)
	da.Stroke(p)
}

// CrossGlyph is a glyph that draws a big X.
type CrossGlyph struct{}

// DrawGlyph implements the Glyph interface.
func (CrossGlyph) DrawGlyph(da *DrawArea, sty GlyphStyle, pt Point) {
	da.SetLineStyle(LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	r := sty.Radius * cosπover4
	var p vg.Path
	p.Move(pt.X-r, pt.Y-r)
	p.Line(pt.X+r, pt.Y+r)
	da.Stroke(p)
	p = vg.Path{}
	p.Move(pt.X-r, pt.Y+r)
	p.Line(pt.X+r, pt.Y-r)
	da.Stroke(p)
}

// MakeDrawArea returns a new DrawArea for a canvas with a
// Size method.
func MakeDrawArea(c interface {
	vg.Canvas
	Size() (vg.Length, vg.Length)
}) DrawArea {
	w, h := c.Size()
	return MakeDrawAreaSize(c, w, h)
}

// MakeDrawAreaSize returns a new DrawArea of the given
// size for a canvas.
func MakeDrawAreaSize(c vg.Canvas, w, h vg.Length) DrawArea {
	return DrawArea{
		Canvas: c,
		Rect:   Rect{Size: Point{w, h}},
	}
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
	return da.ContainsX(p.X) && da.ContainsY(p.Y)
}

// Contains returns true if the DrawArea contains the
// x coordinate.
func (da *DrawArea) ContainsX(x vg.Length) bool {
	return x <= da.Max().X+slop && x >= da.Min.X-slop
}

// ContainsY returns true if the DrawArea contains the
// y coordinate.
func (da *DrawArea) ContainsY(y vg.Length) bool {
	return y <= da.Max().Y+slop && y >= da.Min.Y-slop
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
func (da DrawArea) crop(minx, miny, maxx, maxy vg.Length) DrawArea {
	minpt := Point{
		X: da.Min.X + minx,
		Y: da.Min.Y + miny,
	}
	sz := Point{
		X: da.Max().X + maxx - minpt.X,
		Y: da.Max().Y + maxy - minpt.Y,
	}
	return DrawArea{
		Canvas: vg.Canvas(da),
		Rect:   Rect{Min: minpt, Size: sz},
	}
}

// SetLineStyle sets the current line style
func (da *DrawArea) SetLineStyle(sty LineStyle) {
	da.SetColor(sty.Color)
	da.SetLineWidth(sty.Width)
	var dashDots []vg.Length
	for _, dash := range sty.Dashes {
		dashDots = append(dashDots, dash)
	}
	da.SetLineDash(dashDots, sty.DashOffs)
}

// StrokeLines draws a line connecting a set of points
// in the given DrawArea.
func (da *DrawArea) StrokeLines(sty LineStyle, lines ...[]Point) {
	if len(lines) == 0 {
		return
	}

	da.SetLineStyle(sty)

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		var p vg.Path
		p.Move(l[0].X, l[0].Y)
		for _, pt := range l[1:] {
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
		ls := clipLine(isLeft, Point{da.Max().X, da.Min.Y}, Point{-1, 0}, line)
		lines1 = append(lines1, ls...)
	}
	for _, line := range lines1 {
		ls := clipLine(isRight, Point{da.Min.X, da.Min.Y}, Point{1, 0}, line)
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
		ls := clipLine(isAbove, Point{da.Min.X, da.Min.Y}, Point{0, -1}, line)
		lines1 = append(lines1, ls...)
	}
	for _, line := range lines1 {
		ls := clipLine(isBelow, Point{da.Min.X, da.Max().Y}, Point{0, 1}, line)
		clipped = append(clipped, ls...)
	}
	return
}

// clipLine performs clipping of a line by a single
// clipping line specified by the norm, clip point,
// and in function.
func clipLine(in func(Point, Point) bool, clip, norm Point, pts []Point) (lines [][]Point) {
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

// FillPolygon fills a polygon with the given color.
func (da *DrawArea) FillPolygon(clr color.Color, pts []Point) {
	if len(pts) == 0 {
		return
	}

	da.SetColor(clr)
	var p vg.Path
	p.Move(pts[0].X, pts[0].Y)
	for _, pt := range pts[1:] {
		p.Line(pt.X, pt.Y)
	}
	p.Close()
	da.Fill(p)
}

// ClipPolygonXY returns a slice of lines that
// represent the given polygon clipped in both
// X and Y directions.
func (da *DrawArea) ClipPolygonXY(pts []Point) []Point {
	return da.ClipPolygonY(da.ClipPolygonX(pts))
}

// ClipPolygonX returns a slice of lines that
// represent the given polygon clipped in the
// X direction.
func (da *DrawArea) ClipPolygonX(pts []Point) []Point {
	return clipPoly(isLeft, Point{da.Max().X, da.Min.Y}, Point{-1, 0},
		clipPoly(isRight, Point{da.Min.X, da.Min.Y}, Point{1, 0}, pts))
}

// ClipPolygonY returns a slice of lines that
// represent the given polygon clipped in the
// Y direction.
func (da *DrawArea) ClipPolygonY(pts []Point) []Point {
	return clipPoly(isBelow, Point{da.Min.X, da.Max().Y}, Point{0, 1},
		clipPoly(isAbove, Point{da.Min.X, da.Min.Y}, Point{0, -1}, pts))
}

// clipPoly performs clipping of a polygon by a single
// clipping line specified by the norm, clip point,
// and in function.
func clipPoly(in func(Point, Point) bool, clip, norm Point, pts []Point) (clipped []Point) {
	for i := 0; i < len(pts); i++ {
		j := i + 1
		if i == len(pts)-1 {
			j = 0
		}
		cur, next := pts[i], pts[j]
		curIn, nextIn := in(cur, clip), in(next, clip)
		switch {
		case curIn && nextIn:
			clipped = append(clipped, cur)

		case curIn && !nextIn:
			clipped = append(clipped, cur, isect(cur, next, clip, norm))

		case !curIn && !nextIn:
			// do nothing

		default: // !curIn && nextIn
			clipped = append(clipped, isect(cur, next, clip, norm))
		}
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

// Rect returns a rectangle giving the bounds of
// this text assuming that it is drawn at 0, 0
func (sty TextStyle) Rect(txt string) Rect {
	return Rect{Size: Point{sty.Width(txt), sty.Height(txt)}}
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

// A Point is a location in 2d space.
//
// Points are used for drawing, not for data.  For
// data, see the XYer interface.
type Point struct {
	X, Y vg.Length
}

// Pt returns a point from x, y coordinates.
func Pt(x, y vg.Length) Point {
	return Point{X: x, Y: y}
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
