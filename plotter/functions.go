// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// Function implements the Plotter interface,
// drawing a line for the given function.
type Function struct {
	F       func(float64) float64
	Samples int
	draw.LineStyle
}

// NewFunction returns a Function that plots F using
// the default line style with 50 samples.
func NewFunction(f func(float64) float64) *Function {
	return &Function{
		F:         f,
		Samples:   50,
		LineStyle: DefaultLineStyle,
	}
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (f *Function) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)

	d := (p.X.Max - p.X.Min) / float64(f.Samples-1)
	line := make([]vg.Point, f.Samples)
	for i := range line {
		x := p.X.Min + float64(i)*d
		line[i].X = trX(x)
		line[i].Y = trY(f.F(x))
	}

	// For every continuous block of non-NaN Y values, stroke lines
	for i := 0; i < len(line); i++ {
		if !math.IsNaN(float64(line[i].Y)) {
			j := i + 1
			for ; j < len(line); j++ {
				if math.IsNaN(float64(line[j].Y)) {
					break
				}
			}
			c.StrokeLines(f.LineStyle, c.ClipLinesXY(line[i:j])...)
			i = j
		}
	}
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the function.
func (f Function) Thumbnail(c *draw.Canvas) {
	y := c.Center().Y
	c.StrokeLine2(f.LineStyle, c.Min.X, y, c.Max.X, y)
}
