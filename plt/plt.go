// The plt package implements low-level plotting functionality.
package plt

import (
	"image/color"
)

// A Plot is a pair of axes, an optional title and a set of
// data elements that can be drawn to a drawArea.
type Plot struct {
	Title      struct {
		Text string
		TextStyle
	}
	XAxis      horizontalAxis
	YAxis      verticalAxis
}

// NewPlot returns a new plot.
func NewPlot() *Plot {
	titleFont, err := MakeFont("Times-Roman", 12)
	if err != nil {
		panic(err)
	}
	p := &Plot{
		XAxis: horizontalAxis{makeAxis()},
		YAxis: verticalAxis{makeAxis()},
	}
	p.Title.TextStyle = TextStyle{
		Color: color.RGBA{A: 255},
		Font:  titleFont,
	}
	return p
}

// draw draws a plot to a drawArea.
func (p *Plot) draw(da *drawArea) {
	da.SetColor(White)
	da.Fill(rectPath(da.rect))
	da.SetColor(Black)
	da.Stroke(rectPath(da.rect))

	if p.Title.Text != "" {
		da.setTextStyle(p.Title.TextStyle)
		fillText(da, da.center().x, da.max().y, -0.5, -1, p.Title.Text)
		da.size.y -= textHeight(p.Title.Font, p.Title.Text)
	}

	ywidth := p.YAxis.size()
	p.XAxis.draw(da.crop(ywidth, 0, 0, 0).squishX(p.XAxis.glyphBoxes()))

	xheight := p.XAxis.size()
	p.YAxis.draw(da.crop(0, xheight, 0, 0).squishY(p.YAxis.glyphBoxes()))
}
