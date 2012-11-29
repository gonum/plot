// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// plotter defines a variety of standard Plotters for the Plotinum
// plot package.		
//		
// Plotters use the primitives provided by the plot package to		
// draw to the data area of a plot.  This package provides		
// some standard data styles such as lines, scatter plots,		
// box plots, labels, and more.
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
		Width:    vg.Points(1),
		Dashes:   []vg.Length{},
		DashOffs: 0,
	}

	// DefaultGlyphStyle is the default style used
	// for gyph marks.
	DefaultGlyphStyle = plot.GlyphStyle{
		Color:  color.Black,
		Radius: vg.Points(2.5),
		Shape:  plot.RingGlyph{},
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
func CopyXYs(data XYer) XYs {
	cpy := make(XYs, data.Len())
	for i := range cpy {
		cpy[i].X, cpy[i].Y = data.XY(i)
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

// XYZer wraps the Len and XYZ methods.
type XYZer interface {
	// Len returns the number of x, y, z triples.
	Len() int

	// XYZ returns an x, y, z triple.
	XYZ(int) (float64, float64, float64)
}

// XYZs implements the XYZer interface using a slice.
type XYZs []struct{ X, Y, Z float64 }

// Len implements the Len method of the XYZer interface.
func (xyz XYZs) Len() int {
	return len(xyz)
}

// XYZ implements the XYZ method of the XYZer interface.
func (xyz XYZs) XYZ(i int) (float64, float64, float64) {
	return xyz[i].X, xyz[i].Y, xyz[i].Z
}

// CopyXYZs copies an XYZer.
func CopyXYZs(data XYZer) XYZs {
	cpy := make(XYZs, data.Len())
	for i := range cpy {
		cpy[i].X, cpy[i].Y, cpy[i].Z = data.XYZ(i)
	}
	return cpy
}

// XYValues implements the XYer interface, returning
// the x and y values from an XYZer.
type XYValues struct{ XYZer }

// XY implements the XY method of the XYer interface.
func (xy XYValues) XY(i int) (float64, float64) {
	x, y, _ := xy.XYZ(i)
	return x, y
}

// Labeller wraps the Label methods.
type Labeller interface {
	// Label returns a label.
	Label(int) string
}

// XErrorer wraps the XError method.
type XErrorer interface {
	// XError returns two error values for X data.
	XError(int) (float64, float64)
}

// Errors is a slice of low and high error values.
type Errors []struct{ Low, High float64 }

// XErrors implements the XErrorer interface.
type XErrors Errors

func (xe XErrors) XError(i int) (float64, float64) {
	return xe[i].Low, xe[i].High
}

// YErrorer wraps the YError method.
type YErrorer interface {
	// YError returns two error values for Y data.
	YError(int) (float64, float64)
}

// YErrors implements the YErrorer interface.
type YErrors Errors

func (ye YErrors) YError(i int) (float64, float64) {
	return ye[i].Low, ye[i].High
}
