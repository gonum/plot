// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
	"math"
)

var (
	// DefaultLineStyle is the default style for drawing
	// lines.
	DefaultLineStyle = plot.LineStyle{
		Color:    color.Black,
		Width:    vg.Points(0.5),
		Dashes:   []vg.Length{},
		DashOffs: 0,
	}

	// DefaultGlyphStyle is the default style used
	// for gyph marks.
	DefaultGlyphStyle = plot.GlyphStyle{
		Color:  color.Black,
		Shape:  plot.RingGlyph,
		Radius: vg.Points(2),
	}
)

// Valuer wraps the Len and Value methods.
type Valuer interface {
	// Len returns the number of values.
	Len() int

	// Value returns a value.
	Value(int) float64
}

// Range returns the minimum and maximum values.
func Range(vs Valuer) (min, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)
	for i := 0; i < vs.Len(); i++ {
		v := vs.Value(i)
		min = math.Min(min, v)
		max = math.Max(max, v)
	}
	return
}

// Values implements the Valuer interface.
type Values []float64

// CopyValues returns a Values that is a copy of the
// values from a Valuer.
func CopyValues(vs Valuer) Values {
	cpy := make(Values, vs.Len())
	for i := 0; i < vs.Len(); i++ {
		cpy[i] = vs.Value(i)
	}
	return cpy
}

func (vs Values) Len() int {
	return len(vs)
}

func (vs Values) Value(i int) float64 {
	return vs[i]
}

// XYer wraps the Len and XY methods.
type XYer interface {
	// Len returns the number of x, y pairs.
	Len() int

	// XY returns an x, y pair.
	XY(int) (x, y float64)
}

// XYRange returns the minimum and maximum
// x and y values.
func XYRange(xys XYer) (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = Range(XValues{xys})
	ymin, ymax = Range(YValues{xys})
	return
}

// XYs implements the XYer interface.
type XYs []struct{ X, Y float64 }

// CopyXYs returns an XYs that is a copy of the
// x and y values from an XYer.
func CopyXYs(xys XYer) XYs {
	cpy := make(XYs, xys.Len())
	for i := 0; i < xys.Len(); i++ {
		x, y := xys.XY(i)
		cpy[i].X = x
		cpy[i].Y = y
	}
	return cpy
}

func (xys XYs) Len() int {
	return len(xys)
}

func (xys XYs) XY(i int) (float64, float64) {
	return xys[i].X, xys[i].Y
}

// XValues implements the Valuer interface,
// returning the x value from an XYer.
type XValues struct {
	XYer
}

func (xs XValues) Value(i int) float64 {
	x, _ := xs.XY(i)
	return x
}

// YValues implements the Valuer interface,
// returning the y value from an XYer.
type YValues struct {
	XYer
}

func (ys YValues) Value(i int) float64 {
	_, y := ys.XY(i)
	return y
}

// Labeller wraps the Len and Label methods.
type Labeller interface {
	// Len returns the number of labels.
	Len() int

	// Label returns a label.
	Label(int) string
}

type ValueLabels []struct {
	Value float64
	Label string
}

func (vs ValueLabels) Len() int {
	return len(vs)
}

func (vs ValueLabels) Value(i int) float64 {
	return vs[i].Value
}

func (vs ValueLabels) Label(i int) string {
	return vs[i].Label
}

type XYLabels []struct {
	X, Y  float64
	Label string
}

func (xys XYLabels) Len() int {
	return len(xys)
}

func (xys XYLabels) XY(i int) (float64, float64) {
	return xys[i].X, xys[i].Y
}

func (xys XYLabels) Label(i int) string {
	return xys[i].Label
}
