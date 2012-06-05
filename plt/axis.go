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

	// Padding between the axis line and the data in inches.
	Padding vg.Length

	Tick struct {
		// Label is the TextStyle on the tick labels.
		Label TextStyle

		// LineStyle is the LineStyle of the tick mark lines.
		LineStyle

		// Length is the length of a major tick mark in inches.
		// Minor tick marks are half of the length of major
		// tick marks.
		Length vg.Length

		// TickMarker locates the tick marks given the
		// minimum and maximum values.
		TickMarker
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
	a.Tick.TickMarker = DefaultTicks(struct{}{})

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

// size returns the height of the axis in inches.
func (a *horizontalAxis) size() (h vg.Length) {
	if a.Label.Text != "" {
		h += a.Label.Font.Extents().Height
	}
	marks := a.Tick.marks(a.Min, a.Max)
	if len(marks) > 0 {
		h += a.Tick.Length + tickLabelHeight(a.Tick.Label.Font, marks)
	}
	h += a.Width / 2
	h += a.Padding
	return
}

// draw draws the axis onto the given area.
func (a *horizontalAxis) draw(da *drawArea) {
	y := da.min.y
	if a.Label.Text != "" {
		da.setTextStyle(a.Label.TextStyle)
		y -= a.Label.Font.Extents().Descent
		da.text(da.center().x, y, -0.5, 0, a.Label.Text)
		y += a.Label.Font.Extents().Ascent
	}
	marks := a.Tick.marks(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Tick.LineStyle)
		da.setTextStyle(a.Tick.Label)
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.text(a.x(da, t.Value), y, -0.5, 0, t.Label)
		}
		y += tickLabelHeight(a.Tick.Label.Font, marks)

		len := a.Tick.Length
		for _, t := range marks {
			x := a.x(da, t.Value)
			da.line([]point{{x, y + t.lengthOffset(len)}, {x, y + len}})
		}
		y += len
	}
	da.setLineStyle(a.LineStyle)
	da.line([]point{{da.min.x, y}, {da.max().x, y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *horizontalAxis) glyphBoxes() (boxes []glyphBox) {
	for _, t := range a.Tick.marks(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		w := a.Tick.Label.Font.Width(t.Label)
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

// size returns the width of the axis in inches.
func (a *verticalAxis) size() (w vg.Length) {
	if a.Label.Text != "" {
		w += a.Label.Font.Extents().Ascent
	}
	marks := a.Tick.marks(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Tick.Label.Font, marks); lwidth > 0 {
			w += lwidth
			// Add a space after tick labels to separate
			// them from the tick marks
			w += a.Tick.Label.Font.Width(" ")
		}
		w += a.Tick.Length
	}
	w += a.Width / 2
	w += a.Padding
	return
}

// draw draws the axis onto the given area.
func (a *verticalAxis) draw(da *drawArea) {
	x := da.min.x
	if a.Label.Text != "" {
		x += a.Label.Font.Extents().Ascent
		da.setTextStyle(a.Label.TextStyle)
		da.Push()
		da.Rotate(math.Pi / 2)
		da.text(da.center().y, -x, -0.5, 0, a.Label.Text)
		da.Pop()
		x += -a.Label.Font.Extents().Descent
	}
	marks := a.Tick.marks(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Tick.LineStyle)
		da.setTextStyle(a.Tick.Label)
		if lwidth := tickLabelWidth(a.Tick.Label.Font, marks); lwidth > 0 {
			x += lwidth
			x += a.Tick.Label.Font.Width(" ")
		}
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.text(x, a.y(da, t.Value), -1, -0.5, t.Label+" ")
		}
		len := a.Tick.Length
		for _, t := range marks {
			y := a.y(da, t.Value)
			da.line([]point{{x + t.lengthOffset(len), y}, {x + len, y}})
		}
		x += len
	}
	da.setLineStyle(a.LineStyle)
	da.line([]point{{x, da.min.y}, {x, da.max().y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *verticalAxis) glyphBoxes() (boxes []glyphBox) {
	h := a.Tick.Label.Font.Extents().Height
	for _, t := range a.Tick.marks(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		box := glyphBox{
			y:    a.norm(t.Value),
			rect: rect{min: point{y: -h / 2}, size: point{y: h}},
		}
		boxes = append(boxes, box)
	}
	return
}

// A TickMarker returns a slice of ticks between a given
// range of values. 
type TickMarker interface {
	// marks returns a slice of ticks for the given range.
	marks(min, max float64) []Tick
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

// tickLabelHeight returns the label height.
func tickLabelHeight(f vg.Font, ticks []Tick) vg.Length {
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		return f.Extents().Ascent
	}
	return 0
}

// tickLabelWidth returns the label width.
func tickLabelWidth(f vg.Font, ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		w := f.Width(t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

// A DefaultTicks returns a default set of tick marks within
// the given range.
type DefaultTicks struct{}

// Marks implements the TickMarker Marks method.
func (_ DefaultTicks) marks(min, max float64) []Tick {
	return []Tick{
		{Value: min, Label: fmt.Sprintf("%g", min)},
		{Value: min + (max-min)/4},
		{Value: min + (max-min)/2, Label: fmt.Sprintf("%g", min+(max-min)/2)},
		{Value: min + 3*(max-min)/4},
		{Value: max, Label: fmt.Sprintf("%g", max)},
	}
}

// A ConstantTicks always returns the same set of tick marks.
type ConstantTicks []Tick

// Marks implements the TickMarker Marks method.
func (tks ConstantTicks) marks(min, max float64) []Tick {
	return tks
}
