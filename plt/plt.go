// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vg"
	"image/color"
)

// A Plot is a pair of axes, an optional title and a set of
// data elements that can be drawn to a DrawArea.
type Plot struct {
	Title        string
	TitleStyle   TextStyle
	XAxis HorizontalAxis
	YAxis VerticalAxis
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
		XAxis: HorizontalAxis{ makeAxis() },
		YAxis: VerticalAxis{ makeAxis() },
	}
}

// Draw draws a plot to a DrawArea.
func (p *Plot) Draw(da *DrawArea) {
	da.SetColor(White)
	da.Fill(rectPath(da.Rect))
	da.SetColor(Black)
	da.Stroke(rectPath(da.Rect))

	pad := 5.0 / vg.PtInch
	da = da.crop(0, pad, 0, -pad)

	if p.Title != "" {
		da.setTextStyle(p.TitleStyle)
		da.text(da.center().X, da.Max().Y, -0.5, -1, p.Title)
		da.Size.Y -= p.TitleStyle.Font.Extents().Height / vg.PtInch * da.DPI()
	}

	ywidth := p.YAxis.size()
	p.XAxis.draw(da.crop(ywidth, 0, 0, 0).squishX(p.XAxis.glyphBoxes()))

	xheight := p.XAxis.size()
	p.YAxis.draw(da.crop(0, xheight, 0, 0))
}
