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

// Data is an interface that wraps all of the methods required
// to add data elements to a plot.
type Data interface {
	// Plot draws the data to the given DrawArea using
	// the axes from the given Plot.
	Plot(DrawArea, *Plot)

	// Extents returns the minimum and maximum
	// values of the data.
	Extents() (xmin, ymin, xmax, ymax float64)
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
func (p *Plot) AddData(data ...Data) {
	for _, d := range data {
		xmin, ymin, xmax, ymax := d.Extents()
		p.X.Min = math.Min(p.X.Min, xmin)
		p.X.Max = math.Max(p.X.Max, xmax)
		p.Y.Min = math.Min(p.Y.Min, ymin)
		p.Y.Max = math.Max(p.Y.Max, ymax)
	}
	p.data = append(p.data, data...)
}

// Draw draws a plot to a DrawArea.
func (p *Plot) Draw(da *DrawArea) {
	da.SetColor(color.White)
	da.Fill(rectPath(da.Rect))

	if p.Title.Text != "" {
		da.FillText(p.Title.TextStyle, da.Center().X, da.Max().Y, -0.5, -1, p.Title.Text)
		da.Size.Y -= p.Title.Height(p.Title.Text)
	}

	if p.X.Min == p.X.Max {
		p.X.Min -= 1
		p.X.Max += 1
	}
	x := horizontalAxis{p.X}

	if p.Y.Min == p.Y.Max {
		p.Y.Min -= 1
		p.Y.Max += 1
	}
	y := verticalAxis{p.Y}

	ywidth := y.size()
	x.draw(padX(p, da.crop(ywidth, 0, 0, 0)))

	xheight := x.size()
	y.draw(padY(p, da.crop(0, xheight, 0, 0)))

	da = padY(p, padX(p, da.crop(ywidth, xheight, 0, 0)))
	for _, data := range p.data {
		data.Plot(*da, p)
	}
}

// padX returns a new DrawArea that is padded horizontally
// so that glyphs will no be clipped.
func padX(p *Plot, da *DrawArea) *DrawArea {
	glyphs := p.GlyphBoxes(p)
	l := leftMost(da, glyphs)
	xAxis := horizontalAxis{p.X}
	glyphs = append(glyphs, xAxis.GlyphBoxes(p)...)
	r := rightMost(da, glyphs)

	minx := da.Min.X + (da.Min.X - (da.X(l.X) + l.Min.X))
	maxx := da.Max().X - ((da.X(r.X) + r.Min.X + r.Size.X) - da.Max().X)
	lx := vg.Length(l.X)
	rx := vg.Length(r.X)
	n := (lx*maxx - rx*minx) / (lx - rx)
	m := ((lx-1)*maxx - rx*minx + minx) / (lx - rx)
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		Rect: Rect{
			Min:  Point{X: n, Y: da.Min.Y},
			Size: Point{X: m - n, Y: da.Size.Y},
		},
	}
}

// rightMost returns the right-most GlyphBox.
func rightMost(da *DrawArea, boxes []GlyphBox) GlyphBox {
	maxx := da.Max().X
	r := GlyphBox{X: 1}
	for _, b := range boxes {
		if x := da.X(b.X) + b.Min.X + b.Size.X; x > maxx && b.X <= 1 {
			maxx = x
			r = b
		}
	}
	return r
}

// leftMost returns the left-most GlyphBox.
func leftMost(da *DrawArea, boxes []GlyphBox) GlyphBox {
	minx := da.Min.X
	l := GlyphBox{}
	for _, b := range boxes {
		if x := da.X(b.X) + b.Min.X; x < minx && b.X >= 0 {
			minx = x
			l = b
		}
	}
	return l
}

// padY returns a new DrawArea that is padded vertically
// so that glyphs will no be clipped.
func padY(p *Plot, da *DrawArea) *DrawArea {
	glyphs := p.GlyphBoxes(p)
	b := bottomMost(da, glyphs)
	yAxis := verticalAxis{p.Y}
	glyphs = append(glyphs, yAxis.GlyphBoxes(p)...)
	t := topMost(da, glyphs)

	miny := da.Min.Y + (da.Min.Y - (da.Y(b.Y) + b.Min.Y))
	maxy := da.Max().Y - ((da.Y(t.Y) + t.Min.Y + t.Size.Y) - da.Max().Y)
	by := vg.Length(b.Y)
	ty := vg.Length(t.Y)
	n := (by*maxy - ty*miny) / (by - ty)
	m := ((by-1)*maxy - ty*miny + miny) / (by - ty)
	return &DrawArea{
		vg.Canvas: vg.Canvas(da),
		Rect: Rect{
			Min:  Point{Y: n, X: da.Min.X},
			Size: Point{Y: m - n, X: da.Size.X},
		},
	}
}

// topMost returns the top-most GlyphBox.
func topMost(da *DrawArea, boxes []GlyphBox) GlyphBox {
	maxy := da.Max().Y
	t := GlyphBox{Y: 1}
	for _, b := range boxes {
		if y := da.Y(b.Y) + b.Min.Y + b.Size.Y; y > maxy && b.Y <= 1 {
			maxy = y
			t = b
		}
	}
	return t
}

// bottomMost returns the bottom-most GlyphBox.
func bottomMost(da *DrawArea, boxes []GlyphBox) GlyphBox {
	miny := da.Min.Y
	l := GlyphBox{}
	for _, b := range boxes {
		if y := da.Y(b.Y) + b.Min.Y; y < miny && b.Y >= 0 {
			miny = y
			l = b
		}
	}
	return l
}

// glyphBoxer wraps the GlyphBoxes method.
// It should be implemented by things that meet
// the Data interface that draw glyphs so that
// their glyphs are not clipped if drawn near the
// edge of the DrawArea.
type glyphBoxer interface {
	GlyphBoxes(*Plot) []GlyphBox
}

// A GlyphBox describes the location of a glyph
// and the offset/size of its bounding box.
type GlyphBox struct {
	// The glyph location in normalized coordinates.
	X, Y float64

	// Rect is the offset of the glyph's minimum drawing
	// point relative to the glyph location and its size.
	Rect
}

// GlyphBoxes returns the GlyphBoxes for all plot
// data that meet the glyphBoxer interface.
func (p *Plot) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, d := range p.data {
		if gb, ok := d.(glyphBoxer); ok {
			boxes = append(boxes, gb.GlyphBoxes(p)...)
		}
	}
	return
}
