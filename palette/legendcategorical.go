// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"fmt"
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	// DefaultLineStyle is the default style for drawing
	// lines.
	DefaultLineStyle = draw.LineStyle{
		Color:    color.Black,
		Width:    vg.Points(1),
		Dashes:   []vg.Length{},
		DashOffs: 0,
	}
)

// Legend is a plot.Plotter that draws a legend for a Palette.
type Legend struct {
	p []color.Color

	// Width is the width of each color rectangle.
	Width vg.Length

	// LineStyle is the style of the outline of the rectangles.
	draw.LineStyle

	// Offset is added to the X location of each rectangle.
	// When the Offset is zero, the rectangles are drawn
	// centered at their X location.
	Offset vg.Length

	// XMin is the X location of the first rectangle.  XMin
	// can be changed to move the legend
	// down the X axis in order to make grouped
	// legends.
	XMin float64

	// Horizontal determines wether the legend will be
	// plotted horizontally or vertically.
	// The default is false (vertical).
	Horizontal bool
}

// NewLegend creates a new legend plotter.
func NewLegend(p Palette, width vg.Length) *Legend {
	return &Legend{
		p:         p.Colors(),
		Width:     width,
		LineStyle: DefaultLineStyle,
	}
}

// Plot implements the Plot method of the plot.Plotter interface.
func (l *Legend) Plot(c draw.Canvas, plt *plot.Plot) {
	trCat, trVal := plt.Transforms(&c)
	if !l.Horizontal {
		trCat, trVal = trVal, trCat
	}

	for i, rectColor := range l.p {
		catVal := l.XMin + float64(i)
		catMin := trCat(float64(catVal))
		if l.Horizontal {
			if !c.ContainsX(catMin) {
				continue
			}
		} else {
			if !c.ContainsY(catMin) {
				continue
			}
		}
		catMin = catMin - l.Width/2 + l.Offset
		catMax := catMin + l.Width
		valMin := trVal(0)
		valMax := trVal(1)

		var pts []vg.Point
		var poly []vg.Point
		if l.Horizontal {
			pts = []vg.Point{
				{X: catMin, Y: valMin},
				{X: catMin, Y: valMax},
				{X: catMax, Y: valMax},
				{X: catMax, Y: valMin},
			}
			poly = c.ClipPolygonY(pts)
		} else {
			pts = []vg.Point{
				{X: valMin, Y: catMin},
				{X: valMin, Y: catMax},
				{X: valMax, Y: catMax},
				{X: valMax, Y: catMin},
			}
			poly = c.ClipPolygonX(pts)
		}
		c.FillPolygon(rectColor, poly)

		var outline [][]vg.Point
		if l.Horizontal {
			pts = append(pts, vg.Point{X: catMin, Y: valMin})
			outline = c.ClipLinesY(pts)
		} else {
			pts = append(pts, vg.Point{X: valMin, Y: catMin})
			outline = c.ClipLinesX(pts)
		}
		c.StrokeLines(l.LineStyle, outline...)
	}
}

// DataRange implements the plot.DataRanger interface.
func (l *Legend) DataRange() (xmin, xmax, ymin, ymax float64) {
	catMin := l.XMin
	catMax := catMin + float64(len(l.p)-1)

	valMin := 0.0
	valMax := 1.0
	if l.Horizontal {
		return catMin, catMax, valMin, valMax
	}
	return valMin, valMax, catMin, catMax
}

// GlyphBoxes implements the GlyphBoxer interface.
func (l *Legend) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(l.p))
	for i := range boxes {
		cat := l.XMin + float64(i)
		if l.Horizontal {
			boxes[i].X = plt.X.Norm(cat)
			boxes[i].Rectangle = vg.Rectangle{
				Min: vg.Point{X: l.Offset - l.Width/2},
				Max: vg.Point{X: l.Offset + l.Width/2},
			}
		} else {
			boxes[i].Y = plt.Y.Norm(cat)
			boxes[i].Rectangle = vg.Rectangle{
				Min: vg.Point{Y: l.Offset - l.Width/2},
				Max: vg.Point{Y: l.Offset + l.Width/2},
			}
		}
	}
	return boxes
}

// Legend creates a Legend plotter for this StringMap.
func (sm *StringMap) Legend(width vg.Length) *Legend {
	return &Legend{
		p:         sm.Colors,
		Width:     width,
		LineStyle: DefaultLineStyle,
	}
}

// Legend creates a Legend plotter for this IntMap.
func (im *IntMap) Legend(width vg.Length) *Legend {
	return &Legend{
		p:         im.Colors,
		Width:     width,
		LineStyle: DefaultLineStyle,
	}
}

// SetupPlot changes the default settings of p so that
// they are appropriate for plotting a legend.
func (sm *StringMap) SetupPlot(l *Legend, p *plot.Plot) {
	if !l.Horizontal {
		p.HideX()
		p.Y.Padding = 0
		p.NominalY(sm.Categories...)
	} else {
		p.HideY()
		p.X.Padding = 0
		p.NominalX(sm.Categories...)
	}
}

// SetupPlot changes the default settings of p so that
// they are appropriate for plotting a legend.
func (im *IntMap) SetupPlot(l *Legend, p *plot.Plot) {
	cats := make([]string, len(im.Categories))
	for i, c := range im.Categories {
		cats[i] = fmt.Sprintf("%d", c)
	}
	if !l.Horizontal {
		p.HideX()
		p.Y.Padding = 0
		p.NominalY(cats...)
	} else {
		p.HideY()
		p.X.Padding = 0
		p.NominalX(cats...)
	}
}
