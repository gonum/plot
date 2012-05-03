package plotinum

import (
	"plotinum/vecgfx"
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

	area := da.Rect
	ywidth := p.YAxis.Width() * da.DPI()
	da.Min.X += ywidth
	da.Sz.X -= ywidth
	p.XAxis.DrawHoriz(da)

	da.Rect = area
	xheight := p.XAxis.Height() * da.DPI()
	da.Min.Y += xheight
	da.Sz.Y -= xheight
	p.YAxis.DrawVert(da)
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

func (p Point) Dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Point) Plus(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Minus(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}
