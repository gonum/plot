// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Function implements the Plotter interface, drawing a line
// for the given function.
type Function struct {
	F       func(float64) float64
	Samples int
	plot.LineStyle
}

// MakeFunction returns a Function that plots F using
// the default line style with 50 samples.
func MakeFunction(f func(float64) float64) Function {
	return Function{
		F:         f,
		Samples:   50,
		LineStyle: DefaultLineStyle,
	}
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (l Function) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)
	line := make([]plot.Point, l.Samples)
	d := (p.X.Max - p.X.Min) / float64(l.Samples-1)
	for i := 0; i < l.Samples; i++ {
		x := p.X.Min + float64(i)*d
		line[i].X = trX(x)
		line[i].Y = trY(l.F(x))
	}
	da.StrokeLines(l.LineStyle, da.ClipLinesXY(line)...)
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the function.
func (l Function) Thumbnail(da *plot.DrawArea) {
	da.StrokeLine2(l.LineStyle, da.Min.X, da.Center().Y, da.Max().X, da.Center().Y)
}
