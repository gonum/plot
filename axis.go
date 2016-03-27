// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"image/color"
	"math"
	"strconv"

	"github.com/gonum/floats"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// displayPrecision is a sane level of float precision for a plot.
const displayPrecision = 4

// Ticker creates Ticks in a specified range
type Ticker interface {
	// Ticks returns Ticks in a specified range
	Ticks(min, max float64) []Tick
}

// Normalizer rescales values from the data coordinate system to the
// normalized coordinate system.
type Normalizer interface {
	// Normalize transforms a value x in the data coordinate system to
	// the normalized coordinate system.
	Normalize(min, max, x float64) float64
}

// An Axis represents either a horizontal or vertical
// axis of a plot.
type Axis struct {
	// Min and Max are the minimum and maximum data
	// values represented by the axis.
	Min, Max float64

	Label struct {
		// Text is the axis label string.
		Text string

		// TextStyle is the style of the axis label text. For the vertical axis,
		// 90 degrees of rotation will be added to the label text before drawing.
		draw.TextStyle
	}

	// LineStyle is the style of the axis line.
	draw.LineStyle

	// Padding between the axis line and the data.  Having
	// non-zero padding ensures that the data is never drawn
	// on the axis, thus making it easier to see.
	Padding vg.Length

	Tick struct {
		// Label is the TextStyle on the tick labels.
		Label draw.TextStyle

		// LineStyle is the LineStyle of the tick lines.
		draw.LineStyle

		// Length is the length of a major tick mark.
		// Minor tick marks are half of the length of major
		// tick marks.
		Length vg.Length

		// Marker returns the tick marks.  Any tick marks
		// returned by the Marker function that are not in
		// range of the axis are not drawn.
		Marker Ticker
	}

	// Scale transforms a value given in the data coordinate system
	// to the normalized coordinate system of the axis—its distance
	// along the axis as a fraction of the axis range.
	Scale Normalizer
}

// makeAxis returns a default Axis.
//
// The default range is (∞, ­∞), and thus any finite
// value is less than Min and greater than Max.
func makeAxis(vertical bool) (Axis, error) {
	labelFont, err := vg.MakeFont(DefaultFont, vg.Points(12))
	if err != nil {
		return Axis{}, err
	}

	tickFont, err := vg.MakeFont(DefaultFont, vg.Points(10))
	if err != nil {
		return Axis{}, err
	}

	a := Axis{
		Min: math.Inf(1),
		Max: math.Inf(-1),
		LineStyle: draw.LineStyle{
			Color: color.Black,
			Width: vg.Points(0.5),
		},
		Padding: vg.Points(5),
		Scale:   LinearScale{},
	}
	a.Label.TextStyle = draw.TextStyle{
		Color:  color.Black,
		Font:   labelFont,
		XAlign: -0.5,
		YAlign: 0,
	}
	if vertical {
		a.Tick.Label = draw.TextStyle{
			Color:  color.Black,
			Font:   tickFont,
			XAlign: -1,
			YAlign: -0.5,
		}
	} else {
		a.Tick.Label = draw.TextStyle{
			Color:  color.Black,
			Font:   tickFont,
			XAlign: -0.5,
			YAlign: 0,
		}
	}
	a.Tick.LineStyle = draw.LineStyle{
		Color: color.Black,
		Width: vg.Points(0.5),
	}
	a.Tick.Length = vg.Points(8)
	a.Tick.Marker = DefaultTicks{}

	return a, nil
}

// sanitizeRange ensures that the range of the
// axis makes sense.
func (a *Axis) sanitizeRange() {
	if math.IsInf(a.Min, 0) {
		a.Min = 0
	}
	if math.IsInf(a.Max, 0) {
		a.Max = 0
	}
	if a.Min > a.Max {
		a.Min, a.Max = a.Max, a.Min
	}
	if a.Min == a.Max {
		a.Min--
		a.Max++
	}
}

// LinearScale an be used as the value of an Axis.Scale function to
// set the axis to a standard linear scale.
type LinearScale struct{}

var _ Normalizer = LinearScale{}

// Normalize returns the fractional distance of x between min and max.
func (LinearScale) Normalize(min, max, x float64) float64 {
	return (x - min) / (max - min)
}

// LogScale can be used as the value of an Axis.Scale function to
// set the axis to a log scale.
type LogScale struct{}

var _ Normalizer = LogScale{}

// Normalize returns the fractional logarithmic distance of
// x between min and max.
func (LogScale) Normalize(min, max, x float64) float64 {
	logMin := log(min)
	return (log(x) - logMin) / (log(max) - logMin)
}

