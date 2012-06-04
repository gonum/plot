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
	Padding float64

	// Ticks are the tick marks on the axis.
	Ticks TickMarks
}

// makeAxis returns a default Axis.
//
// The default range is (∞, ­∞), and thus any finite
// value is less than Min and greater than Max.
func makeAxis() Axis {
	labelFont, err := MakeFont(DefaultFont, 12)
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
			Width: 1.0 / 64.0,
		},
		Padding: 1.0 / 8.0,
		Ticks:   makeTickMarks(),
	}
}

// X transfroms the data point x to the drawing coordinate
// for the given drawing area.
func (a *Axis) x(da *DrawArea, x float64) float64 {
	return da.x(a.norm(x))
}

// Y transforms the data point y to the drawing coordinate
// for the given drawing area.
func (a *Axis) y(da *DrawArea, y float64) float64 {
	return da.y(a.norm(y))
}

// norm return the value of x, given in the data coordinate
// system, normalized to its distance as a fraction of the
// range of this axis.  For example, if x is a.Min then the return
// value is 0, and if x is a.Max then the return value is 1.
func (a *Axis) norm(x float64) float64 {
	return (x-a.Min)/(a.Max-a.Min)
}

// A HorizantalAxis draws horizontally across the bottom
// of a plot.
type HorizontalAxis struct {
	Axis
}

// size returns the height of the axis in inches.
func (a *HorizontalAxis) size() (h float64) {
	if a.Label != "" {
		h += a.LabelStyle.Font.Extents().Height / vg.PtInch
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
func (a *HorizontalAxis) draw(da *DrawArea) {
	y := da.Min.Y
	if a.Label != "" {
		da.setTextStyle(a.LabelStyle)
		y += -(a.LabelStyle.Font.Extents().Descent / vg.PtInch * da.DPI())
		da.text(da.center().X, y, -0.5, 0, a.Label)
		y += a.LabelStyle.Font.Extents().Ascent / vg.PtInch * da.DPI()
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
		y += a.Ticks.labelHeight(marks) * da.DPI()

		len := a.Ticks.Length * da.DPI()
		for _, t := range marks {
			x := a.x(da, t.Value)
			da.line([]Point{{x, y + t.lengthOffset(len)}, {x, y + len}})
		}
		y += len
	}
	da.setLineStyle(a.AxisStyle)
	da.line([]Point{{da.Min.X, y}, {da.Max().X, y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *HorizontalAxis) glyphBoxes() (boxes []glyphBox) {
	for _, t := range a.Ticks.marks(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		w := a.Ticks.LabelStyle.Font.Width(t.Label)/vg.PtInch
		box := glyphBox{
			Point: Point{ X: a.norm(t.Value) },
			Rect: Rect{ Min: Point{ X: -w/2 }, Size: Point{ X: w } },
		}
		boxes = append(boxes, box)
	}
	return
}

// A VerticalAxis is drawn vertically up the left side of a plot.
type VerticalAxis struct {
	Axis
}

// size returns the width of the axis in inches.
func (a *VerticalAxis) size() (w float64) {
	if a.Label != "" {
		w += a.LabelStyle.Font.Extents().Ascent / vg.PtInch
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		if lwidth := a.Ticks.labelWidth(marks); lwidth > 0 {
			w += lwidth
			// Add a space after tick labels to separate
			// them from the tick marks
			w += a.Ticks.LabelStyle.Font.Width(" ") / vg.PtInch
		}
		w += a.Ticks.Length
	}
	w += a.AxisStyle.Width / 2
	w += a.Padding
	return
}

// draw draws the axis onto the given area.
func (a *VerticalAxis) draw(da *DrawArea) {
	x := da.Min.X
	if a.Label != "" {
		x += a.LabelStyle.Font.Extents().Ascent / vg.PtInch * da.DPI()
		da.setTextStyle(a.LabelStyle)
		da.Push()
		da.Rotate(math.Pi / 2)
		da.text(da.center().Y, -x, -0.5, 0, a.Label)
		da.Pop()
		x += -a.LabelStyle.Font.Extents().Descent / vg.PtInch * da.DPI()
	}
	marks := a.Ticks.marks(a.Min, a.Max)
	if len(marks) > 0 {
		da.setLineStyle(a.Ticks.MarkStyle)
		da.setTextStyle(a.Ticks.LabelStyle)
		if lwidth := a.Ticks.labelWidth(marks); lwidth > 0 {
			x += lwidth * da.DPI()
			x += a.Ticks.LabelStyle.Font.Width(" ") / vg.PtInch * da.DPI()
		}
		for _, t := range marks {
			if t.minor() {
				continue
			}
			da.text(x, a.y(da, t.Value), -1, -0.5, t.Label+" ")
		}
		len := a.Ticks.Length * da.DPI()
		for _, t := range marks {
			y := a.y(da, t.Value)
			da.line([]Point{{x + t.lengthOffset(len), y}, {x + len, y}})
		}
		x += len
	}
	da.setLineStyle(a.AxisStyle)
	da.line([]Point{{x, da.Min.Y}, {x, da.Max().Y}})
}

// glyphBoxes returns normalized glyphBoxes for the glyphs
// representing the tick mark text.
func (a *VerticalAxis) glyphBoxes() (boxes []glyphBox) {
	h := a.Ticks.LabelStyle.Font.Extents().Height/vg.PtInch
	for _, t := range a.Ticks.marks(a.Min, a.Max) {
		if t.minor() {
			continue
		}
		box := glyphBox{
			Point: Point{ Y: a.norm(t.Value) },
			Rect: Rect{ Min: Point{ Y: -h/2 }, Size: Point{ Y: h } },
		}
		boxes = append(boxes, box)
	}
	return
}

// TickMarks specifies the style and location of the tick marks
// on an axis.
type TickMarks struct {
	// LabelStyle is the TextStyle on the tick labels.
	LabelStyle TextStyle

	// MarkStyle is the LineStyle of the tick mark lines.
	MarkStyle LineStyle

	// Length is the length of a major tick mark in inches.
	// Minor tick marks are half of the length of major
	// tick marks.
	Length float64

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
func (t Tick) lengthOffset(len float64) float64 {
	if t.minor() {
		return len / 2
	}
	return 0
}

// makeTickMarks returns a TickMarks using the default style
// and TickMarker.
func makeTickMarks() TickMarks {
	labelFont, err := MakeFont(DefaultFont, 10)
	if err != nil {
		panic(err)
	}
	return TickMarks{
		LabelStyle: TextStyle{
			Color: color.RGBA{A: 255},
			Font:  labelFont,
		},
		MarkStyle: LineStyle{
			Color: color.RGBA{A: 255},
			Width: 1.0 / 64.0,
		},
		Length:     1.0 / 10.0,
		TickMarker: DefaultTicks(struct{}{}),
	}
}

// labelHeight returns the label height in inches.
func (tick TickMarks) labelHeight(ticks []Tick) float64 {
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		font := tick.LabelStyle.Font
		return font.Extents().Ascent / vg.PtInch
	}
	return 0
}

// labelWidth returns the label width in inches.
func (tick TickMarks) labelWidth(ticks []Tick) float64 {
	maxWidth := 0.0
	for _, t := range ticks {
		if t.minor() {
			continue
		}
		w := tick.LabelStyle.Font.Width(t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth / vg.PtInch
}

// A DefalutTicks returns a default set of tick marks within
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
