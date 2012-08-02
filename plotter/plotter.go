// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

// plotter defines a variety of standard Plotters for the Plotinum plot package.
//
// Plotters use the primitives provided by the plot package to
// draw to the data area of a plot.  This package provides
// some standard data styles such as lines, scatter plots,
// box plots, error bars, and labels.
package plotter

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
	"image/color"
	"math"
)

var (
	// DefaultLineStyle is a reasonable default LineStyle
	// for drawing most lines in a plot.
	DefaultLineStyle = plot.LineStyle{
		Width: vg.Points(0.5),
		Color: color.Black,
	}

	// DefaultGlyhpStyle is a reasonable default GlyphStyle
	// for drawing points on a plot.
	DefaultGlyphStyle = plot.GlyphStyle{
		Radius: vg.Points(2),
		Color:  color.Black,
	}
)

const (
	// DefaultFont is the default font name.
	DefaultFont = "Times-Roman"
)

// An XYer wraps methods for getting a set of
// X and Y data values.
type XYer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// X returns an X value
	X(int) float64

	// Y returns a Y value
	Y(int) float64
}

// XYs is a slice of X, Y pairs, implementing the
// XYer interface.
type XYs []struct {
	X, Y float64
}

// Len returns the number of points.
func (p XYs) Len() int {
	return len(p)
}

// Less returns true if the ith X value is less than
// the jth X value.  This implements the Less
// method of sort.Interface, for sorting points by
// increasing X.
func (p XYs) Less(i, j int) bool {
	return p[i].X < p[j].X
}

// Swap swaps the ith and jth points.
func (p XYs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// X returns the ith X value.
func (p XYs) X(i int) float64 {
	return p[i].X
}

// Y returns the ith Y value.
func (p XYs) Y(i int) float64 {
	return p[i].Y
}

// xDataRange returns the minimum and maximum x
// values of all points from the XYer.
func xDataRange(xys XYer) (xmin, xmax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	for i := 0; i < xys.Len(); i++ {
		x := xys.X(i)
		xmin = math.Min(xmin, x)
		xmax = math.Max(xmax, x)
	}
	return
} // A Yer wraps methods for getting a set of Y data values.
type Yer interface {
	// Len returns the number of X and Y values
	// that are available.
	Len() int

	// Y returns a Y value
	Y(int) float64
}

// Ys is a slice of values, implementing the Yer
// interface.
type Ys []float64

// Len returns the number of values.
func (ys Ys) Len() int {
	return len(ys)
}

// Less returns true if the ith Y value is less than
// the jth Y value.
func (ys Ys) Less(i, j int) bool {
	return ys[i] < ys[j]
}

// Swap swaps the ith and jth values.
func (ys Ys) Swap(i, j int) {
	ys[i], ys[j] = ys[j], ys[i]
}

// Y returns the ith Y value.
func (ys Ys) Y(i int) float64 {
	return ys[i]
}

// ySorted implements sort.Interface, sorting a slice
// of indices for the given Yer.
type ySorter struct {
	Yer
	inds []int
}

// Len returns the number of indices.
func (y ySorter) Len() int {
	return len(y.inds)
}

// Less returns true if the Y value at index i
// is less than the Y value at index j.
func (y ySorter) Less(i, j int) bool {
	return y.Y(y.inds[i]) < y.Y(y.inds[j])
}

// Swap swaps the ith and jth indices.
func (y ySorter) Swap(i, j int) {
	y.inds[i], y.inds[j] = y.inds[j], y.inds[i]
}

// yDataRange returns the minimum and maximum x
// values of all points from the XYer.
func yDataRange(ys Yer) (ymin, ymax float64) {
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for i := 0; i < ys.Len(); i++ {
		y := ys.Y(i)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
	}
	return
}

// XYLabeller wraps both XYer and Labeller.
type XYLabeller interface {
	// XYer returns the XY point that is being labelled.
	XYer

	// Label returns the ith label text.
	Label(int) string
}

// XErrorer wraps the XError method.
type XErrorer interface {
	// XError returns the low and high X errors.
	// Both values are added to the corresponding
	// X value to compute the range of error
	// of the X value of the point, so most likely
	// the low value will be negative.
	XError(int) (float64, float64)
}

// YErrorer wraps the YError method.
type YErrorer interface {
	// YError is the same as the XError method
	// of the XErrorer interface, however it
	// applies to the Y values of points instead
	// of the X values.
	YError(int) (float64, float64)
}

// XYLabelErrors implements the XYer, XYLabeller, XErrorer,
// and YErrorer interfaces.
type XYLabelErrors struct {
	XYs
	Labels  []string
	XErrors []struct{ Low, High float64 }
	YErrors []struct{ Low, High float64 }
}

// MakeYXYLabelErrors returns a new XYLabelErrors
// of the given length.
func MakeXYLabelErrors(l int) XYLabelErrors {
	return XYLabelErrors{
		XYs:     make(XYs, l),
		Labels:  make([]string, l),
		XErrors: make([]struct{ Low, High float64 }, l),
		YErrors: make([]struct{ Low, High float64 }, l),
	}
}

// Label implements the XYLabeller interface.
func (xy XYLabelErrors) Label(i int) string {
	return xy.Labels[i]
}

// XError implements the XErrorer interface.
func (xy XYLabelErrors) XError(i int) (float64, float64) {
	return xy.XErrors[i].Low, xy.XErrors[i].High
}

// YError implements the YErrorer interface.
func (xy XYLabelErrors) YError(i int) (float64, float64) {
	return xy.YErrors[i].Low, xy.YErrors[i].High
}
