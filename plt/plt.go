// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vg"
	"image/color"
)

const (
	defaultFont = "Times-Roman"
)

// Plot is the basic type representing a plot.
type Plot struct {
	Title struct {
		Text string
		TextStyle
	}
	X, Y Axis
}

// New returns a new plot.
func New() *Plot {
	titleFont, err := vg.MakeFont(defaultFont, 12)
	if err != nil {
		panic(err)
	}
	p := &Plot{
		X: makeAxis(),
		Y: makeAxis(),
	}
	p.Title.TextStyle = TextStyle{
		Color: color.Black,
		Font:  titleFont,
	}
	return p
}

// draw draws a plot to a drawArea.
func (p *Plot) draw(da *drawArea) {
	da.SetColor(color.White)
	da.Fill(rectPath(da.rect))
	da.SetColor(color.Black)
	da.Stroke(rectPath(da.rect))

	if p.Title.Text != "" {
		da.setTextStyle(p.Title.TextStyle)
		fillText(da, da.center().x, da.max().y, -0.5, -1, p.Title.Text)
		da.size.y -= textHeight(p.Title.Font, p.Title.Text)
	}

	x := horizontalAxis{p.X}
	y := verticalAxis{p.Y}

	ywidth := y.size()
	x.draw(da.crop(ywidth, 0, 0, 0).squishX(x.glyphBoxes()))

	xheight := x.size()
	y.draw(da.crop(0, xheight, 0, 0).squishY(y.glyphBoxes()))
}
