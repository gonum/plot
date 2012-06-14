// The plt package implements low-level plotting functionality.
package plt

import (
	"code.google.com/p/plotinum/vg"
	"image/color"
	"math"
)

const (
	defaultFont = "Times-Roman"
)

// Plot is the basic type representing a plot.
type Plot struct {
	Title struct {
		// Text is the text of the plot title.  If
		// Text is the empty string then the plot
		// will not have a title.
		Text string
		TextStyle
	}

	// X and Y are the horizontal and vertical axes
	// of the plot respectively.
	X, Y Axis

	// data is a slice of all data elements on the plot.
	data []Data
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

// AddData adds Data to the plot and changes the minimum
// and maximum values of the X and Y axes to fit the
// newly added data.
func (p *Plot) AddData(d Data) {
	xmin, ymin, xmax, ymax := d.extents()
	p.X.Min = math.Min(p.X.Min, xmin)
	p.X.Max = math.Max(p.X.Max, xmax)
	p.Y.Min = math.Min(p.Y.Min, ymin)
	p.Y.Max = math.Max(p.Y.Max, ymax)
	p.data = append(p.data, d)
}

// Draw draws a plot to a drawArea.
func (p *Plot) Draw(da *drawArea) {
	da.SetColor(color.White)
	da.Fill(rectPath(da.rect))
	da.SetColor(color.Black)
	da.Stroke(rectPath(da.rect))

	if p.Title.Text != "" {
		da.fillText(p.Title.TextStyle, da.center().x, da.max().y, -0.5, -1, p.Title.Text)
		da.size.y -= p.Title.height(p.Title.Text)
	}

	x := horizontalAxis{p.X}
	y := verticalAxis{p.Y}

	ywidth := y.size()
	x.draw(da.crop(ywidth, 0, 0, 0).squishX(x.glyphBoxes()))

	xheight := x.size()
	y.draw(da.crop(0, xheight, 0, 0).squishY(y.glyphBoxes()))

	da = da.crop(ywidth, xheight, 0, 0)
	da = da.squishX(x.glyphBoxes()).squishY(y.glyphBoxes())
	for _, data := range p.data {
		data.plot(da, p)
	}
}
