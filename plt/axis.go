package plt

import (
	"code.google.com/p/plotinum/vg"
	"fmt"
	"image/color"
	"math"
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
	labelFont, err := vg.MakeFont(defaultFont, vg.Points(12))
	if err != nil {
		panic(err)
	}
	a := Axis{
		Min: math.Inf(1),
		Max: math.Inf(-1),
		LineStyle: LineStyle{
			Color: color.Black,
			Width: vg.Points(1),
		},
		Padding: vg.Points(5),
	}

	a.Label.TextStyle = TextStyle{
		Color: color.Black,
		Font:  labelFont,
	}

	tickFont, err := vg.MakeFont(defaultFont, vg.Points(10))
	if err != nil {
		panic(err)
	}
	a.Tick.Label = TextStyle{
		Color: color.Black,
		Font:  tickFont,
	}
	a.Tick.LineStyle = LineStyle{
		Color: color.Black,
		Width: vg.Points(1),
	}
	a.Tick.Length = vg.Points(8)
	a.Tick.Marker = DefaultTicks

	return a
}

// Norm return the value of x, given in the data coordinate
// system, normalized to its distance as a fraction of the
// range of this axis.  For example, if x is a.Min then the return
// value is 0, and if x is a.Max then the return value is 1.
func (a *Axis) Norm(x float64) float64 {
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
		h -= a.Label.Font.Extents().Descent
		h += a.Label.Height(a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		h += a.Tick.Length + tickLabelHeight(a.Tick.Label, marks)
	}
	h += a.Width / 2
	h += a.Padding
	return
}

// draw draws the axis along the lower edge of the given DrawArea.
func (a *horizontalAxis) draw(da *DrawArea) {
	y := da.Min.Y
	if a.Label.Text != "" {
		y -= a.Label.Font.Extents().Descent
		da.FillText(a.Label.TextStyle, da.Center().X, y, -0.5, 0, a.Label.Text)
		y += a.Label.Height(a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.FillText(a.Tick.Label, da.X(a.Norm(t.Value)), y, -0.5, 0, t.Label)
		}
		y += tickLabelHeight(a.Tick.Label, marks)

		len := a.Tick.Length
		for _, t := range marks {
			x := da.X(a.Norm(t.Value))
			da.StrokeLine2(a.Tick.LineStyle, x, y+t.lengthOffset(len), x, y+len)
		}
		y += len
	}
	da.StrokeLine2(a.LineStyle, da.Min.X, y, da.Max().X, y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels.
func (a *horizontalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		w := a.Tick.Label.Width(t.Label)
		box := GlyphBox{
			X:    a.Norm(t.Value),
			Rect: Rect{Point{X: -w / 2}, Point{X: w}},
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
		w -= a.Label.Font.Extents().Descent
		w += a.Label.Height(a.Label.Text)
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Tick.Label, marks); lwidth > 0 {
			w += lwidth
			w += a.Label.Width(" ")
		}
		w += a.Tick.Length
	}
	w += a.Width / 2
	w += a.Padding
	return
}

// draw draws the axis along the left side of the DrawArea.
func (a *verticalAxis) draw(da *DrawArea) {
	x := da.Min.X
	if a.Label.Text != "" {
		x += a.Label.Height(a.Label.Text)
		da.Push()
		da.Rotate(math.Pi / 2)
		da.FillText(a.Label.TextStyle, da.Center().Y, -x, -0.5, 0, a.Label.Text)
		da.Pop()
		x += -a.Label.Font.Extents().Descent
	}
	marks := a.Tick.Marker(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Tick.Label, marks); lwidth > 0 {
			x += lwidth
		}
		major := false
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.FillText(a.Tick.Label, x, da.Y(a.Norm(t.Value)), -1, -0.5, t.Label)
			major = true
		}
		if major {
			x += a.Tick.Label.Width(" ")
		}
		len := a.Tick.Length
		for _, t := range marks {
			y := da.Y(a.Norm(t.Value))
			da.StrokeLine2(a.Tick.LineStyle, x+t.lengthOffset(len), y, x+len, y)
		}
		x += len
	}
	da.StrokeLine2(a.LineStyle, x, da.Min.Y, x, da.Max().Y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels
func (a *verticalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		h := a.Tick.Label.Height(t.Label)
		box := GlyphBox{
			Y:    a.Norm(t.Value),
			Rect: Rect{Point{Y: -h / 2}, Point{Y: h}},
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
func ConstantTicks(ts []Tick) func(float64, float64) []Tick {
	return func(float64, float64) []Tick {
		return ts
	}
}

// A Tick is a single tick mark on an axis.
type Tick struct {
	// Value is the value denoted by the tick.
	Value float64
	// Label is the text to display at the tick mark.
	// If Label is an empty string then this is a minor
	// tick mark.
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
func tickLabelHeight(sty TextStyle, ticks []Tick) vg.Length {
	maxHeight := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		h := sty.Height(t.Label)
		if h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}

// tickLabelWidth returns the width of the widest tick mark label.
func tickLabelWidth(sty TextStyle, ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		w := sty.Width(t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}
