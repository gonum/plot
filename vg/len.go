// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package vg

const (
	// ptInch is the number of points in an inch.
	ptInch = 72
)

// A Length is a unit-independent representation of length.
// Internally, the length is stored in postscript points.
type Length float64

// Points returns a length for the given number of points.
func Points(pt float64) Length {
	return Length(pt)
}

// Inches returns a length for the given numer of inches.
func Inches(in float64) Length {
	return Length(in * ptInch)
}

// Dots returns the length in dots for the given Canvas.
func (l Length) Dots(c Canvas) float64 {
	return float64(l) / ptInch * c.DPI()
}

// Points returns the length in postscript points.
func (l Length) Points() float64 {
	return float64(l)
}

// Inches returns the length in inches.
func (l Length) Inches() float64 {
	return float64(l) / ptInch
}
