// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vecgfx"
	"image/color"
)

type Plot struct {
	Title string
	TitleStyle TextStyle
	XAxis, YAxis Axis
}

func NewPlot() *Plot {
	titleFont, err := MakeFont("Times-Roman", 12)
	if err != nil {
		panic(err)
	}
	return &Plot{
		TitleStyle: TextStyle {
			Color: color.RGBA{A:255},
			Font: titleFont,
		},
		XAxis: MakeAxis(0, 0),
		YAxis: MakeAxis(0, 0),
	}
}

func (p *Plot) Draw(da *DrawArea) {
	pad := 5.0/vecgfx.PtInch * da.DPI()
	da.Min.X += pad
	da.Min.Y += pad
	da.Sz.X -= 2*pad
	da.Sz.Y -= 2*pad

	if p.Title != "" {
		da.SetTextStyle(p.TitleStyle)
		da.Text(da.Center().X, da.Max().Y, -0.5, -1, p.Title)
		da.Sz.Y -= p.TitleStyle.Font.Extents().Height/vecgfx.PtInch * da.DPI()
	}

	ywidth := p.YAxis.width()
	p.XAxis.drawHoriz(da.crop(ywidth, 0, -ywidth, 0))

	xheight := p.XAxis.height()
	p.YAxis.drawVert(da.crop(0, xheight, 0, -xheight))
}

type Rect struct {
	Min, Sz Point
}

func (r Rect) Max() Point {
	return Point {
		X: r.Min.X + r.Sz.X,
		Y: r.Min.Y + r.Sz.Y,
	}
}

type Point struct {
	X, Y float64
}

// Pt returns a point with the given x and y values.
func Pt(x, y float64) Point {
	return Point{x, y}
}

// dot returns the dot product of two points.
func (p Point) dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

// plus returns the component-wise sum of two points.
func (p Point) plus(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// minus returns the component-wise differenec of two points.
func (p Point) minus(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// scale returns the component-wise product of a point and a scalar.
func (p Point) scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}
