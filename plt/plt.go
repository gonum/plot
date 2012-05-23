// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vecgfx"
	"image/color"
)

// A Plot is a pair of axes, an optional title and a set of
// data elements that can be drawn to a DrawArea.
type Plot struct {
	Title        string
	TitleStyle   TextStyle
	XAxis, YAxis Axis
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
		XAxis: MakeAxis(),
		YAxis: MakeAxis(),
	}
}

// Draw draws a plot to a DrawArea.
func (p *Plot) Draw(da *DrawArea) {
	da.SetColor(White)
	da.Fill(RectPath(da.Rect))

	pad := 5.0 / vecgfx.PtInch
	da = da.crop(pad, pad, -2*pad, -2*pad)

	if p.Title != "" {
		da.SetTextStyle(p.TitleStyle)
		da.Text(da.Center().X, da.Max().Y, -0.5, -1, p.Title)
		da.Sz.Y -= p.TitleStyle.Font.Extents().Height / vecgfx.PtInch * da.DPI()
	}

	ywidth := p.YAxis.width()
	p.XAxis.drawHoriz(da.crop(ywidth, 0, -ywidth, 0))

	xheight := p.XAxis.height()
	p.YAxis.drawVert(da.crop(0, xheight, 0, -xheight))
}
