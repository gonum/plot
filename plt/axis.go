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

	// Label is the axis label.  If the label is the empty string
	//  then no label is dislayed.
	Label string

	// LabelStyle is the text style of the label on the axis.
	LabelStyle TextStyle

	// AxisStyle is the style of the axis's line.
	AxisStyle LineStyle

	// Padding between the axis line and the data in inches.
	Padding vg.Length

	// Ticks are the tick marks on the axis.
	Ticks tickMarks
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
	return Axis{
		Min:   math.Inf(1),
		Max:   math.Inf(-1),
		Label: "",
		LabelStyle: TextStyle{
			Color: Black,
			Font:  labelFont,
		},
		AxisStyle: LineStyle{
			Color: Black,
			Width: vg.Inches(1.0 / 64.0),
		},
		Padding: vg.Inches(1.0 / 8.0),
		Ticks:   maketickMarks(),
	}
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
	if a.Label != "" {
		h += a.LabelStyle.Font.Extents().Height
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		h += a.Ticks.Length + a.Ticks.labelHeight(marks)
	}
	h += a.AxisStyle.Width / 2
	h += a.Padding
	return
}

// draw draws the axis onto the given area.
func (a *horizontalAxis) draw(da *drawArea) {
	y := da.min.y
	if a.Label != "" {
		da.setTextStyle(a.LabelStyle)
		y -= a.LabelStyle.Font.Extents().Descent
		da.text(da.center().x, y, -0.5, 0, a.Label)
		y += a.LabelStyle.Font.Extents().Ascent
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Ticks.MarkStyle)
		da.setTextStyle(a.Ticks.LabelStyle)
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.text(a.x(da, t.Value), y, -0.5, 0, t.Label)
		}
		y += a.Ticks.labelHeight(marks)

		len := a.Ticks.Length
		for _, t := range marks {
			x := a.x(da, t.Value)
			da.line([]point{{x, y + t.lengthOffset(len)}, {x, y + len}})
		}
		y += len
	}
	da.setLineStyle(a.AxisStyle)
	da.line([]point{{da.min.x, y}, {da.max().x, y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *horizontalAxis) glyphBoxes() (boxes []glyphBox) {
	for _, t := range a.Ticks.marks(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		w := a.Ticks.LabelStyle.Font.Width(t.Label)
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
	if a.Label != "" {
		w += a.LabelStyle.Font.Extents().Ascent
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := a.Ticks.labelWidth(marks); lwidth > 0 {
			w += lwidth
			// Add a space after tick labels to separate
			// them from the tick marks
			w += a.Ticks.LabelStyle.Font.Width(" ")
		}
		w += a.Ticks.Length
	}
	w += a.AxisStyle.Width / 2
	w += a.Padding
	return
}

// draw draws the axis onto the given area.
func (a *verticalAxis) draw(da *drawArea) {
	x := da.min.x
	if a.Label != "" {
		x += a.LabelStyle.Font.Extents().Ascent
		da.setTextStyle(a.LabelStyle)
		da.Push()
		da.Rotate(math.Pi / 2)
		da.text(da.center().y, -x, -0.5, 0, a.Label)
		da.Pop()
		x += -a.LabelStyle.Font.Extents().Descent
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Ticks.MarkStyle)
		da.setTextStyle(a.Ticks.LabelStyle)
		if lwidth := a.Ticks.labelWidth(marks); lwidth > 0 {
			x += lwidth
			x += a.Ticks.LabelStyle.Font.Width(" ")
		}
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.text(x, a.y(da, t.Value), -1, -0.5, t.Label+" ")
		}
		len := a.Ticks.Length
		for _, t := range marks {
			y := a.y(da, t.Value)
			da.line([]point{{x + t.lengthOffset(len), y}, {x + len, y}})
		}
		x += len
	}
	da.setLineStyle(a.AxisStyle)
	da.line([]point{{x, da.min.y}, {x, da.max().y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *verticalAxis) glyphBoxes() (boxes []glyphBox) {
	h := a.Ticks.LabelStyle.Font.Extents().Height
	for _, t := range a.Ticks.marks(a.Min, a.Max) {
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

// tickMarks specifies the style and location of the tick marks
// on an axis.
type tickMarks struct {
	// LabelStyle is the TextStyle on the tick labels.
	LabelStyle TextStyle

	// MarkStyle is the LineStyle of the tick mark lines.
	MarkStyle LineStyle

	// Length is the length of a major tick mark in inches.
	// Minor tick marks are half of the length of major
	// tick marks.
	Length vg.Length

	// TickMarker locates the tick marks given the
	// minimum and maximum values.
	TickMarker
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

// maketickMarks returns a tickMarks using the default style
// and TickMarker.
func maketickMarks() tickMarks {
	labelFont, err := MakeFont(DefaultFont, vg.Points(10))
	if err != nil {
		panic(err)
	}
	return tickMarks{
		LabelStyle: TextStyle{
			Color: color.RGBA{A: 255},
			Font:  labelFont,
		},
		MarkStyle: LineStyle{
			Color: color.RGBA{A: 255},
			Width: vg.Inches(1.0 / 64.0),
		},
		Length:     vg.Inches(1.0 / 10.0),
		TickMarker: DefaultTicks(struct{}{}),
	}
}

// labelHeight returns the label height.
func (tick tickMarks) labelHeight(ticks []Tick) vg.Length {
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		font := tick.LabelStyle.Font
		return font.Extents().Ascent
	}
	return 0
}

// labelWidth returns the label width.
func (tick tickMarks) labelWidth(ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		w := tick.LabelStyle.Font.Width(t.Label)
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

// A ConstTicks always returns the same set of tick marks.
type ConstantTicks []Tick

// Marks implements the TickMarker Marks method.
func (tks ConstantTicks) marks(min, max float64) []Tick {
	return tks
}
