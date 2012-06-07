package plt

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
)

const (
	DefaultFont = "Times-Roman"
)

// An Axis represents either a horizontal or vertical
// axis of a plot.
type Axis struct {
	// Min and Max are the minimum and maximum data
	// coordinates on this axis.
	Min, Max float64

	Label struct {
		// Text is the label string.
		Text string
		// TextStyle is the style of the label text.
		TextStyle
	}

	// LineStyle is the style of the axis line.
	LineStyle

	// Padding between the axis line and the data.
	Padding vg.Length

	Tick struct {
		// Label is the TextStyle on the tick labels.
		Label TextStyle

		// LineStyle is the LineStyle of the tick mark lines.
		LineStyle

		// Length is the length of a major tick mark.
		// Minor tick marks are half of the length of major
		// tick marks.
		Length vg.Length

		// Marker returns the tick marks given the minimum
		// and maximum values of the axis.
		Marker func(min, max float64) []Tick
	}
}

// makeAxis returns a default Axis.
//
// The default range is (∞, ­∞), and thus any finite
// value is less than Min and greater than Max.
func makeAxis() Axis {
	labelFont, err := MakeFont(DefaultFont, vg.Points(12))
	if err != nil {
		panic(err)
	}
	a := Axis{
		Min:   math.Inf(1),
		Max:   math.Inf(-1),
		LineStyle: LineStyle{
			Color: Black,
			Width: vg.Inches(1.0 / 64.0),
		},
		Padding: vg.Inches(1.0 / 8.0),
	}

	a.Label.TextStyle = TextStyle{
		Color: Black,
		Font:  labelFont,
	}

	tickFont, err := MakeFont(DefaultFont, vg.Points(10))
	if err != nil {
		panic(err)
	}
	a.Tick.Label = TextStyle{
		Color: color.RGBA{A: 255},
		Font:  tickFont,
	}
	a.Tick.LineStyle = LineStyle{
		Color: color.RGBA{A: 255},
		Width: vg.Inches(1.0 / 64.0),
	}
	a.Tick.Length = vg.Inches(1.0/10.0)
	a.Tick.Marker = DefaultTicks

	return a
}

// X transfroms the data point x to the drawing coordinate
// for the given drawing area.
func (a *Axis) x(da *drawArea, x float64) vg.Length {
	return da.x(a.norm(x))
}

// Y transforms the data point y to the drawing coordinate
// for the given drawing area.
func (a *Axis) y(da *drawArea, y float64) vg.Length {
	return da.y(a.norm(y))
}

// norm return the value of x, given in the data coordinate
// system, normalized to its distance as a fraction of the
// range of this axis.  For example, if x is a.Min then the return
// value is 0, and if x is a.Max then the return value is 1.
func (a *Axis) norm(x float64) float64 {
	return (x - a.Min) / (a.Max - a.Min)
}

// A HorizantalAxis draws horizontally across the bottom
// of a plot.
type horizontalAxis struct {
	Axis
}

