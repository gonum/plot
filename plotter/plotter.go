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
	"sort"
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

// A Valuer wraps methods for getting a set of data values.
type Valuer interface {
	// Len returns the number of values that are available.
	Len() int

	// Value returns a value
	Value(int) float64
}

// indexSorter implements sort.Interface, sorting a slice
// of indices for the given Valuer.
type indexSorter struct {
	Valuer
	inds []int
}

// Len returns the number of indices.
func (s indexSorter) Len() int {
	return len(s.inds)
}

// Less returns true if the value at index i
// is less than the value at index j.
func (s indexSorter) Less(i, j int) bool {
	return s.Value(s.inds[i]) < s.Value(s.inds[j])
}

// Swap swaps the ith and jth indices.
func (s indexSorter) Swap(i, j int) {
	s.inds[i], s.inds[j] = s.inds[j], s.inds[i]
}

func SortedIndices(vs Valuer) []int {
	sorted := make([]int, vs.Len())
	for i := range sorted {
		sorted[i] = i
	}
	sort.Sort(indexSorter{vs, sorted})
	return sorted
}

// Range returns the minimum and maximum
// values in a Valuer.
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

// Values is a slice of values, implementing the
// Valuer interface.
type Values []float64

// Len returns the number of values.
func (vs Values) Len() int {
	return len(vs)
}

// Less returns true if the ith value is less than
// the jth value.
func (vs Values) Less(i, j int) bool {
	return vs[i] < vs[j]
}

// Swap swaps the ith and jth values.
func (vs Values) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

// Value returns the ith value.
func (vs Values) Value(i int) float64 {
	return vs[i]
}

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

// XYRange returns the range of X and Y values.
func XYRange(xys XYer) (xmin, xmax, ymin, ymax float64) {
	xmin, xmax = Range(XValues{xys})
	ymin, ymax = Range(YValues{xys})
	return
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

// XValues is a Valuer that returns X values.
type XValues struct {
	XYer
}

func (vs XValues) Value(i int) float64 {
	return vs.X(i)
}

// YValues is a Vauler that returns Y values.
type YValues struct {
	XYer
}

func (vs YValues) Value(i int) float64 {
	return vs.Y(i)
}

// Labeller wraps the Labeller method.
type Labeller interface {
	// Len returns the number of labels.
	Len() int

	// Label returns the ith label text.
	Label(int) string
}

// XErrorer wraps the XError method.
type XErrorer interface {
	// Len returns the number of errors.
	Len() int

	// XError returns the low and high X errors.
	// Both values are added to the corresponding
	// X value to compute the range of error
	// of the X value of the point, so most likely
	// the low value will be negative.
	XError(int) (float64, float64)
}

// YErrorer wraps the YError method.
type YErrorer interface {
	// Len returns the number of errors.
	Len() int

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