// Norm returns the value of x, given in the data coordinate
// system, normalized to its distance as a fraction of the
// range of this axis.  For example, if x is a.Min then the return
// value is 0, and if x is a.Max then the return value is 1.
func (a *Axis) Norm(x float64) float64 {
	return a.Scale.Normalize(a.Min, a.Max, x)
}

// drawTicks returns true if the tick marks should be drawn.
func (a *Axis) drawTicks() bool {
	return a.Tick.Width > 0 && a.Tick.Length > 0
}

// A horizontalAxis draws horizontally across the bottom
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
	if marks := a.Tick.Marker.Ticks(a.Min, a.Max); len(marks) > 0 {
		if a.drawTicks() {
			h += a.Tick.Length
		}
		h += tickLabelHeight(a.Tick.Label, marks)
	}
	h += a.Width / 2
	h += a.Padding
	return
}

// draw draws the axis along the lower edge of a draw.Canvas.
func (a *horizontalAxis) draw(c draw.Canvas) {
	y := c.Min.Y
	if a.Label.Text != "" {
		y -= a.Label.Font.Extents().Descent
		c.FillText(a.Label.TextStyle, c.Center().X, y, a.Label.Text)
		y += a.Label.Height(a.Label.Text)
	}

	marks := a.Tick.Marker.Ticks(a.Min, a.Max)
	for _, t := range marks {
		x := c.X(a.Norm(t.Value))
		if !c.ContainsX(x) || t.IsMinor() {
			continue
		}
		c.FillText(a.Tick.Label, x, y, t.Label)
	}

	if len(marks) > 0 {
		y += tickLabelHeight(a.Tick.Label, marks)
	} else {
		y += a.Width / 2
	}

	if len(marks) > 0 && a.drawTicks() {
		len := a.Tick.Length
		for _, t := range marks {
			x := c.X(a.Norm(t.Value))
			if !c.ContainsX(x) {
				continue
			}
			start := t.lengthOffset(len)
			c.StrokeLine2(a.Tick.LineStyle, x, y+start, x, y+len)
		}
		y += len
	}

	c.StrokeLine2(a.LineStyle, c.Min.X, y, c.Max.X, y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels.
func (a *horizontalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker.Ticks(a.Min, a.Max) {
		if t.IsMinor() {
			continue
		}
		w := a.Tick.Label.Width(t.Label)
		box := GlyphBox{
			X: a.Norm(t.Value),
			Rectangle: draw.Rectangle{
				Min: draw.Point{X: -w / 2},
				Max: draw.Point{X: w / 2},
			},
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
	if marks := a.Tick.Marker.Ticks(a.Min, a.Max); len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Tick.Label, marks); lwidth > 0 {
			w += lwidth
			w += a.Label.Width(" ")
		}
		if a.drawTicks() {
			w += a.Tick.Length
		}
	}
	w += a.Width / 2
	w += a.Padding
	return
}

// draw draws the axis along the left side of a draw.Canvas.
func (a *verticalAxis) draw(c draw.Canvas) {
	x := c.Min.X
	if a.Label.Text != "" {
		sty := a.Label.TextStyle
		sty.Rotation += math.Pi / 2
		x += a.Label.Height(a.Label.Text)
		c.FillText(sty, x, c.Center().Y, a.Label.Text)
		x += -a.Label.Font.Extents().Descent
	}
	marks := a.Tick.Marker.Ticks(a.Min, a.Max)
	if w := tickLabelWidth(a.Tick.Label, marks); len(marks) > 0 && w > 0 {
		x += w
	}
	major := false
	for _, t := range marks {
		y := c.Y(a.Norm(t.Value))
		if !c.ContainsY(y) || t.IsMinor() {
			continue
		}
		c.FillText(a.Tick.Label, x, y, t.Label)
		major = true
	}
	if major {
		x += a.Tick.Label.Width(" ")
	}
	if a.drawTicks() && len(marks) > 0 {
		len := a.Tick.Length
		for _, t := range marks {
			y := c.Y(a.Norm(t.Value))
			if !c.ContainsY(y) {
				continue
			}
			start := t.lengthOffset(len)
			c.StrokeLine2(a.Tick.LineStyle, x+start, y, x+len, y)
		}
		x += len
	}
	c.StrokeLine2(a.LineStyle, x, c.Min.Y, x, c.Max.Y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels
func (a *verticalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker.Ticks(a.Min, a.Max) {
		if t.IsMinor() {
			continue
		}
		h := a.Tick.Label.Height(t.Label)
		box := GlyphBox{
			Y: a.Norm(t.Value),
			Rectangle: draw.Rectangle{
				Min: draw.Point{Y: -h / 2},
				Max: draw.Point{Y: h / 2},
			},
		}
		boxes = append(boxes, box)
	}
	return
}

// DefaultTicks is suitable for the Tick.Marker field of an Axis,
// it returns a resonable default set of tick marks.
type DefaultTicks struct{}

var _ Ticker = DefaultTicks{}

// Ticks returns Ticks in a specified range
func (DefaultTicks) Ticks(min, max float64) (ticks []Tick) {
	const SuggestedTicks = 3
	if max < min {
		panic("illegal range")
	}
	tens := math.Pow10(int(math.Floor(math.Log10(max - min))))
	n := (max - min) / tens
	for n < SuggestedTicks {
		tens /= 10
		n = (max - min) / tens
	}

	majorMult := int(n / SuggestedTicks)
	switch majorMult {
	case 7:
		majorMult = 6
	case 9:
		majorMult = 8
	}
	majorDelta := float64(majorMult) * tens
	val := math.Floor(min/majorDelta) * majorDelta
	prec := maxInt(precisionOf(min), precisionOf(max))
	for val <= max {
		if val >= min && val <= max {
			ticks = append(ticks, Tick{Value: val, Label: formatFloatTick(val, prec)})
		}
		if math.Nextafter(val, val+majorDelta) == val {
			break
		}
		val += majorDelta
	}

	minorDelta := majorDelta / 2
	switch majorMult {
	case 3, 6:
		minorDelta = majorDelta / 3
	case 5:
		minorDelta = majorDelta / 5
	}

	val = math.Floor(min/minorDelta) * minorDelta
	for val <= max {
		found := false
		for _, t := range ticks {
			if t.Value == val {
				found = true
			}
		}
		if val >= min && val <= max && !found {
			ticks = append(ticks, Tick{Value: val})
		}
		if math.Nextafter(val, val+minorDelta) == val {
			break
		}
		val += minorDelta
	}
	return
}

// LogTicks is suitable for the Tick.Marker field of an Axis,
// it returns tick marks suitable for a log-scale axis.
type LogTicks struct{}

var _ Ticker = LogTicks{}

// Ticks returns Ticks in a specified range
func (LogTicks) Ticks(min, max float64) []Tick {
	var ticks []Tick
	val := math.Pow10(int(math.Floor(math.Log10(min))))
	if min <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}
	prec := precisionOf(max)
	for val < max*10 {
		for i := 1; i < 10; i++ {
			tick := Tick{Value: val * float64(i)}
			if i == 1 {
				tick.Label = formatFloatTick(val*float64(i), prec)
			}
			ticks = append(ticks, tick)
		}
		val *= 10
	}
	tick := Tick{Value: val, Label: formatFloatTick(val, prec)}
	ticks = append(ticks, tick)
	return ticks
}

// ConstantTicks is suitable for the Tick.Marker field of an Axis.
// This function returns the given set of ticks.
type ConstantTicks []Tick

var _ Ticker = ConstantTicks{}

// Ticks returns Ticks in a specified range
func (ts ConstantTicks) Ticks(float64, float64) []Tick {
	return ts
}

// A Tick is a single tick mark on an axis.
type Tick struct {
	// Value is the data value marked by this Tick.
	Value float64

	// Label is the text to display at the tick mark.
	// If Label is an empty string then this is a minor
	// tick mark.
	Label string
}

// IsMinor returns true if this is a minor tick mark.
func (t Tick) IsMinor() bool {
	return t.Label == ""
}

// lengthOffset returns an offset that should be added to the
// tick mark's line to accout for its length.  I.e., the start of
// the line for a minor tick mark must be shifted by half of
// the length.
func (t Tick) lengthOffset(len vg.Length) vg.Length {
	if t.IsMinor() {
		return len / 2
	}
	return 0
}

// tickLabelHeight returns height of the tick mark labels.
func tickLabelHeight(sty draw.TextStyle, ticks []Tick) vg.Length {
	maxHeight := vg.Length(0)
	for _, t := range ticks {
		if t.IsMinor() {
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
func tickLabelWidth(sty draw.TextStyle, ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.IsMinor() {
			continue
		}
		w := sty.Width(t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

func log(x float64) float64 {
	if x <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}
	return math.Log(x)
}

// formatFloatTick returns a g-formated string representation of v
// to the specified precision.
func formatFloatTick(v float64, prec int) string {
	return strconv.FormatFloat(floats.Round(v, prec), 'g', displayPrecision, 64)
}

// precisionOf returns the precision needed to display x without e notation.
func precisionOf(x float64) int {
	return int(math.Max(math.Ceil(-math.Log10(math.Abs(x))), displayPrecision))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