// size returns the height of the axis.
func (a *horizontalAxis) size() (h vg.Length) {
	if a.Label.Text != "" {
		h += textHeight(a.Label.Font, a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		h += a.Tick.Length + tickLabelHeight(a.Tick.Label.Font, marks)
	}
	h += a.Width / 2
	h += a.Padding
	return
}

// draw draws the axis along the lower edge of the given drawArea.
func (a *horizontalAxis) draw(da *drawArea) {
	y := da.min.y
	if a.Label.Text != "" {
		da.setTextStyle(a.Label.TextStyle)
		y -= a.Label.Font.Extents().Descent
		fillText(da, da.center().x, y, -0.5, 0, a.Label.Text)
		y += textHeight(a.Label.Font, a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Tick.LineStyle)
		da.setTextStyle(a.Tick.Label)
		for _, t := range marks {
			if t.minor() {
				continue
			}
			fillText(da, a.x(da, t.Value), y, -0.5, 0, t.Label)
		}
		y += tickLabelHeight(a.Tick.Label.Font, marks)

		len := a.Tick.Length
		for _, t := range marks {
			x := a.x(da, t.Value)
			strokeLine2(da, x, y + t.lengthOffset(len),x, y + len)
		}
		y += len
	}
	da.setLineStyle(a.LineStyle)
	strokeLine2(da, da.min.x, y, da.max().x, y)
}

// glyphBoxes returns glyphBoxes for the glyphs
// representing the tick mark labels.  The location
// of the glyphBox is normalized to the unit range
// based on its distance along the axis.
func (a *horizontalAxis) glyphBoxes() (boxes []glyphBox) {
	for _, t := range a.Tick.Marker(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		w := textWidth(a.Tick.Label.Font, t.Label)
		box := glyphBox{
			x:    a.norm(t.Value),
			rect: rect{min: point{x: -w / 2}, size: point{x: w}},
		}
		boxes = append(boxes, box)
	}
	return
}

// A verticalAxis is drawn vertically up the left side of a plot.
type verticalAxis struct {
	Axis
}

// size returns the width of the axis.
func (a *verticalAxis) size() (w vg.Length) {
	if a.Label.Text != "" {
		w += textHeight(a.Label.Font, a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Tick.Label.Font, marks); lwidth > 0 {
			w += lwidth
		}
		w += a.Tick.Length
	}
	w += a.Width / 2
	w += a.Padding
	return
}

// draw draws the axis along the left side of the drawArea.
func (a *verticalAxis) draw(da *drawArea) {
	x := da.min.x
	if a.Label.Text != "" {
		x += textHeight(a.Label.Font, a.Label.Text)
		da.setTextStyle(a.Label.TextStyle)
		da.Push()
		da.Rotate(math.Pi / 2)
		fillText(da, da.center().y, -x, -0.5, 0, a.Label.Text)
		da.Pop()
		x += -a.Label.Font.Extents().Descent
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Tick.LineStyle)
		da.setTextStyle(a.Tick.Label)
		if lwidth := tickLabelWidth(a.Tick.Label.Font, marks); lwidth > 0 {
			x += lwidth
		}
		for _, t := range marks {
			if t.minor() {
				continue
			}
			fillText(da, x, a.y(da, t.Value), -1, -0.5, t.Label)
		}
		len := a.Tick.Length
		for _, t := range marks {
			y := a.y(da, t.Value)
			strokeLine2(da, x + t.lengthOffset(len), y, x + len, y)
		}
		x += len
	}
	da.setLineStyle(a.LineStyle)
	strokeLine2(da, x, da.min.y, x, da.max().y)
}

// glyphBoxes returns glyphBoxes for the glyphs
// representing the tick mark labels.  The location
// of the glyphBox is normalized to the unit range
// based on its distance along the axis.
func (a *verticalAxis) glyphBoxes() (boxes []glyphBox) {
	for _, t := range a.Tick.Marker(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		h := textHeight(a.Tick.Label.Font, t.Label)
		box := glyphBox{
			y:    a.norm(t.Value),
			rect: rect{min: point{y: -h / 2}, size: point{y: h}},
		}
		boxes = append(boxes, box)
	}
	return
}

// DefaultTicks is suitable for the Marker field of an Axis, it returns
// the default set of tick marks.
func DefaultTicks(min, max float64) []Tick {
	return []Tick{
		{Value: min, Label: fmt.Sprintf("%g", min)},
		{Value: min + (max-min)/4},
		{Value: min + (max-min)/2, Label: fmt.Sprintf("%g", min+(max-min)/2)},
		{Value: min + 3*(max-min)/4},
		{Value: max, Label: fmt.Sprintf("%g", max)},
	}
}

// ConstantTicks returns a function suitable for the Marker field
// of an Axis.  This function returns the given set of ticks.
func ConstantTicks(ts []Tick) func(float64,float64)[]Tick {
	return func(float64,float64) []Tick {
		return ts
	}
}

// A Tick is a single tick mark
type Tick struct {
	Value float64
	Label string
}

// minor returns true if this is a minor tick mark.
func (t Tick) minor() bool {
	return t.Label == ""
}

// lengthOffset returns an offset that should be added to the
// tick mark's line to accout for its length.  I.e., the start of
// the line for a minor tick mark must be shifted by half of
// the length.
func (t Tick) lengthOffset(len vg.Length) vg.Length {
	if t.minor() {
		return len / 2
	}
	return 0
}

// tickLabelHeight returns height of the tick mark labels.
func tickLabelHeight(f vg.Font, ticks []Tick) vg.Length {
	maxHeight := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		h := textHeight(f, t.Label)
		if h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}

// tickLabelWidth returns the width of the widest tick mark label.
func tickLabelWidth(f vg.Font, ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		w := textWidth(f, t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}
