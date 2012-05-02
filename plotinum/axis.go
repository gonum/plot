package plotinum

import (
	"plotinum/vecgfx"
	"image/color"
	"fmt"
)

const (
	DefaultFont = "Times-Roman"
)

type Axis struct{
	// Min and Max are the minimum and maximum data
	// coordinates on this axis.
	Min, Max float64

	// Label is the axis label
	Label string

	// LabelStyle is the text style of the label on the axis.
	LabelStyle TextStyle

	/// AxisStyle is the style of the axis's line.
	AxisStyle LineStyle

	// Tick are the tick marks on the axis.
	Tick TickMarks
}

// MakeAxis returns a default axis with the given data dimensions.
func MakeAxis(min, max float64) Axis {
	labelFont, err := MakeFont(DefaultFont, 12)
	if err != nil {
		panic(err)
	}
	return Axis{
		Min: min,
		Max: max,
		Label: "",
		LabelStyle: TextStyle{
			Color: Black,
			Font: labelFont,
		},
		AxisStyle: LineStyle{
			Color: Black,
			Width: 1.0/64.0,
		},
		Tick: MakeTickMarks(min, max),
	}
}

// X transfroms the data point x to the drawing coordinate
// for the given drawing area.
func (a *Axis) X(da *DrawArea, x float64) float64 {
	p := (x - a.Min) / (a.Max - a.Min)
	return da.Min.X + p*(da.Max.X - da.Min.X)
}

// Y transforms the data point y to the drawing coordinate
// for the given drawing area.
func (a *Axis) Y(da *DrawArea, y float64) float64 {
	p := (y - a.Min) / (a.Max - a.Min)
	return da.Min.Y + p*(da.Max.Y - da.Min.Y)
}

// Height returns the height of the axis in inches
//  if it is drawn as a horizontal axis.
func (a *Axis) Height() (ht float64) {
	if a.Label != "" {
		ht += a.LabelStyle.Font.Extents().Height/vecgfx.PtInch
	}
	if len(a.Tick.Marks) > 0 {
		ht += a.Tick.Length + a.Tick.labelHeight()
	}
	ht += a.AxisStyle.Width/2
	return
}

// DrawHoriz draws the axis onto the given area.
func (a *Axis) DrawHoriz(da *DrawArea, c vecgfx.Canvas) {
	y := da.Min.Y
	if a.Label != "" {
		da.SetTextStyle(a.LabelStyle)
		y += -da.FontDescent()
		da.Text(da.Center().X, y, -0.5, 0, a.Label)
		y += da.FontAscent()
	}
	if len(a.Tick.Marks) > 0 {
		da.SetLineStyle(a.Tick.MarkStyle)
		da.SetTextStyle(a.Tick.LabelStyle)
		for _, t := range a.Tick.Marks {
			if t.Label == "" {
				continue
			}
			da.Text(a.X(da, t.Value), y, -0.5, 0, t.Label)
		}
		y += a.Tick.labelHeight() * da.DPI()

		len := a.Tick.Length*da.DPI()
		for _, t := range a.Tick.Marks {
			x := a.X(da, t.Value)
			y1 := y
			if t.Label == "" {
				y1 = y +  len/2
			}
			da.Line(Line{{x, y1}, {x, y + len}})
		}
		y += len
	}
	da.SetLineStyle(a.AxisStyle)
	da.Line(Line{{da.Min.X, y}, {da.Max.X, y}})
}

// TickMarks is the style and location of a set of tick marks.
type TickMarks struct {
	// LabelStyle is the text style on the tick labels.
	LabelStyle TextStyle

	// MarkStyle is the style of the tick mark lines.
	MarkStyle LineStyle

	// Length is the length of a major tick mark specified
	// in inches.
	Length float64

	// Marks is a slice of the locations of the tick marks.
	Marks []Tick
}

func MakeTickMarks(min, max float64) TickMarks {
	labelFont, err := MakeFont(DefaultFont, 10)
	if err != nil {
		panic(err)
	}
	return TickMarks{
		LabelStyle: TextStyle{
			Color: color.RGBA{A: 255},
			Font: labelFont,
		},
		MarkStyle: LineStyle{
			Color: color.RGBA{A:255},
			Width: 1.0/64.0,
		},
		Length: 1.0/10.0,
		Marks: []Tick{
			{ Value: min, Label: fmt.Sprintf("%g", min) },
			{ Value: min + (max-min)/2 },
			{ Value: max, Label: fmt.Sprintf("%g", max) },
		},
	}
}

// labelHeight returns the label height in inches.
func (tick TickMarks) labelHeight() float64 {
	for _, t := range tick.Marks {
		if t.Label == "" {
			continue
		}
		font := tick.LabelStyle.Font
		return font.Extents().Ascent/vecgfx.PtInch
	}
	return 0
}

// A Tick is a single tick mark
type Tick struct {
	Value float64
	Label string
}
