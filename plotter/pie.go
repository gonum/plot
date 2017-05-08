// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

const (
	defaultPieChartRadius        = 0.8
	defaultPieChartLabelPosition = 0.8
	defaultPieChartFontFamily    = "Times-Roman"
	defaultPieChartFontSize      = 12
)

var defaultPieChartFont vg.Font

func init() {
	var err error
	defaultPieChartFont, err = vg.MakeFont(defaultPieChartFontFamily, defaultPieChartFontSize)
	if err != nil {
		panic(err)
	}
}

// PieChart presents data values as slices in a pie, with area
// proportional to the data values.
type PieChart struct {
	Values
	sumValues float64

	// Radius is the radius of the pie with respect to
	// the chart size. If 0, the pie is a simple dot, if 1, it fills the chart
	// (given that the axes range from -1 to 1).
	Radius float64

	// Color is the fill color of the pie slices.
	Color color.Color

	// LineStyle is the style of the outline of the pie slices.
	// In particular, LineStyle.Width may be set to change the spacing
	// between slices (defaults to vg.Lengch(2)).
	draw.LineStyle

	Offset struct {
		// Value is the value added to the position of each slice in the pie,
		// as if there was an invisible slice with value Value before.
		// When this Value offset is zero (default), the slices are drawn one after
		// the other, counter-clockwise, starting from the right side.
		Value float64

		// X is the length added to the X position of the pie.
		// When the X offset is zero, the pie is centered horizontally..
		X vg.Length

		// Y is the length added to the Y position the pie.
		// When the Y offset is zero, the pie is centered vertically..
		Y vg.Length
	}

	// Total is the total of values required to fill the pie, it defaults to
	// the sum of all Values.
	// If changed to be greater than this sum, the slices will not fill the pie.
	// If less, it will be reset to the default and fill the pie (slices
	// can not overlay each other).
	Total float64

	Labels struct {
		// Show (default: true) is used to determine whether labels should be
		// displayed in each slice of the pie chart
		Show bool

		// Position is the location of the label:
		// 0 is for the center of the pie;
		// 1 is for the side of the slice.
		// A value greater than 1 is outside of the pie (default: 0.8)
		Position float64

		// TextStyle is the style of the slices label text.
		draw.TextStyle

		// Nominal can be used to override the displayed labels (if not used,
		// regular indexes are displayed).
		// Useful in particular in combination with Offset.Value, to offset
		// manually the displayed indexes.
		Nominal []string

		Values struct {
			// Show (default: false) sets whether values should be displayed
			// along with indexes.
			Show bool

			// Percentage (default: false) is to display the percent value
			// of each slice instead of the real value.
			// The percentage is computed with respect to the given Total
			// value (if any).
			Percentage bool
		}
	}
}

// NewPieChart returns a new pie chart with a single slice for each value.
// The slice areas correspond to the values.
//
// Since this chart is not related to axes, this should probably be used
// in conjunction with p.HideAxes().
func NewPieChart(vs Valuer) (*PieChart, error) {
	values, err := CopyValues(vs)
	if err != nil {
		return nil, err
	}

	pie := &PieChart{
		Values:    values,
		Radius:    defaultPieChartRadius,
		sumValues: sumValues(values),
	}
	// default: separate slices with width-2 white lines
	pie.LineStyle.Color = color.White
	pie.LineStyle.Width = vg.Length(2)
	// default: display values indexes as labels
	pie.Labels.Show = true
	pie.Labels.Position = defaultPieChartLabelPosition
	pie.Labels.TextStyle.XAlign = draw.XCenter
	pie.Labels.TextStyle.YAlign = draw.YCenter
	pie.Labels.TextStyle.Font = defaultPieChartFont

	return pie, nil
}

// Plot implements the plot.Plotter interface.
func (p *PieChart) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	// initialize drawing parameters
	x0 := trX(0) + p.Offset.X
	y0 := trY(0) + p.Offset.Y
	origin := vg.Point{X: x0, Y: y0}
	radius := vg.Length(math.Min(float64(trX(1)-trX(0)), float64(trY(1)-trY(0))) * // take the smallest of x-y scale
		math.Max(0, p.Radius)) // and scale to stored Radius
	labelOrigin := origin.Add(vg.Point{X: radius * vg.Length(p.Labels.Position), Y: 0})
	totalValues := math.Max(p.Total, p.sumValues)
	startAngle := computeValueAngle(p.Offset.Value, totalValues)

	// draw all slices
	var labelPoint vg.Point
	c.SetLineStyle(p.LineStyle)
	for i, v := range p.Values {
		if v <= 0.01 {
			continue
		}
		arcAngle := computeValueAngle(v, totalValues)

		var path vg.Path
		path.Move(origin)                              // move to center
		path.Arc(origin, radius, startAngle, 0)        // move to start of arc
		path.Arc(origin, radius, startAngle, arcAngle) // trace arc
		path.Close()                                   // close path to complete slice

		// fill slice
		c.SetColor(p.Color)
		c.Fill(path)

		// stroke slice
		c.SetColor(p.LineStyle.Color)
		c.Stroke(path)

		// write label at mid-arcAngle
		labelText := ""
		if p.Labels.Nominal == nil {
			labelText = strconv.Itoa(i + 1)
		} else if i < len(p.Labels.Nominal) {
			labelText = p.Labels.Nominal[i]
		}
		if labelText != "" {
			if p.Labels.Values.Show {
				labelText += ": "
				if p.Labels.Values.Percentage {
					labelText += strconv.Itoa(int((v*100/totalValues)+0.5)) + "%"
				} else {
					labelText += strings.TrimRight(strconv.FormatFloat(v, 'f', 2, 64), "0.")
				}
			}
			labelPoint = labelOrigin.Rotate(origin, startAngle+arcAngle/2)
			c.FillText(p.Labels.TextStyle, labelPoint, labelText)
		}

		// next
		startAngle += arcAngle
	}
}

func sumValues(vs Values) float64 {
	sum := 0.
	for _, v := range vs {
		sum += v
	}
	return sum
}

func computeValueAngle(v, total float64) float64 {
	return 2 * math.Pi * v / total
}

// DataRange implements the plot.DataRanger interface.
//
// A pie chart is always defined in the range (-1, -1) to (1, 1).
// If something different is required, change the X/YOffset or Radius
// attributes of the PieChart to move or resize it.
func (p *PieChart) DataRange() (float64, float64, float64, float64) {
	return -1, 1, -1, 1
}

// Thumbnail fulfills the plot.Thumbnailer interface.
func (p *PieChart) Thumbnail(c *draw.Canvas) {
	pts := []vg.Point{
		{X: c.Min.X, Y: c.Min.Y},
		{X: c.Min.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Max.Y},
		{X: c.Max.X, Y: c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(p.Color, poly)

	// force LineStyle.Width to 1, if greater, to preserve readability
	lw := p.LineStyle.Width
	if lw > 1 {
		p.LineStyle.Width = 1
	}
	pts = append(pts, vg.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	c.StrokeLines(p.LineStyle, outline...)
	p.LineStyle.Width = lw
}
