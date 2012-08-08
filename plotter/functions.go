// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
)

// Function implements the Plotter interface,
// drawing a line for the given function.
type Function struct {
	F       func(float64) float64
	Samples int
	plot.LineStyle
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
func (f *Function) Plot(da plot.DrawArea, p *plot.Plot) {
	trX, trY := p.Transforms(&da)

	d := (p.X.Max - p.X.Min) / float64(f.Samples-1)
	line := make([]plot.Point, f.Samples)
	for i := range line {
		x := p.X.Min + float64(i)*d
		line[i].X = trX(x)
		line[i].Y = trY(f.F(x))
	}
	da.StrokeLines(f.LineStyle, da.ClipLinesXY(line)...)
}

// Thumbnail draws a line in the given style down the
// center of a DrawArea as a thumbnail representation
// of the LineStyle of the function.
func (f Function) Thumbnail(da *plot.DrawArea) {
	y := da.Center().Y
	da.StrokeLine2(f.LineStyle, da.Min.X, y, da.Max().X, y)
}
