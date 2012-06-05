// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vg"
	"image/color"
)

// A Plot is a pair of axes, an optional title and a set of
// data elements that can be drawn to a drawArea.
type Plot struct {
	Title      string
	TitleStyle TextStyle
	XAxis      horizontalAxis
	YAxis      verticalAxis
}

// NewPlot returns a new plot.
func NewPlot() *Plot {
	titleFont, err := MakeFont("Times-Roman", 12)
	if err != nil {
		panic(err)
	}
	return &Plot{
		TitleStyle: TextStyle{
			Color: color.RGBA{A: 255},
			Font:  titleFont,
		},
		XAxis: horizontalAxis{makeAxis()},
		YAxis: verticalAxis{makeAxis()},
	}
}

// draw draws a plot to a drawArea.
func (p *Plot) draw(da *drawArea) {
	da.SetColor(White)
	da.Fill(rectPath(da.rect))
	da.SetColor(Black)
	da.Stroke(rectPath(da.rect))

	pad := vg.Points(5.0)
	da = da.crop(0, pad, 0, -pad)

	if p.Title != "" {
		da.setTextStyle(p.TitleStyle)
		da.text(da.center().x, da.max().y, -0.5, -1, p.Title)
		da.size.y -= p.TitleStyle.Font.Extents().Height
	}

	ywidth := p.YAxis.size()
	p.XAxis.draw(da.crop(ywidth, 0, 0, 0).squishX(p.XAxis.glyphBoxes()))

	xheight := p.XAxis.size()
	p.YAxis.draw(da.crop(0, xheight, 0, 0).squishY(p.YAxis.glyphBoxes()))
}
